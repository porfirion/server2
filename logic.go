package main

import (
	"log"
	"math/rand"
	"time"
	"github.com/porfirion/server2/world"
	"github.com/porfirion/server2/network"
)

const (
	SendObjectsTimeout time.Duration = time.Second * 1
)



type Logic struct {
	IncomingMessages network.UserMessagesChannel
	OutgoingMessages network.ServerMessagesChannel
	Users            map[uint64]*network.User
	mWorldMap        *world.WorldMap
}

func (logic Logic) GetIncomingMessagesChannel() network.UserMessagesChannel {
	return logic.IncomingMessages
}
func (logic *Logic) SetIncomingMessagesChannel(channel network.UserMessagesChannel) {
	logic.IncomingMessages = channel
}
func (logic Logic) GetOutgoingMessagesChannel() network.ServerMessagesChannel {
	return logic.OutgoingMessages
}
func (logic *Logic) SetOutgoingMessagesChannel(channel network.ServerMessagesChannel) {
	logic.OutgoingMessages = channel
}

func (logic *Logic) GetUserList(exceptId uint64) []network.User {
	userlist := []network.User{}
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, network.User{Id: user.Id, Name: user.Name})
		}
	}

	log.Printf("Userlist: %#v\n", userlist)

	return userlist
}

/**
 * отправляет сообщение. Первый массив обозначает список целей кому передавать. Второй массив обозначает кому не передавать.
 * @param  {[type]} logic *Logic) SendMessage(msg Message, targets ...[]int [description]
 * @return {[type]} [description]
 */
func (logic *Logic) SendMessage(msg interface{}, targets ...network.UsersList) {
	serverMessage := network.ServerMessage{Data: msg}

	// real targets
	if len(targets) > 0 {
		serverMessage.Targets = targets[0]
	}

	// except this users
	if len(targets) > 1 {
		serverMessage.Except = targets[1]
	}

	logic.OutgoingMessages <- serverMessage
}

func (logic *Logic) SendTextMessage(text string, sender uint64) {
	logic.SendMessage(network.TextMessage{Text: text, Sender: sender})
}

func (logic *Logic) SendTextMessageToUser(text string, sender uint64, userId uint64) {
	logic.SendMessage(network.TextMessage{Text: text, Sender: sender}, network.UsersList{userId})
}

func (logic *Logic) AddUser(id uint64, name string) *network.User {
	user := &network.User{Id: id, Name: name}
	logic.Users[id] = user

	pos := world.Point2D{X: rand.Float64()*1000 - 500, Y: rand.Float64()*1000 - 500}
	logic.mWorldMap.AddUser(user.Id, pos)

	return user
}

func (logic *Logic) RemoveUser(id uint64) {
	delete(logic.Users, id)
	logic.mWorldMap.RemoveUser(id)
}

func (logic *Logic) sendSyncMessage() {
	var currentTime int64 = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond);
	logic.SendMessage(network.SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
}

// Возвращает true, если нужно синхронизировать положение объектов заново
func (logic *Logic) ProcessActionMessage(userId uint64, msg *network.ActionMessage) (needSync bool) {
	log.Println("UNIMPLEMENTED!")
	needSync = false
	switch msg.ActionType {
	case "move":
		user := logic.Users[userId]
		userObject := logic.mWorldMap.UsersObjects[userId]

		x, okX := msg.ActionData["x"].(float64)
		y, okY := msg.ActionData["y"].(float64)

		if okX && okY {
			userObject.StartMoveTo(world.Point2D{X: x, Y: y})
			log.Printf("user #%d try to move it's object #%d to (%f:%f)\n", userId, userObject.Id, x, y)
		} else {
			log.Println("can't get x and y")
		}

		log.Println(user.Name+" is moving to", msg.ActionData)
		needSync = true
	default:
		log.Println("Unknown action type: ", msg.ActionType)
	}

	return
}

func (logic *Logic) ProcessMessage(message network.UserMessage) (needSync bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in %#v\n", r)
		}
	}()

	needSync = false

	switch msg := message.Data.(type) {
	case *network.DataMessage:
		log.Println("Logic: Data message received: ", message)
	case *network.TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case *network.LoginMessage:
		log.Println("Logic: Login message received")

		user := logic.AddUser(msg.Id, msg.Name)
		logic.SendTextMessageToUser("Logic: Wellcome, "+user.Name, 0, user.Id)

		logic.SendMessage(network.UserLoggedinMessage{Id: user.Id, Name: user.Name}, network.UsersList{}, network.UsersList{user.Id})

		logic.SendMessage(network.UserListMessage{logic.GetUserList(user.Id)}, network.UsersList{user.Id})
		log.Println("sent. sync next")
		var currentTime int64 = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond);
		logic.SendMessage(network.SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
	case *network.LogoutMessage:
		log.Println("Logic: Logout message", msg.Id)
		logic.RemoveUser(msg.Id)
		logic.SendMessage(network.UserLoggedoutMessage{Id: msg.Id})
	case *network.ActionMessage:
		needSync = logic.ProcessActionMessage(message.Source, msg)
	default:
		log.Printf("Logic: Unknown message type %#v from %d\n", message.Data, message.Source)
	}

	return
}

func (logic *Logic) Start() {
	rand.Seed(int64(time.Now().Nanosecond()))

	logic.Users = make(map[uint64]*network.User)

	logic.mWorldMap = world.NewWorldMap()

	// стартуем симуляцию
	logic.mWorldMap.ProcessSimulationStep()
	var simulationTimer *time.Timer = time.NewTimer(logic.mWorldMap.TimeToNextStep())
	var sendTimer *time.Timer = time.NewTimer(time.Second * 0)

	log.Println("Logic: started")
	for {
		select {
		case _ = <-simulationTimer.C:
			//log.Println("Logic: simulation step")

			// пока есть что симулировать - симулируем
			simulated := true
			var changed, globallyChanged bool = false, false;
			for simulated {
				// пока нам говорят, что симуляция прошла успешно - делаем очередной шаг
				simulated, changed = logic.mWorldMap.ProcessSimulationStep()
				// при это запоминаем, не поменялось ли чего
				globallyChanged = globallyChanged || changed
			}

			if (globallyChanged) {
				// если что-то изменилось, надо об этом всем рассказать
				logic.sendSyncMessage()
			}

			simulationTimer.Reset(logic.mWorldMap.TimeToNextStep())
		case _ = <-sendTimer.C:
			// дополнительно рассылаем всем уведомления по таймеру
			// по идее потом это можно будет убрать
			logic.sendSyncMessage()
			sendTimer.Reset(SendObjectsTimeout)
		case msg := <-logic.IncomingMessages:
			log.Println("Logic: message received")
			if needSync := logic.ProcessMessage(msg); needSync {
				logic.sendSyncMessage()
			}
			log.Println("Logic: message processed")
		}
	}

	log.Println("Logic: finished")
}

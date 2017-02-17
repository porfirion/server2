package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	SendObjectsTimeout time.Duration = time.Second * 1
)

type LogicInterface interface {
	getIncomingMessagesChannel() UserMessagesChannel
	setIncomingMessagesChannel(channel UserMessagesChannel)
	getOutgoingMessagesChannel() ServerMessagesChannel
	setOutgoingMessagesChannel(channel ServerMessagesChannel)
	Start()
}

type Logic struct {
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[uint64]*User
	mWorldMap        *WorldMap
}

func (logic *Logic) getIncomingMessagesChannel() UserMessagesChannel {
	return logic.IncomingMessages
}
func (logic *Logic) setIncomingMessagesChannel(channel UserMessagesChannel) {
	logic.IncomingMessages = channel
}
func (logic *Logic) getOutgoingMessagesChannel() ServerMessagesChannel {
	return logic.OutgoingMessages
}
func (logic *Logic) setOutgoingMessagesChannel(channel ServerMessagesChannel) {
	logic.OutgoingMessages = channel
}

func (logic *Logic) GetUserList(exceptId uint64) []User {
	userlist := []User{}
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, User{Id: user.Id, Name: user.Name})
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
func (logic *Logic) SendMessage(msg interface{}, targets ...UsersList) {
	serverMessage := ServerMessage{Data: msg}

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
	logic.SendMessage(TextMessage{Text: text, Sender: sender})
}

func (logic *Logic) SendTextMessageToUser(text string, sender uint64, userId uint64) {
	logic.SendMessage(TextMessage{Text: text, Sender: sender}, UsersList{userId})
}

func (logic *Logic) AddUser(id uint64, name string) *User {
	user := &User{Id: id, Name: name}
	logic.Users[id] = user

	pos := Position{X: rand.Float64()*1000 - 500, Y: rand.Float64()*1000 - 500}
	logic.mWorldMap.AddUser(user, pos)

	return user
}

func (logic *Logic) RemoveUser(id uint64) {
	delete(logic.Users, id)
	logic.mWorldMap.RemoveUser(id)
}

func (logic *Logic) ProcessActionMessage(userId uint64, msg *ActionMessage) {
	log.Println("UNIMPLEMENTED!")

	switch msg.ActionType {
	case "move":
		user := logic.Users[userId]
		userObject := logic.mWorldMap.UsersObjects[userId]

		x, okX := msg.ActionData["x"].(float64)
		y, okY := msg.ActionData["y"].(float64)

		if okX && okY {
			userObject.MoveTo(Position{X: x, Y: y})
			log.Printf("user #%d try to move it's object #%d to (%f:%f)\n", userId, userObject.Id, x, y)
		} else {
			log.Println("can't get x and y")
		}

		log.Println(user.Name+" is moving to", msg.ActionData)
	default:
		log.Println("Unknown action type: ", msg.ActionType)
	}
}

func (logic *Logic) ProcessMessage(message UserMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in %#v\n", r)
		}
	}()

	switch msg := message.Data.(type) {
	case *DataMessage:
		log.Println("Logic: Data message received: ", message)
	case *TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case *LoginMessage:
		log.Println("Logic: Login message received")

		user := logic.AddUser(msg.Id, msg.Name)
		logic.SendTextMessageToUser("Logic: Wellcome, "+user.Name, 0, user.Id)

		logic.SendMessage(UserLoggedinMessage{Id: user.Id, Name: user.Name}, UsersList{}, UsersList{user.Id})

		logic.SendMessage(UserListMessage{logic.GetUserList(user.Id)}, UsersList{user.Id})
		log.Println("sent. sync next")
		var currentTime int64 = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond);
		logic.SendMessage(SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
	case *LogoutMessage:
		log.Println("Logic: Logout message", msg.Id)
		logic.RemoveUser(msg.Id)
		logic.SendMessage(UserLoggedoutMessage{Id: msg.Id})
	case *ActionMessage:
		logic.ProcessActionMessage(message.Source, msg)
	default:
		log.Printf("Logic: Unknown message type %#v from %d\n", message.Data, message.Source)
	}
}

func (logic *Logic) Start() {
	rand.Seed(int64(time.Now().Nanosecond()))

	logic.Users = make(map[uint64]*User)

	logic.mWorldMap = NewWorldMap()

	// стартуем симуляцию
	logic.mWorldMap.ProcessSimulationStep()
	var simulationTimer *time.Timer = time.NewTimer(logic.mWorldMap.TimeToNextStep())
	var sendTimer *time.Timer = time.NewTimer(time.Second * 0)

	log.Println("Logic: started")
	for {
		select {
		case _ = <-simulationTimer.C:
			log.Println("Logic: simulation step")

			// пока есть что симулировать - симулируем
			for logic.mWorldMap.ProcessSimulationStep() {
			}

			simulationTimer.Reset(logic.mWorldMap.TimeToNextStep())
		case _ = <-sendTimer.C:
			log.Println("")
			var currentTime int64 = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond);
			logic.SendMessage(SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
			sendTimer.Reset(SendObjectsTimeout)
		case msg := <-logic.IncomingMessages:
			log.Println("Logic: message received")
			logic.ProcessMessage(msg)
			log.Println("Logic: message processed")
		}
	}

	log.Println("Logic: finished")
}

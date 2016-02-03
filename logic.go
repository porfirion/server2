package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

var SimulationTime int

type Logic struct {
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[int]*User
	UsersPositions   map[int]Position
	UsersObjects     map[int]*MapObject
	EventDispatcher  *EventDispatcher
	Objects          []MapObject
}

func (logic *Logic) GetUserList(exceptId int) []User {
	userlist := []User{}
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, User{Id: user.Id, Name: user.Name})
		}
	}

	log.Printf("Userlist: %#v\n", userlist)

	return userlist
}

func (logic *Logic) GetUsersPositions() map[string]Position {
	res := make(map[string]Position)
	for id, pos := range logic.UsersPositions {
		res[strconv.Itoa(id)] = pos
	}

	log.Printf("Users positions: %#v\n", res)

	return res
}

/**
 * отправляет сообщение. Первый массив обозначает список целей кому передавать. Второй массив обозначает кому не передавать.
 * @param  {[type]} logic *Logic) SendMessage(msg Message, targets ...[]int [description]
 * @return {[type]} [description]
 */
func (logic *Logic) SendMessage(msg interface{}, targets ...[]int) {
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

func (logic *Logic) SendTextMessage(text string, sender int) {
	logic.SendMessage(TextMessage{Text: text, Sender: sender})
}

func (logic *Logic) SendTextMessageToUser(text string, sender int, userId int) {
	logic.SendMessage(TextMessage{Text: text, Sender: sender}, []int{userId})
}

func (logic *Logic) AddUser(id int, name string) *User {
	user := &User{Id: id, Name: name}
	pos := Position{X: rand.Int63n(1000) - int64(500), Y: rand.Int63n(1000) - int64(500)}
	obj := &MapObject{Pos: pos, User: user, Speed: 0.0}
	logic.Users[id] = user
	logic.UsersPositions[id] = pos
	logic.UsersObjects[id] = obj

	return user
}

func (logic *Logic) ProcessActionMessage(userId int, msg *ActionMessage) {
	switch msg.ActionType {
	case "move":
		user := logic.Users[userId]
		obj := logic.UsersObjects[userId]

		x, okX := msg.ActionData["x"].(int64)
		y, okY := msg.ActionData["y"].(int64)
		if okX && okY {
			obj.AdjustPosition()
			obj.MoveTo(Position{X: x, Y: y}, 1.0)
		} else {
			log.Println("can't get x and y")
		}

		log.Println(user.Name+" is moving to", msg.ActionData)
	default:
		log.Println("Unknown action type: ", msg.ActionType)
	}

	log.Println("UNIMPLEMENTED!")
}

func (logic *Logic) RemoveUser(id int) {
	delete(logic.Users, id)
	delete(logic.UsersPositions, id)
	delete(logic.UsersObjects, id)
}

func (logic *Logic) ProcessMessage(message UserMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
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

		logic.SendMessage(UserLoggedinMessage{Id: user.Id, Name: user.Name}, []int{}, []int{user.Id})

		logic.SendMessage(UserListMessage{logic.GetUserList(user.Id)}, []int{user.Id})
		logic.SendMessage(SyncPositionsMessage{logic.GetUsersPositions()})
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
	logic.EventDispatcher = &EventDispatcher{}
	logic.EventDispatcher.Init()

	logic.Users = make(map[int]*User)
	logic.UsersPositions = make(map[int]Position)

	log.Println("Logic: started")
	for {
		select {
		case msg := <-logic.IncomingMessages:
			log.Println("Logic: message received")
			logic.ProcessMessage(msg)
			log.Println("Logic: message processed")
		}
	}

	log.Println("Logic: finished")
}

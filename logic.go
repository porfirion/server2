package main

import (
	"fmt"
	"log"
	"strconv"
)

type Logic struct {
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[int]*User
	UsersPositions   map[int]Position
	EventDispatcher  *EventDispatcher
}

func (logic *Logic) GetUserList(exceptId int) []User {
	userlist := []User{}
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, User{Id: user.Id, Name: user.Name})
		}
	}

	log.Println(fmt.Sprintf("Userlist: %#v", userlist))

	return userlist
}

func (logic *Logic) GetUsersPositions() map[string]Position {
	res := make(map[string]Position)
	for id, pos := range logic.UsersPositions {
		res[strconv.Itoa(id)] = pos
	}

	log.Println(fmt.Sprintf("Users positions: %#v", res))

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
	logic.Users[id] = user
	logic.UsersPositions[id] = Position{X: 0, Y: 0}
	return user
}

func (logic *Logic) ProcessMessage(message UserMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	switch msg := message.Data.(type) {
	case *DataMessage:
		log.Println("Data message received: ", message)
	case *TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case *LoginMessage:
		log.Println("Login message received")

		user := logic.AddUser(msg.Id, msg.Name)
		logic.SendTextMessageToUser("Wellcome, "+user.Name, 0, user.Id)

		logic.SendMessage(UserLoggedinMessage{Id: user.Id, Name: user.Name}, []int{}, []int{user.Id})
		logic.SendMessage(UserListMessage{logic.GetUserList(user.Id)}, []int{user.Id})
		logic.SendMessage(SyncPositionsMessage{logic.GetUsersPositions()})
	case *LogoutMessage:
		log.Println("Logout message", msg.Id)
		delete(logic.Users, msg.Id)
		logic.SendMessage(UserLoggedoutMessage{Id: msg.Id})
	default:
		fmt.Printf("Logic: Unknown message type %#v\n", message)
	}
}

func (logic *Logic) Start() {
	logic.EventDispatcher = &EventDispatcher{}
	logic.EventDispatcher.Init()

	logic.Users = make(map[int]*User)
	logic.UsersPositions = make(map[int]Position)

	log.Println("Logic started")
	select {
	case msg := <-logic.IncomingMessages:
		logic.ProcessMessage(msg)
	}
}

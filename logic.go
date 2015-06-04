package main

import (
	"fmt"
	"log"
)

type Logic struct {
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[int]*User
}

func (logic *Logic) GetUserList(exceptId int) []struct {
	Id   int
	Name string
} {
	userlist := []struct {
		Id   int
		Name string
	}{}
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, struct {
				Id   int
				Name string
			}{Id: user.Id, Name: user.Name})
		}
	}

	return userlist
}

// отправляет сообщение всем
func (logic *Logic) SendMessage(msg Message, targets ...[]int) {
	serverMessage := ServerMessage{Data: msg}
	if len(targets) > 0 {
		serverMessage.Targets = targets[0]
	}
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

func (logic *Logic) ProcessMessage(message UserMessage) {
	switch msg := message.Data.(type) {
	case DataMessage:
		log.Println("Data message received: ", message)
	case TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case LoginMessage:
		log.Println("Login message received")

		user := &User{Id: msg.Id, Name: msg.Name}
		logic.Users[msg.Id] = user
		logic.SendTextMessageToUser("Wellcome, "+user.Name, 0, user.Id)

		logic.SendMessage(UserLoggedinMessage{Id: user.Id, Name: user.Name}, []int{}, []int{user.Id})
		logic.SendMessage(UserListMessage{logic.GetUserList(user.Id)}, []int{user.Id})
	case LogoutMessage:
		log.Println("Logout message", msg.Id)
		delete(logic.Users, msg.Id)
		logic.SendMessage(UserLoggedoutMessage{Id: msg.Id})
	default:
		fmt.Printf("Unknown message type %#v\n", message)
	}
}

func (logic *Logic) Start() {
	logic.Users = make(map[int]*User)
	log.Println("Logic started")
	for msg := range logic.IncomingMessages {
		logic.ProcessMessage(msg)
	}
}

package main

import (
	"log"
)

type Logic struct {
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[int]User
}

// отправляет сообщение всем
func (logic *Logic) SendMessage(msg Message, targets []int) {
	logic.OutgoingMessages <- ServerMessage{Data: msg, Targets: targets}
}

func (logic *Logic) SendTextMessage(text string, sender int, targets []int) {
	logic.SendMessage(TextMessage{Text: text, Sender: sender}, targets)
}

func (logic *Logic) ProcessMessage(message UserMessage) {
	switch msg := message.Data.(type) {
	case DataMessage:
		log.Println("Data message received: ", message)
	case TextMessage:
		log.Println("Text message received: ", message)
		logic.SendTextMessage("User "+logic.Users[message.Source].Name+" says "+msg.Text, logic.Users[message.Source].Id, []int{})
	case LoginMessage:
		log.Println("Login message received", msg.User)

		logic.Users[msg.User.Id] = msg.User

		logic.SendTextMessage("Logged "+msg.Name, 0, []int{})
		logic.SendTextMessage("Wellcome, "+msg.Name, 0, []int{message.Source})
	default:
		log.Println("Unknown message type")
	}
}

func (logic *Logic) Start() {
	logic.Users = make(map[int]User)
	log.Println("Logic started")
	for msg := range logic.IncomingMessages {
		logic.ProcessMessage(msg)
	}
}

package main

import (
	"fmt"
	"log"
)

type Logic struct {
	IncomingMessages MessagesChannel
}

// отправляет сообщение всем
func (logic *Logic) SendMessageAll(msg Message) {
	// for _, ch := range logic.chanByUserId {
	// 	ch <- msg
	// }
}

// отправляет сообщение только одному адресату
func (logic *Logic) SendMessageId(msg Message, connId uint64) {
	// logic.chanByUserId[connId] <- msg
}

// отпрвляет сообщение нескольких определённым адресатам
func (logic *Logic) SendMessageMultiple(msg Message, targets []uint64) {

}

// отправлет сообщение всем кроме
func (logic *Logic) SendMessageExcept(msg Message, unwanted []uint64) {

}

func (logic *Logic) ProcessMessage(message Message) {
	fmt.Println("Processing message")
	fmt.Printf(format, message)
	switch msg := message.(type) {
	case DataMessage:
		fmt.Println("Data message received", msg.Data)
	case TextMessage:
		fmt.Println("Text message received", msg.Text)
	case AuthMessage:
		fmt.Println("Auth message received", msg.Uuid)
	default:
		log.Println("Unknown message type")
	}
}

func (logic *Logic) run() {
	for {
		msg := <-logic.IncomingMessages
		log.Println("Message received!")
		logic.ProcessMessage(msg)
	}
}

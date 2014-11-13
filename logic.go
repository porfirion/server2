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
	switch msg := message.(type) {
	case JoinMessage:

	case LeaveMessage:
		// delete(logic.chanByUserId, msg.GetUserId())
	case PingMessage:

	case TextMessage:
		fmt.Println(msg.text)

	default:
		log.Println("Unknown message type")
	}
}

func (logic *Logic) run() {
	select {
	case msg := <-logic.IncomingMessages:
		log.Println("Message received!")
		logic.ProcessMessage(msg)
	}
}

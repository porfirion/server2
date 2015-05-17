package main

import (
	"log"
)

type Logic struct {
	IncomingMessages MessagesChannel
	OutgoingMessages MessagesChannel
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
	case DataMessage:
		log.Println("Data message received", msg.Data)
	case TextMessage:
		log.Println("Text message received", msg.Text)
	case AuthMessage:
		log.Println("Auth message received", msg.Uuid)
		logic.OutgoingMessages <- &TextMessage{Text: "hello!"}
	default:
		log.Println("Unknown message type")
	}
}

func (logic *Logic) Start() {
	log.Println("Logic started")
	for msg := range logic.IncomingMessages {
		logic.ProcessMessage(msg)
	}
}

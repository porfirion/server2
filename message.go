package main

type Message interface {
}

type TextMessage struct {
	Text string
}

type DataMessage struct {
	Data []byte
}

type AuthMessage struct {
	Uuid string
}

type ServerMessage struct {
	Targets []int
	Data    Message
}

type MessagesChannel chan Message

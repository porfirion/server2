package main

type Message interface {
	GetUserId() uint64
	GetResponseChannel() MessagesChannel
}

type LoginMessage struct {
	Message
}

type LogoutMessage struct {
	Message
}

type Ping struct {
	Message
}

type Pong struct {
	Message
}

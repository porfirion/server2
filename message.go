package main

type Message interface {
	GetUserId() uint64
	GetResponseChannel() MessageChannel
}

type MessageChannel chan Message

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

package main

import (
	"time"
)

type Message struct {
	Sender       chan Message
	TimeReceived time.Time
	Recipient    int64
}

// ???
type Response struct {
}

type LoginMessage struct {
	Message
	ClientId        int64
	ResponseChannel chan Message
}

type LogoutMessage struct {
	Message
	ClientId int64
}

type Ping struct {
	Message
}
type Pong struct {
	Message
}

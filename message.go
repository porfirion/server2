package main

type Message interface {
	GetUserId() uint64
	GetResponseChannel() MessagesChannel
	SetResponseChannel(MessagesChannel)
}

type BaseMessage struct {
	userId          uint64
	responseChannel MessagesChannel
}

func (msg *BaseMessage) GetUserId() uint64 {
	return msg.userId
}

func (msg *BaseMessage) GetResponseChannel() MessagesChannel {
	return msg.responseChannel
}

func (msg *BaseMessage) SetResponseChannel(channel MessagesChannel) {
	msg.responseChannel = channel
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

type DataMessage struct {
	Message
	data []byte
}

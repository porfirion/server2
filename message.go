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

type JoinMessage struct {
	Message
}

type LeaveMessage struct {
	Message
}

type PingMessage struct {
	Message
}

type PongMessage struct {
	Message
}

type TextMessage struct {
	Message
	text string
}

type DataMessage struct {
	Message
	data []byte
}

type AuthMessage struct {
	Message
	name string
}

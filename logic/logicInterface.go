package logic

type LogicInterface interface {
	GetIncomingMessagesChannel() UserMessagesChannel
	GetOutgoingMessagesChannel() ServerMessagesChannel
}

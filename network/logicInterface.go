package network

type LogicInterface interface {
	GetIncomingMessagesChannel() UserMessagesChannel
	GetOutgoingMessagesChannel() ServerMessagesChannel
}

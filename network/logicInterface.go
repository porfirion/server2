package network

type LogicInterface interface {
	GetIncomingMessagesChannel() UserMessagesChannel
	SetIncomingMessagesChannel(channel UserMessagesChannel)
	GetOutgoingMessagesChannel() ServerMessagesChannel
	SetOutgoingMessagesChannel(channel ServerMessagesChannel)
}

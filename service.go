package main

import "github.com/porfirion/server2/network"

type Service interface {
	GetRequiredMessagesTypes() []int
	GetIncomingChannel() network.UserMessagesChannel
	GetOutgoingChannel() network.ServerMessagesChannel
}

type MessageBroker struct {
	Services []Service
	IncomingMessagesChannel network.UserMessagesChannel
	OutgoingMessagesChannel network.ServerMessagesChannel
}
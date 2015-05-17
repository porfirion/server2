package main

type Connection interface {
	StartReading(ch MessagesChannel)
	Close()
	GetResponseChannel() MessagesChannel
}

type ConnectionsChannel chan Connection

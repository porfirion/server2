package main

type Gate interface {
	Start(ConnectionsChannel, MessagesChannel)
}

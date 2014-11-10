package main

import (
	"fmt"
)

type ConnectionsPool struct {
	peer2sserver        MessagesChannel    // сообщения, которые отправляются на обработку в сервер
	server2peer         MessagesChannel    // сообщения, которые приходят на отправку из сервера
	incomingConnections ConnectionsChannel // входящие соединения
}

func (pool *ConnectionsPool) processConnection(connection Connection) {

}

func (pool *ConnectionsPool) start(incomingConnections ConnectionsChannel) {
	pool.incomingConnections = incomingConnections

	select {
	case connection := <-incomingConnections:
		fmt.Println("Connection received", connection)
	}
}

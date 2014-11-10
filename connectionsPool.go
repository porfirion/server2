package main

import (
	"log"
)

type ConnectionsPool struct {
	logic               *Logic
	incomingConnections ConnectionsChannel // входящие соединения
	connections         []Connection
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	pool.connections = append(pool.connections, connection)
	connection.StartReading(pool.logic.IncomingMessages)
}

func (pool *ConnectionsPool) Start() {
	pool.connections = make([]Connection, 100)

	select {
	case connection := <-pool.incomingConnections:
		log.Println("Connection received: ", connection)

		pool.processConnection(connection)
	}
}

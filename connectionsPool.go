package main

import (
	"log"
)

type ConnectionsPool struct {
	logic               *Logic
	incomingConnections ConnectionsChannel // входящие соединения
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	log.Println("Receiving auth")

	go connection.StartReading(pool.logic.IncomingMessages)
}

func (pool *ConnectionsPool) Start() {
	for {
		connection := <-pool.incomingConnections
		log.Println("Connection received: ", connection)

		pool.processConnection(connection)
	}

	log.Println("Connections pool finished")
}

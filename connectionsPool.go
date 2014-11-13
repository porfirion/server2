package main

import (
	"log"
)

type ConnectionsPool struct {
	logic               *Logic
	incomingConnections ConnectionsChannel // входящие соединения
	players             []*Player
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	newPlayer := connection.GetAuth()
	if newPlayer != nil {
		pool.players = append(pool.players, newPlayer)

		go connection.StartReading(pool.logic.IncomingMessages)
	} else {
		log.Println("error receiving auth. closing connection")
		connection.Close()
	}
}

func (pool *ConnectionsPool) Start() {
	pool.players = make([]*Player, 100)

	select {
	case connection := <-pool.incomingConnections:
		log.Println("Connection received: ", connection)

		pool.processConnection(connection)
	}
}

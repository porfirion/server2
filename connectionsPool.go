package main

import (
	"log"
)

type ConnectionsPool struct {
	logic                 *Logic
	incomingConnections   ConnectionsChannel // входящие соединения
	ConnectionsEnumerator chan int
	Connections           map[int]Connection
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	go connection.StartReading(pool.logic.IncomingMessages)

	var connectionId = <-pool.ConnectionsEnumerator
	pool.Connections[connectionId] = connection
}

func (pool *ConnectionsPool) InitEnumerator() {
	pool.ConnectionsEnumerator = make(chan int, 1)

	go func() {
		var connectionId int = 1
		for {
			pool.ConnectionsEnumerator <- connectionId
			connectionId++
		}
	}()
}

func (pool *ConnectionsPool) Start() {
	log.Println("Connections pool started")

	pool.InitEnumerator()

	pool.Connections = make(map[int]Connection)

	for {
		select {
		case connection := <-pool.incomingConnections:
			log.Println("CPool: connection received", connection)

			pool.processConnection(connection)
		case message := <-pool.logic.OutgoingMessages:
			for _, conn := range pool.Connections {
				conn.GetResponseChannel() <- message
			}
			log.Println("Outgoing message", message)
		}
	}

	log.Println("Connections pool finished")
}

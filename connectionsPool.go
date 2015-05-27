package main

import (
	"log"
)

type ConnectionsPool struct {
	logic                 *Logic
	incomingConnections   ConnectionsChannel // входящие соединения
	ConnectionsEnumerator chan int
	Connections           map[int]Connection
	ClosingChannel        chan int
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	connection.StartReading(pool.logic.IncomingMessages)

	var connectionId = <-pool.ConnectionsEnumerator
	pool.Connections[connectionId] = connection
	connection.SetId(connectionId)
	connection.SetClosingChannel(pool.ClosingChannel)
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

func (pool *ConnectionsPool) RemoveConnection(connectionId int) {
	log.Println("Removing connection", connectionId)
	delete(pool.Connections, connectionId)
}

func (pool *ConnectionsPool) DispathMessage(msg *ServerMessage) {
	if len(msg.Targets) == 0 {
		for _, conn := range pool.Connections {
			conn.GetResponseChannel() <- msg.Data
		}
	} else {
		for _, connectionId := range msg.Targets {
			pool.Connections[connectionId].GetResponseChannel() <- msg.Data
		}
	}
}

func (pool *ConnectionsPool) Start() {
	log.Println("Connections pool started")

	pool.InitEnumerator()

	pool.Connections = make(map[int]Connection)
	pool.ClosingChannel = make(chan int)

	for {
		select {
		case connection := <-pool.incomingConnections:
			log.Println("CPool: connection received", connection)

			pool.processConnection(connection)
		case message := <-pool.logic.OutgoingMessages:
			pool.DispathMessage(message)
			log.Println("Outgoing message", message)
		case connectionId := <-pool.ClosingChannel:
			pool.RemoveConnection(connectionId)
		}
	}

	log.Println("Connections pool finished")
}

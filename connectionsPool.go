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
	go func() {
		authMessage, err := connection.GetAuth()
		if err != nil {
			log.Println("Error authorization", err)
			/**
			 * TODO по идее нужно прибить это соединение, а то оно так и будет крутиться
			 * Также неплохо бы отправить ответ в соединение, что не прошла авторизация по таким-то причинам
			 */
			connection.GetResponseChannel() <- ErrorMessage{Code: 0, Description: "Authorization failed"}
		} else {
			log.Println("Authorization successful: ", authMessage)

			var connectionId = <-pool.ConnectionsEnumerator

			connection.SetId(connectionId)
			connection.SetClosingChannel(pool.ClosingChannel)

			// извещаем клиента о том, что он подключился
			connection.GetResponseChannel() <- WellcomeMessage{Id: connectionId}

			pool.Connections[connectionId] = connection
			pool.logic.IncomingMessages <- UserMessage{Data: &LoginMessage{Id: connectionId, Name: authMessage.Name}, Source: connectionId}

			connection.StartReading(pool.logic.IncomingMessages)
		}
	}()
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

	pool.logic.IncomingMessages <- UserMessage{Data: &LogoutMessage{connectionId}, Source: connectionId}
}

func (pool *ConnectionsPool) DispathMessage(msg ServerMessage) {
	if len(msg.Targets) == 0 {
	AllConnectionsLoop:
		for _, conn := range pool.Connections {
			if len(msg.Except) > 0 {
				for _, id := range msg.Except {
					if conn.GetId() == id {
						continue AllConnectionsLoop
					}
				}
			}
			conn.GetResponseChannel() <- msg.Data
		}
	} else {
	TargetConnectionsLoop:
		for _, connectionId := range msg.Targets {
			if len(msg.Except) > 0 {
				for _, id := range msg.Except {
					if connectionId == id {
						continue TargetConnectionsLoop
					}
				}
			}
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
			// log.Println("Outgoing message", message)
		case connectionId := <-pool.ClosingChannel:
			log.Println("Closing connection", connectionId)
			pool.RemoveConnection(connectionId)
		}
	}

	log.Println("Connections pool finished")
}

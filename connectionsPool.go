package main

import (
	"log"
)

type ConnectionsPool struct {
	logic                 *Logic
	incomingConnections   ConnectionsChannel // входящие соединения
	ConnectionsEnumerator chan uint64
	Connections           map[uint64]Connection
	ClosingChannel        chan uint64
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
			//log.Println("Authorization successful: ", authMessage)

			var connectionId uint64 = <-pool.ConnectionsEnumerator

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
	pool.ConnectionsEnumerator = make(chan uint64, 1)

	go func() {
		var connectionId uint64 = 1
		for {
			pool.ConnectionsEnumerator <- connectionId
			connectionId++
		}
	}()
}

func (pool *ConnectionsPool) RemoveConnection(connectionId uint64) {
	log.Println("CPool: Removing connection", connectionId)

	delete(pool.Connections, connectionId)

	log.Println("CPool: sending message to logic")
	pool.logic.IncomingMessages <- UserMessage{Data: &LogoutMessage{connectionId}, Source: connectionId}
	log.Println("CPool: message in sent to logic")
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

	pool.Connections = make(map[uint64]Connection)
	pool.ClosingChannel = make(chan uint64)

	for {
		select {
		case connection := <-pool.incomingConnections:
			log.Println("CPool: connection received", connection)
			pool.processConnection(connection)
			log.Println("CPool connection processed")
		case message := <-pool.logic.OutgoingMessages:
			log.Printf("CPool: Outgoing message %T\n", message)
			pool.DispathMessage(message)
			log.Println("CPool: Message is sent")
		case connectionId := <-pool.ClosingChannel:
			log.Println("CPool: Closing connection", connectionId)
			pool.RemoveConnection(connectionId)
			log.Println("CPool: connection closed")
		}
	}

	log.Println("Connections pool finished")
}

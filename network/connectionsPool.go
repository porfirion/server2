package network

import (
	"log"
)

type SearchableArray []uint64

func (arr SearchableArray) indexOf(value uint64) (bool, int) {
	if len(arr) == 0 {
		return false, -1
	}
	for ind, val := range arr {
		if val == value {
			return true, ind
		}
	}
	return false, -1
}

type ConnectionsPool struct {
	Logic                 LogicInterface
	IncomingConnections   ConnectionsChannel // входящие соединения
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
			connection.GetResponseChannel() <- WelcomeMessage{Id: connectionId}

			pool.Connections[connectionId] = connection
			pool.Logic.GetIncomingMessagesChannel() <- UserMessage{Data: &LoginMessage{Id: connectionId, Name: authMessage.Name}, Source: connectionId}

			connection.StartReading(pool.Logic.GetIncomingMessagesChannel())
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
	pool.Logic.GetIncomingMessagesChannel() <- UserMessage{Data: &LogoutMessage{connectionId}, Source: connectionId}
	log.Println("CPool: message in sent to logic")
}

func (pool *ConnectionsPool) DispathMessage(msg ServerMessage) {
	except := SearchableArray(msg.Except)
	if len(msg.Targets) == 0 {
		for _, conn := range pool.Connections {
			if exists, _ := except.indexOf(conn.GetId()); !exists {
				conn.GetResponseChannel() <- msg.Data
			}
		}
	} else {
		for _, connectionId := range msg.Targets {
			if exists, _ := except.indexOf(connectionId); !exists {
				pool.Connections[connectionId].GetResponseChannel() <- msg.Data
			}
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
		case connection := <-pool.IncomingConnections:
			log.Println("CPool: connection received", connection)
			pool.processConnection(connection)
			log.Println("CPool connection processed")
		case message := <-pool.Logic.GetOutgoingMessagesChannel():
			//log.Printf("CPool: Outgoing message %T\n", message)
			pool.DispathMessage(message)
			//log.Println("CPool: Message is sent")
		case connectionId := <-pool.ClosingChannel:
			log.Println("CPool: Closing connection", connectionId)
			pool.RemoveConnection(connectionId)
			log.Println("CPool: connection closed")
		}
	}

	log.Println("Connections pool finished")
}

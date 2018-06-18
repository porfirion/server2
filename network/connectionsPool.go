package network

import (
	"log"
)

type MessageFromClient struct {
	ClientId    uint64
	MessageType uint64
	Data        []byte
}

type MessageForClient struct {
	MessageType uint64
	Data        []byte
}

type MessageForClientWrapper struct {
	Targets []uint64 // send only to
	Except  []uint64 // do not send to
	Data    MessageForClient
}

type ConnectionsPool struct {
	IncomingConnections   chan Connection // входящие соединения
	ConnectionsEnumerator chan uint64
	Connections           map[uint64]Connection
	ClosingChannel        chan uint64                  // в этот канал приходят id соединений, которые закрываются
	IncomingMessages      chan MessageFromClient       // сюда приходят сообщения от клиентов
	OutgoingMessages      chan MessageForClientWrapper // канал сообщений для клиентов
}

func (pool *ConnectionsPool) processConnection(connection Connection) {
	go func() {
		pool.Connections[connection.GetId()] = connection
	}()
}

func (pool *ConnectionsPool) RemoveConnection(connectionId uint64) {
	log.Println("CPool: Removing connection", connectionId)
	conn := pool.Connections[connectionId]
	conn.Close("close reason not implemented")
	delete(pool.Connections, connectionId)
}

func (pool *ConnectionsPool) DispatchMessage(msg MessageForClientWrapper) {
	except := SearchableArray(msg.Except)
	if len(msg.Targets) == 0 {
		for _, conn := range pool.Connections {
			if exists, _ := except.indexOf(conn.GetId()); !exists {
				conn.WriteMessage(msg.Data.MessageType, msg.Data.Data)
			}
		}
	} else {
		for _, connectionId := range msg.Targets {
			if exists, _ := except.indexOf(connectionId); !exists {
				pool.Connections[connectionId].WriteMessage(msg.Data.MessageType, msg.Data.Data)
			}
		}
	}
}

func (pool *ConnectionsPool) Start() {
	log.Println("Connections pool started")

	for {
		select {
		case connection := <-pool.IncomingConnections:
			log.Println("CPool: connection received", connection)
			pool.processConnection(connection)
			log.Println("CPool connection processed")
		case message := <-pool.OutgoingMessages:
			//log.Printf("CPool: Outgoing message %T\n", message)
			pool.DispatchMessage(message)
			//log.Println("CPool: Message is sent")
		case connectionId := <-pool.ClosingChannel:
			log.Println("CPool: Closing connection", connectionId)
			pool.RemoveConnection(connectionId)
			log.Println("CPool: connection closed")
		}
	}

	log.Println("Connections pool finished")
}

func NewConnectionsPool() *ConnectionsPool {
	pool := &ConnectionsPool{
		IncomingConnections:   make(chan Connection),
		ConnectionsEnumerator: make(chan uint64),
		Connections:           make(map[uint64]Connection),
		ClosingChannel:        make(chan uint64),
		IncomingMessages:      make(chan MessageFromClient),
		OutgoingMessages:      make(chan MessageForClientWrapper),
	}

	go func() {
		var connectionId uint64 = 1
		for {
			pool.ConnectionsEnumerator <- connectionId
			connectionId++
		}
	}()

	return pool
}

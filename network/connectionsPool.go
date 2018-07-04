package network

import (
	"log"
	"github.com/porfirion/server2/service"
	"encoding/binary"
)

type MessageFromClient struct {
	ClientId    uint64
	Data        service.TypedMessage
}

type MessageForClient struct {
	Targets     []uint64 // send only to
	Data        service.TypedMessage
}

type TypedBytesMessage []byte

func (t TypedBytesMessage) GetType() uint64 {
	return binary.BigEndian.Uint64(t[:8])
}

type ConnectionsPool struct {
	IncomingConnections   chan Connection // входящие соединения
	ConnectionsEnumerator chan uint64
	Connections           map[uint64]Connection
	ClosingChannel        chan uint64            // в этот канал приходят id соединений, которые закрываются
	IncomingMessages      chan MessageFromClient // в этот канал мы пишем сообщения от клиентов
	OutgoingMessages      chan MessageForClient  // канал сообщений для клиентов
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

func (pool *ConnectionsPool) DispatchMessage(msg MessageForClient) {
	if len(msg.Targets) == 0 {
		for _, conn := range pool.Connections {
			conn.WriteMessage(msg.Data)
		}
	} else {
		for _, connectionId := range msg.Targets {
			pool.Connections[connectionId].WriteMessage(msg.Data)
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

			// входящие сообщения должен принимать тот, кто создал пул
			//case message := <-pool.IncomingMessages:
			//	fmt.Printf("Incoming message %v\n", message)
		case connectionId := <-pool.ClosingChannel:
			log.Println("CPool: Closing connection", connectionId)
			pool.RemoveConnection(connectionId)
			log.Println("CPool: connection closed")
		}
	}

	log.Println("Connections pool finished")
}

func NewConnectionsPool(incoming chan MessageFromClient) *ConnectionsPool {
	pool := &ConnectionsPool{
		IncomingConnections:   make(chan Connection),
		ConnectionsEnumerator: make(chan uint64),
		Connections:           make(map[uint64]Connection),
		ClosingChannel:        make(chan uint64),
		IncomingMessages:      incoming,
		OutgoingMessages:      make(chan MessageForClient),
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

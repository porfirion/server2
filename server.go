package main

import (
	"log"
	"net"
)

type MessagesChannel chan Message
type ConnectionsChannel chan Connection

var ControlChannel chan int = make(chan int, 10)

type Connection interface {
	StartReading(ch MessagesChannel)
	WriteMessage(msg Message)
	GetAuth() *Player
	Close()
}

func main() {
	var incomingConnections ConnectionsChannel = make(ConnectionsChannel, 10)
	var incomingMessages MessagesChannel = make(MessagesChannel, 100)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic := &Logic{IncomingMessages: incomingMessages}
	go logic.run()

	log.Println("starting websocket gate")
	ws := &WebSocketGate{":8080", incomingConnections}
	ws.Start()
	log.Println("websocket gate started")

	log.Println("starting tcp gate")
	ts := &TcpGate{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001}, incomingConnections}
	go ts.Start()
	log.Println("tcp gate started")

	log.Println("starting pool")
	pool := &ConnectionsPool{logic: logic, incomingConnections: incomingConnections}
	go pool.Start()
	log.Println("pool started")

	log.Println("running")

	for {
		signal := <-ControlChannel
		log.Println("signal received ", signal)
	}

	log.Println("exit")
}

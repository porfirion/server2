package main

import (
	"log"
	"net"
)

type MessagesChannel chan Message
type ConnectionsChannel chan *Connection

func main() {
	var incomingConnections ConnectionsChannel = make(ConnectionsChannel, 10)
	var incomingMessages MessagesChannel = make(MessagesChannel, 100)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic := new(Logic)
	logic.incomingMessages = incomingMessages
	go logic.run()

	log.Println("starting websocket gate")
	ws := &WebSocketGate{":8080", incomingConnections, incomingMessages}
	ws.Start()
	log.Println("websocket gate started")

	log.Println("starting tcp gate")
	ts := &TcpGate{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001}, incomingConnections, incomingMessages}
	go ts.Start()
	log.Println("tcp gate started")

	log.Println("running")

	for {

	}

	log.Println("exit")
}

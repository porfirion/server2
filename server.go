package main

import (
	"log"
	"net"
)

var ControlChannel chan int = make(chan int, 10)

const format = "%T(%v)\n"

func main() {
	var incomingConnections ConnectionsChannel = make(ConnectionsChannel, 100)

	var incomingMessages UserMessagesChannel = make(UserMessagesChannel)
	var outgoingMessages ServerMessagesChannel = make(ServerMessagesChannel)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic := &Logic{IncomingMessages: incomingMessages, OutgoingMessages: outgoingMessages}
	go logic.Start()

	pool := &ConnectionsPool{logic: logic, incomingConnections: incomingConnections}
	go pool.Start()

	wsGate := &WebSocketGate{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}, incomingConnections}
	go wsGate.Start()

	tcpGate := &TcpGate{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001}, incomingConnections}
	go tcpGate.Start()

	log.Println("Running")

	for {
		signal := <-ControlChannel
		log.Println("signal received ", signal)
	}

	log.Println("exit")
}

package main

import (
	"log"
)

const (
	HTTP_HOST string = ""
	HTTP_PORT string = "8080"
)

func main() {
	var incomingConnections chan *Connection = make(chan *Connection, 10)
	var incomingMessages MessageChannel = make(chan Message, 100)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic := new(Logic)
	logic.incomingMessages = incomingMessages
	go logic.run()

	var ws WebSocketGate
	ws.Start(incomingConnections)

	log.Println("running")

	for {

	}

	log.Println("exit")
}

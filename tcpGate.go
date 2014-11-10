package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type TcpConnection struct {
	conn net.Conn
}

func (connection *TcpConnection) StartReading(ch MessagesChannel) {
	log.Println("starting reading")
	defer connection.conn.Close()

	var buffer []byte

	for {
		buffer = make([]byte, 1024)
		if n, err := connection.conn.Read(buffer); err != nil && err != io.EOF {
			log.Println("error reading from connection")
			fmt.Println(err)
			break
		} else {
			log.Println("read bytes: ", n)
		}

		ch <- DataMessage{&BaseMessage{0, nil}, buffer}
	}
}
func (connection *TcpConnection) WriteMessage(msg Message) {}

func NewTcpConnection(conn net.Conn) Connection {
	connection := &TcpConnection{conn}
	return connection
}

type TcpGate struct {
	addr                *net.TCPAddr
	incomingConnections ConnectionsChannel
}

func (gate *TcpGate) Start() {

	listener, err := net.ListenTCP("tcp4", gate.addr)

	if err != nil {
		log.Println("Error opening listener: ", err)
		return
	}

	log.Println("Listening tcp..", gate.addr)

	// main loop
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("Error: ", err)
		}

		connection := NewTcpConnection(conn)

		log.Println("Connected! ", conn.RemoteAddr())

		gate.incomingConnections <- connection
	}
	log.Println("finished")
}

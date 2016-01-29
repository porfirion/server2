package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type TcpConnection struct {
	*BasicConnection
	socket net.Conn
}

func (connection *TcpConnection) StartReading(ch UserMessagesChannel) {
	go func() {
		log.Println("starting reading")
		defer connection.Close(0, "Unimplemented")

		var buffer []byte

		for {
			buffer = make([]byte, 1024)
			if n, err := connection.socket.Read(buffer); err != nil && err != io.EOF {
				log.Println("error reading from connection")
				fmt.Println(err)
				break
			} else {
				log.Println("read bytes: ", n)
			}

			ch <- UserMessage{connection.id, DataMessage{buffer}}
		}

		log.Println("Reading finished")
	}()
}
func (connection *TcpConnection) StartWriting() {
	log.Println("Not implemented")
}

func (connection *TcpConnection) Close(code int, message string) {
	connection.socket.Close()
}

func (connection *TcpConnection) GetAuth() (*AuthMessage, error) {
	log.Println("TcpConnection.GetAuth is not implemented")
	return nil, errors.New("Not implemented")
}

func NewTcpConnection(socket net.Conn) Connection {
	connection := &TcpConnection{socket: socket}
	connection.StartWriting()
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

	log.Println("Listening tcp:", gate.addr)

	// main loop
	defer listener.Close()

	for {
		socket, err := listener.AcceptTCP()
		if err != nil {
			log.Println("Error: ", err)
		}

		connection := NewTcpConnection(socket)

		log.Println("Connected tcp from ", socket.RemoteAddr())

		gate.incomingConnections <- connection
	}
	log.Println("finished")
}

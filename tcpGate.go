package main

import (
	"fmt"
	"net"
)

type TcpGate struct {
	addr                *net.TCPAddr
	incomingConnections ConnectionsChannel
	incomingMessages    MessagesChannel
}

func (gate *TcpGate) Start() {

	listener, err := net.ListenTCP("tcp4", gate.addr)

	if err != nil {
		fmt.Println("Error opening listener: ", err)
		return
	}

	fmt.Println("Listening tcp..", gate.addr)

	// main loop
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
		}

		connection := NewTcpConnection(conn)

		fmt.Println("Connected! (", conn.RemoteAddr(), ")")

		gate.incomingConnections <- &connection
	}
	fmt.Println("finished")
}

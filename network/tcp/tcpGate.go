package tcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"github.com/porfirion/server2/network"
)

type TcpConnection struct {
	*network.BasicConnection
	socket net.Conn
}

func (connection *TcpConnection) StartReading(ch network.UserMessagesChannel) {
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

			ch <- network.UserMessage{connection.Id, network.DataMessage{buffer}}
		}

		log.Println("Reading finished")
	}()
}
func (connection *TcpConnection) StartWriting() {
	log.Println("Not implemented")
}

func (connection *TcpConnection) Write(msg interface{}) {
	log.Println("Not implemented")
}

func (connection *TcpConnection) Close(code int, message string) {
	connection.socket.Close()
}

func (connection *TcpConnection) GetAuth() (*network.AuthMessage, error) {
	log.Println("TcpConnection.GetAuth is not implemented")
	return nil, errors.New("Not implemented")
}

func NewTcpConnection(socket net.Conn) network.Connection {
	connection := &TcpConnection{socket: socket}
	connection.StartWriting()
	return connection
}

type TcpGate struct {
	Addr                *net.TCPAddr
	Pool                *network.ConnectionsPool
}

func (gate *TcpGate) Start() error {

	listener, err := net.ListenTCP("tcp4", gate.Addr)

	if err != nil {
		log.Println("Error opening listener: ", err)
		return err
	}

	log.Println("Listening tcp:", gate.Addr)

	// main loop
	defer listener.Close()

	for {
		socket, err := listener.AcceptTCP()
		if err != nil {
			log.Println("Error: ", err)
		}

		connection := NewTcpConnection(socket)

		log.Println("Connected tcp from ", socket.RemoteAddr())

		gate.Pool.IncomingConnections <- connection
	}
	log.Println("finished")

	return nil
}

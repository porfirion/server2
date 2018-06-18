package tcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"github.com/porfirion/server2/network"
	"encoding/binary"
)

type TcpConnection struct {
	*network.BasicConnection
	socket net.Conn
}

func (connection *TcpConnection) Close(message string) {
	panic("implement me")
}

func (connection *TcpConnection) WriteMessage(msg interface{}) {
	panic("implement me")
}

func (connection *TcpConnection) StartReading(ch chan network.MessageFromClient) {
	go func() {
		log.Println("starting reading")
		defer func() {
			connection.NotifyPoolWeAreClosing()
		}()

		var buffer []byte

		for {
			buffer = make([]byte, 1024)
			if n, err := connection.socket.Read(buffer); err != nil && err != io.EOF {
				log.Println("error reading from connection")
				fmt.Println(err)
				break
			} else {
				log.Println("read bytes: ", n)
				ch <- network.MessageFromClient{
					connection.Id,
					binary.BigEndian.Uint64(buffer[:8]),
					buffer[8:n],
				}
			}
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

func NewTcpConnection(id uint64, incoming chan network.MessageFromClient, closing chan uint64, socket net.Conn) network.Connection {
	connection := &TcpConnection{
		BasicConnection: network.NewBasicConnection(id, incoming, closing),
		socket:          socket,
	}
	connection.StartWriting()
	return connection
}

type TcpGate struct {
	Addr *net.TCPAddr
	Pool *network.ConnectionsPool
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

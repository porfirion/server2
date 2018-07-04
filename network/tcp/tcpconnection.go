package tcp

import (
	"io"
	"fmt"
	"github.com/porfirion/server2/network"
	"net"
	"log"
	"github.com/porfirion/server2/service"
)

type TcpConnection struct {
	*network.BasicConnection
	socket net.Conn
}

func (connection *TcpConnection) WriteMessage(messageData service.TypedMessage) {
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
					network.TypedBytesMessage(buffer),
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

func (connection *TcpConnection) Close(message string) {
	connection.socket.Close()
}

func NewTcpConnection(
	id uint64,
	incoming chan network.MessageFromClient,
	closing chan uint64,
	socket net.Conn) network.Connection {
	connection := &TcpConnection{
		BasicConnection: network.NewBasicConnection(
			id,
			incoming,
			closing,
		),
		socket: socket,
	}
	connection.StartWriting()
	return connection
}

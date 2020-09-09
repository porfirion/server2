package tcp

import (
	"fmt"
	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/network/pool"
	"github.com/porfirion/server2/service"
	"io"
	"log"
	"net"
)

type TcpConnection struct {
	*pool.BasicConnection
	socket net.Conn
}

func (connection *TcpConnection) WriteMessage(messageData service.TypedMessage) {
	panic("implement me")
}

func (connection *TcpConnection) StartReading(ch chan pool.MessageFromClient) {
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
				if msg, err := messages.DeserializeFromBinary(buffer); err == nil {
					ch <- pool.MessageFromClient{
						ClientId: connection.Id,
						Data:     msg,
					}
				} else {
					log.Println("Error parsing binary message", err)
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
	incoming chan pool.MessageFromClient,
	closing chan uint64,
	socket net.Conn) pool.Connection {
	connection := &TcpConnection{
		BasicConnection: pool.NewBasicConnection(
			id,
			incoming,
			closing,
		),
		socket: socket,
	}
	connection.StartWriting()
	return connection
}

package tcp

import (
	"log"
	"net"
	"github.com/porfirion/server2/network"
)

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

		connection := NewTcpConnection(
			<-gate.Pool.ConnectionsEnumerator,
			gate.Pool.IncomingMessages,
			gate.Pool.ClosingChannel,
			socket,
		)

		log.Println("Connected tcp from ", socket.RemoteAddr())

		gate.Pool.IncomingConnections <- connection
	}
	log.Println("finished")

	return nil
}

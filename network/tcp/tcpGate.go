package tcp

import (
	"github.com/porfirion/server2/network/pool"
	"log"
	"net"
)

type TcpGate struct {
	Addr *net.TCPAddr
	Pool *pool.ConnectionsPool
}

func (gate *TcpGate) Start() {

	listener, err := net.ListenTCP("tcp4", gate.Addr)

	if err != nil {
		log.Fatal("Error opening listener: ", err)
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

	log.Println("TCP Gate finished")
}

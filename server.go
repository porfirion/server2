package main

import (
	"log"
	"net"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/network/ws"
	"github.com/porfirion/server2/network/tcp"
	"time"
)

var (
	ControlChannel chan int = make(chan int, 10)
	logic          *Logic
)

const format = "%T(%v)\n"

func main() {
	// log.SetFlags(log.Ltime | log.Lshortfile) - may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	var incomingConnections network.ConnectionsChannel = make(network.ConnectionsChannel)

	var incomingMessages network.UserMessagesChannel = make(network.UserMessagesChannel)
	var outgoingMessages network.ServerMessagesChannel = make(network.ServerMessagesChannel)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	logic = &Logic{}
	logic.SetParams(LogicParams{
		SimulateByStep:           false,                  // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
		SimulationStepTime:       100 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
		SimulationStepRealTime:   100 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
		SendObjectsTimeout:       time.Second * 1,
		MaxSimulationStepsAtOnce: 10,
	})
	logic.SetIncomingMessagesChannel(incomingMessages)
	logic.SetOutgoingMessagesChannel(outgoingMessages)
	go logic.Start()

	pool := &network.ConnectionsPool{Logic: logic, IncomingConnections: incomingConnections}
	go pool.Start()

	wsGate := &ws.WebSocketGate{Addr: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 8080}, IncomingConnections: incomingConnections}
	go wsGate.Start()

	tcpGate := &tcp.TcpGate{Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001}, IncomingConnections: incomingConnections}
	go tcpGate.Start()

	log.Println("Running")

	for {
		signal := <-ControlChannel
		log.Println("signal received ", signal)
	}

	log.Println("exit")
}

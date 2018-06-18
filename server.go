package main

import (
	"log"
	"net"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/network/ws"
	"github.com/porfirion/server2/network/tcp"
)

var (
	ControlChannel = make(chan int, 10)
	//logic          *Logic
)

func main() {
	// log.SetFlags(log.Ltime | log.Lshortfile) - may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	//logic = &Logic{
	//	IncomingMessages: make(network.UserMessagesChannel, 10),
	//	OutgoingMessages: make(network.ServerMessagesChannel, 10),
	//	Params: LogicParams{
	//		SimulateByStep:           true,                  // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
	//		SimulationStepTime:       500 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
	//		SimulationStepRealTime:   500 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
	//		SendObjectsTimeout:       time.Millisecond * 500,
	//		MaxSimulationStepsAtOnce: 10,
	//	},
	//}
	//go logic.Start()

	pool := network.NewConnectionsPool()
	go pool.Start()

	wsGate := &ws.WebSocketGate{
		Addr: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 8080},
		Pool: pool,
	}
	go wsGate.Start()

	tcpGate := &tcp.TcpGate{
		Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001},
		Pool: pool,
	}
	go tcpGate.Start()

	log.Println("Running")

	for {
		signal := <-ControlChannel
		log.Println("signal received ", signal)
	}

	log.Println("exit")
}

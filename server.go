package main

import (
	"log"
	"net"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/network/ws"
	"github.com/porfirion/server2/network/tcp"
	"github.com/porfirion/server2/logic"
	"time"
)

var (
	ControlChannel = make(chan int, 10)
	lg          *logic.GameLogic
)

func main() {
	// log.SetFlags(log.Ltime | log.Lshortfile) - may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	broker := logic.NewBroker()
	broker.Start()

	chat := logic.NewChat(logic.NewBasicService(1))
	chat.Start()
	broker.RegisterService(chat)

	logicMessages := make(logic.ServerMessagesChannel, 10)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	lg := &logic.GameLogic{
		IncomingMessages: make(logic.UserMessagesChannel, 10),
		OutgoingMessages: logicMessages,
		Params: logic.LogicParams{
			SimulateByStep:           true,                  // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
			SimulationStepTime:       500 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
			SimulationStepRealTime:   500 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
			SendObjectsTimeout:       time.Millisecond * 500,
			MaxSimulationStepsAtOnce: 10,
		},
	}
	go lg.Start()

	logicSvc := &logic.GameLogicService{
		BasicService:          logic.NewBasicService(1),
		Logic:                 lg,
		LogicOutgoingMessages: logicMessages,
	}
	logicSvc.Start()

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

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
)

func main() {
	// log.SetFlags(log.Ltime | log.Lshortfile) - may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	broker := logic.NewBroker()
	go broker.Start()

	chat := logic.NewChat(logic.NewBasicService(1))
	go chat.Start()
	broker.RegisterService(chat)

	logicMessages := make(logic.ServerMessagesChannel, 10)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	lg := &logic.GameLogic{
		IncomingMessages: make(logic.UserMessagesChannel, 10),
		OutgoingMessages: logicMessages,
		Params: logic.LogicParams{
			SimulateByStep:           true,                   // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
			SimulationStepTime:       500 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
			SimulationStepRealTime:   500 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
			SendObjectsTimeout:       time.Millisecond * 500,
			MaxSimulationStepsAtOnce: 10,
		},
	}
	go lg.Start()
	log.Println("GameLogic started")

	logicSvc := &logic.GameLogicService{
		BasicService:          logic.NewBasicService(logic.SERVICE_TYPE_LOGIC),
		Logic:                 lg,
		LogicOutgoingMessages: logicMessages,
	}
	go logicSvc.Start()
	log.Println("LogicService started")

	networkSvc := &logic.NetworkService{
		BasicService: logic.NewBasicService(logic.SERVICE_TYPE_NETWORK),
	}
	go networkSvc.Start()
	log.Println("Network started")
	broker.RegisterService(networkSvc)

	pool := network.NewConnectionsPool(make(chan network.MessageFromClient))
	go pool.Start()
	log.Println("Pool started")
	networkSvc.SetPool(pool)

	wsGate := &ws.WebSocketGate{
		Addr: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 8080},
		Pool: pool,
	}
	go wsGate.Start()
	log.Println("WsGate started")

	tcpGate := &tcp.TcpGate{
		Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 25001},
		Pool: pool,
	}
	go tcpGate.Start()
	log.Println("TcpGate started")


	log.Println("Running")

	for {
		signal := <-ControlChannel
		log.Println("signal received ", signal)
	}

	log.Println("exit")
}

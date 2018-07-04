package main

import (
	"log"
	"net"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/network/ws"
	"github.com/porfirion/server2/network/tcp"
	"github.com/porfirion/server2/service"
	"time"
	"os"
	"os/signal"
	"fmt"
	"github.com/porfirion/server2/game"
	"github.com/porfirion/server2/chat"
	"github.com/porfirion/server2/messages"
)

func main() {
	//log.SetFlags(log.Ltime | log.Lshortfile) //may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	broker := service.NewBroker()
	go broker.Start()

	chat := chat.NewChat(service.NewBasicService(service.TypeChat))
	go chat.Start()
	broker.RegisterService(chat)

	logicMessages := make(messages.ServerMessagesChannel, 10)

	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	lg := &game.GameLogic{
		IncomingMessages: make(messages.UserMessagesChannel, 10),
		OutgoingMessages: logicMessages,
		Params: game.LogicParams{
			SimulateByStep:           true,                   // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
			SimulationStepTime:       500 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
			SimulationStepRealTime:   500 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
			SendObjectsTimeout:       time.Millisecond * 500,
			MaxSimulationStepsAtOnce: 10,
		},
	}
	go lg.Start()
	log.Println("GameLogic started")

	logicSvc := &game.GameLogicService{
		BasicService:          service.NewBasicService(service.TypeLogic),
		Logic:                 lg,
		LogicOutgoingMessages: logicMessages,
	}
	go logicSvc.Start()
	log.Println("LogicService started")
	broker.RegisterService(logicSvc)

	networkSvc := &network.NetworkService{
		BasicService: service.NewBasicService(service.TypeNetwork),
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	select {
	case sig := <-interrupt:
		fmt.Printf("Got signal \"%s\"\n", sig.String())
	}

	fmt.Println("FINISH MAIN LOOP")
}

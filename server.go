package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/porfirion/server2/auth"
	"github.com/porfirion/server2/chat"
	"github.com/porfirion/server2/game/next"
	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/service"
)

func main() {
	wsPort := flag.Int("wsport", 18080, "port to listen for WebSocket connections")
	tcpPort := flag.Int("tcpport", 25001, "port to listen to TCP connections")
	httpPort := flag.Int("httpport", 2018, "")
	noStatic := flag.Bool("nostatic", false, "disables serving static files")

	flag.Parse()

	log.SetFlags(log.Ltime | log.Lshortfile) //may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	broker := service.NewBroker(func(message service.ServiceMessage) service.ServiceType {
		switch message.MessageData.(type) {
		case *messages.TextMessage:
			if message.SourceServiceType == service.TypeNetwork {
				return service.TypeChat
			}
		case *messages.AuthMessage:
			return service.TypeAuth
		case *messages.SimulateMessage,
			*messages.ChangeSimulationMode:
			return service.TypeLogic
		}

		return 0
	})
	go broker.Start()

	authService := auth.NewService()
	go authService.Start()
	broker.RegisterService(authService)

	chatService := chat.NewService()
	go chatService.Start()
	broker.RegisterService(chatService)

	logicService := next.NewService()
	go logicService.Start()
	broker.RegisterService(logicService)

	networkService := network.NewService(*wsPort, *tcpPort, *httpPort, *noStatic)
	go networkService.Start()
	broker.RegisterService(networkService)

	log.Println("All services started")

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	sig := <-interruptChan
	fmt.Printf("Got signal \"%s\"\n", sig.String())
	fmt.Println("FINISH MAIN LOOP")
}

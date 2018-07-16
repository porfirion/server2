package main

import (
	"log"
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/service"
	"os"
	"os/signal"
	"fmt"
	"github.com/porfirion/server2/game"
	"github.com/porfirion/server2/chat"
	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/auth"
)

func main() {
	//log.SetFlags(log.Ltime | log.Lshortfile) //may be very useful to know where print was called
	log.SetFlags(log.Lmicroseconds)

	broker := service.NewBroker(func(message service.ServiceMessage) uint64 {
		switch message.MessageData.(type) {
		case *messages.TextMessage:
			if message.SourceServiceType == service.TypeNetwork {
				return service.TypeChat
			}
		case *messages.AuthMessage:
			return service.TypeAuth
		}

		return 0
	})
	go broker.Start()

	auth := auth.NewService()
	go auth.Start()
	broker.RegisterService(auth)

	chat := chat.NewService()
	go chat.Start()
	broker.RegisterService(chat)

	logicSvc := game.NewService()
	go logicSvc.Start()
	broker.RegisterService(logicSvc)

	networkSvc := network.NewService()
	go networkSvc.Start()
	broker.RegisterService(networkSvc)

	log.Println("All services started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	select {
	case sig := <-interrupt:
		fmt.Printf("Got signal \"%s\"\n", sig.String())
	}

	fmt.Println("FINISH MAIN LOOP")
}

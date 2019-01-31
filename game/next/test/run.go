package main

import (
	"github.com/porfirion/server2/game/next"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Creating")
	logic := next.NewLogic(next.SimulationModeContinuous, time.Second, time.Second)
	log.Println("Starting")
	logic.Start()
	log.Println("Started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
		case <-interrupt:
			logic.Stop()
	}
}
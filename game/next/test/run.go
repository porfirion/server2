package main

import (
	"github.com/porfirion/server2/game/next"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	control := make(chan next.ControlMessage)
	input := make(chan next.PlayerInput)
	log.Println("Creating")
	logic := next.NewLogic(control, input, next.SimulationModeContinuous, time.Second, time.Second)
	log.Println("Starting")
	logic.Start()
	log.Println("Started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	<-logic.Stop()
	log.Println("Stopped")
}

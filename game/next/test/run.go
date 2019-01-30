package main

import (
	"github.com/porfirion/server2/game/next"
	"log"
	"time"
)

func main() {
	log.Println("Creating")
	logic := next.NewLogic(next.SimulationModeContinuous, time.Second, time.Second)
	log.Println("Starting")
	logic.Start()
	log.Println("Started")
	select {}
}
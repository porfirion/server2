package next

import (
	"fmt"
	"testing"
	"time"
)

func TestLogicImplContinuous(t *testing.T) {
	controlChan := make(chan ControlMessage)
	inputChan := make(chan PlayerInput)
	monitorChan := make(chan GameState)

	l := NewLogic(controlChan, inputChan, SimulationModeContinuous, time.Second, time.Second)
	l.SetMonitorChan(monitorChan)

	go func() {
		for msg := range monitorChan {
			fmt.Println("State: %v", msg)
		}
	}()

	go func() {
		for msg := range l.outputChan {
			fmt.Println("Output %v", msg)
		}
	}()

	log.Println(l)
	l.Start()

	//for i := 0; i < 10; i++ {
	//	controlChan <- ControlMessageSimulate
	//}

	log.Println(l)
	l.Stop()
}

// Each simulate command should result in message from monitor chan
func TestLogicImplStepByStep(t *testing.T) {
	controlChan := make(chan ControlMessage)
	inputChan := make(chan PlayerInput)
	monitorChan := make(chan GameState)

	l := NewLogic(controlChan, inputChan, SimulationModeStepByStep, time.Second, time.Second)
	l.SetMonitorChan(monitorChan)

	go func() {
		for msg := range monitorChan {
			fmt.Println("State: %v", msg)
		}
	}()

	go func() {
		for msg := range l.outputChan {
			fmt.Println("Output %v", msg)
		}
	}()

	log.Println(l)
	l.Start()

	for i := 0; i < 10; i++ {
		controlChan <- ControlMessageSimulate
		<-monitorChan
	}

	log.Println(l)
	l.Stop()
}

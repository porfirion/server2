package next

import (
	"fmt"
	"testing"
	"time"
)

func drain(ch chan interface{}) {
	go func() {
		for range ch {}
	}()
}

func drainStates(ch chan GameState) <-chan int {
	drained := 0
	resp := make(chan int, 1)
	go func() {
		for range ch {
			drained++
		}

		resp <- drained
		close(resp)
	}()

	return resp
}

func TestLogicImplContinuous(t *testing.T) {
	controlChan := make(chan ControlMessage)
	inputChan := make(chan PlayerInput)

	l := NewLogic(controlChan, inputChan, SimulationModeContinuous, time.Second, time.Second)

	drain(l.outputChan)
	l.Start()
	<-l.Stop()
}

func TestLogicImplStep(t *testing.T) {
	controlChan := make(chan ControlMessage)
	inputChan := make(chan PlayerInput)
	monitorChan := make(chan GameState)

	l := NewLogic(controlChan, inputChan, SimulationModeStepByStep, time.Second, time.Second)
	l.SetMonitorChan(monitorChan)

	iterationsCount := 10

	drain(l.outputChan)
	dr := drainStates(monitorChan)

	l.Start()

	for i := 0; i < iterationsCount; i++ {
		fmt.Println("sending simulate")
		controlChan <- ControlMessageSimulate
	}

	<-l.Stop()
	statesCount := <- dr

	if iterationsCount != statesCount {
		t.Fatalf("iterations count was %d and received %d states", iterationsCount, statesCount)
	} else {
		t.Logf("processed %d iterations", iterationsCount)
	}
}

func TestLogicImplReplay(t *testing.T) {
	controlChan := make(chan ControlMessage)
	inputChan := make(chan PlayerInput)
	monitorChan := make(chan GameState)

	l := NewLogic(controlChan, inputChan, SimulationModeReplay, time.Second, time.Second)
	l.SetMonitorChan(monitorChan)

	iterationsCount := 10

	drain(l.outputChan)
	dr := drainStates(monitorChan)

	l.Start()

	for i := 0; i < iterationsCount; i++ {
		fmt.Println("sending simulate")
		controlChan <- ControlMessageSimulate
	}

	<-l.Stop()
	statesCount := <- dr

	if iterationsCount != statesCount {
		t.Fatalf("iterations count was %d and received %d states", iterationsCount, statesCount)
	} else {
		t.Logf("processed %d steps", iterationsCount)
	}
}

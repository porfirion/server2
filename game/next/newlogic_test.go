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

func drainStates(ch chan GameState) <-chan ([]GameState) {
	drainedCount := 0
	resp := make(chan []GameState, 1)
	states := make([]GameState, 0)

	go func() {
		for state := range ch {
			states = append(states, state)
			drainedCount++
		}

		resp <- states
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
	states := <- dr

	if iterationsCount != len(states) {
		t.Fatalf("iterations count was %d and received %d states", iterationsCount, len(states))
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
	states := <- dr

	if iterationsCount != len(states) {
		t.Fatalf("iterations count was %d and received %d states", iterationsCount, len(states))
	} else {
		t.Logf("processed %d steps", iterationsCount)
	}
}

package next

import (
	"log"
	"testing"
	"time"
)

func TestLogicImpl(t *testing.T) {
	l := NewLogic(SimulationModeStepByStep, time.Second, time.Second)
	log.Println(l)
	l.Start()
	log.Println(l)
}

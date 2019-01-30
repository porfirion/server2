package next

import (
	"log"
	"testing"
	"time"
)

func TestLogicImpl(t *testing.T) {
	l := NewLogic(SimulationModeContinuous, time.Second, time.Second)
	log.Println(l)
	l.Start()
	log.Println(l)
	l.Stop()
}

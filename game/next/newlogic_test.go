package next

import (
	"log"
	"testing"
)

func TestLogicImpl(t *testing.T) {
	l := NewLogic()
	log.Println(l)
	l.mainStep()
	log.Println(l)
}

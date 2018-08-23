package game

import (
	"testing"
	"log"
)

func TestLogicImpl(t *testing.T) {
	l := &LogicImpl{
		gameTick:      0,
		playersInputs: []PlayerInput{{}, {}, {}, {gameTick:1}},
	}

	log.Println(l)

	l.mainStep()

	log.Println(l)
}

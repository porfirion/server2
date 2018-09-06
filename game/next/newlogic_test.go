package next

import (
	"testing"
	"log"
)

func TestLogicImpl(t *testing.T) {
	l := &LogicImpl{
		gameTick:             0,
		pendingPlayersInputs: []PlayerInput{{}, {}, {}, {gameTick: 1}},
	}

	log.Println(l)

	l.mainStep()

	log.Println(l)
}

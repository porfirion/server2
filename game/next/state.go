package next

import (
	"github.com/porfirion/server2/world"
	"io"
	"time"
)

type GameStateImpl struct {
	*world.WorldMap
}

func (st *GameStateImpl) ProcessSimulationStep(time.Duration) GameState {
	return st
}

func (st *GameStateImpl) Copy() GameState {
	return st
}

func (st *GameStateImpl) GetTickAndTime() (int64, time.Time) {
	panic("implement me")
}

func (st *GameStateImpl) SetTickAndTime(int64, time.Time) {
	panic("implement me")
}

func (st *GameStateImpl) Serialize(writer io.Writer) {
	panic("implement me")
}

func NewGameState() GameState {
	return &GameStateImpl{ world.NewWorldMap() };
}
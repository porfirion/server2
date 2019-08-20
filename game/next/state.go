package next

import (
	"encoding/gob"
	"github.com/porfirion/server2/world"
	"io"
	"time"
)

type GameStateImpl struct {
	*world.WorldMap
}

func (st *GameStateImpl) GetPlayerState(playerId uint) PlayerState {
	// взять вьюпорт пользователя
	// найти объекты, которые в него попадают
	// отправить пользователю найденные объекты пользователю
	// TODO как быть с теми объектами, которые раньше пользователю отправляли а теперь они исчезли?

	log.Println("Creating player state stub")
	return PlayerState{}
}

func (st *GameStateImpl) ProcessSimulationStep(time.Duration) GameState {
	return st
}

func (st *GameStateImpl) Copy() GameState {
	return st
}

func (st *GameStateImpl) Serialize(writer io.Writer) {
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(st); err != nil {
		log.Fatal("Can't encode current state")
	}
}

func NewGameState() GameState {
	return &GameStateImpl{ world.NewWorldMap(10000, 10000) };
}
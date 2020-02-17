// GameState - all the information about game
// Some entities belong to player
// Some player entities have physical representation (map objects)
// Player's input is applied only to it's own entity
// Player's field of view equals players map objects field of view
// List of player's entities and visible map objects (they are entities too) is synchronized with client
// There are some global entities (time of day, leaderboard, etc.) - they are synchronized for each player
// On each sync we send global entities, private entities and visible entities
// Each entity has it's own sync rate
// WorldMap - is a system!
// +HealthSystem
package next

import (
	"encoding/gob"
	"github.com/porfirion/server2/world"
	"io"
	"time"
)

type PhysicalObjectComponent struct {
	Component
}
type HealthComponent struct {
	Component
}

type Component struct {
	id uint64
	tp uint64
}

type Entity struct {
	owner uint64 // playerId
	visibility uint // private (player's private state), public (visible for other players in field of view), global (visible for everyone)?
	components []Component // list of components for this entity // maybe should be fixes size to decrease allocations?

	children []*Component // ? child entities for complex cases
}

type GameStateImpl struct {
	*world.WorldMap // SYSTEM!

	entities []*Entity // all entities in the world
	playersEntities map[uint64]*Entity // entities bound to players (copy of links from entities)
	globalEntities []*Entity // list of entities visible for all
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
	return &GameStateImpl{
		WorldMap: world.NewWorldMap(10000, 10000),
	};
}
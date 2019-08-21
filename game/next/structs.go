package next

import (
	"github.com/porfirion/server2/network/pool"
	"github.com/porfirion/server2/service"
	"io"
	"time"
)

type InputKey int

const (
	LeftMouseButton   InputKey = 1
	MiddleMouseButton InputKey = 2
	RightMouseButton  InputKey = 3
	SpaceButton       InputKey = 4
)

type Vector struct {
	X float32
	Y float32
}

type SimulationMode int8

const (
	// обычный режим - запускаем следующий шаг симуляции по времени
	SimulationModeContinuous SimulationMode = 1
	// пошаговый режим - запускаем следующий шаг только при приходе соответствующей команды
	SimulationModeStepByStep SimulationMode = 2
	// непрерывная симуляция для "проигрывания" записи. Фактически это аналог по шагам,
	// только тут в логику будут передавать предварительно записанные инпуты, а потом посылать simulate_message.
	// Вопрос как поставить такое воспроизведение на паузу?
	SimulationModeReplay SimulationMode = 3
)

type ControlMessage int16

const (
	ControlMessageStop                 ControlMessage = 1
	ControlMessageSimulate             ControlMessage = 2
	ControlMessageChangeModeContinuous ControlMessage = 4 | 1
	ControlMessageChangeModeStepByStep ControlMessage = 4 | 2
	ControlMessageChangeModeReplay     ControlMessage = 4 | 3
)

type PlayerAction int32

const (
	PlayerActionMove    = 1
	PlayerActionAbility = 2
)

type PlayerInput struct {
	PlayerId    uint
	gameTick    uint // tick of the game when input was received
	Action      PlayerAction
	actionData  interface{}
	PressedKeys []InputKey
	MouseVector Vector // position of mouse relative to screen center (aka viewport position/player object position)
}

// Maybe some time we will adjust GameTick with player latency - so we should use this getter instead of field itself
func (i PlayerInput) GetGameTick() uint {
	return i.gameTick
}

// Описание видимого мира для конкретного игрока
type PlayerState struct {

}

type Player struct {
	Id                   uint
	prevStates           map[uint]PlayerState
	playerObjectId       uint64
	additionalObjectsIds []uint64
	conn                 pool.Connection // непосредственное соединение с игроком
}

func (player Player) SendState(state PlayerState) {
	// count diff with prev state
	// send diff to player
	player.conn.WriteMessage(service.TypedMessage())
}

// GameState should contain all information about players, objects, etc.
// Consists of parts:
//   1) list of PlayerState (can contain link to physical object)
//   2) list of physical objects (can contain additional attributes such as playerId, etc)
//   3) current time and all settings (does it differs from Logic itself?)
type GameState interface {
	ProcessSimulationStep(time.Duration) GameState
	Copy() GameState
	Serialize(writer io.Writer)
	GetPlayerState(playerId uint) PlayerState
}

type HistoryEntry struct {
	state    GameState
	tick     uint64
	gameTime time.Time //(tick * stepDuration)
	realTime time.Time
}

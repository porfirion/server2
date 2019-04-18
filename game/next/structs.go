package next

import (
	"github.com/porfirion/server2/network/pool"
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
	conn pool.Connection // непосредственное соединение с игроком
}

type Player struct {
	prevStates        map[uint]PlayerState
	playerObjectId    uint64
	additionalObjects []uint64
}

func (player Player) SendState(state PlayerState) {
	// count diff with prev state
	// send diff to player
}

// GameState should contain all information about players, objects, etc.
type GameState interface {
	ProcessSimulationStep(time.Duration) GameState

	Copy() GameState

	Serialize(writer io.Writer)
}

type HistoryEntry struct {
	state    GameState
	tick     uint64
	gameTime time.Time //(tick * stepDuration)
	realTime time.Time
}

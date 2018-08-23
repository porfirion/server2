package game

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
	// непрерывная симуляция для "проигрывания" записи
	SimulationModeReplay     SimulationMode = 3
)

type PlayerInput struct {
	PlayerId    uint
	gameTick    uint // tick of the game when input was received
	PressedKeys []InputKey
	MouseVector Vector // position of mouse relative to screen center (aka viewport position/player object position)
}

// Maybe some time we will adjust GameTick with player latency - so we should use this getter instead of field itself
func (i PlayerInput) GetGameTick() uint {
	return i.gameTick
}

type LogicImpl struct {
	gameTick      uint
	playersInputs []PlayerInput
}

func (l *LogicImpl) GetCurrentGameTick() uint {
	return l.gameTick
}

func (l *LogicImpl) processPlayerInput(input PlayerInput) {}

func (l LogicImpl) applyPlayersInputsUpTo(currentTick uint) {
	var ind int
	for ind = 0; ind < len(l.playersInputs) && l.playersInputs[ind].GetGameTick() <= currentTick; ind++ {
		l.processPlayerInput(l.playersInputs[ind])
	}

	l.playersInputs = l.playersInputs[ind:]
}

func (l *LogicImpl) ShouldSimulate() bool {

}

func (l *LogicImpl) main() {
	for l.ShouldSimulate() {
		l.simulateStep()
	}

}
func (l *LogicImpl) simulateStep() {
	l.gameTick++

	l.applyPlayersInputsUpTo(l.gameTick)

}

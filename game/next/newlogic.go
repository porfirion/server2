package next

import (
	"log"
	"math"
	"time"

	"github.com/porfirion/server2/world"
)

type LogicImpl struct {
	gameTick     uint64
	prevTickTime time.Time

	controlChan chan ControlMessage
	inputChan   chan PlayerInput

	players map[uint]Player

	worldMap *world.WorldMap

	Mode               SimulationMode // режим логики
	flagShouldStop     bool           // говорит что нужно остановить логику
	flagShouldSimulate bool           // сбрасывается при начале mainStep

	simulationStepTime     time.Duration // сколько виртуального времени проходит за один тик
	simulationStepRealTime time.Duration // сколько реального времени проходит за один тик

	prevStates map[uint64]GameState
	state      GameState
}

func (l *LogicImpl) NextSimulationTime() time.Time {
	if l.Mode == SimulationModeContinuous {
		return l.prevTickTime.Add(l.simulationStepRealTime)
	} else {
		// never
		return time.Unix(math.MaxInt64, 0)
	}
}

// Говорит наступило ли время для симуляции
func (l *LogicImpl) ShouldSimulate() bool {
	return l.NextSimulationTime().Before(time.Now())
}

func (l *LogicImpl) receiveInputsUntilShouldSimulate() []PlayerInput {
	//listen to control chan in parallel.

	var inputsBuffer []PlayerInput

	var timeout int64

	if (l.Mode == SimulationModeContinuous) {
		// receive everything from input chan until next simulation time
		var timeout = l.NextSimulationTime().Sub(time.Now())
		if timeout <= 0 {
			// время уже прошло!
			return inputsBuffer;
		}
	} else if (l.Mode == SimulationModeStepByStep) {
		// receive everything from input chang until simulateMessage
		timeout = math.MaxInt64
	} else if (l.Mode == SimulationModeReplay) {
		// get everything from replay source (aka input chan) and/not wait for simulate message

		// будем предполагать, что воспроизведение - это будет обёртка над логикой, а логика об этом не будет ничего знать
		// логика будет забирать инпут из обычного канала и ждать сигнала к симуляции.
		// обёртка же будет скармливать в логику весь нужный инпут и потом выдавать сигнал к симуляции
		// потому для логики поведение будет таким же, как для воспроизведения по шагам
		timeout = math.MaxInt64
	}

	timer := time.NewTimer(time.Duration(timeout))
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	countReceived := 0

	var stopReceiving bool

	for !stopReceiving && !l.ShouldSimulate() {
		select {
		case msg := <-l.inputChan:
			inputsBuffer = append(inputsBuffer, msg)
			log.Println("Logic: received message ", msg)
			countReceived++
		case ctrl := <-l.controlChan:
			switch ctrl {
			case ControlMessageStop:
				l.flagShouldStop = true
				stopReceiving = true
			case ControlMessageSimulate:
				l.flagShouldSimulate = true
				stopReceiving = true
			default:
				log.Printf("Unknown control message %d", ctrl)
			}
			log.Println("Logic: SHOULD STOP!!!")
		case <-timer.C:
			// пришло время, больше ничего читать не будем
			stopReceiving = true
		}
	}

	log.Printf("received %d inputs\n", countReceived)

	return inputsBuffer
}

func (l *LogicImpl) applyPlayerInputs(state GameState, inputs []PlayerInput) GameState {
	log.Println("apply inputs stub")
	for range inputs {

	}
	return state
}

func (l *LogicImpl) sendStateToPlayers(state GameState) {
	for range l.players {
		// взять вьюпорт пользователя
		// найти объекты, которые в него попадают
		// отправить пользователю найденные объекты пользователю
		// TODO как быть с теми объектами, которые раньше пользователю отправляли а теперь они исчезли?
	}
}

func (l *LogicImpl) simulateStep(state GameState) GameState {
	log.Println("simulate step stub")
	return state;
}

// Main step of server. Apply pending players inputs
// run physics engine simulation,
// calculate states for all players,
// sending states to respective players
func (l *LogicImpl) mainStep(inputs []PlayerInput) {
	// ждём когда придёт время симулировать
	if l.ShouldSimulate() {
		l.gameTick++
		l.state = l.applyPlayerInputs(l.state, inputs)
		l.state = l.applyQueuedEvents(l.state)
		l.state = l.simulateStep(l.state)
	}
}

// основной цикл логики
// получаем весь инпут, складываем его в очередь
// как только настаёт время - выполняем основной шаг (mainStep)
func (l *LogicImpl) mainLoop() {
	for l.flagShouldStop {
		// вычитываем инпут и кладём в очередь
		// до тех пор, пока не придёт время симулировать
		inputsBuffer := l.receiveInputsUntilShouldSimulate()
		l.mainStep(inputsBuffer)
	}

}

func (l *LogicImpl) Start() {
	go l.mainLoop()
}

// Говорит логике остановиться
// В реальности логика остановится только тогда, когда она проверит этот канал,
// а это происходит тогда, когда она принимает инпуты
func (l *LogicImpl) Stop() {
	l.controlChan <- ControlMessageStop
}

func (l *LogicImpl) applyQueuedEvents(state GameState) GameState {
	return state
}

func NewLogic(mode SimulationMode, stepTime, stepRealTime time.Duration) *LogicImpl {
	logic := &LogicImpl{
		gameTick:               0,
		prevTickTime:           time.Time{},
		controlChan:            make(chan ControlMessage),
		inputChan:              make(chan PlayerInput),
		players:                make(map[uint]Player),
		worldMap:               world.NewWorldMap(),
		Mode:                   mode,
		flagShouldStop:         false,
		flagShouldSimulate:     false,
		simulationStepTime:     stepTime,
		simulationStepRealTime: stepRealTime,
		prevStates:             make(map[uint64]GameState),
		state:                  NewGameState(),
	}
	return logic
}

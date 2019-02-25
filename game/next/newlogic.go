package next

import (
	logMod "log"
	"math"
	"os"
	"time"

	"github.com/porfirion/server2/world"
)

var log = logMod.New(os.Stdout, "Newlogic: ", logMod.Ltime | logMod.Lshortfile)

type LogicImpl struct {
	gameTick         uint64
	prevTickRealTime time.Time

	gameTime time.Time

	controlChan chan ControlMessage
	inputChan   chan PlayerInput

	players map[uint]Player

	worldMap *world.WorldMap

	Mode               SimulationMode // режим логики
	flagShouldStop     bool           // говорит что нужно остановить логику
	flagShouldSimulate bool           // сбрасывается при начале mainStep

	simulationStepTime     time.Duration // сколько виртуального времени проходит за один тик
	simulationStepRealTime time.Duration // сколько реального времени проходит за один тик

	prevStates []GameState
	state      GameState
}

func (l *LogicImpl) NextSimulationTime() time.Time {
	if l.Mode == SimulationModeContinuous {
		return l.prevTickRealTime.Add(l.simulationStepRealTime)
	} else {
		// never
		return time.Unix(math.MaxInt64, 0)
	}
}

// Говорит наступило ли время для симуляции
func (l *LogicImpl) ShouldSimulate() bool {
	return l.flagShouldSimulate || l.NextSimulationTime().Before(time.Now())
}

func (l *LogicImpl) receiveInputsUntilShouldSimulate() []PlayerInput {
	//listen to control chan in parallel.
	log.Println("receiving inputs")

	var inputsBuffer []PlayerInput

	var timeout time.Duration

	if (l.Mode == SimulationModeContinuous) {
		// receive everything from input chan until next simulation time
		timeout = l.NextSimulationTime().Sub(time.Now())
		if timeout <= 0 {
			log.Println("time to tick has already come!")
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
	timerFired := false
	defer func() {
		if !timerFired && !timer.Stop() {
			// мы вышли не по таймеру (т.е. не вычитали из него) и таймер уже сработал.. Очистим канал
			<-timer.C
		}
	}()

	countReceived := 0

	var stopReceiving bool

	for !stopReceiving && !l.ShouldSimulate() {
		select {
		case msg := <-l.inputChan:
			inputsBuffer = append(inputsBuffer, msg)
			log.Println("received message ", msg)
			countReceived++
		case ctrl := <-l.controlChan:
			switch ctrl {
			case ControlMessageStop:
				l.flagShouldStop = true
				stopReceiving = true
			case ControlMessageSimulate:
				// если у нас обычная симуляция, то мы просто заставим шаг произойти преждевременно
				l.flagShouldSimulate = true
				stopReceiving = true
			case ControlMessageChangeModeContinuous:
				l.Mode = SimulationModeContinuous
				stopReceiving = true
			case ControlMessageChangeModeStepByStep:
				l.Mode = SimulationModeStepByStep
				stopReceiving = true
			case ControlMessageChangeModeReplay:
				l.Mode = SimulationModeReplay
				stopReceiving = true
			default:
				log.Printf("Unknown control message %d", ctrl)
			}
			log.Println("SHOULD STOP!!!")
		case <-timer.C:
			// пришло время, больше ничего читать не будем
			log.Println("timer fired")
			stopReceiving = true
			timerFired = true
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

// Main step of server. Apply pending players inputs
// run physics engine simulation,
// calculate states for all players,
// sending states to respective players
func (l *LogicImpl) mainStep(inputs []PlayerInput) {
	log.Println("main step started")

	l.prevStates = append(l.prevStates, l.state)

	l.state = l.state.Copy()
	tick, tm := l.state.GetTickAndTime()
	l.state.SetTickAndTime(tick + 1, tm.Add(l.simulationStepTime))

	l.state = l.applyPlayerInputs(l.state, inputs)
	l.state = l.applyQueuedEvents(l.state)

	l.prevTickRealTime = time.Now()
	l.flagShouldSimulate = false

	l.state.ProcessSimulationStep(l.simulationStepTime)

	log.Println("main step finished")
}

// основной цикл логики
// получаем весь инпут, складываем его в очередь
// как только настаёт время - выполняем основной шаг (mainStep)
func (l *LogicImpl) mainLoop() {
	log.Println("starting main loop")
	var inputsBuffer []PlayerInput

	for !l.flagShouldStop {
		// Read the inputs and put into queue until simulation time comes.
		// But reading can break in case of receiving control message. So before simulation
		// we should check if simulation time has come really
		inputsBuffer := append(inputsBuffer, l.receiveInputsUntilShouldSimulate()...)

		if l.ShouldSimulate() {
			l.mainStep(inputsBuffer)

			// Clear inputs buffer only after simulation,
			// because receiving could be aborted before simulation time
			inputsBuffer = inputsBuffer[:0]
		}
	}
	log.Println("main loop stopped")
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

// предполагаем, что события будут инициироваться не только игроками, но и самой логикой.
// Они будут складываться в очередь и срабатывать в назначенное время.
func (l *LogicImpl) applyQueuedEvents(state GameState) GameState {
	return state
}

func NewLogic(mode SimulationMode, stepTime, stepRealTime time.Duration) *LogicImpl {
	logic := &LogicImpl{
		gameTick:               0,
		prevTickRealTime:       time.Time{},
		controlChan:            make(chan ControlMessage),
		inputChan:              make(chan PlayerInput),
		players:                make(map[uint]Player),
		Mode:                   mode,
		flagShouldStop:         false,
		flagShouldSimulate:     false,
		simulationStepTime:     stepTime,
		simulationStepRealTime: stepRealTime,
		prevStates:             make([]GameState, 0, 10),
		state:                  NewGameState(),
	}
	return logic
}

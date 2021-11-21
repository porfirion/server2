package next

import (
	logMod "log"
	"math"
	"os"
	"time"
)

var log = logMod.New(os.Stdout, "NewLogic: ", logMod.Lmicroseconds|logMod.Lshortfile)

// Logic describes logic of receiving inputs, processing simulation step, etc.
type Logic struct {
	gameTick         uint64
	prevTickRealTime time.Time

	gameTime time.Time

	// канал, в который поступают управляющие события
	controlChan chan ControlMessage

	// канал, в который поступают действия игроков
	inputChan <-chan PlayerInput

	// заготовка под канал по которому будет отправляться информация для игроков
	outputChan chan interface{}

	// канал в окторый мы будем писать стейт целиком. Можно использовать для отладки
	monitorChan chan<- *GameState

	players map[uint]Player

	Mode               SimulationMode // режим логики (нормальный/пошаговый/воспроизведение)
	flagShouldStop     bool           // говорит что нужно остановить логику
	flagShouldSimulate bool           // сбрасывается при начале mainStep

	simulationStepDuration     time.Duration // сколько виртуального времени проходит за один тик
	simulationStepRealDuration time.Duration // сколько реального времени проходит за один тик

	history      []HistoryEntry
	State        *GameState
	finishedChan chan bool // канал, единственная функция которого - быть открытым или закрытым
}

func (l *Logic) NextSimulationTime() time.Time {
	if l.Mode == SimulationModeContinuous {
		return l.prevTickRealTime.Add(l.simulationStepRealDuration)
	} else {
		// never
		return time.Unix(math.MaxInt64, 0)
	}
}

// Говорит наступило ли время для симуляции
func (l *Logic) ShouldSimulate() bool {
	return l.flagShouldSimulate ||
		(l.Mode == SimulationModeContinuous && l.NextSimulationTime().Before(time.Now()))
}

func (l *Logic) receiveInputsUntilShouldSimulate() []PlayerInput {
	// listen to control chan in parallel.
	log.Println("receiving inputs")

	inputsReceived := 0
	controlsReceived := 0
	defer func() {
		log.Printf("received %d inputs, %d controls\n", inputsReceived, controlsReceived)
	}()

	var inputsBuffer []PlayerInput

	var timeout time.Duration

	if l.Mode == SimulationModeContinuous {
		// receive everything from input chan until next simulation time

		timeout = time.Until(l.NextSimulationTime())
		if timeout <= 0 {
			log.Println("time to tick has already come!")
			// время уже прошло!
			return inputsBuffer
		}
	} else if l.Mode == SimulationModeStepByStep {
		// receive everything from input chang until simulateMessage
		timeout = math.MaxInt64
	} else if l.Mode == SimulationModeReplay {
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

	stopReceiving := false

	if l.ShouldSimulate() {
		log.Println("strange - we already should simulate")
		stopReceiving = true
	}

	for !stopReceiving && !l.ShouldSimulate() {
		select {
		case msg := <-l.inputChan:
			inputsBuffer = append(inputsBuffer, msg)
			log.Printf("input message %v\n", msg)
			inputsReceived++
		case ctrl := <-l.controlChan:
			controlsReceived++
			log.Printf("control message %v\n", ctrl)
			switch ctrl {
			case ControlMessageStop:
				l.flagShouldStop = true
				stopReceiving = true
			case ControlMessageSimulate:
				// если у нас обычная симуляция, то мы просто заставим шаг произойти преждевременно
				l.flagShouldSimulate = true
				stopReceiving = true
			case ControlMessageChangeModeContinuous, ControlMessageChangeModeStepByStep, ControlMessageChangeModeReplay:
				var newMode SimulationMode
				switch ctrl {
				case ControlMessageChangeModeContinuous:
					newMode = SimulationModeContinuous
				case ControlMessageChangeModeStepByStep:
					newMode = SimulationModeStepByStep
				case ControlMessageChangeModeReplay:
					newMode = SimulationModeReplay
				}
				if l.Mode != newMode {
					l.Mode = newMode
					stopReceiving = true
				}
			default:
				log.Printf("Unknown control message %d", ctrl)
			}
			log.Println("stopping wait loop")
		case <-timer.C:
			// пришло время, больше ничего читать не будем
			log.Println("timer fired")
			stopReceiving = true
			timerFired = true
		}
	}

	return inputsBuffer
}

func (l *Logic) applyPlayerInputs(state *GameState, inputs []PlayerInput) *GameState {
	log.Println("apply inputs stub")
	for _, inp := range inputs {
		// TODO
		switch inp.Action {
		case PlayerActionMove:
			log.Println("Player actions not implemented")
		case PlayerActionAbility:
			log.Println("Player abilities not implemented")
		}
	}
	return state
}

func (l *Logic) sendStateToPlayers(state *GameState) {
	for _, player := range l.players {
		player.SendState(state.GetPlayerState(player.Id))
	}

	// send state if available
	select {
	case l.monitorChan <- state:
	default:
	}
}

// Main step of server. Apply pending players inputs
// run physics engine simulation,
// calculate states for all players,
// sending states to respective players
func (l *Logic) mainStep(inputs []PlayerInput) {
	log.Println("main step started")

	l.State = l.State.Copy()
	l.gameTick++
	l.gameTime.Add(l.simulationStepDuration)

	l.State = l.applyPlayerInputs(l.State, inputs)
	l.State = l.applyQueuedEvents(l.State)
	l.State = l.State.ProcessSimulationStep(l.simulationStepDuration)

	l.prevTickRealTime = time.Now()
	l.flagShouldSimulate = false

	l.history = append(l.history, HistoryEntry{
		state:    l.State,
		tick:     l.gameTick,
		gameTime: l.gameTime,
		realTime: time.Now(),
	})

	l.sendStateToPlayers(l.State)

	log.Println("main step finished")
}

// основной цикл логики
// receive inputs, put them into queue
// when time has come - simulate next step
func (l *Logic) mainLoop() {
	log.Println("starting main loop")
	var inputsBuffer []PlayerInput

	for !l.flagShouldStop {
		// Read the inputs and put into queue until simulation time comes.
		// But reading can break in case of receiving control message. So before simulation
		// we should check if simulation time has come really
		inputsBuffer = append(inputsBuffer, l.receiveInputsUntilShouldSimulate()...)

		if l.ShouldSimulate() {
			l.mainStep(inputsBuffer)

			// Clear inputs buffer only after simulation,
			// because receiving could be aborted before simulation time
			inputsBuffer = inputsBuffer[:0]
		}
	}
	log.Println("main loop stopped")
	l.stop()
}

func (l *Logic) Start() {
	go l.mainLoop()
}

// Говорит логике остановиться
// В реальности логика остановится только тогда, когда она проверит этот канал,
// а это происходит тогда, когда она принимает инпуты
func (l *Logic) stop() {
	if l.monitorChan != nil {
		close(l.monitorChan)
	}
	if l.outputChan != nil {
		close(l.outputChan)
	}
	close(l.finishedChan)
}

// предполагаем, что события будут инициироваться не только игроками, но и самой логикой.
// Они будут складываться в очередь и срабатывать в назначенное время.
func (l *Logic) applyQueuedEvents(state *GameState) *GameState {
	// TODO
	return state
}

func (l *Logic) SetMonitorChan(ch chan<- *GameState) {
	l.monitorChan = ch
}

// синхронная операция - ждёт пока логика действительно не остановится
func (l *Logic) Stop() chan bool {
	l.controlChan <- ControlMessageStop
	// ждём пока этот канал не закроют
	return l.finishedChan
}

func NewLogic(controlChan chan ControlMessage, inputChan <-chan PlayerInput, mode SimulationMode, stepTime, stepRealTime time.Duration) *Logic {
	logic := &Logic{
		gameTick:                   0,
		prevTickRealTime:           time.Time{},
		gameTime:                   time.Time{},
		controlChan:                controlChan,
		inputChan:                  inputChan,
		outputChan:                 make(chan interface{}, 10),
		monitorChan:                nil,
		players:                    make(map[uint]Player),
		Mode:                       mode,
		flagShouldStop:             false,
		flagShouldSimulate:         false,
		simulationStepDuration:     stepTime,
		simulationStepRealDuration: stepRealTime,
		history:                    make([]HistoryEntry, 0, 10),
		State:                      NewGameState(),
		finishedChan:               make(chan bool),
	}

	return logic
}

package next

import (
	"log"
	"time"

	"github.com/porfirion/server2/world"
)

type LogicImpl struct {
	gameTick             uint
	pendingPlayersInputs []PlayerInput

	stopChan  chan bool
	inputChan chan interface{}

	players map[uint]Player

	worldMap *world.WorldMap
}

func (l *LogicImpl) applyPlayerInput(input PlayerInput) {}

func (l LogicImpl) applyPlayersInputsUpTo(currentTick uint) {
	var ind int
	for ind = 0; ind < len(l.pendingPlayersInputs) && l.pendingPlayersInputs[ind].GetGameTick() <= currentTick; ind++ {
		l.applyPlayerInput(l.pendingPlayersInputs[ind])
	}

	l.pendingPlayersInputs = nil
}

func (l *LogicImpl) NextSimulationTime() time.Time {
	// TODO заглушка
	log.Println("Error: nextSimulationTime stub")
	return time.Now().Add(time.Hour)
}

// Говорит наступило ли время для симуляции
func (l *LogicImpl) ShouldSimulate() bool {
	return l.NextSimulationTime().Before(time.Now())
}

func (l *LogicImpl) receiveInputsUntil(until time.Time) {
	var timeout = until.Sub(time.Now())

	if timeout <= 0 {
		// время уже прошло!
		return
	}

	timer := time.NewTimer(timeout)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	// сначала забираем весь пользовательский инпут
	// теоретически пользователи могут заспамить сообщениями там,
	// что логика не будет успевать отработать
	// потому нужно ограничивать максимальнео количество
	countReceived := 0

receivingInputs:
	for !l.ShouldSimulate() {
		select {
		case msg := <-l.inputChan:
			log.Println("Logic: received message ", msg)
			countReceived++
		case <-l.stopChan:
			log.Println("Logic: SHOULD STOP!!!")
		case <-timer.C:
			// пришло время, больше ничего читать не будем
			break receivingInputs
		}
	}
}

// Calculate visible objects
func (l *LogicImpl) calculateUsersStates() map[uint]PlayerState {
	states := make(map[uint]PlayerState)

	return states
}

// Run physics engine simulation
func (l *LogicImpl) simulateStep() {

}

// Main step of server. Apply pending players inputs
// run physics engine simulation,
// calculate states for all players,
// sending states to respective players
func (l *LogicImpl) mainStep() {
	// ждём когда придёт время симулировать
	if l.ShouldSimulate() {
		l.gameTick++

		l.applyPlayersInputsUpTo(l.gameTick)

		l.simulateStep()

		states := l.calculateUsersStates()

		for playerId, state := range states {
			player := l.players[playerId]
			player.SendState(state)
		}
	}
}

// основной цикл логики
// получаем весь инпут, складываем его в очередь
// как только настаёт время - выполняем основной шаг (mainStep)
func (l *LogicImpl) mainLoop() {
	var shouldExit = false

	for !shouldExit {
		// вычитываем инпут и кладём в очередь
		// до тех пор, пока не придёт время симулировать
		l.receiveInputsUntil(l.NextSimulationTime())

		l.mainStep()
	}

}

func (l *LogicImpl) Start() {
	go l.mainLoop()
}
func (l *LogicImpl) Stop() {

}

func NewLogic() *LogicImpl {
	logic := &LogicImpl{
		gameTick:             0,
		pendingPlayersInputs: nil,
		stopChan:             make(chan bool),
		inputChan:            make(chan interface{}, 10),
		players:              make(map[uint]Player),
		worldMap:             world.NewWorldMap(),
	}
	return logic
}

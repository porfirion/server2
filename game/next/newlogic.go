package next

import (
	"time"
)

type LogicImpl struct {
	gameTick             uint
	pendingPlayersInputs []PlayerInput

	stopChan  chan bool
	inputChan chan interface{}
}

func (l *LogicImpl) applyPlayerInput(input PlayerInput) {}

func (l LogicImpl) applyPlayersInputsUpTo(currentTick uint) {
	var ind int
	for ind = 0; ind < len(l.pendingPlayersInputs) && l.pendingPlayersInputs[ind].GetGameTick() <= currentTick; ind++ {
		l.applyPlayerInput(l.pendingPlayersInputs[ind])
	}

	l.pendingPlayersInputs = l.pendingPlayersInputs[ind:]
}

func (l *LogicImpl) NextSimulationTime() time.Time {
	// TODO заглушка
	return time.Now().Add(time.Hour)
}

// Говорит наступило ли время для симуляции
func (l *LogicImpl) ShouldSimulate() bool {
	return l.NextSimulationTime().Before(time.Now());
}

func (l *LogicImpl) simulateStep() {
	l.gameTick++
	l.applyPlayersInputsUpTo(l.gameTick)
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
	const MAX_RECEIVE = 50
	countReceived := 0

receivingInputs:
	for countReceived < MAX_RECEIVE && !l.ShouldSimulate() {
		select {
		case msg := <-l.inputChan:
			println(msg)
			countReceived++
		case <- timer.C:
			// пришло время, больше ничего читать не будем
			break receivingInputs
		}
	}
}

func (l *LogicImpl) mainLoop() {
	var shouldExit = false

	for !shouldExit {
		// вычитываем инпут и кладём в очередь
		l.receiveInputsUntil(l.NextSimulationTime())

		// ждём когда придёт время симулировать
		if l.ShouldSimulate() {
			// симулировать нужно уже сейчас!
			l.simulateStep()
		} else {
			// нет, время ещё не пришло - дождёмся
		}

		for l.ShouldSimulate() {
			l.simulateStep()
		}
	}

}

func (l *LogicImpl) Start() {
	go l.mainLoop()
}

func (l *LogicImpl) Stop() {

}

func NewLogic() {}

package game

import (
	"github.com/porfirion/server2/service"
	"github.com/porfirion/server2/messages"
	"log"
	"time"
)

type GameLogicService struct {
	*service.BasicService
	logic                 *GameLogic
}

func (s *GameLogicService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (s *GameLogicService) Start() {
	go s.startReading()
}

func (s *GameLogicService) startReading() {
	s.WaitForRegistration()

	// ура! теперь нам есть куда писать!!!
	go s.startWriting()

	for msg := range s.IncomingMessages {
		// TODO переделать!
		// пока просто прокидываем сообщения внутрь логики

		log.Printf("Logic: Can't cast service message to UserMessage %v\n", msg.MessageData)
	}
}

func (s *GameLogicService) startWriting() {
	for msg := range s.logic.OutgoingMessages {
		log.Println("Message from service to pass to broker", msg)
		// TODO переделать!
		// пока тупо прокидываем сообщения из логики в брокер (но он их не поймёт)
		// TODO FORTEST ONLY
		//s.OutgoingMessages <- msg
	}
}

func NewService() *GameLogicService {
	logic := &GameLogic{
		IncomingMessages: make(messages.UserMessagesChannel, 10),
		OutgoingMessages: make(messages.ServerMessagesChannel, 10),
		Params: LogicParams{
			SimulateByStep:           true,                   // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
			SimulationStepTime:       500 * time.Millisecond, // сколько виртуального времени проходит за один шаг симуляции
			SimulationStepRealTime:   500 * time.Millisecond, // сколько реального времени проходит за один шаг симуляции
			SendObjectsTimeout:       time.Millisecond * 500,
			MaxSimulationStepsAtOnce: 10,
		},
	}
	// стартуем логику. она готова, чтобы принимать и обрабатывать соощения
	go logic.Start()

	return &GameLogicService{
		BasicService: service.NewBasicService(service.TypeLogic),
		logic:        logic,
	}
}

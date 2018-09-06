package game

import (
	"github.com/porfirion/server2/service"
	"log"
)

type LogicService struct {
	*service.BasicService
	logic *GameLogic
}

func (s *LogicService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (s *LogicService) Start() {
	go func() {
		// пока мы не зарегистрируемся в брокере - читать неоткуда и писать некуда
		s.WaitForRegistration()

		go s.startReading()
		go s.startWriting()
	}()
}

func (s *LogicService) startReading() {
	for msg := range s.IncomingMessages {
		// TODO переделать!
		// пока просто прокидываем сообщения внутрь логики

		log.Printf("Logic: Can't cast service message to UserMessage %v\n", msg.MessageData)
	}
}

func (s *LogicService) startWriting() {
	for msg := range s.logic.OutgoingMessages {
		log.Println("Message from service to pass to broker", msg)
		// TODO переделать!
		// пока тупо прокидываем сообщения из логики в брокер (но он их не поймёт)
		// TODO FORTEST ONLY
		//s.OutgoingMessages <- msg
	}
}

func NewService() *LogicService {
	logic := NewGameLogic()

	return &LogicService{
		BasicService: service.NewBasicService(service.TypeLogic),
		logic:        logic,
	}
}

package service

import "fmt"

type GameLogicService struct {
	*BasicService
	Logic                 *GameLogic
	LogicOutgoingMessages ServerMessagesChannel
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

		fmt.Printf("Can't cast service message to UserMessage %v\n", msg.MessageData)
	}
}

func (s *GameLogicService) startWriting() {
	for msg := range s.LogicOutgoingMessages {
		fmt.Println("Message from service to pass to broker", msg)
		// TODO переделать!
		// пока тупо прокидываем сообщения из логики в брокер (но он их не поймёт)
		// TODO FORTEST ONLY
		//s.OutgoingMessages <- msg
	}
}

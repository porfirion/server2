package logic

import "fmt"

type GameLogicService struct {
	*BasicService
	Logic                 *GameLogic
	LogicOutgoingMessages ServerMessagesChannel
}

func (s *GameLogicService) Start() {
	go s.startReading()
}

func (s *GameLogicService) startReading() {
	regMsg := <-s.IncomingMessages
	dt := regMsg.MessageData.(BrokerRegisterServiceResponse)
	s.Id = dt.Id
	s.OutgoingMessages = dt.Ch

	// ура! теперь нам есть куда писать!!!
	go s.startWriting()

	for msg := range s.IncomingMessages {
		// TODO переделать!
		// пока просто прокидываем сообщения внутрь логики
		if data, ok := msg.MessageData.(UserMessage); ok {
			s.Logic.IncomingMessages <- data
		} else {
			fmt.Printf("Can't cast to UserMessage %#v\n", msg.MessageData)
		}
	}
}

func (s *GameLogicService) startWriting() {
	for msg := range s.LogicOutgoingMessages {
		fmt.Println("Message from logic to pass to broker", msg)
		// TODO переделать!
		// пока тупо прокидываем сообщения из логики в брокер (но он их не поймёт)
		// TODO FORTEST ONLY
		//s.OutgoingMessages <- msg
	}
}

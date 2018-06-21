package logic

import "fmt"

type ChatService struct {
	*BasicService
}

// async start method
func (s *ChatService) Start() {
	go s.StartReading()
}

func (s *ChatService) StartReading() {
	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
	regMsg := <-s.IncomingMessages
	dt := regMsg.MessageData.(BrokerRegisterServiceResponse)
	s.Id = dt.Id
	s.OutgoingMessages = dt.Ch

	for msg := range s.IncomingMessages {
		fmt.Println(msg)
	}
}

func NewChat(bs *BasicService) *ChatService {
	return &ChatService{
		bs,
	}
}

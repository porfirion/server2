package chat

import (
	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/service"
	"log"
)

type ChatService struct {
	*service.BasicService
}

// async start method
func (s *ChatService) Start() {
	go s.StartReading()
}

func (s *ChatService) StartReading() {
	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		log.Printf("ChatService: %#v", msg.MessageData)
		switch msg.MessageData.(type) {
		case *messages.TextMessage, messages.TextMessage:
			s.SendMessageToBroker(msg.MessageData, 0, service.TypeNetwork, 0, nil)
		default:
			log.Printf("Chat: unexpected message type %T\n", msg.MessageData)
		}

	}
}

func (s *ChatService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func NewService() *ChatService {
	return &ChatService{
		service.NewBasicService(service.TypeChat),
	}
}

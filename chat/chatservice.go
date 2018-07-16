package chat

import (
	"github.com/porfirion/server2/service"
	"log"
	"github.com/porfirion/server2/messages"
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
		log.Println("ChatService: ", msg.MessageData)
		switch msg.MessageData.(type) {
		case *messages.TextMessage:
			s.SendMessage(msg.MessageData, 0, service.TypeNetwork, 0, nil)
		default:
			log.Printf("Chat: unexpected message type %T\n", msg.MessageData)
		}

	}
}

func (s *ChatService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func NewChat(bs *service.BasicService) *ChatService {
	return &ChatService{
		bs,
	}
}

package chat

import (
	"fmt"
	"github.com/porfirion/server2/service"
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
		fmt.Println("CHAT: ", msg)
		s.SendMessage(msg.MessageData, 0, service.TypeNetwork, 0, nil)
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

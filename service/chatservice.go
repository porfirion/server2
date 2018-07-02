package service

import (
	"fmt"
)

type ChatService struct {
	*BasicService
}

// async start method
func (s *ChatService) Start() {
	go s.StartReading()
}

func (s *ChatService) StartReading() {
	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		fmt.Println(msg)
	}
}

func (s *ChatService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func NewChat(bs *BasicService) *ChatService {
	return &ChatService{
		bs,
	}
}

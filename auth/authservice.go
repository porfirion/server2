package auth

import (
	"github.com/porfirion/server2/service"
	"log"
)

type AuthService struct {
	*service.BasicService
}

func (auth *AuthService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (auth *AuthService) Start() {
	go auth.startReading()
}

func (auth *AuthService) startReading() {
	auth.WaitForRegistration()

	for msg := range auth.IncomingMessages {
		log.Printf("AuthService: %#v\n", msg.MessageData)
	}
}

func NewService() *AuthService {
	return &AuthService{
		service.NewBasicService(service.TypeAuth),
	}
}

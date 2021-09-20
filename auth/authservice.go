package auth

import (
	"log"

	"github.com/porfirion/server2/messages"
	"github.com/porfirion/server2/service"
)

type Service struct {
	*service.BasicService
	nextId uint64
}

func (auth *Service) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (auth *Service) Start() {
	go auth.startReading()
}

func (auth *Service) startReading() {
	auth.WaitForRegistration()

	for msg := range auth.IncomingMessages {
		log.Printf("AuthService: %#v\n", msg.MessageData)

		switch data := msg.MessageData.(type) {
		case *messages.AuthMessage:
			id := auth.nextId
			auth.nextId++
			auth.SendMessageToBroker(msg.MessageData, 0, service.TypeChat, 0, nil)
			auth.SendMessageToBroker(messages.LoginMessage{
				Id:   id,
				Name: data.Name,
			}, 0, service.TypeLogic, 0, nil)
		default:
			log.Printf("unknown message type %+v", msg.MessageData)
		}

	}
}

func NewService() *Service {
	return &Service{
		BasicService: service.NewBasicService(service.TypeAuth),
		nextId:       1,
	}
}

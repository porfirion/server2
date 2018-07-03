package service

import (
	"github.com/porfirion/server2/network"
	"fmt"
	"encoding/json"
)

type NetworkService struct {
	*BasicService
	pool *network.ConnectionsPool
}

func (s *NetworkService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (s *NetworkService) Start() {
	s.startReceivingFromBroker()
}

func (s *NetworkService) SetPool(pool *network.ConnectionsPool) {
	s.pool = pool
	go s.startReadingFromClients()
}

func (s *NetworkService) startReceivingFromBroker() {
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		fmt.Println(msg)
		if bytes, err := json.Marshal(msg.MessageData); err == nil {
			s.pool.OutgoingMessages <- network.MessageForClient{
				Targets:     msg.DestinationServiceClients,
				MessageType: msg.MessageType,
				Data:        bytes,
			}
		} else {
			fmt.Printf("Error matshalling message %v\n", msg)
		}
	}
}

func (s *NetworkService) startReadingFromClients() {
	// TODO здесь ещё нужна проверка на то, зарегистрировали ли нас и есть ли нам куда писать
	for msg := range s.pool.IncomingMessages {
		// TODO сейчас в пробкер отправляются сырые байты и никакого парсинга не происходит
		// также не указывается целевой сервис, в который мы отправляем эти данные
		s.SendMessage(msg.MessageType, msg.Data, msg.ClientId, 0, 0, nil)
	}
}

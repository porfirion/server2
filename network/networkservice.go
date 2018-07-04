package network

import (
	"github.com/porfirion/server2/service"
)

type NetworkService struct {
	*service.BasicService
	pool *ConnectionsPool
}

func (s *NetworkService) GetRequiredMessageTypes() []uint {
	return []uint{}
}

func (s *NetworkService) Start() {
	s.startReceivingFromBroker()
}

func (s *NetworkService) SetPool(pool *ConnectionsPool) {
	s.pool = pool
	go s.startReadingFromClients()
}

func (s *NetworkService) startReceivingFromBroker() {
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		s.pool.OutgoingMessages <- MessageForClient{
			Targets: msg.DestinationServiceClients,
			Data:    msg.MessageData,
		}
	}
}

func (s *NetworkService) startReadingFromClients() {
	// TODO здесь ещё нужна проверка на то, зарегистрировали ли нас и есть ли нам куда писать
	for msg := range s.pool.IncomingMessages {
		// TODO сейчас в пробкер отправляются сырые байты и никакого парсинга не происходит
		// также не указывается целевой сервис, в который мы отправляем эти данные
		s.SendMessage(msg.Data, msg.ClientId, service.TypeLogic, 0, nil)
	}
}
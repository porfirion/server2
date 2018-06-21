package logic

import (
	"github.com/porfirion/server2/network"
	"fmt"
	"encoding/json"
)

type NetworkService struct {
	*BasicService
	pool *network.ConnectionsPool
}

func (s *NetworkService) Start() {
	s.startReceivingFromBroker()
}

func (s *NetworkService) SetPool(p *network.ConnectionsPool) {
	s.pool = p
	s.startReadingFromClients()
}

func (s *NetworkService) startReceivingFromBroker() {
	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
	regMsg := <-s.IncomingMessages
	dt := regMsg.MessageData.(BrokerRegisterServiceResponse)
	s.Id = dt.Id
	s.OutgoingMessages = dt.Ch

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
		fmt.Printf("Incoming message from pool %v\n", msg)
		s.SendMessage(msg.MessageType, msg.Data, msg.ClientId, 0, 0, nil)
	}
}

package logic

import "github.com/porfirion/server2/utils"

/**
 * Брокер, который разруливает в какой сервис отправлять сообщение
 */
type MessageBroker interface {
	Send(msg ServiceMessage)      // отправка сообщения в брокер для конкретного получаетля
	Broadcast(msg ServiceMessage) // отправка широковещательного сообщения
	RegisterService(svc Service)  // регистрация нового сервиса в брокере
	Start()
}

// IMPLEMENTATION

type BrokerBroadcastMessage struct {
	Data ServiceMessage
}

type BrokerTargetMessage struct {
	Data ServiceMessage
}

type BrokerRegisterServiceMessage struct {
	Service Service
}

type BrokerImplementation struct {
	utils.IdGenerator
	mainChan chan interface{}
	services map[uint64]Service
}

func (broker *BrokerImplementation) Start() {
	go broker.StartReading()

}

func (broker *BrokerImplementation) StartReading() {
	for untypedMessage := range broker.mainChan {
		switch msg := untypedMessage.(type) {
		case BrokerTargetMessage:

		case BrokerBroadcastMessage:

		case BrokerRegisterServiceMessage:
			nextId := broker.NextId()
			broker.services[nextId] = msg.Service
			msg.Service.Register(nextId, broker.mainChan)
		}
	}
}

func (broker *BrokerImplementation) Send(msg ServiceMessage) {
	broker.mainChan <- msg
}

func (broker *BrokerImplementation) Broadcast(msg ServiceMessage) {
	panic("implement me")
}

func (broker *BrokerImplementation) RegisterService(svc Service) {
	serviceId := broker.NextId()
	broker.services[serviceId] = svc
	svc.Register(serviceId, broker.mainChan)
}

func NewBroker() MessageBroker {
	return &BrokerImplementation{
		utils.NewIdGenerator(1),
		make(chan interface{}),
		make(map[uint64]Service),
	}
}

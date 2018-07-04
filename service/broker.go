package service

import (
	"github.com/porfirion/server2/utils"
	"log"
)


// Брокер, который разруливает в какой сервис отправлять сообщение
type MessageBroker interface {
	Send(msg ServiceMessage)      // отправка сообщения в брокер для конкретного получателя
	RegisterService(svc Service)  // регистрация нового сервиса в брокере
	Start()
}

// IMPLEMENTATION

type BrokerRegisterServiceMessage struct {
	TypedMessageStub
	Service Service
}

type BrokerRegisterServiceResponse struct {
	TypedMessageStub
	Id uint64
	Ch chan ServiceMessage
}

type TypedMessageStub struct {}

func (m TypedMessageStub) GetType() uint64 {
	log.Println("Warning! Using TypedMessageStub")
	return 1
}

type BrokerImplementation struct {
	utils.IdGenerator
	mainChan chan ServiceMessage
	services map[uint64]Service // each service has unique id
	serviceByTypes map[uint64][]Service // many services of the same type can exist
}

func (broker *BrokerImplementation) Start() {
	go broker.StartReading()

}

func (broker *BrokerImplementation) StartReading() {
	for untypedMessage := range broker.mainChan {
		switch msg := untypedMessage.MessageData.(type) {

		case BrokerRegisterServiceMessage:
			nextId := broker.NextId()
			service := msg.Service
			serviceType := service.GetType()

			broker.services[nextId] = service
			if broker.serviceByTypes[serviceType] == nil {
				broker.serviceByTypes[serviceType] = []Service{service}
			} else {
				broker.serviceByTypes[serviceType] = append(broker.serviceByTypes[serviceType], service)
			}

			service.Register(nextId, broker.mainChan)
		default:
			log.Printf("Broker: Unexpected message %T\n", msg)
			if untypedMessage.DestinationServiceId != 0 {
				if dest := broker.services[untypedMessage.DestinationServiceId]; dest != nil {
					dest.Deliver(untypedMessage)
				} else {
					log.Printf("Broker: can't find destination service #%d\n", untypedMessage.DestinationServiceId)
				}
			} else if untypedMessage.DestinationServiceType != 0 {
				if dests := broker.serviceByTypes[untypedMessage.DestinationServiceType]; dests != nil && len(dests) > 0 {
					for _, dest := range dests {
						dest.Deliver(untypedMessage)
					}
				} else {
					log.Printf("Broker: can't find any services with type %d\n", untypedMessage.DestinationServiceType)
				}
			} else {
				for _, service := range broker.services {
					service.Deliver(untypedMessage)
				}
			}
		}
	}
}

func (broker *BrokerImplementation) Send(msg ServiceMessage) {
	broker.mainChan <- msg
}

func (broker *BrokerImplementation) RegisterService(svc Service) {
	broker.Send(ServiceMessage{
		SourceServiceType:         0,
		SourceServiceId:           0,
		SourceServiceClient:       0,
		DestinationServiceType:    0,
		DestinationServiceId:      0,
		DestinationServiceClients: nil,
		MessageData:               BrokerRegisterServiceMessage{Service: svc},
	})
}

func NewBroker() MessageBroker {
	return &BrokerImplementation{
		IdGenerator:    utils.NewIdGenerator(1),
		mainChan:       make(chan ServiceMessage),
		services:       make(map[uint64]Service),
		serviceByTypes: make(map[uint64][]Service),
	}
}

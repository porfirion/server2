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

// фнукция, которая определяет куда направлять сообщение, если адресат не был задан
// Вызывается только когда serviceId и serviceType не указаны
// Если эта функция вернёт 0 - значит сообщение никуда не будет отправлено.
//
// Возвращает тип сервиса, в который нужно передать сообщение
type MessageRouter func(msg ServiceMessage) uint64

type BrokerImplementation struct {
	utils.IdGenerator
	mainChan       chan ServiceMessage
	services       map[uint64]Service // each service has unique id
	serviceByTypes map[uint64]Service // пока упрощённое - по одному сервису каждого типа
	messageRouter  MessageRouter
}

func (broker *BrokerImplementation) Start() {
	go broker.StartReading()

}

func (broker *BrokerImplementation) StartReading() {
	for serviceMessage := range broker.mainChan {

		//log.Println(serviceMessage)

		switch msg := serviceMessage.MessageData.(type) {

		case BrokerRegisterServiceMessage:
			//log.Println("registering service")
			nextId := broker.NextId()
			service := msg.Service
			serviceType := service.GetType()

			broker.services[nextId] = service
			broker.serviceByTypes[serviceType] =  service

			//log.Println("Broker: sending registration to service")
			service.StoreRegisteration(nextId, broker.mainChan)
			//log.Printf("Broker: registered service %d %s\n", service.GetId(), service.GetType())
		default:
			//log.Printf("Broker: Delivering message type %T\n", msg)
			if serviceMessage.DestinationServiceId != 0 {
				if dest := broker.services[serviceMessage.DestinationServiceId]; dest != nil {
					dest.Deliver(serviceMessage)
				} else {
					log.Printf("Broker: can't find destination service #%d\n", serviceMessage.DestinationServiceId)
				}
			} else if serviceMessage.DestinationServiceType != 0 {
				if dest := broker.serviceByTypes[serviceMessage.DestinationServiceType]; dest != nil {
					dest.Deliver(serviceMessage)
				} else {
					log.Printf("Broker: can't find service type %d\n", serviceMessage.DestinationServiceType)
					log.Fatal("Broker: can't find service type %d\n", serviceMessage.DestinationServiceType)
				}
			} else if destinationServiceType := broker.messageRouter(serviceMessage); destinationServiceType != 0 {
				if dest := broker.serviceByTypes[destinationServiceType]; dest != nil {
					dest.Deliver(serviceMessage)
				} else {
					log.Printf("Broker: can't find service type %d\n", destinationServiceType)
				}
			} else {
				log.Printf("Don't know where to deliver mesage type %T\n", serviceMessage.MessageData)
				//log.Fatalf("Don't know where to deliver mesage type %T\n", serviceMessage.MessageData)
				//broker.deliverAll(serviceMessage)
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

func NewBroker(messageRouter MessageRouter) MessageBroker {

	return &BrokerImplementation{
		IdGenerator:    utils.NewIdGenerator(1),
		mainChan:       make(chan ServiceMessage, 100),
		services:       make(map[uint64]Service),
		serviceByTypes: make(map[uint64]Service),
		messageRouter:  messageRouter,
	}
}

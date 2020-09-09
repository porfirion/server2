package service

import (
	"errors"
	"log"
	"time"
)

type ServiceType uint64

// Начитавшись хабра (https://habrahabr.ru/company/mailru/blog/220359/)
// пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте
// Более того - это скорее даже мешает - всё валится в одну кучу.
// Также авторизация остаётся незакрытым вопросом.
type Service interface {
	Deliver(msg ServiceMessage)                           // через этот метод сообщения закидываются в сервис
	StoreRegisteration(id uint64, ch chan ServiceMessage) // уведомляет сервис о том, что он был зарегистрирован
	//GetRequiredMessageTypes() []uint            // отдаёт список ожидаемых сообщений
	Start()
	GetType() ServiceType
	GetId() uint64
}

// сообщение, которое ходит между сервисами
type ServiceMessage struct {
	SourceServiceType   ServiceType
	SourceServiceId     uint64
	SourceServiceClient uint64 // по идее у нас не может быть только один отправитель или не быть его вообще

	DestinationServiceType    ServiceType
	DestinationServiceId      uint64
	DestinationServiceClients []uint64 // зато может быть много получателей

	MessageData TypedMessage
}

type TypedMessage interface {
	GetType() uint64
}

const (
	TypeLogic   ServiceType = 1
	TypeAuth    ServiceType = 2
	TypeNetwork ServiceType = 3
	TypeChat    ServiceType = 4
)

type BasicService struct {
	Id   uint64
	Type ServiceType

	IncomingMessages chan ServiceMessage // это канал для полечения сообщений от брокера
	OutgoingMessages chan ServiceMessage // это канал для отправки и нам должен дать его сам брокер, когда зарегистрирует наш сервис
}

func (service *BasicService) GetType() ServiceType {
	return service.Type
}

func (service *BasicService) GetId() uint64 {
	return service.Id
}

// Через этот метод брокер отправляет сообщения в сервис
func (service *BasicService) Deliver(msg ServiceMessage) {
	service.IncomingMessages <- msg
}

func (service *BasicService) StoreRegisteration(serviceId uint64, out chan ServiceMessage) {
	if (service).IncomingMessages == nil {
		log.Fatal("incoming messages chan is not initialized")
	}
	service.IncomingMessages <- ServiceMessage{
		MessageData: BrokerRegisterServiceResponse{
			Id: serviceId,
			Ch: out,
		},
	}
}

// Отправляет сообщение брокеру
func (service *BasicService) SendMessageToBroker(
	msg TypedMessage,
	sourceClientId uint64,
	targetServiceType ServiceType,
	targetServiceId uint64,
	targets []uint64) error {

	if service.OutgoingMessages == nil {
		return errors.New("no output channel")
	}

	serverMessage := ServiceMessage{
		SourceServiceType:      service.Type,
		SourceServiceId:        service.Id,
		SourceServiceClient:    sourceClientId,
		DestinationServiceType: targetServiceType,
		DestinationServiceId:   targetServiceId,
		MessageData:            msg,
	}

	if len(targets) > 0 {
		serverMessage.DestinationServiceClients = targets
	}

	t := time.NewTimer(time.Millisecond * 100)
	select {
	case service.OutgoingMessages <- serverMessage:
		t.Stop()
		return nil
	case <-t.C:
		return errors.New("output write timeout")
	}
}

// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
func (service *BasicService) WaitForRegistration() {
	//log.Println("BasicService: wating for registration")
	dt := (<-service.IncomingMessages).MessageData.(BrokerRegisterServiceResponse)
	service.Id = dt.Id
	service.OutgoingMessages = dt.Ch
	//log.Println("BasicService: Registration received")
}

func NewBasicService(serviceType ServiceType) *BasicService {
	return &BasicService{
		Id:               0,
		Type:             serviceType,
		IncomingMessages: make(chan ServiceMessage),
		OutgoingMessages: nil,
	}
}

//Пример регистрации сервиса:

//type Svc struct {
//	*BasicService
//}
//
//func (s *Svc) Start() {
//	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
//  s.WaitForRegistration()
//
//	for msg := range s.IncomingMessages {
//		fmt.Println(msg)
//	}
//}
//
//func example() {
//	var broker MessageBroker
//  broker.Start()
//	svc := &Svc{NewBasicService(1)}
//	svc.Start()
//	broker.RegisterService(svc)
//}

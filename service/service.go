package service

import (
	"errors"
	"time"
)

// Начитавшись хабра (https://habrahabr.ru/company/mailru/blog/220359/)
// пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте
// Более того - это скорее даже мешает - всё валится в одну кучу.
// Также авторизация остаётся незакрытым вопросом.
type Service interface {
	Deliver(msg ServiceMessage)                 // через этот метод сообщения закидываются в сервис
	Register(id uint64, ch chan ServiceMessage) // уведомляет сервис о том, что он был зарегистрирован
	GetRequiredMessageTypes() []uint            // отдаёт список ожидаемых сообщений
	Start()
	GetType() uint64
	GetId() uint64
}

// сообщение, которое ходит между сервисами
type ServiceMessage struct {
	SourceServiceType   uint64
	SourceServiceId     uint64
	SourceServiceClient uint64 // по идее у нас не может быть только один отправитель

	DestinationServiceType    uint64
	DestinationServiceId      uint64
	DestinationServiceClients []uint64 // зато может быть много получателей

	// ВАЖНО! это совсем не тот тип, который ходит по сети.
	// Типы сообщений, используемые в сети, являются частью публичного API
	// и могут вообще никак не пересекаться с внутренним обозначением типов
	MessageType uint64
	MessageData interface{}
}

const (
	TypeLogic   uint64 = 1
	TypeAuth           = 2
	TypeNetwork        = 3
	TypeChat           = 4
)

type BasicService struct {
	Id   uint64
	Type uint64

	IncomingMessages chan ServiceMessage // это канал для полечения сообщений от брокера
	OutgoingMessages chan ServiceMessage // это канал для отправки и нам должен дать его сам брокер, когда зарегистрирует наш сервис
}

func (service *BasicService) GetType() uint64 {
	return service.Type
}

func (service *BasicService) GetId() uint64 {
	return service.Id
}

// Через этот метод брокер отправляет сообщения в сервис
func (service *BasicService) Deliver(msg ServiceMessage) {
	service.IncomingMessages <- msg
}

func (service *BasicService) Register(serviceId uint64, out chan ServiceMessage) {
	service.IncomingMessages <- ServiceMessage{
		MessageData: BrokerRegisterServiceResponse{
			serviceId,
			out,
		},
	}
}

// Отправляет сообщение брокеру
func (service *BasicService) SendMessage(
	msgType uint64,
	msg interface{},
	sourceClientId uint64,
	targetServiceType uint64,
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
		MessageType:            msgType,
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
	dt := (<-service.IncomingMessages).MessageData.(BrokerRegisterServiceResponse)
	service.Id = dt.Id
	service.OutgoingMessages = dt.Ch
}

func NewBasicService(serviceType uint64) *BasicService {
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
//	regMsg := <-s.IncomingMessages
//	dt := regMsg.MessageData.(struct {
//		Id uint64
//		Ch chan ServiceMessage
//	})
//	s.Id = dt.Id
//	s.OutgoingMessages = dt.Ch
//
//	for msg := range s.IncomingMessages {
//		fmt.Println(msg)
//	}
//}
//
//func example() {
//	svc := &Svc{NewBasicService(1)}
//	svc.Start()
//	var broker MessageBroker
//	broker.RegisterService(svc)
//}

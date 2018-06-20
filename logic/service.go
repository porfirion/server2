package logic

import (
	"errors"
	"time"
)

/**
 * Начитавшись хабра (https://habrahabr.ru/company/mailru/blog/220359/)
 * пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте
 * Более того - это скорее даже мешает - всё валится в одну кучу.
 * Также авторизация остаётся незакрытым вопросом. Пожалуй стоит оформить каждый из этих фрагментов как отдельный сервис.
 */
type Service interface {
	Deliver(msg ServiceMessage)              // через этот метод сообщения закидываются в сервис
	Register(id uint64, ch chan interface{}) // уведомляет сервис о том, что он был зарегистрирован
	Start()
}

// сообщение, которое ходит между сервисами
type ServiceMessage struct {
	SourceServiceType   uint64
	SourceServiceId     uint64
	SourceServiceClient uint64 // по идее у нас не может быть много отправителей

	DestinationServiceType    uint64
	DestinationServiceId      uint64
	DestinationServiceClients []uint64 // зато может быть много получателей

	MessageType uint64
	MessageData interface{}
}

type BasicService struct {
	Id   uint64
	Type uint64

	IncomingMessages chan ServiceMessage
	OutgoingMessages chan interface{}
}

func (service *BasicService) Deliver(msg ServiceMessage) {
	service.IncomingMessages <- msg
}

func (service *BasicService) Register(serviceId uint64, out chan interface{}) {
	data := struct {
		Id uint64
		Ch chan interface{}
	}{
		serviceId,
		out,
	}
	service.IncomingMessages <- ServiceMessage{MessageData: data}
}

/**
 * отправляет сообщение. Первый массив обозначает список целей кому передавать. Второй массив обозначает кому не передавать.
 * @param  {[type]} logic *Logic) SendMessage(msg Message, targets ...[]int [description]
 * @return {[type]} [description]
 */
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

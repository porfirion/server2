package main

import (
	"testing"
	"github.com/porfirion/server2/service"
)

type OutputService struct {
	*service.BasicService
}

func (s *OutputService) Start() {
	//log.Println("waiting for registration")
	s.WaitForRegistration()

	//log.Println("registration received")

	for _ = range s.IncomingMessages {
		output <- true
	}
}

type EchoService struct {
	*service.BasicService
}

func (s *EchoService) Start() {
	s.WaitForRegistration()

	for msg := range s.IncomingMessages {
		msg.DestinationServiceType = 1
		s.BasicService.OutgoingMessages <- msg
	}
}

type TestMessage struct{}

func (m TestMessage) GetType() uint64 { return 1 }

var (
	broker                     service.MessageBroker
	outputService, echoService service.Service
	output                     = make(chan bool)
	finalRes                   bool
)

func setupBrokerAndServices() {
	broker = service.NewBroker(func(message service.ServiceMessage) uint64 {
		return 0
	})
	go broker.Start()

	outputService = &OutputService{service.NewBasicService(1)}
	go outputService.Start()

	echoService = &EchoService{service.NewBasicService(2)}
	go echoService.Start()

	broker.RegisterService(outputService)
	//broker.RegisterService(echoService)
}

// Проверка скорости отправки сообщений через брокер.
// В данном случае мы отправляем в брокер сообщение, предназначенное для OutputService
// А OutputService пробрасывает эти сообщения в output.
// Соответственное мы меряем время между отправкой в брокер до выхода результата из output
func BenchmarkBroker(b *testing.B) {
	b.StopTimer()

	setupBrokerAndServices()

	t := TestMessage{}
	m := service.ServiceMessage{
		DestinationServiceType: 1,
		MessageData:            t,
	}
	var res bool

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		broker.Send(m)
		res = <-output
	}

	finalRes = res

}

func BenchmarkBrokerEcho(b *testing.B) {
	b.StopTimer()

	setupBrokerAndServices()

	t := TestMessage{}
	m := service.ServiceMessage{
		DestinationServiceType: 2,
		MessageData:            t,
	}
	var res bool

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		broker.Send(m)
		res = <-output
	}

	finalRes = res

}

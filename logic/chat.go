package logic

import "fmt"

type Chat struct {
	*BasicService
}

// async start method
func (s *Chat) Start() {
	go s.StartReading()
}

func (s *Chat) StartReading() {
	// первое сообщение, которое должно придти в канал - это сообщение от брокера о регистрации сервиса
	regMsg := <-s.IncomingMessages
	dt := regMsg.MessageData.(struct {
		Id uint64
		Ch chan interface{}
	})
	s.Id = dt.Id
	s.OutgoingMessages = dt.Ch

	for msg := range s.IncomingMessages {
		fmt.Println(msg)
	}
}

func NewChat(bs *BasicService) *Chat {
	return &Chat{
		bs,
	}
}

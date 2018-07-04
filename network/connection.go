package network

import (
	"time"
	"errors"
	"github.com/porfirion/server2/service"
)

type Connection interface {
	Close(message string)
	WriteMessage(data service.TypedMessage)
	GetId() uint64
}

type BasicConnection struct {
	Id              uint64
	OutgoingChannel chan MessageForClient
	IncomingChannel chan MessageFromClient
	ClosingChannel  chan uint64
}

func (connection *BasicConnection) GetId() uint64 {
	return connection.Id
}

func (connection *BasicConnection) WriteMessage(data service.TypedMessage) {
	connection.OutgoingChannel <- MessageForClient{Data: data}
}

// Отправляет сообщение "наверх" (в пул / сервис / брокер)
func (connection *BasicConnection) Notify(data service.TypedMessage) error {
	t := time.NewTimer(time.Millisecond * 100)
	select {
	case connection.IncomingChannel <- MessageFromClient{
		ClientId:    connection.Id,
		Data:        data,
	}:
		t.Stop()
		return nil
	case <-t.C:
		return errors.New("notify message timeout exceeded")
	}
}

func (connection *BasicConnection) NotifyPoolWeAreClosing() {
	connection.ClosingChannel <- connection.Id
}

func NewBasicConnection(id uint64, incoming chan MessageFromClient, closing chan uint64) *BasicConnection {
	return &BasicConnection{
		Id:              id,
		OutgoingChannel: make(chan MessageForClient),
		IncomingChannel: incoming,
		ClosingChannel:  closing,
	}
}
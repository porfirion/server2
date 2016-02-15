package main

type Connection interface {
	StartReading(ch UserMessagesChannel)
	Close(code int, message string)
	GetResponseChannel() MessagesChannel
	SetId(id uint64)
	GetId() uint64
	IsClosed() bool
	SetClosingChannel(chan uint64)
	GetAuth() (*AuthMessage, error)
}

type BasicConnection struct {
	id              uint64
	closed          bool
	responseChannel MessagesChannel
	closingChannel  chan uint64
}

func (connection *BasicConnection) SetId(id uint64) {
	connection.id = id
}
func (connection *BasicConnection) GetId() uint64 {
	return connection.id
}

func (connection *BasicConnection) IsClosed() bool {
	return connection.closed
}

func (connection *BasicConnection) GetResponseChannel() MessagesChannel {
	if connection.responseChannel == nil {
		connection.responseChannel = make(MessagesChannel)
	}

	return connection.responseChannel
}

func (connection *BasicConnection) SetClosingChannel(closingChannel chan uint64) {
	connection.closingChannel = closingChannel
}

type ConnectionsChannel chan Connection

package main

type Connection interface {
	StartReading(ch UserMessagesChannel)
	Close()
	GetResponseChannel() MessagesChannel
	SetId(id int)
	GetId() int
	IsClosed() bool
	SetClosingChannel(chan int)
	GetAuth() (*AuthMessage, error)
}

type BasicConnection struct {
	id              int
	closed          bool
	responseChannel MessagesChannel
	closingChannel  chan int
}

func (connection *BasicConnection) SetId(id int) {
	connection.id = id
}
func (connection *BasicConnection) GetId() int {
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

func (connection *BasicConnection) SetClosingChannel(closingChannel chan int) {
	connection.closingChannel = closingChannel
}

type ConnectionsChannel chan Connection

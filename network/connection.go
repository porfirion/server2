package network

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
	Id              uint64
	Closed          bool
	ResponseChannel MessagesChannel
	ClosingChannel  chan uint64
}

func (connection *BasicConnection) SetId(id uint64) {
	connection.Id = id
}
func (connection *BasicConnection) GetId() uint64 {
	return connection.Id
}

func (connection *BasicConnection) IsClosed() bool {
	return connection.Closed
}

func (connection *BasicConnection) GetResponseChannel() MessagesChannel {
	if connection.ResponseChannel == nil {
		connection.ResponseChannel = make(MessagesChannel)
	}

	return connection.ResponseChannel
}

func (connection *BasicConnection) SetClosingChannel(closingChannel chan uint64) {
	connection.ClosingChannel = closingChannel
}

type ConnectionsChannel chan Connection

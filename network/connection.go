package network

type Connection interface {
	Close(code int, message string)
	WriteMessage(msg interface{})
	GetAuth() (*AuthMessage, error)
	GetId() uint64
}

type ConnectionsChannel chan Connection

type BasicConnection struct {
	Id              uint64
	OutgoingChannel MessagesChannel
	IncomingChannel chan interface{}
	ClosingChannel  chan uint64
}

func (connection *BasicConnection) SetId(id uint64) {
	connection.Id = id
}

func (connection *BasicConnection) GetId() uint64 {
	return connection.Id
}

func (connection *BasicConnection) NotifyPoolWeAreClosing() {
	connection.ClosingChannel <- connection.Id
}

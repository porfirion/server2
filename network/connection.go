package network

type Connection interface {
	Close(message string)
	WriteMessage(msgType uint64, data []byte)
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

func (connection *BasicConnection) WriteMessage(msgType uint64, data []byte) {
	connection.OutgoingChannel <- MessageForClient{MessageType: msgType, Data: data}
}

func (connection *BasicConnection) Notify(msgType uint64, data []byte) {
	connection.IncomingChannel <- MessageFromClient{
		ClientId: connection.Id,
		MessageType: msgType,
		Data: data,
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
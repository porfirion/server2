package main

type Logic struct {
	chanByClientId map[int]chan Message

	incomingMessages chan Message
}

// отправляет сообщение всем
func (logic *Logic) SendMessageAll(msg Message) {
	for _, ch := range logic.chanByConnId {
		ch <- msg
	}
}

// отправляет сообщение только одному адресату
func (logic *Logic) SendMessageId(connId int, msg Message) {
	logic.chanByConnId[connId] <- msg
}

// отпрвляет сообщение нескольких определённым адресатам
func (logic *Logic) SendMessageMultiple(msg Message, targets []int) {

}

// отправлет сообщение всем кроме
func (logic *Logic) SendMessageExcept(msg Message, unwanted []int) {

}

func (logic *Logic) ProcessMessage(Message msg) {
	switch msg.(type) {
	case LoginMessage:
		chanByClientId[msg.ClientId] = msg.ResponseChannel
	case LogoutMessage:
		delete(chanByClientId, msg.ClientId)
	case Ping:

	default:

	}
}

func (logic *Logic) run() {
	select {
	case msg := <-logic.incomingMessages:
		logic.ProcessMessage(msg)
	}
}

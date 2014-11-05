package main

type Logic struct {
	chanByUserId map[uint64]MessageChannel

	incomingMessages MessageChannel
}

// отправляет сообщение всем
func (logic *Logic) SendMessageAll(msg Message) {
	for _, ch := range logic.chanByUserId {
		ch <- msg
	}
}

// отправляет сообщение только одному адресату
func (logic *Logic) SendMessageId(msg Message, connId uint64) {
	logic.chanByUserId[connId] <- msg
}

// отпрвляет сообщение нескольких определённым адресатам
func (logic *Logic) SendMessageMultiple(msg Message, targets []uint64) {

}

// отправлет сообщение всем кроме
func (logic *Logic) SendMessageExcept(msg Message, unwanted []uint64) {

}

func (logic *Logic) ProcessMessage(msg Message) {
	switch msg.(type) {
	case LoginMessage:
		logic.chanByUserId[msg.GetUserId()] = msg.GetResponseChannel()
	case LogoutMessage:
		delete(logic.chanByUserId, msg.GetUserId())
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

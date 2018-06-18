package logic

import (
	"log"
)

/**
 * Начитавшись хабра (https://habrahabr.ru/company/mailru/blog/220359/)
 * пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте
 * Более того - это скорее даже мешает - всё валится в одну кучу.
 * Также авторизация остаётся незакрытым вопросом. Пожалуй стоит оформить каждый из этих фрагментов как отдельный сервис.
 */
type ServiceInterface interface {
	Deliver(msg interface{}) // через этот метод сообщения закидываются в сервис
	SetBroker(broker MessageBroker) // внедрение зависимости брокера
}

type BasicService struct {
	OutgoingMessages chan interface{}
	IncomingMessages chan interface{}
}

func (service *BasicService) Deliver(msg interface{}) {
	panic("implement me")
}

func (service *BasicService) SetBroker(broker MessageBroker) {
	panic("implement me")
}

/**
 * отправляет сообщение. Первый массив обозначает список целей кому передавать. Второй массив обозначает кому не передавать.
 * @param  {[type]} logic *Logic) SendMessage(msg Message, targets ...[]int [description]
 * @return {[type]} [description]
 */
func (service *BasicService) SendMessage(msg interface{}, targets ...[]uint64) {
	serverMessage := ServerMessage{Data: msg}

	// real targets
	if len(targets) > 0 {
		serverMessage.Targets = targets[0]
	}

	// except this users
	if len(targets) > 1 {
		serverMessage.Except = targets[1]
	}

	select {
	case service.OutgoingMessages <- serverMessage:
	default:
		log.Println("busy outgoing messages chan")
	}
}

func (service *BasicService)

func (service *BasicService) SendTextMessage(text string, sender uint64) {
	service.SendMessage(TextMessage{Text: text, Sender: sender})
}

func (service *BasicService) SendTextMessageToUser(text string, sender uint64, userId uint64) {
	service.SendMessage(TextMessage{Text: text, Sender: sender}, []uint64{userId})
}
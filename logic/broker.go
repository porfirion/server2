package logic

/**
 * Брокер, который разруливает в какой сервис отправлять сообщение
 */
type MessageBroker interface {
	Send(msg interface{})
	Broadcast(msg interface{})
	AddService(svc ServiceInterface)
}

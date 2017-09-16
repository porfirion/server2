package main

import "github.com/porfirion/server2/network"

/**
 * Начитавшись хабра (https://habrahabr.ru/company/mailru/blog/220359/)
 * пришёл к выводу, что чат с игровой механикой не стоит держать в одном месте
 * Более того - это скорее даже мешает - всё валится в одну кучу.
 * Также авторизация остаётся незакрытым вопросом. Пожалуй стоит оформить каждый из этих фрагментов как отдельный сервис.
 */
type Service interface {
	GetRequiredMessagesTypes() []int
	GetIncomingChannel() network.UserMessagesChannel
	GetOutgoingChannel() network.ServerMessagesChannel
}

/**
 * Брокер, который разруливает в какой сервис отправлять сообщение
 */
type MessageBroker struct {
	Services []*Service
	IncomingMessagesChannel network.UserMessagesChannel
	OutgoingMessagesChannel network.ServerMessagesChannel
}
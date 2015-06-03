package main

type Message interface {
}

/**
 * При получении сервером ретранслируется всем адресатам
 */
type TextMessage struct {
	Sender int
	Text   string
}

type DataMessage struct {
	Data []byte
}

/**
 * Посылается пользователм на сервер для прохождения авторизации
 */
type AuthMessage struct {
	Uuid string
}

/**
 * Посылается пулом соединений для извещения о входе
 */
type LoginMessage struct {
	User
}

/**
 * Посылается пулом сообщений для извещения о выходе
 */
type LogoutMessage struct {
	User
}

/**
 * Используется для синронизации списка пользователей с клиентом
 */
type UserListMessage struct {
	Users []User
}

type ServerMessage struct {
	Targets []int
	Data    Message
}

type UserMessage struct {
	Source int
	Data   Message
}

type MessagesChannel chan Message

type ServerMessagesChannel chan ServerMessage
type UserMessagesChannel chan UserMessage

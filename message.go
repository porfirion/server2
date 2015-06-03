package main

type Message interface {
}

type ErrorMessage struct {
	Code        int
	Description string
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
 * Посылается клиенту, чтобы сообщить, что он успешно подключился и сказать ему его id
 */
type WellcomeMessage struct {
	Id int
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

/* SPECIAL STRUCTURES */

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

func GetMessageTypeId(msg Message) int {
	var res int = 0
	/*
		1     - 99    Initial messages
		100   - 999   Errors
		1000  - 9999  Information messages
		10000 - 99999 Complex information messages
	*/
	switch msg.(type) {
	case AuthMessage:
		res = 1
	case WellcomeMessage:
		res = 2
	case LoginMessage:
		res = 10
	case LogoutMessage:
		res = 11
	case ErrorMessage:
		res = 100
	case DataMessage:
		res = 1000
	case TextMessage:
		res = 1001
	case UserListMessage:
		res = 10000
	default:
		// Unknown message type
		res = 0
	}

	return res
}

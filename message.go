package main

import (
	"fmt"
	"reflect"
)

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
	Name string
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
	Id   int
	Name string
}

/**
 * Посылается пулом сообщений для извещения о выходе
 */
type LogoutMessage struct {
	Id int
}

/**
 * Используется для синронизации списка пользователей с клиентом
 */
type UserListMessage struct {
	Users []struct {
		Id   int
		Name string
	}
}
type UserLoggedinMessage struct {
	Id   int
	Name string
}

type UserLoggedoutMessage struct {
	Id int
}

type SyncPositionsMessage struct {
	Positions map[string]Position
}

type SyncTimeMessage struct {
	Time int
}

/* SPECIAL STRUCTURES */

type ServerMessage struct {
	Data    Message
	Targets []int // send only to
	Except  []int // do not send to
}

type UserMessage struct {
	Source int
	Data   Message
}

type MessagesChannel chan Message

type ServerMessagesChannel chan ServerMessage
type UserMessagesChannel chan UserMessage

/*func GetMessageTypeId(msg Message) int {
	var res int = 0

	//	1     - 99    Initial messages
	//	100   - 999   Errors
	//	1000  - 9999  Information messages
	//	10000 - 99999 Complex information messages

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
	case UserLoggedinMessage:
		res = 10001
	case UserLoggedoutMessage:
		res = 10002
	case SyncPositionsMessage:
		res = 10003
	case SyncTimeMessage:
		res = 10004
	default:
		// Unknown message type
		fmt.Printf("Unknown message type %#v\n", msg)
		res = 0
	}

	return res
}*/

var dict map[reflect.Type]int = map[reflect.Type]int{
	reflect.TypeOf(AuthMessage{}):          1,
	reflect.TypeOf(WellcomeMessage{}):      2,
	reflect.TypeOf(LoginMessage{}):         10,
	reflect.TypeOf(LogoutMessage{}):        11,
	reflect.TypeOf(ErrorMessage{}):         100,
	reflect.TypeOf(DataMessage{}):          1000,
	reflect.TypeOf(TextMessage{}):          1001,
	reflect.TypeOf(UserListMessage{}):      10000,
	reflect.TypeOf(UserLoggedinMessage{}):  10001,
	reflect.TypeOf(UserLoggedoutMessage{}): 10002,
	reflect.TypeOf(SyncPositionsMessage{}): 10003,
	reflect.TypeOf(SyncTimeMessage{}):      10004,
}

func GetMessageTypeId(value interface{}) int {
	if id, ok := dict[reflect.TypeOf(value)]; ok {
		return id
	} else {
		return -1
	}
}

func GetValueByTypeId(typeId int) interface{} {
	for typeDec, id := range dict {
		fmt.Printf("Type: %#v Id: %#v", typeDec, id)
		if id == typeId {
			return reflect.New(typeDec).Interface()
		}
	}
	fmt.Println("Unknown message type", typeId)
	return nil
}

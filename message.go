package main

import (
	"fmt"
	"reflect"
)

type ErrorMessage struct {
	Code        int    `json: "code"`
	Description string `json: "description"`
}

/**
 * При получении сервером ретранслируется всем адресатам
 */
type TextMessage struct {
	Sender int    `json:"sender"`
	Text   string `json:"text"`
}

type DataMessage struct {
	Data []byte `json:"data"`
}

/**
 * Посылается пользователм на сервер для прохождения авторизации
 */
type AuthMessage struct {
	Name string `json:"name"`
}

/**
 * Посылается клиенту, чтобы сообщить, что он успешно подключился и сказать ему его id
 */
type WellcomeMessage struct {
	Id int `json:"id"`
}

/**
 * Посылается пулом соединений для извещения о входе
 */
type LoginMessage struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

/**
 * Посылается пулом сообщений для извещения о выходе
 */
type LogoutMessage struct {
	Id int `json:"id"`
}

/**
 * Используется для синронизации списка пользователей с клиентом
 */
type UserListMessage struct {
	Users []User `json:"users"`
}
type UserLoggedinMessage struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserLoggedoutMessage struct {
	Id int `json:"id"`
}

type SyncPositionsMessage struct {
	Positions map[string]Position `json:"positions"`
}

type SyncTimeMessage struct {
	Time int64 `json:"time"`
}

/* SPECIAL STRUCTURES */

type ServerMessage struct {
	Data    interface{}
	Targets []int // send only to
	Except  []int // do not send to
}

type UserMessage struct {
	Source int
	Data   interface{}
}

type MessagesChannel chan interface{}

type ServerMessagesChannel chan ServerMessage
type UserMessagesChannel chan UserMessage

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
		fmt.Println("Type is not presented in list")
		return -1
	}
}

func GetValueByTypeId(typeId int) interface{} {
	for typeDec, id := range dict {
		if id == typeId {
			return reflect.New(typeDec).Interface()
		}
	}
	fmt.Println("Can't get value. Unknown message type", typeId)
	return nil
}

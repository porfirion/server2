package service

import (
	"fmt"
	"reflect"
	"github.com/porfirion/server2/world"
)

type ErrorMessage struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// При получении сервером ретранслируется всем адресатам
type TextMessage struct {
	Sender uint64 `json:"sender"`
	Text   string `json:"text"`
}

type DataMessage struct {
	Data []byte `json:"data"`
}

// Посылается пользователм на сервер для прохождения авторизации
type AuthMessage struct {
	Name string `json:"name"`
}

// Посылается клиенту, чтобы сообщить, что он успешно подключился и сказать ему его id
type WelcomeMessage struct {
	Id uint64 `json:"id"`
}

// Посылается пулом соединений для извещения о входе
type LoginMessage struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

// Посылается пулом сообщений для извещения о выходе
type LogoutMessage struct {
	Id uint64 `json:"id"`
}

// Используется для синронизации списка пользователей с клиентом
type UserListMessage struct {
	Users []User `json:"users"`
}

type UserLoggedinMessage struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserLoggedoutMessage struct {
	Id uint64 `json:"id"`
}

// Отправляет на клиент список объектов с координатами
type SyncPositionsMessage struct {
	Positions map[string]world.MapObjectDescription `json:"positions"` // список объектов
	Time      uint64                                `json:"time"`      // время по серверу
}

type ServerStateMessage struct {
	SimulationByStep       bool   `json:"simulation_by_step"`        // идёт ли симуляция по шагам или по времени
	SimulationStepTime     uint64 `json:"simulation_step_time"`      // мс, сколько времени симулируется за раз
	SimulationStepRealTime uint64 `json:"simulation_step_real_time"` // мс, сколько времени реально проходит между симуляциями
	SimulationTime         uint64 `json:"simulation_time"`           // мс игровое время (не привязано ни к какой точке отсчёта, имеет смысл лишь рассматривать изменение этой величины)
	ServerTime             uint64 `json:"server_time"`               // мс, текущее серверное время
}

// Служебное сообщение для ыравнивания времени на сервере и клиенте
type SyncTimeMessage struct {
	Time int64 `json:"time"`
}

// Действие пользователя (двигаться, остановиться, ...)
type ActionMessage struct {
	ActionType string                 `json:"action_type"`
	ActionData map[string]interface{} `json:"action_data"`
}

// Выполнить симуляцию определённого количества шагов (для режима отладки)
type SimulateMessage struct {
	Steps int `json:"steps"`
}

// Сменить режим симуляции на пощаговый или непрерывный
type ChangeSimulationMode struct {
	StepByStep bool `json:"step_by_step"`
}

/* SPECIAL STRUCTURES */

type ServerMessage struct {
	Data    interface{}
	Targets []uint64 // send only to
	Except  []uint64 // do not send to
}

type UserMessage struct {
	Source uint64
	Data   interface{}
}

type MessagesChannel chan interface{}

type ServerMessagesChannel chan ServerMessage
type UserMessagesChannel chan UserMessage

var dict = map[reflect.Type]int{
	reflect.TypeOf(AuthMessage{}):          1,
	reflect.TypeOf(WelcomeMessage{}):       2,
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
	reflect.TypeOf(ServerStateMessage{}):   10005,

	reflect.TypeOf(ActionMessage{}):        1000000,
	reflect.TypeOf(SimulateMessage{}):      1000001,
	reflect.TypeOf(ChangeSimulationMode{}): 1000002,
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

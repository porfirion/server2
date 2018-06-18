package logic

import (
	"fmt"
	"reflect"
	"github.com/porfirion/server2/world"
)

/**
 * Отправляет на клиент список объектов с координатами
 */
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

/**
 * Служебное сообщение для ыравнивания времени на сервере и клиенте
 */
type SyncTimeMessage struct {
	Time int64 `json:"time"`
}

/**
 * Выполнить симуляцию определённого количества шагов (для режима отладки)
 */
type SimulateMessage struct {
	Steps int `json:"steps"`
}

/**
 * Сменить режим симуляции на пощаговый или непрерывный
 */
type ChangeSimulationMode struct {
	StepByStep bool `json:"step_by_step"`
}

/* SPECIAL STRUCTURES */

var dict = map[reflect.Type]int{
	reflect.TypeOf(SyncPositionsMessage{}): 10003,
	reflect.TypeOf(SyncTimeMessage{}):      10004,
	reflect.TypeOf(ServerStateMessage{}):   10005,

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

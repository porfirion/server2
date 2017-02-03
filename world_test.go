package main

import (
	"encoding/json"
	"testing"
	"time"
)

var now time.Time

func CreateAndAdjust(t1, t2 time.Time, t *testing.T) MapObject {
	obj := MapObject{
		CurrentPosition:             Position{X: 0, Y: 0},
		StartTime:                   t1,
		DestinationPosition:         Position{X: 100, Y: 55},
		DestinationTime:             t2,
		Speed:                       10,
	}

	t.Logf("Position: %v, destination: %v\n", obj.CurrentPosition, obj.DestinationPosition)
	obj.AdjustPosition()
	t.Logf("Position: %v, destination: %v\n", obj.CurrentPosition, obj.DestinationPosition)

	return obj
}

func TestAdjustPosition(t *testing.T) {
	now = time.Now()
	var obj MapObject

	obj = CreateAndAdjust(now.Add(-7 * time.Second), now.Add(3 * time.Second), t) // ещё идти 3 секунды
	// с этим тестом иногда случаются проблемы -
	// почему-то разница между взятием текущего времени и вызовом метода Adjust составляет около 1мс,
	// тогда и сравнение наше идёт прахом
	if obj.CurrentPosition.DistanceTo(Position{70, 38.5}) > 0.0001 {
		t.Log(obj.CurrentPosition, obj.CurrentPosition.DistanceTo(Position{70, 38.5}), time.Now().Sub(now), obj.CurrentPosition.DistanceTo(Position{70, 38.5})/float64(time.Now().Sub(now)))
		t.Error("Wrong adjustion for future time")
	} else {
		t.Logf("Computing error: %v\n", obj.CurrentPosition.DistanceTo(Position{70.0, 38.5}))
	}
	obj = CreateAndAdjust(now.Add(-10 * time.Second), now, t) // только что пришли
	if obj.CurrentPosition.DistanceTo(Position{100.0, 55.0}) > 0.0001 {
		t.Log(obj.CurrentPosition)
		t.Error("Wrong adjustion for current time")
	} else {
		t.Logf("Computing error: %v\n", obj.CurrentPosition.DistanceTo(Position{100.0, 55.0}))
	}
	obj = CreateAndAdjust(now.Add(-11 * time.Second), now.Add(-1 * time.Second), t) // пришли ещё секунду назад
	if obj.CurrentPosition.DistanceTo(Position{100.0, 55.0}) > 0.0001 {
		t.Log(obj.CurrentPosition)
		t.Error("Wrong adjustion for previous time")
	} else {
		t.Logf("Computing error: %v\n", obj.CurrentPosition.DistanceTo(Position{100.0, 55.0}))
	}
}

/*func TestRound(t *testing.T) {
	t.Log("0.2 => ", Round(0.2))
	t.Log("0.7 => ", Round(0.7))
	t.Log("0.0 => ", Round(0.0))
	t.Log("-0.2 => ", Round(-0.2))
	t.Log("-0.7 => ", Round(-0.7))
	t.Log("-1.0 => ", Round(-1.0))

	if Round(0.2) != 0.0 {
		t.Error("Error rounding 0.2")
	}
	if Round(0.7) != 1.0 {
		t.Error("Error rounding 0.7")
	}
	if Round(0.0) != 0.0 {
		t.Error("Error rounding 0.0")
	}
	if Round(-0.2) != 0.0 {
		t.Error("Error rounding -0.2")
	}
	if Round(-0.7) != -1.0 {
		t.Error("Error rounding -0.7")
	}
	if Round(-1.0) != -1.0 {
		t.Error("Error rounding -1.0")
	}
}*/

func TestMarshaling(t *testing.T) {
	obj := MapObjectDescription{
		Id:                                  1,
		ObjectType:                          MapObjectTypeObstacle,
		Position:                            Position{10, 20},
		DestinationPosition:                 Position{30, 40},
		DestinationTime:                     time.Now().Unix(),
		Speed:                               1.0,
		StartTime:                           time.Now().Add(time.Minute * -1).Unix(),
		UserId:                              123,
	}

	_, err := json.Marshal(obj)

	if err == nil {
		// нормально сериализовалось
	} else {
		t.Error("Error", err)
	}
}

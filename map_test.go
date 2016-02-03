package main

import (
	"testing"
	"time"
)

func CreateAndAdjust(t1, t2 time.Time, t *testing.T) MapObject {
	obj := MapObject{
		Pos:             Position{X: 0, Y: 0},
		StartTime:       t1,
		Destination:     Position{X: 100, Y: 55},
		DestinationTime: t2,
		Speed:           10,
	}

	t.Log("%+v\n", obj)
	obj.AdjustPosition()
	t.Log("%+v\n", obj)

	return obj
}

func TestAdjustPosition(t *testing.T) {
	now := time.Now()
	var obj MapObject

	obj = CreateAndAdjust(now.Add(-7*time.Second), now.Add(3*time.Second), t) // ещё идти 2 секунды
	if obj.Pos.X != 70 || obj.Pos.Y != 39 {
		t.Error("Wrong adjustion for future time")
	}
	obj = CreateAndAdjust(now.Add(-10*time.Second), now, t) // уже пришли
	if obj.Pos.X != 100 || obj.Pos.Y != 55 {
		t.Error("Wrong adjustion for current time")
	}
	obj = CreateAndAdjust(now.Add(-11*time.Second), now.Add(-1*time.Second), t) // пришли ещё секунду назад
	if obj.Pos.X != 100 || obj.Pos.Y != 55 {
		t.Error("Wrong adjustion for previous time")
	}
}

func TestRound(t *testing.T) {
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
}

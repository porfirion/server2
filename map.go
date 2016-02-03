package main

import (
	"math"
	"time"
)

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MapObject struct {
	Pos             Position
	Destination     Position
	DestinationTime time.Time
	Speed           float64
	StartTime       time.Time
	User            *User
}

func (obj *MapObject) AdjustPosition() {
	if obj.Speed > 0 {

		deltaTime := time.Now().Sub(obj.StartTime).Seconds()            // сколько прошло времени с начала движения
		assumedTime := obj.DestinationTime.Sub(obj.StartTime).Seconds() // сколько времени должно пройти до окончания

		coeff := (float64)(deltaTime / assumedTime)
		if coeff >= 1 {
			// мы уже пришли на место
			obj.Pos = obj.Destination
			obj.StartTime = obj.DestinationTime

			obj.Destination = Position{}
			obj.DestinationTime = time.Time{}
			obj.Speed = 0

		} else {
			// мы ещё не пришли. Рассчитываем текущее положение и записываем его в качестве стартового
			dst := obj.Destination
			src := obj.Pos

			obj.Pos = Position{X: src.X + Round((float64)(dst.X-src.X)*coeff), Y: src.Y + Round((float64)(dst.Y-src.Y)*coeff)}
			obj.StartTime = time.Now()
		}

	}
}

func (obj *MapObject) MoveTo(pos Position, speed float64) {

}

func Round(f float64) int64 {
	return (int64)(math.Floor(f + .5))
}

type Map struct {
	Objects []MapObject
	Width   int
	Height  int
}

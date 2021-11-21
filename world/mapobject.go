package world

import (
	"log"
)

type MapObject struct {
	Id                  uint64   // id объекта
	Speed               Vector2D // speed pixels/second
	CurrentPosition     Point2D  // текущее положение объекта
	DestinationPosition Point2D  // точка, к которой движется объект
	Size                float64  // размер объекта в пикселях (пока оперируем только с кругами)
	Mass                uint16   // Масса объекта
}

type ByLeft []*MapObject

func (a ByLeft) Len() int           { return len(a) }
func (a ByLeft) Less(i, j int) bool { return a[i].CurrentPosition.X < a[j].CurrentPosition.X }
func (a ByLeft) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type MapObjectCollision struct {
	obj1 *MapObject
	obj2 *MapObject
}

func (obj *MapObject) StartMoveTo(dest Point2D) {
	obj.Speed = obj.CurrentPosition.VectorTo(dest).Unit().Mult(ObjectSpeed)
	obj.DestinationPosition = dest
	log.Printf("pos %#v dest %#v speed %#v", obj.CurrentPosition, dest, obj.Speed)
}

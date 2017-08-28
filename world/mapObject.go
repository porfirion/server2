package world

import (
	"log"
)

type MapObjectType int

const (
	MapObjectTypeObstacle MapObjectType = 1   // неподвижный объект, который не может изменять своё положение
	MapObjectTypeMovable                = 10  // объект, который может изменять своё положение, но не управляется пользователем
	MapObjectTypeUser                   = 100 // объект, который может изменять своё положение и принадлежащий какому-либо пользователю
)

type MapObjectDescription struct {
	Id         uint64        `json:"id"`
	ObjectType MapObjectType `json:"objectType"`
	Speed      Vector2D      `json:"speed"`
	Position   Point2D       `json:"position"`
	UserId     uint64        `json:"userId"`
	Size       uint64        `json:"size"`
}

type MapObject struct {
	Id                  uint64        // id объекта
	ObjectType          MapObjectType // тип обхекта. Задаётся константами типа MapObjectType
	UserId              uint64        // ссылка на обхект пользователя, если это пользовательский объект
	Speed               Vector2D      // speed pixels/second
	CurrentPosition     Point2D       // текущее положение обхекта
	DestinationPosition Point2D       // точка, к которой движется объект
	Size                float64       // размер объекта в пикселях (пока оперируем только с кругами)
	Mass                uint16        // Масса объекта
}

type ByLeft []*MapObject

func (a ByLeft) Len() int           { return len(a) }
func (a ByLeft) Less(i, j int) bool { return a[i].CurrentPosition.X < a[j].CurrentPosition.X }
func (a ByLeft) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type MapObjectCollision struct {
	obj1 *MapObject
	obj2 *MapObject
}

func (obj *MapObject) GetDescription() MapObjectDescription {
	description := MapObjectDescription{
		Id:         obj.Id,
		ObjectType: obj.ObjectType,
		Position:   obj.CurrentPosition,
		Speed:      obj.Speed,
		//StartPosition:       obj.StartPosition,
		//StartTime:           obj.StartTime.UnixNano() / int64(time.Millisecond),
		//DestinationPosition: obj.DestinationPosition,
		//DestinationTime:     obj.DestinationTime.UnixNano() / int64(time.Millisecond),
		//Direction:           direction.Modulus(),
	}

	if obj.UserId != 0 {
		description.UserId = obj.UserId
	} else {
		description.UserId = 0
	}

	return description
}

func (obj *MapObject) StartMoveTo(dest Point2D) {
	obj.Speed = obj.CurrentPosition.VectorTo(dest).Modulus(ObjectSpeed)
	obj.DestinationPosition = dest
	log.Printf("pos %#v dest %#v speed %#v", obj.CurrentPosition, dest, obj.Speed)
}

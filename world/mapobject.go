package world

import (
	"log"
)

type MapObjectType int

const (
	MapObjectTypeObstacle MapObjectType = 1   // неподвижный объект, который не может изменять своё положение
	MapObjectTypeMovable  MapObjectType = 10  // объект, который может изменять своё положение, но не управляется пользователем
	MapObjectTypeUser     MapObjectType = 100 // объект, который может изменять своё положение и принадлежащий какому-либо пользователю
)

type MapObjectDTO struct {
	Id          uint64        `json:"id"`
	ObjectType  MapObjectType `json:"objectType"`
	Speed       Vector2D      `json:"speed"`
	Position    Point2D       `json:"position"`
	UserId      uint64        `json:"userId"`
	Size        float64       `json:"size"`
	Destination Point2D       `json:"destination"`
}

func CreateDTOFromMapObject(obj *MapObject) MapObjectDTO {
	dto := MapObjectDTO{
		Id:          obj.Id,
		ObjectType:  obj.ObjectType,
		Position:    obj.CurrentPosition,
		Speed:       obj.Speed,
		Size:        obj.Size,
		Destination: obj.DestinationPosition,
		//StartPosition:       obj.StartPosition,
		//StartTime:           obj.StartTime.UnixNano() / int64(time.Millisecond),
		//DestinationPosition: obj.DestinationPosition,
		//DestinationTime:     obj.DestinationTime.UnixNano() / int64(time.Millisecond),
		//Direction:           direction.Modulus(),
	}

	if obj.UserId != 0 {
		dto.UserId = obj.UserId
	} else {
		dto.UserId = 0
	}

	return dto
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

func (obj *MapObject) StartMoveTo(dest Point2D) {
	obj.Speed = obj.CurrentPosition.VectorTo(dest).Unit().Mult(ObjectSpeed)
	obj.DestinationPosition = dest
	log.Printf("pos %#v dest %#v speed %#v", obj.CurrentPosition, dest, obj.Speed)
}

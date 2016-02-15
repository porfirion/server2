package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (pos Position) Distance(dest Position) float64 {
	return math.Sqrt(math.Pow(dest.X-pos.X, 2) + math.Pow(dest.Y-pos.Y, 2))
}

type MapObjectType int

const (
	MapObjectTypeObstacle MapObjectType = 1
	MapObjectTypeNPC                    = 10
	MapObjectTypeUser                   = 100
)

type MapObject struct {
	Id              uint64
	ObjectType      MapObjectType
	Pos             Position
	Destination     Position
	DestinationTime time.Time
	Speed           float64
	StartTime       time.Time
	User            *User
}

// Обновление положения объекта
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

			obj.Pos = Position{X: src.X + (dst.X-src.X)*coeff, Y: src.Y + (dst.Y-src.Y)*coeff}
			obj.StartTime = time.Now()
		}

	}
}

// Получение текущего положения объекта
func (obj *MapObject) getPosition() Position {
	if obj.Speed > 0 {

		deltaTime := time.Now().Sub(obj.StartTime).Seconds()            // сколько прошло времени с начала движения
		assumedTime := obj.DestinationTime.Sub(obj.StartTime).Seconds() // сколько времени должно пройти до окончания

		coeff := (float64)(deltaTime / assumedTime)
		if coeff >= 1 {
			// мы уже пришли на место
			return obj.Destination

		} else {
			// мы ещё не пришли. Рассчитываем текущее положение и записываем его в качестве стартового
			dst := obj.Destination
			src := obj.Pos

			return Position{X: src.X + (dst.X-src.X)*coeff, Y: src.Y + (dst.Y-src.Y)*coeff}
		}

	} else {
		return obj.Pos
	}
}

// Начало движения в указанную точку
func (obj *MapObject) MoveTo(pos Position, speed float64) {
	obj.AdjustPosition()
	obj.Destination = pos
	obj.Speed = speed

	// duration is counted in nanoseconds
	duration := time.Duration((obj.Pos.Distance(pos) / speed) * 1000000000)

	obj.DestinationTime = time.Now().Add(duration)
}

type WorldMap struct {
	Objects      map[uint64]*MapObject
	UsersObjects map[uint64]*MapObject
	Width        float64
	Height       float64

	NextObjectId uint64
}

func NewWorldMap() *WorldMap {
	var world *WorldMap = new(WorldMap)
	world.Objects = make(map[uint64]*MapObject)
	world.UsersObjects = make(map[uint64]*MapObject)

	for i := 0; i < 10; i++ {
		world.NewObject(Position{X: rand.Float64()*world.Width*2 - world.Width,
			Y: rand.Float64()*world.Height*2 - world.Height,
		},
			MapObjectTypeObstacle)
	}

	return world
}

func (world *WorldMap) NewObject(pos Position, objectType MapObjectType) *MapObject {
	world.NextObjectId++
	return &MapObject{Pos: pos, ObjectType: objectType, Id: world.NextObjectId}
}

func (world *WorldMap) AddObject(obj *MapObject) {
	world.Objects[obj.Id] = obj
}

func (world *WorldMap) AddUser(user *User, pos Position) {
	obj := world.NewObject(pos, MapObjectTypeUser)
	world.AddObject(obj)

	world.UsersObjects[user.Id] = obj
}

func (world *WorldMap) RemoveObject(obj *MapObject) {
	delete(world.Objects, obj.Id)
}

func (world *WorldMap) RemoveUser(userId uint64) {
	obj := world.UsersObjects[userId]
	delete(world.UsersObjects, userId)
	world.RemoveObject(obj)
}

func (world *WorldMap) GetObjectsPositions() map[string]Position {
	res := make(map[string]Position)
	for id, obj := range world.Objects {
		res[strconv.FormatUint(id, 10)] = obj.getPosition()
	}

	log.Printf("Map: users positions %#v\n", res)

	return res
}

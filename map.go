package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	SimulationStepTime time.Duration = 1 * time.Second
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2D) Unit() Vector2D {
	modulo := v.Length()
	return Vector2D{X: v.X / modulo, Y: v.Y / modulo}
}

func (pos Position) Distance(dest Position) float64 {
	return math.Sqrt(math.Pow(dest.X-pos.X, 2) + math.Pow(dest.Y-pos.Y, 2))
}

type MapObjectType struct {
	IsMovable     bool
	BelongsToUser bool
}

var (
	MapObjectTypeObstacle MapObjectType = MapObjectType{IsMovable: false}                      // неподвижный объект, который не может изменять своё положение
	MapObjectTypeMovable                = MapObjectType{IsMovable: true, BelongsToUser: false} // объект, который может изменять своё положение, но не управляется пользователем
	MapObjectTypeUser                   = MapObjectType{}                                      // объект, принадлежащий какому-либо пользователю
)

type MapObjectDescription struct {
	Id                  uint64        `json:"id"`
	ObjectType          MapObjectType `json:"objectType"`
	StartPosition       Position      `json:"startPosition"`
	StartTime           int64         `json:"startTime"`
	DestinationPosition Position      `json:"destinationPosition"`
	DestinationTime     int64         `json:"destinationTime"`
	Speed               float64       `json:"speed"`
	Position            Position      `json:"position"`
	UserId              uint64        `json:"userId"`
	Direction           Vector2D      `json:"direction"`
}

type MapObject struct {
	Id              uint64        // id объекта
	ObjectType      MapObjectType // тип обхекта. Задаётся константами типа MapObjectType
	User            *User         // ссылка на обхект пользователя, если это пользовательский обхект
	Speed           float64       // speed pixels/second
	Acceration      Vector2D      // текущее ускорение, которое действует на объект
	CurrentPosition Position      // текущее положение обхекта

	StartPosition       Position
	StartTime           time.Time
	DestinationPosition Position
	DestinationTime     time.Time
}

// Обновление положения объекта
func (obj *MapObject) AdjustPosition() {
	if obj.Speed > 0 {

		deltaTime := time.Now().Sub(obj.StartTime).Seconds()            // сколько прошло времени с начала движения
		assumedTime := obj.DestinationTime.Sub(obj.StartTime).Seconds() // сколько времени должно пройти до окончания

		coeff := (float64)(deltaTime / assumedTime)
		if coeff >= 1 {
			// мы уже пришли на место
			obj.StartPosition = obj.DestinationPosition
			obj.StartTime = obj.DestinationTime

			obj.DestinationPosition = Position{}
			obj.DestinationTime = time.Time{}
			obj.Speed = 0

		} else {
			// мы ещё не пришли. Рассчитываем текущее положение и записываем его в качестве стартового
			dst := obj.DestinationPosition
			src := obj.StartPosition

			obj.StartPosition = Position{X: src.X + (dst.X-src.X)*coeff, Y: src.Y + (dst.Y-src.Y)*coeff}
			obj.StartTime = time.Now()
		}

	}
}

// Получение текущего положения объекта
func (obj *MapObject) getPosition() (Position, float64) {
	if obj.Speed > 0 {

		deltaTime := time.Now().Sub(obj.StartTime).Seconds()            // сколько прошло времени с начала движения
		assumedTime := obj.DestinationTime.Sub(obj.StartTime).Seconds() // сколько времени должно пройти до окончания

		coeff := (float64)(deltaTime / assumedTime)
		if coeff >= 1 {
			// мы уже пришли на место
			return obj.DestinationPosition, 0.0

		} else {
			// мы ещё не пришли. Рассчитываем текущее положение и записываем его в качестве стартового
			dst := obj.DestinationPosition
			src := obj.StartPosition

			return Position{X: src.X + (dst.X-src.X)*coeff, Y: src.Y + (dst.Y-src.Y)*coeff}, obj.Speed
		}

	} else {
		return obj.StartPosition, 0.0
	}
}

// Начало движения в указанную точку
func (obj *MapObject) MoveTo(pos Position, speed float64) {
	obj.AdjustPosition()
	obj.DestinationPosition = pos
	obj.Speed = speed

	// duration is counted in nanoseconds
	duration := time.Duration((obj.StartPosition.Distance(pos) / speed) * 1000000000)

	obj.DestinationTime = time.Now().Add(duration)

}

func (obj *MapObject) GetDescription() MapObjectDescription {
	direction := Vector2D{X: obj.DestinationPosition.X - obj.StartPosition.X, Y: obj.DestinationPosition.Y - obj.StartPosition.Y}
	position, speed := obj.getPosition() // server calculated position

	description := MapObjectDescription{
		Id:                  obj.Id,
		ObjectType:          obj.ObjectType,
		StartPosition:       obj.StartPosition,
		StartTime:           obj.StartTime.UnixNano() / int64(time.Millisecond),
		DestinationPosition: obj.DestinationPosition,
		DestinationTime:     obj.DestinationTime.UnixNano() / int64(time.Millisecond),
		Speed:               speed,
		Position:            position,
		UserId:              0,
		Direction:           direction.Unit(),
	}

	if obj.User != nil {
		description.UserId = obj.User.Id
	}

	return description
}

type WorldMap struct {
	Objects      map[uint64]*MapObject
	UsersObjects map[uint64]*MapObject
	Width        float64
	Height       float64

	NextObjectId uint64

	// max uint64 - 18446744073709500000
	// если предположить, что шаг симуляции - наносекунда, то
	//1,84467E+19	наносекунд влезает в uint64
	//18446744074	секунд
	//5124095,576	часов
	//213503,9823	дней
	//584,9424174	лет
	SimulationStep uint64    // номер последнего рассчитанного шага симуляции
	StartTime      time.Time // время начала симуляции (отсчитывается от первого вызова simulationStep)
	NextStepTime   time.Time
}

func NewWorldMap() *WorldMap {
	var world *WorldMap = new(WorldMap)
	world.Width = 10000
	world.Height = 10000
	world.Objects = make(map[uint64]*MapObject)
	world.UsersObjects = make(map[uint64]*MapObject)
	world.SimulationStep = 0
	for i := 0; i < 10; i++ {
		obj := world.NewObject(
			Position{
				X: rand.Float64()*300 - 150,
				Y: rand.Float64()*300 - 150,
			},
			MapObjectTypeObstacle)

		world.AddObject(obj)
	}

	log.Println("world created")
	return world
}

func (world *WorldMap) NewObject(pos Position, objectType MapObjectType) *MapObject {
	world.NextObjectId++
	return &MapObject{Id: world.NextObjectId, ObjectType: objectType, StartPosition: pos, StartTime: time.Now(), Speed: 0.0}
}

func (world *WorldMap) AddObject(obj *MapObject) {
	world.Objects[obj.Id] = obj
}

func (world *WorldMap) AddUser(user *User, pos Position) {
	obj := world.NewObject(pos, MapObjectTypeUser)
	world.AddObject(obj)
	obj.User = user

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

func (world *WorldMap) GetObjectsPositions() map[string]MapObjectDescription {
	res := make(map[string]MapObjectDescription)
	for id, obj := range world.Objects {
		res[strconv.FormatUint(id, 10)] = obj.GetDescription()
	}

	// log.Printf("Map: users positions %#v\n", res)

	return res
}

func (world *WorldMap) TimeToNextStep() time.Duration {
	if world.NextStepTime.After(time.Now()) {
		return world.NextStepTime.Sub(time.Now())
	} else {
		return 0
	}
}

func (world *WorldMap) ProcessSimulationStep() bool {
	if world.SimulationStep == 0 {
		world.StartTime = time.Now()
		world.NextStepTime = time.Now()
	} else {
		if time.Now().Before(world.NextStepTime) {
			// время ещё не пришло
			return false
		}
	}
	world.SimulationStep++
	log.Println("Simulation step ", world.SimulationStep)
	world.NextStepTime = world.NextStepTime.Add(SimulationStepTime)

	for id, obj := range world.Objects {
		log.Println("id, obj", id, obj)
	}

	return true
}

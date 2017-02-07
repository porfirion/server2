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
	ObjectSpeed        float64       = 10.0
	MinDistance                      = 1.0
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	NilPosition Position = Position{X: math.MaxFloat64, Y: math.MaxFloat64}
)

/**
 * Расстояние между точками
 */
func (pos Position) DistanceTo(dest Position) float64 {
	return math.Sqrt(math.Pow(dest.X-pos.X, 2) + math.Pow(dest.Y-pos.Y, 2))
}

type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// длина вектора
func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// нормированный вектор
func (v Vector2D) Modulus() Vector2D {
	modulo := v.Length()
	return Vector2D{X: v.X / modulo, Y: v.Y / modulo}
}

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
	Position   Position      `json:"position"`
	UserId     uint64        `json:"userId"`
	//Direction           Vector2D      `json:"direction"`
	//StartPosition       Position      `json:"startPosition"`
	//StartTime           int64         `json:"startTime"`
	//DestinationPosition Position      `json:"destinationPosition"`
	//DestinationTime     int64         `json:"destinationTime"`
}

type MapObject struct {
	Id              uint64        // id объекта
	ObjectType      MapObjectType // тип обхекта. Задаётся константами типа MapObjectType
	User            *User         // ссылка на обхект пользователя, если это пользовательский обхект
	Speed           Vector2D      // speed pixels/second
	CurrentPosition Position      // текущее положение обхекта

	// unused properties
	//Acceration          Vector2D
	//StartPosition       Position
	//StartTime           time.Time
	DestinationPosition Position
	//DestinationTime     time.Time
}

func (obj *MapObject) GetDescription() MapObjectDescription {

	description := MapObjectDescription{
		Id:         obj.Id,
		ObjectType: obj.ObjectType,
		Position:   obj.CurrentPosition,
		//StartPosition:       obj.StartPosition,
		//StartTime:           obj.StartTime.UnixNano() / int64(time.Millisecond),
		//DestinationPosition: obj.DestinationPosition,
		//DestinationTime:     obj.DestinationTime.UnixNano() / int64(time.Millisecond),
		//Speed:               obj.Speed,
		//Direction:           direction.Modulus(),
	}

	if obj.User != nil {
		description.UserId = obj.User.Id
	} else {
		description.UserId = 0;
	}

	return description
}

func (obj *MapObject) MoveTo(pos Position) {
	obj.DestinationPosition = pos;
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
	SimulationTime time.Time // время, в которое, по идее, произошёл текущий шаг симуляции
	StartTime      time.Time // время начала симуляции (отсчитывается от первого вызова simulationStep)
	NextStepTime   time.Time // время, в которое должен произойти следующий шаг симуляции
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
	return &MapObject{Id: world.NextObjectId, ObjectType: objectType, CurrentPosition: pos, DestinationPosition: NilPosition}
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

func (world *WorldMap) GetStepTime(step int) time.Time {
	return world.StartTime.Add(SimulationStepTime * time.Duration(step));
}

func (world *WorldMap) ProcessSimulationStep() bool {
	if world.SimulationStep == 0 {
		// это наш первый шаг симуляции, запоминаем когда стартовали
		world.StartTime = time.Now()
		world.NextStepTime = time.Now()
	} else {
		if time.Now().Before(world.NextStepTime) {
			// время ещё не пришло
			return false
		}
	}
	world.SimulationStep++
	world.SimulationTime = world.StartTime.Add((time.Duration)(world.SimulationStep) * SimulationStepTime);
	log.Println("Simulation step ", world.SimulationStep)
	world.NextStepTime = world.NextStepTime.Add(SimulationStepTime)

	for id, obj := range world.Objects {
		if obj.DestinationPosition != NilPosition {
			log.Println("moving ", id)
			distance := obj.CurrentPosition.DistanceTo(obj.DestinationPosition);
			if distance <= ObjectSpeed {
				obj.CurrentPosition = obj.DestinationPosition;
				obj.DestinationPosition = NilPosition
			} else {
				dx := obj.DestinationPosition.X - obj.CurrentPosition.X;
				dy := obj.DestinationPosition.Y - obj.CurrentPosition.Y;
				obj.CurrentPosition.X += dx / distance * ObjectSpeed;
				obj.CurrentPosition.Y += dy / distance * ObjectSpeed;
			}
		}
		//log.Println("id, obj", id, obj)
	}

	return true
}

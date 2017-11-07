package world

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const (
	SimulationStepTime time.Duration = 100 * time.Millisecond // сколько виртуального времени проходит за один шаг симуляции
	ObjectSpeed        float64       = 50.0                   // скорость объекта по умолчанию
)

type WorldMap struct {
	ObjectsById  map[uint64]*MapObject // список объектов по id
	Objects      []*MapObject          // список объектов, отсоритрованный по левой границе
	UsersObjects map[uint64]*MapObject // список пользовательских объектов
	Width        float64               // ширина карты
	Height       float64               // высота карты

	NextObjectId uint64

	// max uint64 - 18446744073709500000
	// если предположить, что шаг симуляции - наносекунда, то
	//1,84467E+19	наносекунд влезает в uint64
	//18446744074	секунд
	//5124095,576	часов
	//213503,9823	дней
	//584,9424174	лет
	SimulationStep      uint64    // номер последнего рассчитанного шага симуляции
	SimulationTime      time.Time // текущее игровое время
	SimulationStartTime time.Time // время начала симуляции (по идее оно чисто условное и может начинаться с любого времени)
}

func NewWorldMap() *WorldMap {
	var world *WorldMap = new(WorldMap)
	world.Width = 10000
	world.Height = 10000
	world.ObjectsById = make(map[uint64]*MapObject)
	world.UsersObjects = make(map[uint64]*MapObject)
	world.SimulationStep = 0
	for i := 0; i < 10; i++ {
		obj := world.NewObject(
			Point2D{
				X: rand.Float64()*200 - 100,
				Y: rand.Float64()*200 - 100,
			},
			MapObjectTypeObstacle)

		world.AddObject(obj)
	}

	// TODO FORTEST
	world.SimulationStartTime = time.Unix(0, 0)

	log.Println("world created")
	return world
}

func (world *WorldMap) NewObject(pos Point2D, objectType MapObjectType) *MapObject {
	world.NextObjectId++
	return &MapObject{Id: world.NextObjectId, ObjectType: objectType, CurrentPosition: pos, DestinationPosition: NilPosition, Mass: 10, Size: 10}
}

func (world *WorldMap) AddObject(obj *MapObject) {
	world.ObjectsById[obj.Id] = obj
	world.Objects = append(world.Objects, obj)
	sort.Sort(ByLeft(world.Objects))
}

func (world *WorldMap) AddUser(userId uint64, pos Point2D) {
	obj := world.NewObject(pos, MapObjectTypeUser)
	world.AddObject(obj)
	obj.UserId = userId

	world.UsersObjects[userId] = obj
}

func (world *WorldMap) RemoveObject(obj *MapObject) {
	delete(world.ObjectsById, obj.Id)
}

func (world *WorldMap) RemoveUser(userId uint64) {
	obj := world.UsersObjects[userId]
	delete(world.UsersObjects, userId)
	world.RemoveObject(obj)
}

func (world *WorldMap) GetUserObject(userId uint64) *MapObject {
	return world.UsersObjects[userId]
}

func (world *WorldMap) GetObjectsPositions() map[string]MapObjectDescription {
	res := make(map[string]MapObjectDescription)
	for id, obj := range world.ObjectsById {
		res[strconv.FormatUint(id, 10)] = obj.GetDescription()
	}
	//log.Printf("Map: users positions %#v\n", res)

	return res
}

/**
 * Возвращает игровое время для указанного шага симуляции
 */
func (world *WorldMap) GetStepTime(step int) time.Time {
	return world.SimulationStartTime.Add(SimulationStepTime * time.Duration(step))
}

// Выполнение симуляции.
// Первый возвращаемый результат - была ли выполнены симуляция
// Второй возвращаемый результат - произошли ли какие-то существенные изменения
func (world *WorldMap) ProcessSimulationStep() (somethingChanged bool) {
	somethingChanged = false

	world.SimulationStep++
	world.SimulationTime = world.SimulationStartTime.Add(time.Duration(world.SimulationStep) * SimulationStepTime)

	var passedTime float64 = float64(SimulationStepTime) / float64(time.Second)

	for _, obj := range world.ObjectsById {
		if obj.DestinationPosition != NilPosition {
			//log.Println("moving ", id)
			distance := obj.CurrentPosition.DistanceTo(obj.DestinationPosition)
			if distance <= ObjectSpeed*passedTime {
				obj.CurrentPosition = obj.DestinationPosition
				obj.DestinationPosition = NilPosition
				obj.Speed = Vector2D{}
				somethingChanged = true
			} else {
				//dx := obj.DestinationPosition.X - obj.CurrentPosition.X
				//dy := obj.DestinationPosition.Y - obj.CurrentPosition.Y
				//obj.CurrentPosition.X += dx / distance * ObjectSpeed
				//obj.CurrentPosition.Y += dy / distance * ObjectSpeed
				//log.Printf("Position: %#v Speed %#v passedTime %f", obj.CurrentPosition, obj.Speed, passedTime)

				obj.CurrentPosition.X += obj.Speed.X * passedTime
				obj.CurrentPosition.Y += obj.Speed.Y * passedTime

				//log.Printf("Position: %#v Speed %#v", obj.CurrentPosition, obj.Speed)
			}
		}
		//log.Println("id, obj", id, obj)
	}

	if collisions := world.detectCollisions(); len(collisions) > 0 {
		world.resolveCollisions(collisions)
	}

	return
}

func (world *WorldMap) detectCollisions() []MapObjectCollision {
	collisions := make([]MapObjectCollision, 0)
	for i := 0; i < len(world.Objects); i++ {
		obj1 := world.Objects[i]
		if (i < len(world.Objects) - 1) {
			// это не послдений объект в списке
			for j := i + 1; j < len(world.Objects); j++ {
				obj2 := world.Objects[j]

				if (obj1.Id == obj2.Id) {
					log.Println("WARNING! the same objects!");
					continue
				}

				id1, id2 := obj1.Id, obj2.Id

				//if obj1.ObjectType == MapObjectTypeUser || obj2.ObjectType == MapObjectTypeUser {
				//	log.Printf("%d -- %d dist %f > %f + %f", id1, id2, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
				//}
				if obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition) < float64(obj1.Size+obj2.Size) {
					log.Printf("collide %d VS %d \n", id1, id2)
					collisions = append(collisions, MapObjectCollision{obj1, obj2})
				}
			}
		}
	}

	return collisions
}

func (world *WorldMap) resolveCollisions(collisions []MapObjectCollision) {
	log.Printf("Resolving %d collisions", len(collisions))
	for _, collision := range collisions {
		resolver := GetResolver(collision.obj1, collision.obj2)
		resolver.resolve(collision.obj1, collision.obj2)
	}
	// TODO introduce collisions resolver
}

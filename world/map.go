package world

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const (
	ObjectSpeed = 50.0 // скорость объекта по умолчанию

	// Чтобы из-за точности float не происходило лишних расчётов столкновений,
	// когда объекты сталкиваются на 1.1368683772161603e-13, вводим эту константу
	CollisionThreshold = 0.001 // минимальное расстояние, которое считается столкновением
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
	SimulationTime      time.Duration // сколько прошло игрового времени с начала симуляции
	SimulationStartTime time.Time     // время начала симуляции (по идее оно чисто условное и может начинаться с любого времени)
}

func NewWorldMap() *WorldMap {
	var world = new(WorldMap)
	world.Width = 10000
	world.Height = 10000
	world.ObjectsById = make(map[uint64]*MapObject)
	world.UsersObjects = make(map[uint64]*MapObject)
	for i := 0; i < 10; i++ {
		obj := world.NewObject(
			Point2D{
				X: rand.Float64()*200 - 100,
				Y: rand.Float64()*200 - 100,
			},
			MapObjectTypeObstacle)
		world.AddObject(obj)
	}

	world.SimulationStartTime = time.Now()
	world.SimulationTime = 0 //time.Duration(10000000000 * rand.Float32())

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

// возвращает текущее время в миллисекундах
// (начиная с неизвестно чего). Имеет смысл ориентироваться только на разницу во времени,
// а не на его абсолютное значение
func (world *WorldMap) GetCurrentTimeMillis() uint64 {
	return uint64(world.SimulationTime / time.Millisecond)
}

//Выполнение симуляции.
//return произошли ли какие-то существенные изменения
func (world *WorldMap) ProcessSimulationStep(passedTimeDur time.Duration) (somethingChanged bool) {
	somethingChanged = false

	world.SimulationTime += passedTimeDur

	// время, которое прошло за шаг, секунды
	var dt = float64(passedTimeDur) / float64(time.Second)

	for _, obj := range world.ObjectsById {
		if obj.DestinationPosition != NilPosition {
			//log.Println("moving ", id)
			distance := obj.CurrentPosition.Distance2To(obj.DestinationPosition)
			if distance <= (ObjectSpeed*dt)*(ObjectSpeed*dt) {
				obj.CurrentPosition = obj.DestinationPosition
				obj.DestinationPosition = NilPosition
				obj.Speed = Vector2D{}
				somethingChanged = true
			} else {
				//dx := obj.DestinationPosition.X - obj.CurrentPosition.X
				//dy := obj.DestinationPosition.Y - obj.CurrentPosition.Y
				//obj.CurrentPosition.X += dx / distance * ObjectSpeed
				//obj.CurrentPosition.Y += dy / distance * ObjectSpeed
				//log.Printf("Position: %#v Speed %#v dt %f", obj.CurrentPosition, obj.Speed, dt)

				obj.Speed = obj.CurrentPosition.VectorTo(obj.DestinationPosition).Modulus(ObjectSpeed)

				obj.CurrentPosition.X += obj.Speed.X * dt
				obj.CurrentPosition.Y += obj.Speed.Y * dt

				//log.Printf("Position: %#v Speed %#v", obj.CurrentPosition, obj.Speed)
			}
		}
		//log.Println("id, obj", id, obj)
	}

	if collisions := world.detectPossibleCollisions(); len(collisions) > 0 {
		res := world.resolveCollisions(collisions)
		somethingChanged = somethingChanged || res
	}

	return
}

/**
 * Ищет возможные коллизии.
 * TODO здесь надо бы переделать на bounding box
 * Wide phase
 */
func (world *WorldMap) detectPossibleCollisions() []MapObjectCollision {
	collisions := make([]MapObjectCollision, 0)
	for i := 0; i < len(world.Objects); i++ {
		obj1 := world.Objects[i]
		if i < len(world.Objects)-1 {
			// это не послдений объект в списке
			for j := i + 1; j < len(world.Objects); j++ {
				obj2 := world.Objects[j]

				if obj1.Id == obj2.Id {
					log.Println("WARNING! the same objects!")
					continue
				}

				id1, id2 := obj1.Id, obj2.Id

				//if obj1.ObjectType == MapObjectTypeUser || obj2.ObjectType == MapObjectTypeUser {
				//	log.Printf("%d -- %d dist %f > %f + %f", id1, id2, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
				//}

				// маленький хак - используем не расстояние, а его квадрат, чтобы не извлекать корень
				distance := obj1.CurrentPosition.Distance2To(obj2.CurrentPosition)
				minimum := float64(obj1.Size+obj2.Size) * float64(obj1.Size+obj2.Size)
				if distance < minimum && minimum-distance > CollisionThreshold {
					log.Printf("collide %d VS %d (%f < %f)\n", id1, id2, distance, minimum)
					collisions = append(collisions, MapObjectCollision{obj1, obj2})
				}
			}
		}
	}

	return collisions
}

func (world *WorldMap) resolveCollisions(collisions []MapObjectCollision) bool {
	changed := false
	log.Printf("Resolving %d collisions", len(collisions))
	for _, collision := range collisions {
		res := GetResolver(collision.obj1, collision.obj2).resolve(collision.obj1, collision.obj2)
		changed = changed || res
	}
	return changed
}

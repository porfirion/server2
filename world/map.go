package world

import (
	"log"
	"math/rand"
	"sort"
	"time"
)

const (
	ObjectSpeed = 50.0 // скорость объекта по умолчанию

	// Чтобы из-за точности float не происходило лишних расчётов столкновений,
	// когда объекты сталкиваются на 1.1368683772161603e-13, вводим эту константу
	CollisionThreshold = 0.001 // минимальное расстояние, которое считается столкновением
)

type MapLayer int64 // слои на карте, чтобы считать коллизии

const (
	LayerDefault     = 1
	LayerUnderground = 1 << 2
	LayerAir         = 1 << 3
	// и т.д. до 63
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
	SimulationTime time.Duration // сколько прошло игрового времени с начала симуляции
}

func NewWorldMap(width, height float64) *WorldMap {
	var world = new(WorldMap)
	world.Width = width
	world.Height = height
	world.ObjectsById = make(map[uint64]*MapObject)
	world.UsersObjects = make(map[uint64]*MapObject)

	world.SimulationTime = 0 //time.Duration(10000000000 * rand.Float32())

	log.Println("world created")
	return world
}

func (m *WorldMap) TestFill() {
	for i := 0; i < 10; i++ {
		m.NewObject(
			Point2D{
				X: rand.Float64()*200 - 100,
				Y: rand.Float64()*200 - 100,
			},
			MapObjectTypeObstacle)
	}
}

func (m *WorldMap) NewObject(pos Point2D, objectType MapObjectType) *MapObject {
	m.NextObjectId++
	obj := &MapObject{Id: m.NextObjectId, ObjectType: objectType, CurrentPosition: pos, DestinationPosition: NilPosition, Mass: 10, Size: 10}
	m.ObjectsById[obj.Id] = obj
	m.Objects = append(m.Objects, obj)
	sort.Sort(ByLeft(m.Objects))

	return obj
}

// should be removed from map
func (m *WorldMap) AddUser(userId uint64, pos Point2D) *MapObject {
	obj := m.NewObject(pos, MapObjectTypeUser)
	obj.UserId = userId

	m.UsersObjects[userId] = obj

	return obj
}

func (m *WorldMap) RemoveObject(obj *MapObject) {
	delete(m.ObjectsById, obj.Id)
}

func (m *WorldMap) RemoveUser(userId uint64) {
	obj := m.UsersObjects[userId]
	delete(m.UsersObjects, userId)
	m.RemoveObject(obj)
}

func (m *WorldMap) GetUserObject(userId uint64) *MapObject {
	return m.UsersObjects[userId]
}

// получает положение всех видимых объектов
// TODO добавить дополнительные параметры - может ли игрок видеть невидимые объекты; какие слои видит
func (world *WorldMap) GetObjectsPositions(center Point2D, radius float64) map[uint64]MapObjectDTO {
	res := make(map[uint64]MapObjectDTO)
	for id, obj := range world.ObjectsById {
		res[id] = CreateDTOFromMapObject(obj)
	}
	//log.Printf("Map: users positions %#v\n", res)

	return res
}

func (*WorldMap) GetObjectsInRadius(center Point2D, radius float64) []*MapObject {
	panic("implement me!")
}

func (*WorldMap) GetObjectsInRect(leftTop Point2D, rightBottom Point2D) []*MapObject {
	panic("implement me!")
}

// получение размеров карты
func (m WorldMap) GetSize() Point2D {
	return Point2D{m.Width, m.Height}
}

// возвращает текущее время в миллисекундах
// (начиная с неизвестно чего). Имеет смысл ориентироваться только на разницу во времени,
// а не на его абсолютное значение
func (m *WorldMap) GetCurrentTimeMillis() uint64 {
	return uint64(m.SimulationTime / time.Millisecond)
}

//Выполнение симуляции.
//return произошли ли какие-то существенные изменения
func (m *WorldMap) ProcessSimulationStep(passedTimeDur time.Duration) {

	m.SimulationTime += passedTimeDur

	// время, которое прошло за шаг, секунды
	var dt = float64(passedTimeDur) / float64(time.Second)

	for _, obj := range m.ObjectsById {
		if obj.DestinationPosition != NilPosition {
			//log.Println("moving ", id)
			distance := obj.CurrentPosition.Distance2To(obj.DestinationPosition)
			if distance <= (ObjectSpeed*dt)*(ObjectSpeed*dt) {
				obj.CurrentPosition = obj.DestinationPosition
				obj.DestinationPosition = NilPosition
				obj.Speed = Vector2D{}
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

	if collisions := m.detectPossibleCollisions(); len(collisions) > 0 {
		m.resolveCollisions(collisions)
	}

	return
}

// Ищет возможные коллизии.
// TODO здесь надо бы переделать на bounding box
// Wide phase
func (m *WorldMap) detectPossibleCollisions() []MapObjectCollision {
	collisions := make([]MapObjectCollision, 0)
	for i := 0; i < len(m.Objects); i++ {
		obj1 := m.Objects[i]
		if i < len(m.Objects)-1 {
			// это не послдений объект в списке
			for j := i + 1; j < len(m.Objects); j++ {
				obj2 := m.Objects[j]

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

func (m *WorldMap) resolveCollisions(collisions []MapObjectCollision) bool {
	changed := false
	log.Printf("Resolving %d collisions", len(collisions))
	for _, collision := range collisions {
		res := GetResolver(collision.obj1, collision.obj2).resolve(collision.obj1, collision.obj2)
		changed = changed || res
	}
	return changed
}

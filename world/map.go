package world

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	SimulationStepTime time.Duration = 100 * time.Millisecond
	ObjectSpeed        float64       = 50.0
)

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
			Point2D{
				X: rand.Float64()*300 - 150,
				Y: rand.Float64()*300 - 150,
			},
			MapObjectTypeObstacle)

		world.AddObject(obj)
	}

	log.Println("world created")
	return world
}

func (world *WorldMap) NewObject(pos Point2D, objectType MapObjectType) *MapObject {
	world.NextObjectId++
	return &MapObject{Id: world.NextObjectId, ObjectType: objectType, CurrentPosition: pos, DestinationPosition: NilPosition, Mass: 10, Size: 10}
}

func (world *WorldMap) AddObject(obj *MapObject) {
	world.Objects[obj.Id] = obj
}

func (world *WorldMap) AddUser(userId uint64, pos Point2D) {
	obj := world.NewObject(pos, MapObjectTypeUser)
	world.AddObject(obj)
	obj.UserId = userId

	world.UsersObjects[userId] = obj
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
	//log.Printf("Map: users positions %#v\n", res)

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
	return world.StartTime.Add(SimulationStepTime * time.Duration(step))
}

// Выполнение симуляции.
// Первый возвращаемый результат - была ли выполнены симуляция
// Второй возвращаемый результат - произошли ли какие-то существенные изменения
func (world *WorldMap) ProcessSimulationStep() (simulationPassed bool, somethingChanged bool) {
	if world.SimulationStep == 0 {
		// это наш первый шаг симуляции, запоминаем когда стартовали
		world.StartTime = time.Now()
		world.NextStepTime = time.Now()
	} else {
		if time.Now().Before(world.NextStepTime) {
			// время ещё не пришло
			return false, false
		}
		// ага, время уже настало. Симулируем
	}

	simulationPassed = true
	somethingChanged = false

	world.SimulationStep++
	world.SimulationTime = world.StartTime.Add((time.Duration)(world.SimulationStep) * SimulationStepTime)
	//log.Println("Simulation step ", world.SimulationStep)
	world.NextStepTime = world.NextStepTime.Add(SimulationStepTime)
	var passedTime float64 = float64(SimulationStepTime) / float64(time.Second)

	for id, obj := range world.Objects {
		if obj.DestinationPosition != NilPosition {
			log.Println("moving ", id)
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

	collisions := world.detectCollisions()

	world.resolveCollisions(collisions)

	return
}

func (world *WorldMap) detectCollisions() []MapObjectCollision {
	collisions := make([]MapObjectCollision, 0)
	for id1, obj1 := range world.Objects {
		for id2, obj2 := range world.Objects {
			if id1 != id2 {
				if obj1.ObjectType == MapObjectTypeUser || obj2.ObjectType == MapObjectTypeUser {
					log.Printf("%d -- %d dist %f > %f + %f", id1, id2, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
				}
				if obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition) < float64(obj1.Size+obj2.Size) {
					log.Printf("collide %d VS %d \n", id1, id2);
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

package service

import (
	"github.com/porfirion/server2/world"
	"log"
	"math/rand"
	"time"
)

const (
	MAX_SYNC_TIMEOUT = 100 * time.Millisecond
)

/**
 * соотношением SimulationStepTime / SimulationStepRealTime можно регулировать скорость игрового сервера
 */
type LogicParams struct {
	SimulateByStep           bool          // если выставить этот флаг, то симуляция запускается не по таймеру, а по приходу события Simulate
	SimulationStepTime       time.Duration // сколько виртуального времени проходит за один шаг симуляции
	SimulationStepRealTime   time.Duration // сколько реального времени проходит за один шаг симуляции
	SendObjectsTimeout       time.Duration // частота отправки состояний объектов клиентам
	MaxSimulationStepsAtOnce int           // максимальнео количество симуляций подряд.
}

type SendStub struct {
}

func (st *SendStub) SendMessage(args ...interface{})           {}
func (st *SendStub) SendTextMessage(args ...interface{})       {}
func (st *SendStub) SendTextMessageToUser(args ...interface{}) {}

type GameLogic struct {
	*BasicService
	*SendStub
	Params           LogicParams
	IncomingMessages UserMessagesChannel
	OutgoingMessages ServerMessagesChannel
	Users            map[uint64]*User

	mWorldMap *world.WorldMap

	forceSimulationChannel      chan int  // отправка сообщения в этот канал инициирует новый шаг симуляции
	changeSimulationModeChannel chan bool // отправка сообщения в этот канал инициирует изменение режима симуляции
	prevSyncTime                time.Time

	StartTime    time.Time // время начала симуляции (отсчитывается от первого вызова simulationStep)
	NextStepTime time.Time // время, в которое должен произойти следующий шаг симуляции
	PrevStepTime time.Time // время, в которое произошёл предыдущий шаг симуляции
}

func (logic GameLogic) GetIncomingMessagesChannel() UserMessagesChannel {
	return logic.IncomingMessages
}

func (logic GameLogic) GetOutgoingMessagesChannel() ServerMessagesChannel {
	return logic.OutgoingMessages
}

func (logic *GameLogic) GetUserList(exceptId uint64) []User {
	var userlist []User
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, User{Id: user.Id, Name: user.Name})
		}
	}

	log.Printf("Userlist: %#v\n", userlist)

	return userlist
}

func (logic *GameLogic) AddUser(id uint64, name string) *User {
	user := &User{Id: id, Name: name}
	logic.Users[id] = user

	pos := world.Point2D{X: rand.Float64()*600 - 300, Y: rand.Float64()*600 - 300}
	logic.mWorldMap.AddUser(user.Id, pos)

	return user
}

func (logic *GameLogic) RemoveUser(id uint64) {
	delete(logic.Users, id)
	logic.mWorldMap.RemoveUser(id)
}

func (logic *GameLogic) sendSyncPositionMessage() {
	logic.SendStub.SendMessage(
		SyncPositionsMessage{
			logic.mWorldMap.GetObjectsPositions(),
			logic.mWorldMap.GetCurrentTimeMillis(),
		},
	)
	logic.prevSyncTime = time.Now()
}

// Возвращает true, если нужно синхронизировать положение объектов заново
func (logic *GameLogic) ProcessActionMessage(userId uint64, msg *ActionMessage) (needSync bool) {
	log.Println("UNIMPLEMENTED action message processing")
	needSync = false
	switch msg.ActionType {
	case "move":
		user := logic.Users[userId]
		userObject := logic.mWorldMap.GetUserObject(userId)

		x, okX := msg.ActionData["x"].(float64)
		y, okY := msg.ActionData["y"].(float64)

		if okX && okY {
			userObject.StartMoveTo(world.Point2D{X: x, Y: y})
			log.Printf("user #%d try to move it's object #%d to (%f:%f)\n", userId, userObject.Id, x, y)
		} else {
			log.Println("can't get x and y")
		}

		log.Println(user.Name+" is moving to", msg.ActionData)
		needSync = true
	default:
		log.Println("Unknown action type: ", msg.ActionType)
	}

	return
}

func (logic *GameLogic) ProcessMessage(message UserMessage) (needSync bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in %#v\n", r)
		}
	}()

	needSync = false

	switch msg := message.Data.(type) {
	case *TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case *LoginMessage:
		log.Println("GameLogic: Login message received")

		user := logic.AddUser(msg.Id, msg.Name)
		logic.SendStub.SendTextMessageToUser("GameLogic: Welcome, "+user.Name, 0, user.Id)
		logic.SendStub.SendMessage(logic.getServerStateMessage(), UsersList{user.Id})
		logic.SendStub.SendMessage(UserLoggedinMessage{Id: user.Id, Name: user.Name}, UsersList{}, UsersList{user.Id})
		logic.SendStub.SendMessage(UserListMessage{logic.GetUserList(user.Id)}, UsersList{user.Id})
		log.Println("sent. sync next")
		logic.SendStub.SendMessage(SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), logic.mWorldMap.GetCurrentTimeMillis()})
	case *LogoutMessage:
		log.Println("GameLogic: Logout message", msg.Id)
		logic.RemoveUser(msg.Id)
		logic.SendStub.SendMessage(UserLoggedoutMessage{Id: msg.Id})
	case *ActionMessage:
		needSync = logic.ProcessActionMessage(message.Source, msg)
	case *SimulateMessage:
		if logic.Params.SimulateByStep {
			select {
			case logic.forceSimulationChannel <- msg.Steps:
			default:
				log.Println("Already busy force chan")
			}
		} else {
			log.Println("We are not in step by step mode")
		}
	case *ChangeSimulationMode:
		newValue := msg.StepByStep
		if logic.Params.SimulateByStep != newValue {
			select {
			case logic.changeSimulationModeChannel <- newValue:
			default:
				log.Println("Already busy mode chan")
			}
		} else {
			log.Printf("Simulation already in mode %v\n", newValue)
		}
	default:
		log.Printf("GameLogic: Unknown message type %#v from %d\n", message.Data, message.Source)
	}

	return
}

func (logic *GameLogic) getServerStateMessage() ServerStateMessage {
	return ServerStateMessage{
		SimulationByStep:       logic.Params.SimulateByStep,
		SimulationStepTime:     (uint64)(logic.Params.SimulationStepTime / time.Millisecond),
		SimulationStepRealTime: (uint64)(logic.Params.SimulationStepRealTime / time.Millisecond),
		ServerTime:             (uint64)(time.Now().UnixNano() / int64(time.Millisecond)),
		SimulationTime:         logic.mWorldMap.GetCurrentTimeMillis(),
	}
}

func (logic *GameLogic) executeSimulation(dt time.Duration) (changed bool) {
	changed = logic.mWorldMap.ProcessSimulationStep(logic.Params.SimulationStepTime)
	logic.PrevStepTime = time.Now()
	logic.NextStepTime = logic.NextStepTime.Add(logic.Params.SimulationStepRealTime)
	return
}

func (logic *GameLogic) Start() {
	rand.Seed(int64(time.Now().Nanosecond()))

	logic.Users = make(map[uint64]*User)
	logic.changeSimulationModeChannel = make(chan bool, 2)
	logic.forceSimulationChannel = make(chan int, 1)

	logic.mWorldMap = world.NewWorldMap()

	logic.NextStepTime = time.Now()
	logic.prevSyncTime = time.Unix(0, 0)
	// таймер, который инициализирует симуляцию
	var simulationTimer = time.NewTimer(0)

	if logic.Params.SimulateByStep {
		if !simulationTimer.Stop() {
			<-simulationTimer.C
		}
	}

	log.Println("GameLogic: started")
	for {
		// сначала пытаемся вычитаь все входящие соединения.
		// потом надо будет учесть вариант, что нас могут заспамить и симуляция вообще не произойдёт.
		select {
		case msg := <-logic.IncomingMessages:
			//log.Println("GameLogic: message received")
			if needSync := logic.ProcessMessage(msg); needSync {
				logic.sendSyncPositionMessage()
			}
			//log.Println("GameLogic: message processed")
			continue
		default:
		}

		select {
		case mode := <-logic.changeSimulationModeChannel:
			// обработка смены режима симуляции
			if logic.Params.SimulateByStep != mode {
				logic.Params.SimulateByStep = mode

				if logic.Params.SimulateByStep {
					log.Println("Simulation mode changed to STEP_BY_STEP")
					// симуляция по шагам
					if !simulationTimer.Stop() {
						<-simulationTimer.C
					}
					// обнулили таймер и просто ждём что нам скажут делать дальше:
					// либо симулировать очередной шаг по команде,
					// либо вернуть непрерывную симуляцию
				} else {
					log.Println("Simulation mode changed to CONTINIOUS")
					// непрерывная симуляция
					logic.NextStepTime = time.Now()
					log.Println("stopping timer")

					// если попытаться по-новой остановить уже остановленный таймер - он вернёт false.
					// но при этом в C будет пусто и мы просто навечно залокируемся тут
					if !simulationTimer.Stop() {
						select {
						case <-simulationTimer.C:
						default:
							log.Println("timer channel is already empty")
						}
					}

					log.Println("resetting timer")
					simulationTimer.Reset(0)
				}

				log.Println("sending server state message")
				// а теперь уведомляем всех об изменившемся режиме
				logic.SendStub.SendMessage(logic.getServerStateMessage())
				logic.sendSyncPositionMessage()
				log.Println("server state message sent")
			} else {
				log.Printf("Simulation mode already was %v\n", mode)
			}
		case _ = <-simulationTimer.C:
			// по идее уже пора выполнять очередной шаг симуляции
			//log.Println("Timer fired!")
			//log.Printf("Now %v next %v\n", time.Now(), service.NextStepTime)

			stepsCount := 0
			globallyChanged := false

			if (!logic.NextStepTime.Equal(time.Now()) && !logic.NextStepTime.Before(time.Now())) {
				log.Println("WARNING! simulation timer fired before next step!")
			}

			startTime := time.Now()

			// если уже пора симулировать, то симулируем, н оне больше 10 шагов
			for (logic.NextStepTime.Equal(time.Now()) || logic.NextStepTime.Before(time.Now())) && stepsCount < logic.Params.MaxSimulationStepsAtOnce {
				// если что-то изменилось - нужно разослать всем уведомления
				changed := logic.executeSimulation(logic.Params.SimulationStepTime)
				globallyChanged = globallyChanged || changed
				stepsCount++

				if stepsCount > 1 {
					log.Printf("step %d\n", stepsCount)
				}
			}

			passedTime := time.Now().Sub(startTime)
			log.Printf("Simulated %d steps (%d mcs): world time %d ms\n", stepsCount, passedTime/time.Microsecond, logic.mWorldMap.GetCurrentTimeMillis())

			if globallyChanged || time.Now().Sub(logic.prevSyncTime) > MAX_SYNC_TIMEOUT {
				logic.sendSyncPositionMessage()
			}

			timeToNextStep := logic.NextStepTime.Sub(time.Now())
			//log.Printf("Delaying timer for %v nanoseconds\n", timeToNextStep.Nanoseconds())
			simulationTimer.Reset(timeToNextStep)
		case _ = <-logic.forceSimulationChannel:
			// нас попросили выполнить очередной шаг симуляции
			if logic.Params.SimulateByStep {
				log.Println("Simulating!")
				startTime := time.Now()
				logic.executeSimulation(logic.Params.SimulationStepTime)
				passedTime := time.Now().Sub(startTime)
				log.Printf("Simulated 1 step (%d mcs): world time %d ms", passedTime/time.Microsecond, logic.mWorldMap.GetCurrentTimeMillis())
				logic.sendSyncPositionMessage()
			} else {
				log.Println("Not in step by step mode")
			}
		case msg := <-logic.IncomingMessages:
			//log.Println("GameLogic: message received")
			if needSync := logic.ProcessMessage(msg); needSync {
				logic.sendSyncPositionMessage()
			}
			//log.Println("GameLogic: message processed")
		}
	}

	log.Println("GameLogic: finished")
}

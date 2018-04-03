package main

import (
	"github.com/porfirion/server2/network"
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

type Logic struct {
	params           LogicParams
	IncomingMessages network.UserMessagesChannel
	OutgoingMessages network.ServerMessagesChannel
	Users            map[uint64]*network.User

	mWorldMap *world.WorldMap

	forceSimulationChannel      chan int  // отправка сообщения в этот канал инициирует новый шаг симуляции
	changeSimulationModeChannel chan bool // отправка сообщения в этот канал инициирует изменение режима симуляции
	prevSyncTime                time.Time

	StartTime    time.Time // время начала симуляции (отсчитывается от первого вызова simulationStep)
	NextStepTime time.Time // время, в которое должен произойти следующий шаг симуляции
	PrevStepTime time.Time // время, в которое произошёл предыдущий шаг симуляции
}

func (logic *Logic) SetParams(params LogicParams) {
	logic.params = params
}

func (logic Logic) GetIncomingMessagesChannel() network.UserMessagesChannel {
	return logic.IncomingMessages
}
func (logic *Logic) SetIncomingMessagesChannel(channel network.UserMessagesChannel) {
	logic.IncomingMessages = channel
}
func (logic Logic) GetOutgoingMessagesChannel() network.ServerMessagesChannel {
	return logic.OutgoingMessages
}
func (logic *Logic) SetOutgoingMessagesChannel(channel network.ServerMessagesChannel) {
	logic.OutgoingMessages = channel
}

func (logic *Logic) GetUserList(exceptId uint64) []network.User {
	var userlist []network.User
	for userId, user := range logic.Users {
		if userId != exceptId {
			userlist = append(userlist, network.User{Id: user.Id, Name: user.Name})
		}
	}

	log.Printf("Userlist: %#v\n", userlist)

	return userlist
}

/**
 * отправляет сообщение. Первый массив обозначает список целей кому передавать. Второй массив обозначает кому не передавать.
 * @param  {[type]} logic *Logic) SendMessage(msg Message, targets ...[]int [description]
 * @return {[type]} [description]
 */
func (logic *Logic) SendMessage(msg interface{}, targets ...network.UsersList) {
	serverMessage := network.ServerMessage{Data: msg}

	// real targets
	if len(targets) > 0 {
		serverMessage.Targets = targets[0]
	}

	// except this users
	if len(targets) > 1 {
		serverMessage.Except = targets[1]
	}

	select {
	case logic.OutgoingMessages <- serverMessage:
	default:
		log.Println("busy outgoing messages chan")
	}

}

func (logic *Logic) SendTextMessage(text string, sender uint64) {
	logic.SendMessage(network.TextMessage{Text: text, Sender: sender})
}

func (logic *Logic) SendTextMessageToUser(text string, sender uint64, userId uint64) {
	logic.SendMessage(network.TextMessage{Text: text, Sender: sender}, network.UsersList{userId})
}

func (logic *Logic) AddUser(id uint64, name string) *network.User {
	user := &network.User{Id: id, Name: name}
	logic.Users[id] = user

	pos := world.Point2D{X: rand.Float64()*600 - 300, Y: rand.Float64()*600 - 300}
	logic.mWorldMap.AddUser(user.Id, pos)

	return user
}

func (logic *Logic) RemoveUser(id uint64) {
	delete(logic.Users, id)
	logic.mWorldMap.RemoveUser(id)
}

func (logic *Logic) sendSyncPositionMessage() {
	logic.SendMessage(
		network.SyncPositionsMessage{
			logic.mWorldMap.GetObjectsPositions(),
			logic.mWorldMap.GetCurrentTimeMillis(),
		},
	)
	logic.prevSyncTime = time.Now()
}

// Возвращает true, если нужно синхронизировать положение объектов заново
func (logic *Logic) ProcessActionMessage(userId uint64, msg *network.ActionMessage) (needSync bool) {
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

func (logic *Logic) ProcessMessage(message network.UserMessage) (needSync bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in %#v\n", r)
		}
	}()

	needSync = false

	switch msg := message.Data.(type) {
	case *network.DataMessage:
		log.Println("Logic: Data message received: ", message)
	case *network.TextMessage:
		// log.Println("Text message received: ", message)
		logic.SendTextMessage(msg.Text, logic.Users[message.Source].Id)
	case *network.LoginMessage:
		log.Println("Logic: Login message received")

		user := logic.AddUser(msg.Id, msg.Name)
		logic.SendTextMessageToUser("Logic: Wellcome, "+user.Name, 0, user.Id)
		logic.SendMessage(logic.getServerStateMessage(), network.UsersList{user.Id})
		logic.SendMessage(network.UserLoggedinMessage{Id: user.Id, Name: user.Name}, network.UsersList{}, network.UsersList{user.Id})
		logic.SendMessage(network.UserListMessage{logic.GetUserList(user.Id)}, network.UsersList{user.Id})
		log.Println("sent. sync next")
		logic.SendMessage(network.SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), logic.mWorldMap.GetCurrentTimeMillis()})
	case *network.LogoutMessage:
		log.Println("Logic: Logout message", msg.Id)
		logic.RemoveUser(msg.Id)
		logic.SendMessage(network.UserLoggedoutMessage{Id: msg.Id})
	case *network.ActionMessage:
		needSync = logic.ProcessActionMessage(message.Source, msg)
	case *network.SimulateMessage:
		if logic.params.SimulateByStep {
			select {
			case logic.forceSimulationChannel <- msg.Steps:
			default:
				log.Println("Already busy force chan")
			}
		} else {
			log.Println("We are not in step by step mode")
		}
	case *network.ChangeSimulationMode:
		newValue := msg.StepByStep
		if logic.params.SimulateByStep != newValue {
			select {
			case logic.changeSimulationModeChannel <- newValue:
			default:
				log.Println("Already busy mode chan")
			}
		} else {
			log.Printf("Simulation already in mode %v\n", newValue)
		}
	default:
		log.Printf("Logic: Unknown message type %#v from %d\n", message.Data, message.Source)
	}

	return
}

func (logic *Logic) getServerStateMessage() network.ServerStateMessage {
	return network.ServerStateMessage{
		SimulationByStep:       logic.params.SimulateByStep,
		SimulationStepTime:     (uint64)(logic.params.SimulationStepTime / time.Millisecond),
		SimulationStepRealTime: (uint64)(logic.params.SimulationStepRealTime / time.Millisecond),
		ServerTime:             (uint64)(time.Now().UnixNano() / int64(time.Millisecond)),
		SimulationTime:         logic.mWorldMap.GetCurrentTimeMillis(),
	}
}

func (logic *Logic) executeSimulation(dt time.Duration) (changed bool) {
	changed = logic.mWorldMap.ProcessSimulationStep(logic.params.SimulationStepTime)
	logic.PrevStepTime = time.Now()
	logic.NextStepTime = logic.NextStepTime.Add(logic.params.SimulationStepRealTime)
	return
}

func (logic *Logic) Start() {
	rand.Seed(int64(time.Now().Nanosecond()))

	logic.Users = make(map[uint64]*network.User)
	logic.changeSimulationModeChannel = make(chan bool, 2)
	logic.forceSimulationChannel = make(chan int, 1)

	logic.mWorldMap = world.NewWorldMap()

	logic.NextStepTime = time.Now()
	logic.prevSyncTime = time.Unix(0, 0)
	// таймер, который инициализирует симуляцию
	var simulationTimer = time.NewTimer(0)

	if logic.params.SimulateByStep {
		if !simulationTimer.Stop() {
			<-simulationTimer.C
		}
	}

	log.Println("Logic: started")
	for {
		// сначала пытаемся вычитаь все входящие соединения.
		// потом надо будет учесть вариант, что нас могут заспамить и симуляция вообще не произойдёт.
		select {
			case msg := <-logic.IncomingMessages:
				//log.Println("Logic: message received")
				if needSync := logic.ProcessMessage(msg); needSync {
					logic.sendSyncPositionMessage()
				}
				//log.Println("Logic: message processed")
				continue
			default:
		}

		select {
		case mode := <-logic.changeSimulationModeChannel:
			// обработка смены режима симуляции
			if logic.params.SimulateByStep != mode {
				logic.params.SimulateByStep = mode

				if logic.params.SimulateByStep {
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
				logic.SendMessage(logic.getServerStateMessage())
				logic.sendSyncPositionMessage()
				log.Println("server state message sent")
			} else {
				log.Printf("Simulation mode already was %v\n", mode)
			}
		case _ = <-simulationTimer.C:
			// по идее уже пора выполнять очередной шаг симуляции
			//log.Println("Timer fired!")
			//log.Printf("Now %v next %v\n", time.Now(), logic.NextStepTime)

			stepsCount := 0
			globallyChanged := false

			if (!logic.NextStepTime.Equal(time.Now()) && !logic.NextStepTime.Before(time.Now())) {
				log.Println("WARNING! simulation timer fired before next step!")
			}

			startTime := time.Now()

			// если уже пора симулировать, то симулируем, н оне больше 10 шагов
			for (logic.NextStepTime.Equal(time.Now()) || logic.NextStepTime.Before(time.Now())) && stepsCount < logic.params.MaxSimulationStepsAtOnce {
				// если что-то изменилось - нужно разослать всем уведомления
				changed := logic.executeSimulation(logic.params.SimulationStepTime)
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
			if logic.params.SimulateByStep {
				log.Println("Simulating!")
				startTime := time.Now()
				logic.executeSimulation(logic.params.SimulationStepTime)
				passedTime := time.Now().Sub(startTime)
				log.Printf("Simulated 1 step (%d mcs): world time %d ms", passedTime/time.Microsecond, logic.mWorldMap.GetCurrentTimeMillis())
				logic.sendSyncPositionMessage()
			} else {
				log.Println("Not in step by step mode")
			}
		case msg := <-logic.IncomingMessages:
			//log.Println("Logic: message received")
			if needSync := logic.ProcessMessage(msg); needSync {
				logic.sendSyncPositionMessage()
			}
			//log.Println("Logic: message processed")
		}
	}

	log.Println("Logic: finished")
}

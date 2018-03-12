package main

import (
	"github.com/porfirion/server2/network"
	"github.com/porfirion/server2/world"
	"log"
	"math/rand"
	"time"
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

	logic.OutgoingMessages <- serverMessage
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

func (logic *Logic) sendSyncMessage() {
	var currentTime = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond)
	logic.SendMessage(network.SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
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
		var currentTime = logic.mWorldMap.SimulationTime.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond)
		logic.SendMessage(network.SyncPositionsMessage{logic.mWorldMap.GetObjectsPositions(), currentTime})
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
				log.Println("Pushed to force simulation chan")
			default:
				log.Println("Already busy force chan")
			}
		} else {
			log.Println("We are not in step by step mode")
		}
	case *network.ChangeSimulationMode:
		newValue := msg.StepByStep
		if logic.params.SimulateByStep != newValue {

			log.Printf("Changing simulation mode from %b to %b\n", logic.params.SimulateByStep, newValue)
			logic.params.SimulateByStep = newValue

			// а теперь уведомляем всех об изменившемся режиме
			logic.SendMessage(logic.getServerStateMessage())
		} else {
			log.Printf("Simulation already in mode %b\n", newValue)
		}
	default:
		log.Printf("Logic: Unknown message type %#v from %d\n", message.Data, message.Source)
	}

	return
}

func (logic *Logic) getServerStateMessage() network.ServerStateMessage {
	return network.ServerStateMessage{
		SimulationByStep:       logic.params.SimulateByStep,
		SimulationStepTime:     (uint64)(logic.params.SimulationStepTime.Nanoseconds() / 1000),
		SimulationStepRealTime: (uint64)(logic.params.SimulationStepRealTime.Nanoseconds() / 1000),
		ServerTime:             (uint64)(time.Now().UnixNano() / 1000),
	}
}

func (logic *Logic) TimeToNextStep() time.Duration {
	if logic.NextSimulationStepTime().After(time.Now()) {
		return logic.NextSimulationStepTime().Sub(time.Now())
	} else {
		return 0
	}
}

func (logic *Logic) NextSimulationStepTime() time.Time {
	if logic.params.SimulateByStep {
		// следующий шаг симуляции не произойдёт никогда! мухахаха
		return time.Date(3000, time.January, 0, 0, 0, 0, 0, time.Local)
	} else {
		return logic.PrevStepTime.Add(logic.params.SimulationStepRealTime)
	}
}
func (logic *Logic) executeSimulation(dt time.Duration) (changed bool) {
	changed = logic.mWorldMap.ProcessSimulationStep(logic.params.SimulationStepTime)
	logic.PrevStepTime = time.Now()
	logic.NextStepTime.Add(logic.params.SimulationStepRealTime)
	return
}

func (logic *Logic) Start() {
	rand.Seed(int64(time.Now().Nanosecond()))

	logic.Users = make(map[uint64]*network.User)

	logic.forceSimulationChannel = make(chan int, 1)

	logic.mWorldMap = world.NewWorldMap()

	// таймер, который инициализирует симуляцию
	var simulationTimer = time.NewTimer(logic.TimeToNextStep())

	if logic.params.SimulateByStep {
		if !simulationTimer.Stop() {
			<-simulationTimer.C
		}
	}

	log.Println("Logic: started")
	for {
		select {
		case mode := <-logic.changeSimulationModeChannel:
			if logic.params.SimulateByStep != mode {
				log.Printf("Simulation mode changed to %v\n", mode)
				logic.params.SimulateByStep = mode

				if logic.params.SimulateByStep {
					// симуляция по шагам
					if !simulationTimer.Stop() {
						<-simulationTimer.C
					}
				} else {
					// непрерывная симуляция
					logic.PrevStepTime = time.Now().Add(-1 * logic.params.SimulationStepRealTime - 1)
					logic.NextStepTime = time.Now()
					simulationTimer.Reset(0)
				}
			} else {
				log.Printf("Simulation mode already was %v\n", mode)
			}
		case _ = <-simulationTimer.C:
			//log.Println("Logic: simulation step")

			// ага, уже пора производить симуляцию

			stepsCount := 0
			globallyCahnged := false

			// если уже пора симулировать, то симулируем, н оне больше 10 шагов
			for logic.NextStepTime.Before(time.Now()) && stepsCount < logic.params.MaxSimulationStepsAtOnce {
				// если что-то изменилось - нужно разослать всем уведомления
				changed := logic.executeSimulation(logic.params.SimulationStepTime)
				globallyCahnged = globallyCahnged || changed
				stepsCount++
			}

			if globallyCahnged {
				logic.sendSyncMessage()
			}

			simulationTimer.Reset(logic.TimeToNextStep())
		case _ = <-logic.forceSimulationChannel:
			if logic.params.SimulateByStep {
				log.Println("Simulating!")
				logic.executeSimulation(logic.params.SimulationStepTime)
				logic.sendSyncMessage()
			} else {
				log.Println("Not in step by step mode")
			}
		case msg := <-logic.IncomingMessages:
			log.Println("Logic: message received")
			if needSync := logic.ProcessMessage(msg); needSync {
				logic.sendSyncMessage()
			}
			log.Println("Logic: message processed")
		}
	}

	log.Println("Logic: finished")
}

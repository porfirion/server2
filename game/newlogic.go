package game

type Logic interface {
	SetOutputChan(chan interface{})
	AddMessage(message interface{})
	Start()
}

type LogicImpl struct {

}

package next

import (
	"time"

	"github.com/porfirion/server2/service"
)

type svc struct {
	*service.BasicService
	logic *Logic
}

func (s *svc) Start() {
	for msg := range s.IncomingMessages {
		log.Printf("NewLogic: %+v", msg)
	}
}

func NewService() service.Service {
	return &svc{
		BasicService: service.NewBasicService(service.TypeLogic),
		logic: NewLogic(
			make(chan ControlMessage),
			make(chan PlayerInput),
			SimulationModeContinuous,
			time.Second,
			time.Second,
		),
	}
}

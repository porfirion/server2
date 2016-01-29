package main

import (
	"fmt"
	"log"
	"time"
)

type EventDispatcher struct {
	OutgoigEvents  chan Event
	IncomingEvents chan Event
	events         []Event
	nearestTime    int64
}

func (ed *EventDispatcher) Init() {
	ed.OutgoigEvents = make(chan Event)
	ed.IncomingEvents = make(chan Event)
	ed.events = make([]Event, 0)

	go func() {
		defer log.Println("EventDispatcher terminated")

		log.Println("EventDispatcher started")
		for {
			select {
			case event := <-ed.IncomingEvents:
				log.Println(fmt.Sprintf("Appending event %#v", event))
				ed.events = append(ed.events, event)
				if ed.nearestTime != 0 && event.Time < ed.nearestTime {
					ed.nearestTime = event.Time
				}
			case <-time.Tick(10 * time.Second):
				log.Println("EvDis. Tick!Tock!")
			}
		}
	}()
}

type Event struct {
	Time int64
}

type MoveEvent struct {
	Event
	Position Position
}

type TickEvent struct {
	Event
}

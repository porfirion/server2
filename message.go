package main

type Message interface {
}

type TextMessage struct {
	Text string
}

type DataMessage struct {
	Data []byte
}

type AuthMessage struct {
	Uuid string
}

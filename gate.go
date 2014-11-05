package main

type Gate interface {
	Start(chan MessageChannel)
}

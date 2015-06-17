package main

type Position struct {
	X int
	Y int
}

type MapObject struct {
	Pos         Position
	Destination Position
	Speed       int
	StartTime   int
	User        *User
}

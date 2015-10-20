package main

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MapObject struct {
	Pos         Position
	Destination Position
	Speed       int
	StartTime   int
	User        *User
}

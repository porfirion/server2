package main

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MapObject struct {
	Pos         Position
	Destination Position
	Speed       int
	StartTime   int
	User        *User
}

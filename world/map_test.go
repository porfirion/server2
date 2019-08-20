package world

import (
	"math/rand"
	"testing"
)

func TestNewWorldMap(t *testing.T) {
	m := NewWorldMap(10000, 10000)
	if m == nil {
		t.Fail()
	}
}

func TestWorldMap_NewObject(t *testing.T) {
	m := NewWorldMap(10000, 10000)
	m.NewObject(Point2D{0, 0}, MapObjectTypeObstacle)

	if len(m.Objects) != 1 { t.Error("There should be only one object") }
}

func TestWorldMap_RemoveObject(t *testing.T) {
	m := NewWorldMap(10000, 10000)
	ob := m.NewObject(Point2D{0, 0}, MapObjectTypeObstacle)
	m.RemoveObject(ob)

	if len(m.Objects) > 0 {
		t.Error("There should be no objects")
	}
}

func TestWorldMap_AddUser(t *testing.T) {
	m := NewWorldMap(10000, 10000)
	m.AddUser(rand.Uint64(), Point2D{0, 0})


}
package world

import (
	"log"
)

type CollisionResolver interface {
	resolve(obj1 *MapObject, obj2 *MapObject) bool
}

var (
	simpleResolver = &SimpleResolver{}
)

type SimpleResolver struct {
}

func (r SimpleResolver) resolve(obj1 *MapObject, obj2 *MapObject) bool {
	log.Printf("============================================\n")
	// строим линию между центрами
	line1To2 := LineByPoints(obj1.CurrentPosition, obj2.CurrentPosition)

	// берём от неё направляющий вектор
	directing1To2 := obj1.CurrentPosition.VectorTo(obj2.CurrentPosition)
	// и делаем из него единичный
	vector1to2 := directing1To2.Unit()

	distance := obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition)
	penetration := obj1.Size + obj2.Size - distance

	log.Printf("\n\t%6d: %v size %.2f\n\t%6d: %v size %.2f\n\t  line: %v\n\tdirect: %v\n\t  vect: %v\n\t  dist: %f\n\tpenetr: %f\n",
		obj1.Id,
		obj1.CurrentPosition,
		obj1.Size,
		obj2.Id,
		obj2.CurrentPosition,
		obj2.Size,
		line1To2,
		directing1To2,
		vector1to2,
		distance,
		penetration)

	if penetration < 0 {
		log.Printf("No penetration! %v %v dist %f (%f + %f)\n", obj1.CurrentPosition, obj2.CurrentPosition, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
		// это какая-то ерунда - получается, что соприкосновения и не было!
		return false
	}

	sumMass := (float64)(obj1.Mass + obj2.Mass)
	var obj1Proportion = float64(obj1.Mass) / sumMass
	var obj2Proportion = 1 - obj1Proportion
	log.Printf("mass proportions: 1 - %f  2 - %f \n", obj1Proportion, obj2Proportion)

	prev1, prev2 := obj1.CurrentPosition, obj2.CurrentPosition

	obj1.CurrentPosition = obj1.CurrentPosition.Move(vector1to2.Revers().Mult(penetration * obj1Proportion))
	obj2.CurrentPosition = obj2.CurrentPosition.Move(vector1to2.Mult(penetration * obj2Proportion))

	log.Printf("\n\t%6d: {%f, %f} -> {%f, %f}\n\t%6 d: {%f, %f} -> {%f, %f}\n",
		obj1.Id,
		prev1.X, prev1.Y,
		obj1.CurrentPosition.X, obj1.CurrentPosition.Y,
		obj2.Id,
		prev2.X, prev2.Y,
		obj2.CurrentPosition.X, obj2.CurrentPosition.Y,
	)

	//center := Point2D{obj1.X + (obj2.X - obj1.X) * obj1Proportion, obj1.Y + (obj2.Y - obj1.Y) * obj1Proportion};

	return true
}

func GetResolver(obj1 *MapObject, obj2 *MapObject) CollisionResolver {
	// temporary hack - always returning simple resolver
	return simpleResolver
}

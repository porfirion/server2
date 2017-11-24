package world

import (
	"fmt"
)

type CollisionResolver interface {
	resolve(obj1 *MapObject, obj2 *MapObject);
}

var (
	simpleResolver *SimpleResolver = &SimpleResolver{}
)

type SimpleResolver struct {
}

func (r SimpleResolver) resolve(obj1 *MapObject, obj2 *MapObject) {
	fmt.Printf("============================================\n")
	// строим линию между центрами
	line1To2 := LineByPoints(obj1.CurrentPosition, obj2.CurrentPosition)

	// берём от неё направляющий вектор
	directing1To2 := obj1.CurrentPosition.VectorTo(obj2.CurrentPosition)
	// и делаем из него единичный
	vector1to2 := directing1To2.Unit()

	distance := obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition)
	penetration := obj1.Size + obj2.Size - distance;

	fmt.Printf("%d: %v size %f\n%d: %v size %f\nline: %v\ndirect: %v\nvect: %v\ndist: %f\npenetr: %f\n",
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

	if (penetration < 0) {
		fmt.Printf("No penetration! %v %v dist %f (%f + %f)", obj1.CurrentPosition, obj2.CurrentPosition, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
		// это какая-то ерунда - получается, что соприкосновения и не было!
		return;
	}

	sumMass := (float64)(obj1.Mass + obj2.Mass);
	var obj1Proportion = float64(obj1.Mass) / sumMass;
	var obj2Proportion = 1 - obj1Proportion;
	fmt.Printf("mass proportions: 1 - %f  2 - %f \n", obj1Proportion, obj2Proportion)

	prev1, prev2 := obj1.CurrentPosition, obj2.CurrentPosition

	obj1.CurrentPosition = obj1.CurrentPosition.Move(vector1to2.Revers().Mult(penetration * obj1Proportion));
	obj2.CurrentPosition = obj2.CurrentPosition.Move(vector1to2.Mult(penetration * obj2Proportion))

	fmt.Printf("{%f, %f} -> {%f, %f} and\n{%f, %f} -> {%f, %f}\n",
		prev1.X, prev1.Y,
		obj1.CurrentPosition.X, obj1.CurrentPosition.Y,
		prev2.X, prev2.Y,
		obj2.CurrentPosition.X, obj2.CurrentPosition.Y,
	)

	//center := Point2D{obj1.X + (obj2.X - obj1.X) * obj1Proportion, obj1.Y + (obj2.Y - obj1.Y) * obj1Proportion};
}

func GetResolver(obj1 *MapObject, obj2 *MapObject) CollisionResolver {
	// temporary hack - always returning simple resolver
	return simpleResolver;
}

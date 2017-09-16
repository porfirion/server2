package world

import (
	"log"
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
	log.Printf("Resolving %d VS %d\n", obj1.Id, obj2.Id)
	// строим линию между центрами
	line1to2 := LineByPoints(obj1.CurrentPosition, obj2.CurrentPosition).Directing().Unit()
	fmt.Println(obj1.CurrentPosition, obj2.CurrentPosition, line1to2)

	penetration := (-1)*(obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition) - obj1.Size - obj2.Size);
	if (penetration < 0) {
		log.Printf("No penetration! %v %v dist %f (%f + %f)", obj1.CurrentPosition, obj2.CurrentPosition, obj1.CurrentPosition.DistanceTo(obj2.CurrentPosition), obj1.Size, obj2.Size)
		// это какая-то ерунда - получается, что соприкосновения и не было!
		return;
	}

	sumMass := (float64)(obj1.Mass + obj2.Mass);
	var obj1Proportion float64 = float64(obj1.Mass) / sumMass;
	var obj2Proportion float64 = 1 - obj1Proportion;

	prev1, prev2 := obj1.CurrentPosition, obj2.CurrentPosition

	obj1.CurrentPosition = obj1.CurrentPosition.Move(line1to2.Opposite().Mult(penetration * obj1Proportion));
	obj2.CurrentPosition = obj2.CurrentPosition.Move(line1to2.Mult(penetration * obj2Proportion))

	log.Printf("{%f, %f} -> {%f, %f} and {%f, %f} -> {%f, %f}\n",
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

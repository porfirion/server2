package world

import "math"

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	NilPosition Point2D = Point2D{X: math.MaxFloat64, Y: math.MaxFloat64}
)

/**
 * Расстояние между точками
 */
func (pos Point2D) DistanceTo(dest Point2D) float64 {
	return math.Sqrt(math.Pow(dest.X-pos.X, 2) + math.Pow(dest.Y-pos.Y, 2))
}

func (pos Point2D) VectorTo(dest Point2D) Vector2D {
	return Vector2D{dest.X - pos.X, dest.Y - pos.Y}
}

type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// длина вектора
func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// приведение длины вектора
func (v Vector2D) Modulus(base float64) Vector2D {
	modulo := v.Length() / base
	return Vector2D{v.X / modulo, v.Y / modulo}
}

func (v Vector2D) Plus(v2 Vector2D) Vector2D {
	return Vector2D{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2D) Minus(v2 Vector2D) Vector2D {
	return Vector2D{v.X - v2.X, v.Y - v2.Y};
}

func (v Vector2D) Devide(devider float64) Vector2D {
	return Vector2D{v.X / devider, v.Y / devider}
}
func (v Vector2D) Mult(multiplier float64) Vector2D {
	return Vector2D{v.X * multiplier, v.Y * multiplier}
}
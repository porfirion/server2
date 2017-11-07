package world

import (
	"math"
	"fmt"
)

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	NilPosition = Point2D{X: 0, Y: 0}
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

func (pos Point2D) Move(v Vector2D) Point2D {
	return Point2D{pos.X + v.X, pos.Y + v.Y}
}

func (pos Point2D) String() string {
	return fmt.Sprintf("{%f, %f}", pos.X, pos.Y)
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
// сложение векторов
func (v Vector2D) Plus(v2 Vector2D) Vector2D {
	return Vector2D{v.X + v2.X, v.Y + v2.Y}
}
// вычитание векторов
func (v Vector2D) Minus(v2 Vector2D) Vector2D {
	return Vector2D{v.X - v2.X, v.Y - v2.Y};
}
// деление вектора на число
func (v Vector2D) Devide(devider float64) Vector2D {
	return Vector2D{v.X / devider, v.Y / devider}
}
// умножение вектора на число
func (v Vector2D) Mult(multiplier float64) Vector2D {
	return Vector2D{v.X * multiplier, v.Y * multiplier}
}
// единичный вектор, колинеарный данному
func (v Vector2D) Unit() Vector2D {
	var sum float64 = 1 / math.Sqrt(v.X * v.X + v.Y * v.Y);
	return Vector2D{v.X * sum, v.Y * sum}
}
// обратный вектор
func (v Vector2D) Revers() Vector2D {
	return Vector2D{-v.X, -v.Y}
}

func (v Vector2D) String() string {
	return fmt.Sprintf("{{%f, %f}}", v.X, v.Y)
}

type Line2D struct {
	// Ax + By + C = 0
	A, B, C float64

    // (-B; A) - направляющий вектор
	// (A; B) - нормаль
}

func LineByPoints(p1, p2 Point2D) Line2D {
	// X(y1 - y2) + Y (x2 - x1) + x1y1 - x2y1 = 0
	var res Line2D = Line2D{};
	res.A = p1.Y - p2.Y;
	res.B = p2.X - p1.X;
	res.C = p1.X * p2.Y - p2.X * p1.Y;
	return res;
}

// направляющий вектор (он же колинеарный) к прямой
func (l Line2D) Directing() Vector2D {
	return Vector2D{-l.B, l.A}
}

// нормаль/нормальный (перпендикулярный) вектор к прямой
func (l Line2D) Normal() Vector2D {
	return Vector2D{l.A, l.B}
}

// уравнение перпендикулярной прямой через точку
func (l Line2D) Perpendicular(p Point2D) Line2D {
	// https://www.desmos.com/calculator/ywbgran4rg
	return Line2D{-l.B, l.A, p.X * l.B - p.Y * l.A}
}

type Circle struct {
	Center Point2D
	Radius float64
}

type Rectangle struct {
	Center Point2D
	Height, Width float64
	Angle float64
}

type HairLine struct {
	Start Point2D
	End Point2D
}
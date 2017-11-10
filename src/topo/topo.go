package topo

type Shape interface {
	Area() float64
	Perimeter() float64
	isInside(Point) bool
}

package topo

import (
	"math"
)

type Circle struct {
	center Point
	radius float64
}

func GetCircle(center Point, radius float64) Circle {
	return Circle{center,radius}
}

func (c Circle) GetCenter() Point {
	return c.center
}

func (c Circle) GetRadius() float64 {
	return c.radius
}

func (c Circle) Area() float64 {
	return math.Pi*math.Pow(c.radius, 2)
}

func (c Circle) Perimeter() float64 {
	return 2*math.Pi*c.radius
}

func (c Circle) IsInside(p Point) bool {
	if(c.center.DistanceFrom(p) <= c.radius){
		return true
	}
	return false
}
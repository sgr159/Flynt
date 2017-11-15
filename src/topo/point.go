package topo

import (
	"math"
)

const TOLERANCE = 0.001

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceFrom(p2 Point) float64 {
	return math.Sqrt((math.Pow(p.X - p2.X, 2) + math.Pow(p.Y - p2.Y, 2)))
}

func Distancebetween(p,p2 Point) float64 {
	return math.Sqrt((math.Pow(p.X - p2.X, 2) + math.Pow(p.Y - p2.Y, 2)))
}

func (p Point) IsEqual(p2 Point) bool {
	if(math.Abs(p.X - p2.X) < TOLERANCE && math.Abs(p.Y - p2.Y) < TOLERANCE){
		return true
	}
	return false
}
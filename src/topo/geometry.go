package topo

import (
	"math"
)

type Degrees float64
type Radians float64

 func GetEquiGeoCoordinates(numOfVertices uint64, distanceFromCenter float64) []Point {
 	var points []Point
 	points = make([]Point,numOfVertices)
 	//centered around origin, take the first coordinate on x axis
 	for i,_ := range points {
 		if i==0 {
 			points[i] = Point{distanceFromCenter,0}
 			continue
 		}
 		points[i] = rotateLeft(Point{0,0}, points[i-1], Radians(2*math.Pi/float64(numOfVertices)))
 	}
 	return points
 }
 
 func rotateLeft(anchor Point, end Point, angleRad Radians) Point {
 	/*
 	Matrix multiplication
 	[x3,y3]=[cosθ -sinθ, sinθ cosθ][x2−x1,y2−y1]+[x1,y1]
 	*/
 	angle := float64(angleRad)
 	var x,y,xdiff,ydiff float64
 	xdiff = end.X - anchor.X
 	ydiff = end.Y - anchor.Y
 	x = math.Cos(angle)*xdiff - math.Sin(angle)*ydiff + anchor.X
 	y = math.Cos(angle)*ydiff + math.Sin(angle)*xdiff + anchor.Y
 	return Point{x,y}
 }
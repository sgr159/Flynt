package main

import (
	"fmt"
//	"topo"
	"drone"
)

var IDVAR uint64 = 1

func main() {
	var field = GetNewField(50, 50)
	fmt.Println("Enter number of drones: ")
	var numOfDrones uint64
	fmt.Scan(&numOfDrones)
	for i:=0;i<5;i++ {
		field.AddDrone(drone.GetDrone(IDVAR,100,3))
		IDVAR++
	}
	field.ArrangeDrones()
	field.PlotField()
	/*
	fmt.Println("Enter number of drones: ")
	var numOfDrones uint64
	fmt.Scan(&numOfDrones)
	
	fmt.Println("Enter distance from center: ")
	var distanceFromCenter float64
	fmt.Scan(&distanceFromCenter)
	fmt.Println("Distance from center:",distanceFromCenter)
	var positions []topo.Point
	positions = append(positions,topo.Point{0,0})
	positions = append(positions,topo.GetEquiGeoCoordinates(uint64(numOfDrones-1), distanceFromCenter)...)
	fmt.Println("Initial positions of drones:")
	for _,point := range positions {
		fmt.Printf("{%.2f, %.2f} ",point.X,point.Y)
	}
	fmt.Printf("\n")
	userPos := []topo.Point{topo.Point{3,4},topo.Point{-1,4},topo.Point{-1,-2}}
	plotMap(positions,userPos)
	
	*/
	fmt.Println("Check out points.png yo!")
}

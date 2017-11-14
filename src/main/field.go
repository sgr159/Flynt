package main

import (
	"drone"
	"user"
	"topo"
)

type Field struct {
	//length and breadth of the field. Drones cannot be deployed outside this limit
	length float64
	breadth float64
	droneGroups []*drone.DroneGroup
	users []*user.User
	numOfDrones uint64
}

func GetNewField(length, breadth float64) *Field {
	drone.XMax = length/2
	drone.XMin = length/2 * -1
	drone.YMax = breadth/2
	drone.YMin = breadth/2 * -1
	
	return &Field{length,breadth,nil,nil,0}
}

func (f *Field)AddDrone(d *drone.Drone) {
	droneGrp := drone.GetDroneGroup(IDVAR)
	IDVAR++
	droneGrp.AddDrone(d)
	f.droneGroups = append(f.droneGroups,droneGrp)
	f.numOfDrones++
}

func (f *Field)AddUser(u *user.User) {
	f.users = append(f.users,u)
}

func (f *Field)GetDronePositions() []topo.Point {
	var dronePoints []topo.Point
	
	for _,dg := range f.droneGroups {
		drones :=dg.GetDrones()
		
		for _,dr := range drones {
			dronePoints = append(dronePoints,dr.GetCurrentPosition())
		}
	}
	return dronePoints
}

func (f *Field)GetUserPositions() []topo.Point {
	var userPoints []topo.Point
	
	for _,u := range f.users {
		userPoints = append(userPoints,u.GetCurrentPosition())
	}
	
	return userPoints
}

func (f *Field)PlotField() {
	var dronePoints,userPoints []topo.Point
	
	dronePoints = f.GetDronePositions()
	userPoints = f.GetUserPositions()
	
	plotMap(dronePoints, userPoints)
}

func (f *Field)ArrangeDrones() {
	equiPoints := topo.GetEquiGeoCoordinates(f.numOfDrones - 1, 5)
	for i:=1;i<=len(equiPoints);i++ {
		f.droneGroups[i].GetDrones()[0].MoveTo(equiPoints[i-1])
	}
}

func (f *Field)GetClosestDroneGroup (p topo.Point) *drone.DroneGroup {
	return nil
}
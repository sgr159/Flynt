package main

import (
	"fmt"
	"topo"
	"drone"
	"user"
)

var IDVAR uint64 = 0

func getId() uint64 {
	IDVAR++
	return IDVAR
}
func main() {
	var field = GetNewField(200, 200)
	userPoints := []topo.Point{
		topo.Point{1,1},
		topo.Point{12,1},
		topo.Point{-4,-9},
		topo.Point{-7,-5},
//		topo.Point{-43,65},
	}
	
	for _,p := range userPoints {
		field.AddUser(user.GetNewUser(getId(), p))
	}
	
	fmt.Println("Enter number of drones: ")
	var numOfDrones int
	fmt.Scan(&numOfDrones)
	for i:=0;i<numOfDrones;i++ {
		field.AddDrone(drone.GetDrone(getId(),100,3))
	}
	clustpos := field.ClusterUsers()
	field.ArrangeDrones()
	_,userPointsconf := field.GetUsersAndPositions()
	_,dronePoints := field.GetDronesAndDronePositions()
	fmt.Println("User Points:",userPointsconf)
	for _,ug := range field.GetUserGroups() {
		fmt.Println("center",ug.GetCenterofGroup(),"points",ug.GetUserPoints())
	}
	fmt.Println("CLUSTER COMP:",clustpos)
	fmt.Println("drone Points:",dronePoints)
	field.PlotField()

	fmt.Println("Check out points.png yo!")
}

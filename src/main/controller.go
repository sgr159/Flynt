package main

import (
	"fmt"
	"topo"
	"drone"
	"user"
	"math/rand"
)

var IDVAR uint64 = 0

func getId() uint64 {
	IDVAR++
	return IDVAR
}
func main() {
	var field = GetNewField(200, 200)
	fmt.Println("Enter number of users: ")
	var numOfUsers int
	fmt.Scan(&numOfUsers)
	
	for i:=0;i<numOfUsers;i++ {
		field.AddUser(user.GetNewUser(getId(), topo.Point{(rand.Float64() -0.5)*200,(rand.Float64() -0.5)*200}))
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
	field.PlotField2()

	fmt.Println("Check out points.png yo!")
}

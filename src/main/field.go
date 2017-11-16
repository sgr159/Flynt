package main

import (
	"drone"
	"user"
	"topo"
	"fmt"
)

type Field struct {
	//length and breadth of the field. Drones cannot be deployed outside this limit
	length float64
	breadth float64
	droneGroups []*drone.DroneGroup
	userGroups []*user.UserGroup
	numOfDrones uint64
	numOfUsers uint64
}

func GetNewField(length, breadth float64) *Field {
	drone.XMax = length/2
	drone.XMin = length/2 * -1
	drone.YMax = breadth/2
	drone.YMin = breadth/2 * -1
	
	return &Field{length,breadth,nil,nil,0,0}
}

func (f *Field)AddDrone(d *drone.Drone) {
	droneGrp := drone.GetDroneGroup(getId())
	droneGrp.AddDrone(d)
	f.droneGroups = append(f.droneGroups,droneGrp)
	f.numOfDrones++
}

func (f *Field)AddUser(u *user.User) {
	ug := user.GetUserGroup(0)
	ug.AddUser(u)
	f.userGroups = append(f.userGroups,ug)
	f.numOfUsers++
}

func (f *Field)GetUserGroups() []*user.UserGroup{
	return f.userGroups
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

func (f *Field) GetUsersAndPositions() ([]*user.User,[]topo.Point) {
	var userPoints []topo.Point
	var users []*user.User
	
	for _,ug := range f.userGroups {
		userPoints = append(userPoints,ug.GetUserPoints()...)
		users = append(users,ug.GetUsers()...)
	}
	
	return users,userPoints
}

func (f *Field) GetDronesAndDronePositions() ([]*drone.Drone,[]topo.Point) {
	var dronePoints []topo.Point
	var drones []*drone.Drone
	
	for _,dg := range f.droneGroups {
		dronePoints = append(dronePoints,dg.GetDronePoints()...)
		drones = append(drones,dg.GetDrones()...)
	}
	
	return drones,dronePoints
}

func (f *Field)PlotField() {
	dronePoints := f.GetDronePositions()
	_,userPoints := f.GetUsersAndPositions()
	
	plotMap(dronePoints, userPoints)
}

func (f *Field) ClusterUsers () []int {
	users,userPoints := f.GetUsersAndPositions()
	var numOfClusters uint64
	if f.numOfDrones-1 < f.numOfUsers {
		numOfClusters = f.numOfDrones-1
	} else {
		numOfClusters = f.numOfUsers
	}
	clusterPos := cluster(userPoints, int(numOfClusters))
	for i,n := range clusterPos {
		if users[i].GetGroup().GetClusterNumber() == uint64(n) {
			continue;
		}
		if ug:=f.GetUserGroup(uint64(n)); ug != nil {
			fmt.Println("cluster",n,"center before",ug.GetCenterofGroup())
			users[i].ChangeGroupTo(ug)
			fmt.Println("cluster",n,"center after",ug.GetCenterofGroup())
			continue;
		}
		//find available cluster
		if ug:=f.GetUserGroup(uint64(0)); ug != nil {
			ug.SetClusterNumber(uint64(n))
			fmt.Println("cluster",n,"center before",ug.GetCenterofGroup())
			users[i].ChangeGroupTo(ug)
			fmt.Println("cluster",n,"center after",ug.GetCenterofGroup())
			continue;
		}
	}
	return clusterPos
}

func (f *Field) GetUserGroup (cluster uint64) *user.UserGroup {
	for _,ug := range f.userGroups {
		if ug.GetClusterNumber() == cluster {
			return ug
		}
	}
	return nil
}

func (f *Field)ArrangeDrones() {
	equiPoints := topo.GetEquiGeoCoordinates(f.numOfDrones - 1, 5)
	f.droneGroups[0].GetDrones()[0].Serve() //anchor node
	for i:=1;i<=len(equiPoints);i++ {
		f.droneGroups[i].GetDrones()[0].MoveTo(equiPoints[i-1])
	}
	
	for _,ug := range f.userGroups {
		if len(ug.GetUsers()) == 0 {
			continue;
		}
		dg := f.GetClosestDroneGroup(ug.GetCenterofGroup()).GetDrones()
		dg[0].MoveTo(ug.GetCenterofGroup())
		dg[0].Serve()
		fmt.Println("Moving drone: ",dg[0].GetId(),"to cluster",ug.GetClusterNumber(),"drone pos",dg[0].GetCurrentPosition())
	}
}

func (f *Field)GetClosestDroneGroup (p topo.Point) *drone.DroneGroup {
	var minDist float64
	var minDistDg *drone.DroneGroup
	init := false
	for _,dg := range f.droneGroups {
		if dg.GetDrones()[0].GetStatus() != drone.Free {
			continue;
		}
		if !init {
			minDist = dg.GetCenterOfGroup().DistanceFrom(p)
			minDistDg = dg
			init = true
		}
		dist := dg.GetCenterOfGroup().DistanceFrom(p)
		if dist < minDist {
			minDist = dist
			minDistDg = dg
		}
	}
	return minDistDg
}
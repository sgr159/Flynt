package main

import (
	"drone"
	"fmt"
	"topo"
	"user"
	"math"
	"errors"
)

type Field struct {
	//length and breadth of the field. Drones cannot be deployed outside this limit
	length      float64
	breadth     float64
	freeDrones  []*drone.Drone
	droneGroups []*drone.DroneGroup
	userGroups  []*user.UserGroup
	numOfDrones uint64
	numOfUsers  uint64
}

func GetNewField(length, breadth float64) *Field {
	drone.XMax = length / 2
	drone.XMin = length / 2 * -1
	drone.YMax = breadth / 2
	drone.YMin = breadth / 2 * -1

	return &Field{length, breadth, nil, nil, nil, 0, 0}
}

func (f *Field) AddDrone(d *drone.Drone) {
	if f.droneGroups == nil {
		f.droneGroups = append(f.droneGroups, drone.GetDroneGroup(getId()))
		d.SetRole(drone.Anchor)
		d.Serve(nil)
		f.droneGroups[0].AddDrone(d)
	} else {
		d.SetParent(f.droneGroups[0].GetDrones()[0])
		f.freeDrones = append(f.freeDrones, d)
	}
	f.numOfDrones++
}

func (f *Field) DetachDrone(d *drone.Drone) {
	var index int
	var dg *drone.DroneGroup
	for index, dg = range f.droneGroups {
		if dg == d.GetGroup() {
			break
		}
	}
	if index == 0 {
		fmt.Println("panic!, trying to detach from dg:", d.GetGroup().GetId(), "not found in field.droneGroups")
	}
	d.GetGroup().RemoveDrone(d)
	if len(dg.GetDrones()) == 0 {
		f.droneGroups[index] = f.droneGroups[len(f.droneGroups)-1]
		f.droneGroups[len(f.droneGroups)-1] = nil
		f.droneGroups = f.droneGroups[:len(f.droneGroups)-1]
	}
}

func (f *Field) Serve(d *drone.Drone, users []*user.User) {
	for i,fd := range f.freeDrones {
		if fd == d {
			f.freeDrones[i] = f.freeDrones[len(f.freeDrones)-1]
			f.freeDrones = f.freeDrones[:len(f.freeDrones)-1]
			break;
		}
	}
	d.Serve(users)
	if d.GetGroup() != nil {
		return
	}
	droneGrp := drone.GetDroneGroup(getId())
	droneGrp.AddDrone(d)
	f.droneGroups = append(f.droneGroups, droneGrp)
}

func (f *Field) AddUser(u *user.User) {
	if f.userGroups == nil {
		//usergroup for anchor
		ug := user.GetUserGroup(1)
		f.userGroups = append(f.userGroups, ug)
	}
	ug := user.GetUserGroup(0)
	ug.AddUser(u)
	f.userGroups = append(f.userGroups, ug)
	f.numOfUsers++
}

func (f *Field) GetAnchor() *drone.Drone {
	return f.droneGroups[0].GetDrones()[0]
}

func (f *Field) GetUserGroups() []*user.UserGroup {
	return f.userGroups
}

func (f *Field) GetDronePositions() []topo.Point {
	var dronePoints []topo.Point

	for _, dg := range f.droneGroups {
		drones := dg.GetDrones()

		for _, dr := range drones {
			dronePoints = append(dronePoints, dr.GetCurrentPosition())
		}
	}
	return dronePoints
}

func (f *Field) GetFreeDrones(numOfDrones int) ([]*drone.Drone, error) {
	if len(f.freeDrones) < numOfDrones {
		return nil, errors.New("do not have requested num of free drones")
	}
	return f.freeDrones[:numOfDrones], nil
}

func (f *Field) GetUsersAndPositions() ([]*user.User, []topo.Point) {
	var userPoints []topo.Point
	var users []*user.User

	for _, ug := range f.userGroups {
		userPoints = append(userPoints, ug.GetUserPoints()...)
		users = append(users, ug.GetUsers()...)
	}

	return users, userPoints
}

func (f *Field) GetDronesAndDronePositions() ([]*drone.Drone, []topo.Point) {
	var dronePoints []topo.Point
	var drones []*drone.Drone

	for _, dg := range f.droneGroups {
		dronePoints = append(dronePoints, dg.GetDronePoints()...)
		drones = append(drones, dg.GetDrones()...)
	}

	return drones, dronePoints
}

func (f *Field) PlotField() {
	dronePoints := f.GetDronePositions()
	_, userPoints := f.GetUsersAndPositions()

	plotMap(dronePoints, userPoints)
}

func (f *Field) ClusterUsers() []int {
	users, userPoints := f.GetUsersAndPositions()
	var numOfClusters uint64
	if (f.numOfDrones-1)/2 < f.numOfUsers {
		numOfClusters = (f.numOfDrones - 1)/2
	} else {
		numOfClusters = f.numOfUsers
	}
	clusterPos := cluster(userPoints, int(numOfClusters))
	for i, p := range userPoints {
		if p.DistanceFrom(topo.Point{0,0}) < drone.DroneRange {
			clusterPos[i] = 1
		} else {
			clusterPos[i]++
		}
	}
	for i, n := range clusterPos {
		n = n+1
		if users[i].GetGroup().GetClusterNumber() == uint64(n) {
			continue
		}
		if ug := f.GetUserGroup(uint64(n)); ug != nil {
			users[i].ChangeGroupTo(ug)
			continue
		}
		//find available cluster
		if ug := f.GetUserGroup(uint64(0)); ug != nil {
			ug.SetClusterNumber(uint64(n))
			users[i].ChangeGroupTo(ug)
			continue
		}
	}
	return clusterPos
}

func (f *Field) GetUserGroup(cluster uint64) *user.UserGroup {
	for _, ug := range f.userGroups {
		if ug.GetClusterNumber() == cluster {
			return ug
		}
	}
	return nil
}

func (f *Field) ArrangeDrones() {
	/*	equiPoints := topo.GetEquiGeoCoordinates(f.numOfDrones-1, 5)
		for i := 1; i <= len(equiPoints); i++ {
			f.droneGroups[i].GetDrones()[0].MoveTo(equiPoints[i-1])
		}
	*/
	for _, ug := range f.userGroups {
		if len(ug.GetUsers()) == 0 {
			continue
		}
		if ug.GetClusterNumber() == 1 {
			f.droneGroups[0].GetDrones()[0].Serve(ug.GetUsers())
			continue
		}
		f.processUserGroup(ug)
	}
}

func (f *Field) GetClosestDrone(p topo.Point) *drone.Drone {
	var minDist float64
	var minDistD *drone.Drone
	init := false
	for _, d := range f.freeDrones {
		if d.GetStatus() != drone.Free {
			continue
		}
		if !init {
			minDist = d.GetCurrentPosition().DistanceFrom(p)
			minDistD = d
			init = true
		}
		dist := d.GetCurrentPosition().DistanceFrom(p)
		if dist < minDist {
			minDist = dist
			minDistD = d
		}
	}
	return minDistD
}

func (f *Field) GetClosestServingDrone(p topo.Point) (float64, *drone.Drone){
	var minDist float64 = -1
	var minDistDrone *drone.Drone
	for _,dg := range f.droneGroups {
		for _,d := range dg.GetDrones() {
			if d.GetStatus() != drone.Serving {
				continue
			}
			if minDist == -1 {
				minDist = d.GetCurrentPosition().DistanceFrom(p)
				minDistDrone = d
			} else if d.GetCurrentPosition().DistanceFrom(p) < minDist {
				minDist = d.GetCurrentPosition().DistanceFrom(p)
				minDistDrone = d
			}
		}
	}
	return minDist, minDistDrone
}

func (f *Field) processUserGroup(ug *user.UserGroup) {
	/*
		u1,u2,dist := ug.GetMostDistantUsers()
		if(dist > 2 * drone.DroneRange) {
			fmt.Println("UG subclusters needed")
		}
	*/

	_, maxDist := ug.GetMostDistantUserFrom(topo.Point{0, 0})
	if maxDist < drone.DroneRange {
		//assign to anchor
		fmt.Println("cluster",ug.GetClusterNumber(),"assigned to anchor")
		f.Serve(f.droneGroups[0].GetDrones()[0], ug.GetUsers())
		return
	}
	
	numOfConnectorNodes := int(math.Ceil(ug.GetCenterofGroup().DistanceFrom(topo.Point{0,0})/(2*drone.DroneRange)))

	var prevD *drone.Drone = f.droneGroups[0].GetDrones()[0]
	fmt.Println("num of connectors:",numOfConnectorNodes)
	for i:=1;i<numOfConnectorNodes;i++ {
		ratio := float64(i)/float64(numOfConnectorNodes)
		p := topo.Point{ug.GetCenterofGroup().X*ratio, ug.GetCenterofGroup().Y*ratio}
		d := f.GetClosestDrone(p)
		d.MoveTo(p)
		f.Serve(d,nil)
		d.SetParent(prevD)
		prevD = d
	}
	d := f.GetClosestDrone(ug.GetCenterofGroup())
	d.MoveTo(ug.GetCenterofGroup())
	f.Serve(d,ug.GetUsers())
	d.SetParent(prevD)
	
	fmt.Println("Moving drone: ", d.GetId(), "to cluster", ug.GetClusterNumber(), "drone pos", d.GetCurrentPosition())
	//	numOfDrones := math.Ceil(maxDist/(drone.DroneRange*2))

}

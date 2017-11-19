package main

import (
	"drone"
	"fmt"
	"topo"
	"user"
)

type Field struct {
	//length and breadth of the field. Drones cannot be deployed outside this limit
	length      float64
	breadth     float64
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

	return &Field{length, breadth, nil, nil, 0, 0}
}

func (f *Field) AddDrone(d *drone.Drone) {
	droneGrp := drone.GetDroneGroup(getId())
	droneGrp.AddDrone(d)
	f.droneGroups = append(f.droneGroups, droneGrp)
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

func (f *Field) Serve (d *drone.Drone) {
	d.Serve()
	if d.GetGroup() != nil {
		return
	}
	droneGrp := drone.GetDroneGroup(getId())
	droneGrp.AddDrone(d)
	f.droneGroups = append(f.droneGroups, droneGrp)
}

func (f *Field) AddUser(u *user.User) {
	ug := user.GetUserGroup(0)
	ug.AddUser(u)
	f.userGroups = append(f.userGroups, ug)
	f.numOfUsers++
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
	if f.numOfDrones-1 < f.numOfUsers {
		numOfClusters = f.numOfDrones - 1
	} else {
		numOfClusters = f.numOfUsers
	}
	clusterPos := cluster(userPoints, int(numOfClusters))
	for i, n := range clusterPos {
		if users[i].GetGroup().GetClusterNumber() == uint64(n) {
			continue
		}
		if ug := f.GetUserGroup(uint64(n)); ug != nil {
			fmt.Println("cluster", n, "center before", ug.GetCenterofGroup())
			users[i].ChangeGroupTo(ug)
			fmt.Println("cluster", n, "center after", ug.GetCenterofGroup())
			continue
		}
		//find available cluster
		if ug := f.GetUserGroup(uint64(0)); ug != nil {
			ug.SetClusterNumber(uint64(n))
			fmt.Println("cluster", n, "center before", ug.GetCenterofGroup())
			users[i].ChangeGroupTo(ug)
			fmt.Println("cluster", n, "center after", ug.GetCenterofGroup())
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
	equiPoints := topo.GetEquiGeoCoordinates(f.numOfDrones-1, 5)
	f.droneGroups[0].GetDrones()[0].Serve() //anchor node
	f.droneGroups[0].GetDrones()[0].SetRole(drone.Anchor)
	for i := 1; i <= len(equiPoints); i++ {
		f.droneGroups[i].GetDrones()[0].MoveTo(equiPoints[i-1])
	}

	for _, ug := range f.userGroups {
		if len(ug.GetUsers()) == 0 {
			continue
		}
		d := f.GetClosestDrone(ug.GetCenterofGroup())
		f.DetachDrone(d)
		d.MoveTo(ug.GetCenterofGroup())
		f.Serve(d)
		fmt.Println("Moving drone: ", d.GetId(), "to cluster", ug.GetClusterNumber(), "drone pos", d.GetCurrentPosition())
	}
}

func (f *Field) GetClosestDrone(p topo.Point) *drone.Drone {
	var minDist float64
	var minDistD *drone.Drone
	init := false
	for _, dg := range f.droneGroups {
		for _, d := range dg.GetDrones() {
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
	}
	return minDistD
}

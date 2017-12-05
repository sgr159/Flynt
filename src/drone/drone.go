package drone

import (
	"topo"
	"fmt"
	"user"
)

const (
	DroneRange float64 = 40
)

type Role uint8

const (
	Unknown Role = iota
	Anchor 
	Edge
	Parent
)

type Status uint8

const (
	Free Status = iota
	Moving
	Serving
	Reserved
)

var XMax float64 = 0
var XMin float64 = 0
var YMax float64 = 0
var YMin float64 = 0

type DroneGroup struct {
	groupId uint64
	centerOfGroup topo.Point
	drones []*Drone
}

func GetDroneGroup(groupId uint64) *DroneGroup {
	return &DroneGroup{groupId,topo.Point{0,0},nil}
}

func (dg *DroneGroup) AddDrone(drone *Drone) {
	sumX := dg.centerOfGroup.X*float64(len(dg.drones)) +drone.GetCurrentPosition().X
	sumY := dg.centerOfGroup.Y*float64(len(dg.drones)) +drone.GetCurrentPosition().Y
	
	dg.centerOfGroup.X = sumX/float64(len(dg.drones)+1)
	dg.centerOfGroup.Y = sumY/float64(len(dg.drones)+1)
	
	dg.drones = append(dg.drones,drone)
	drone.droneGroup = dg
}


func (dg *DroneGroup) RemoveDrone(d *Drone) {
	if len(dg.drones) == 1 {
		dg.centerOfGroup.X, dg.centerOfGroup.Y = 0, 0
		dg.drones = dg.drones[:0]
	} else {
		var index int
		found := false

		for i, mem := range dg.drones {
			if mem == d {
				index = i
				found = true
				break
			}
		}
		if !found {
			fmt.Println("user", d.id, "you,re trying to delete is not found in dronegroup", dg.groupId)
			return
		}

		sumX := dg.centerOfGroup.X*float64(len(dg.drones)) - d.GetCurrentPosition().X
		sumY := dg.centerOfGroup.Y*float64(len(dg.drones)) - d.GetCurrentPosition().Y

		dg.centerOfGroup.X = sumX / float64(len(dg.drones)-1)
		dg.centerOfGroup.Y = sumY / float64(len(dg.drones)-1)

		dg.drones[index] = dg.drones[len(dg.drones)-1]
		dg.drones[len(dg.drones)-1] = nil //prevent memleak
		dg.drones = dg.drones[:len(dg.drones)-1]
	}
	d.droneGroup = nil
}

func (dg *DroneGroup) GetId() uint64 {
	return dg.groupId
}

func (dg *DroneGroup) GetDrones() []*Drone {
	return dg.drones
}

func (dg *DroneGroup) GetCenterOfGroup() topo.Point {
	return dg.centerOfGroup
}

func (dg *DroneGroup) GetDronePoints() []topo.Point {
	var dronePoints []topo.Point
	for _,d := range dg.drones {
		dronePoints = append(dronePoints,d.GetCurrentPosition())
	}
	return dronePoints
}

type Drone struct {
	id uint64
	currentPosition topo.Point
	destinationPosition topo.Point
	capacity uint64
	signalRange uint64
	droneGroup *DroneGroup
	role Role
	status Status
	parent *Drone
	users []*user.User
}

func GetDrone(id uint64, capacity uint64, signalRange uint64) *Drone {
	return &Drone {id,
		topo.Point{0,0},
		topo.Point{0,0},
		capacity,
		signalRange,
		nil,
		Unknown,
		Free,
		nil,
		nil}
}

func (d *Drone) GetCurrentPosition() topo.Point {
	return d.currentPosition
}

func (d *Drone) GetStatus() Status {
	return d.status
}

func (d *Drone) GetGroup() *DroneGroup {
	return d.droneGroup
}

func (d *Drone) GetId() uint64 {
	return d.id
}

func (d *Drone) GetUsers() []*user.User {
	return d.users
}

func (d *Drone) Serve (users []*user.User) {
	d.status = Serving
	d.users = append(d.users,users...)
}

func (d *Drone) SetStatus (s Status) {
	d.status = s
}

func (d *Drone) SetRole (r Role) {
	d.role = r
}

func (d *Drone) SetParent (p *Drone) {
	d.parent = p
}

func (d *Drone) Reserve () {
	d.status = Reserved
}

func (d *Drone) GetParent() *Drone {
	return d.parent
}

func (d *Drone) MoveTo(p topo.Point) bool {
	/*
	if !isInLimit(p) {
		return false
	}
	*/
	d.currentPosition.X = p.X
	d.currentPosition.Y = p.Y
	return true
}

func isInLimit(p topo.Point) bool {
	if XMax == 0 && XMin == 0 && YMax == 0 && YMin == 0 {
		return true
	}
	if XMax > p.X && XMin < p.X && YMax > p.Y && YMin < p.Y {
		return true
	}
	return false
}


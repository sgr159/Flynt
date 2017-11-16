package drone

import (
	"topo"
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
		nil}
}

func (d *Drone) GetCurrentPosition() topo.Point {
	return d.currentPosition
}

func (d *Drone) GetStatus() Status {
	return d.status
}

func (d *Drone) GetId() uint64 {
	return d.id
}

func (d *Drone) Serve () {
	d.status = Serving
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


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
	Uninitialized Status = iota
	Moving
	Serving
)

type DroneGroup struct {
	groupId uint64
	drones []*Drone
}

func GetDroneGroup(groupId uint64) *DroneGroup {
	return &DroneGroup{groupId,make([]*Drone,5)}
}

func (dg *DroneGroup) AddDrone(drone *Drone) {
	dg.drones = append(dg.drones,drone)
	drone.groupId = dg.groupId
}

type Drone struct {
	name string
	currentPosition topo.Point
	destinationPosition topo.Point
	capacity uint64
	signalRange uint64
	groupId uint64
	role Role
	status Status
	parent *Drone
}

func GetDrone(name string, currentPosition topo.Point, capacity uint64, signalRange uint64) *Drone {
	return &Drone {name,
		currentPosition,
		currentPosition,
		capacity,
		signalRange,
		0,
		Unknown,
		Uninitialized,
		nil}
}

func (d *Drone) GetCurrentPosition() topo.Point {
	return d.currentPosition
}


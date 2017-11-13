package user

import (
	"topo"
)

type User struct {
	id string
	position topo.Point
	droneId uint64 
	droneGroupId uint64
}
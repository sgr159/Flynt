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

func GetNewUser(id string, position topo.Point) *User {
	return &User{id, position,0,0}
}

func (u *User) GetCurrentPosition() topo.Point {
	return u.position
}
package user

import (
	"drone"
	"fmt"
	"topo"
)

type User struct {
	id       uint64
	position topo.Point
	drone    *drone.Drone
	group    *UserGroup
}

func GetNewUser(id uint64, position topo.Point) *User {
	return &User{id, position, nil, nil}
}

func (u *User) GetCurrentPosition() topo.Point {
	return u.position
}

func (u *User) GetGroup() *UserGroup {
	return u.group
}

func (u *User) ChangeGroupTo(ug *UserGroup) {
	u.group.RemoveUser(u)
	ug.AddUser(u)
}

type UserGroup struct {
	cluster       uint64
	users         []*User
	droneGroupId  uint64
	centerOfGroup topo.Point
}

func GetUserGroup(cluster uint64) *UserGroup {
	return &UserGroup{cluster, nil, 0, topo.Point{0, 0}}
}

func (ug *UserGroup) GetClusterNumber() uint64 {
	return ug.cluster
}

func (ug *UserGroup) GetCenterofGroup() topo.Point {
	return ug.centerOfGroup
}

func (ug *UserGroup) SetClusterNumber(num uint64) {
	ug.cluster = num
}

func (ug *UserGroup) AddUser(u *User) {
	sumX := ug.centerOfGroup.X*float64(len(ug.users)) + u.GetCurrentPosition().X
	sumY := ug.centerOfGroup.Y*float64(len(ug.users)) + u.GetCurrentPosition().Y

	ug.centerOfGroup.X = sumX / float64(len(ug.users)+1)
	ug.centerOfGroup.Y = sumY / float64(len(ug.users)+1)

	ug.users = append(ug.users, u)
	u.group = ug
}

func (ug *UserGroup) RemoveUser(u *User) {
	if len(ug.users) == 1 {
		ug.centerOfGroup.X, ug.centerOfGroup.Y = 0, 0
		ug.users = ug.users[:0]
	} else {
		var index int
		found := false

		for i, mem := range ug.users {
			if mem == u {
				index = i
				found = true
				break
			}
		}
		if !found {
			fmt.Println("user", u.id, "you,re trying to delete is not found in usergroup", ug.cluster)
			return
		}

		sumX := ug.centerOfGroup.X*float64(len(ug.users)) - u.GetCurrentPosition().X
		sumY := ug.centerOfGroup.Y*float64(len(ug.users)) - u.GetCurrentPosition().Y

		ug.centerOfGroup.X = sumX / float64(len(ug.users)-1)
		ug.centerOfGroup.Y = sumY / float64(len(ug.users)-1)

		ug.users[index] = ug.users[len(ug.users)-1]
		ug.users[len(ug.users)-1] = nil //prevent memleak
		ug.users = ug.users[:len(ug.users)-1]
	}
	u.group = nil
}

func (ug *UserGroup) GetUsers() []*User {
	return ug.users
}

func (ug *UserGroup) GetUserPoints() []topo.Point {
	var userPoints []topo.Point
	for _, u := range ug.users {
		userPoints = append(userPoints, u.GetCurrentPosition())
	}
	return userPoints
}

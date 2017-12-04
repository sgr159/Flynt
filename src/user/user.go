package user

import (
	"fmt"
	"topo"
)

type User struct {
	id       uint64
	position topo.Point
	group    *UserGroup
}

func GetNewUser(id uint64, position topo.Point) *User {
	return &User{id, position, nil}
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


func (ug *UserGroup) GetMostDistantUsers() (*User,*User,float64) {
	var maxXu,maxYu,minXu,minYu *User
	var maxX,maxY,minX,minY,maxDist float64 = 0,0,-1,-1,0
	for _,u := range ug.GetUsers() {
		if u.GetCurrentPosition().X > maxX {
			maxX = u.GetCurrentPosition().X
			maxXu = u
		}
		if u.GetCurrentPosition().X < minX || minX == -1 {
			minX = u.GetCurrentPosition().X
			minXu = u
		}
		if u.GetCurrentPosition().Y > maxY {
			maxY = u.GetCurrentPosition().Y
			maxYu = u
		}
		if u.GetCurrentPosition().Y < minY || minY == -1 {
			minY = u.GetCurrentPosition().Y
			minYu = u
		}
	}
	var borderUsers []*User
	var maxdistu1,maxdistu2 *User
	borderUsers = append(borderUsers,maxXu,maxYu,minXu,minYu)
	for _,u1 := range borderUsers {
		for _,u2 := range borderUsers {
			if u1 == u2 {
				continue;
			}
			dist := u1.GetCurrentPosition().DistanceFrom(u2.GetCurrentPosition())
			if maxDist < dist {
				maxDist = dist
				maxdistu1 = u1
				maxdistu2 = u2
			}
		}
	}
	return maxdistu1,maxdistu2,maxDist
}

func (ug *UserGroup) GetMostDistantUserFrom(p topo.Point) (*User,float64) {
	var maxDist float64 = -1
	var maxDistUser *User
	
	for _,u := range ug.GetUsers() {
		if maxDist < u.GetCurrentPosition().DistanceFrom(p) {
			maxDist = u.GetCurrentPosition().DistanceFrom(p)
			maxDistUser = u
		}
	}
	
	return maxDistUser,maxDist
}
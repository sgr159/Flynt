package drone

import (
    "testing"
    "topo"
    "fmt"
)


func TestCenterOfGroup(t *testing.T) {
	drone := GetDrone(1, 1, 1)
	drone.MoveTo(topo.Point{1,1})
	dg := GetDroneGroup(2)
	dg.AddDrone(drone)
	if !dg.centerOfGroup.IsEqual(drone.GetCurrentPosition()) {
		t.Error(
		        "For single drone case, expected the center of group to be the same as drone position, observed:",dg.centerOfGroup)
	}
	drone2 := GetDrone(0, 1, 1)
	drone2.MoveTo(topo.Point{-1,-1})
	dg.AddDrone(drone2)
	if !dg.centerOfGroup.IsEqual(topo.Point{0,0}) {
		t.Error(
		        "expected center of position to be 0,0, observed:",dg.centerOfGroup)
	}
}


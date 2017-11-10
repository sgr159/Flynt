package topo

import (
    "testing"
    "math"
    "fmt"
)

func TestArea(t *testing.T) {
	var index = Point{0,0}
	var r = Rectangle{index,1,1}
	var res = r.Area()
	if res != 1 {
		t.Errorf("SGRDBG area: %f",res)
	}
}

type testDistancePair struct {
	p1 Point
	p2 Point
	distance float64
}

func TestDistance(t *testing.T) {
	var points = []testDistancePair{
		{p1:Point{0,0},p2:Point{0,1},distance:1},
		{p1:Point{1,0},p2:Point{0,0},distance:1},
	}
	for _,pair := range points {
		v := pair.p1.DistanceFrom(pair.p2)
		if pair.distance != v {
			t.Error(
		        "For", pair.p1, "and",pair.p2,
		        "expected", pair.distance,
		        "got", v,
	      )
		}
	}
}

type testRotatePair struct {
	p1,p2,result Point
	angle Radians
}

func TestRotateLeft(t *testing.T) {
	var points = []testRotatePair{
		{p1:Point{0,0},p2:Point{1,0},angle:Radians(math.Pi/2),result:Point{0,1}},
		{p1:Point{1,1},p2:Point{2,1},angle:Radians(math.Pi/2),result:Point{1,2}},
	}
	for _,pair := range points {
		ans := rotateLeft(pair.p1, pair.p2, pair.angle)
		if !pair.result.IsEqual(ans) {
			t.Error(
		        "For", pair.p1, "and",pair.p2,
		        "expected", pair.result,
		        "got", ans,
	      )
		}
	}
}

type testEquiCoordSet struct {
	numOfVertices uint64
	distanceFromOrigin float64
	results []Point
}

func TestEquilateralCoordinates(t *testing.T) {
	var cases = []testEquiCoordSet{
		{numOfVertices:4,distanceFromOrigin:1,results:[]Point{Point{1,0},Point{0,1},Point{-1,0},Point{0,-1}}},
	}
	for _,set := range cases {
		ans := GetEquiGeoCoordinates(set.numOfVertices, set.distanceFromOrigin)
		if len(ans) != len(set.results){
			t.Error(
		        "For case", set,
		        "got", ans,
			)
			return;
		}
		for i := range ans{
			if !ans[i].IsEqual(set.results[i]) {
				t.Error(
			        "For case", set,
			        "got", ans,
				)
			}
		}
	}
}



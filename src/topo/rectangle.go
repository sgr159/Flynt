package topo

import (

)

/* Rectangle Strct */

type Rectangle struct {
	index  Point
	height float64
	width  float64
}

func GetRectangle(p Point, h float64, w float64) Rectangle{
	return Rectangle{p,h,w}
}

func (r Rectangle) Area() float64 {
	return r.height * r.width
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.height + r.width)
}

func (r Rectangle) IsInside(p Point) bool {
	if p.X < r.index.X && p.X > r.height+r.index.X {
		return false
	}
	if p.Y < r.index.Y && p.Y > r.width+r.index.Y {
		return false
	}
	return true
}
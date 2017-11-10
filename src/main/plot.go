package main

import (
	"fmt"
	"math/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"topo"
)

func plotsgr(points []topo.Point) {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Drone Positions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err = plotutil.AddScatters(p,
		"Drones", ptsFromPoints(points))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
	fmt.Println("err yo:",err)
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

// randomPoints returns some random x, y points.
func ptsFromPoints(points []topo.Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := range pts {
		pts[i].X = points[i].X
		pts[i].Y = points[i].Y
	}
	return pts
}
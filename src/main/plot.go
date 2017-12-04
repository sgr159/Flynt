package main

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"math/rand"
	"topo"
)

func plotMap(points, users []topo.Point) {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Drone Positions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.Add(plotter.NewGrid())

	err = plotutil.AddScatters(p,
		"Drones", ptsFromPoints(points),
		"Users", ptsFromPoints(users))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func (f *Field) PlotField2() {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Drone Positions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.Add(plotter.NewGrid())

	for _, dg := range f.droneGroups {
		for _, d := range dg.GetDrones() {
			if d.GetParent() != nil {
				var linePointsData []topo.Point
				linePointsData = append(linePointsData, d.GetCurrentPosition(), d.GetParent().GetCurrentPosition())
				// Make a line plotter and set its style.
				l, err := plotter.NewLine(ptsFromPoints(linePointsData))
				if err != nil {
					panic(err)
				}
				l.LineStyle.Width = vg.Points(1)
				l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
				l.LineStyle.Color = color.RGBA{R: 255, A: 255}
				p.Add(l)
			}
			for _, u := range d.GetUsers() {
				var linePointsData []topo.Point
				linePointsData = append(linePointsData, d.GetCurrentPosition(), u.GetCurrentPosition())
				// Make a line plotter and set its style.
				l, err := plotter.NewLine(ptsFromPoints(linePointsData))
				if err != nil {
					panic(err)
				}
				l.LineStyle.Width = vg.Points(1)
				//		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
				l.LineStyle.Color = color.RGBA{B: 255, A: 255}
				p.Add(l)

			}
		}
	}
	dronePoints := f.GetDronePositions()
	_, userPoints := f.GetUsersAndPositions()

	d, err := plotter.NewScatter(ptsFromPoints(dronePoints))
	if err != nil {
		panic(err)
	}
	d.Color = color.RGBA{R: 255, A: 255}
	d.Shape = draw.PyramidGlyph{}
	
	p.Add(d)
	
	u, err := plotter.NewScatter(ptsFromPoints(userPoints))
	if err != nil {
		panic(err)
	}
	u.Color = color.RGBA{G: 255, A: 255}
	u.Shape = draw.CircleGlyph{}
	
	p.Add(u)
	
	/*
	err = plotutil.AddScatters(p,
		"Drones", ptsFromPoints(dronePoints),
		"Users", ptsFromPoints(userPoints))
	if err != nil {
		panic(err)
	}
	*/

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points2.png"); err != nil {
		panic(err)
	}
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

func ptsFromPoints(points []topo.Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := range pts {
		pts[i].X = points[i].X
		pts[i].Y = points[i].Y
	}
	return pts
}

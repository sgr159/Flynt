package main

import (
	"github.com/mpraski/clusters"
	"topo"
)

func cluster(points []topo.Point, numOfClusters int) []int {
	var data [][]float64
//	var observation []float64
	
	data = make([][]float64,len(points))
	for i,p := range points {
		data[i] = []float64{p.X,p.Y}
	}
	
	// Create a new KMeans++ clusterer with 1000 iterations,
	// 3 clusters and a distance measurement function of type func([]float64, []float64) float64).
	// Pass nil to use clusters.EuclideanDistance
	c, e := clusters.KMeans(1000, numOfClusters, clusters.EuclideanDistance)
	if e != nil {
		panic(e)
	}
	
	// Use the data to train the clusterer
	if e = c.Learn(data); e != nil {
		panic(e)
	}
	
	/*
	fmt.Printf("Clustered data set into %d\n", c.Sizes())

	fmt.Printf("Assigned observation %v to cluster %d\n", observation, c.Predict(observation))

	for index, number := range c.Guesses() {
		fmt.Printf("Assigned data point %v to cluster %d\n", data[index], number)
	}
	*/
	
	return c.Guesses()
}

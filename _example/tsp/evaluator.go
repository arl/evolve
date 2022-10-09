package main

import (
	"math"
)

type routeEvaluator struct {
	cities []point
	dists  [][]float64
}

func newRouteEvaluator(cities []point) *routeEvaluator {
	// Create a matrix to hold distances between every city pair.
	ncities := len(cities)

	// TODO(arl): see if a 1D array could help
	dists := make([][]float64, ncities)
	for i := 0; i < ncities; i++ {
		dists[i] = make([]float64, ncities)
	}

	// Compute all distances.
	for i := 0; i < ncities; i++ {
		for j := 0; j < ncities; j++ {
			hypot := math.Hypot(cities[i].X-cities[j].X, cities[i].Y-cities[j].Y)
			dists[i][j] = hypot
			dists[j][i] = hypot
		}
	}

	return &routeEvaluator{
		dists: dists,
	}
}

func (e *routeEvaluator) Fitness(ind []int, pop [][]int) float64 {
	// Compute perimeter of the polygon formed by the closed path.
	var tot float64
	for i := 0; i < len(ind)-1; i++ {
		tot += e.dists[ind[i]][ind[i+1]]
	}
	tot += e.dists[0][len(ind)-1]

	return tot
}

func (e *routeEvaluator) IsNatural() bool {
	// TSP optimizes for shorter total route.
	return false
}

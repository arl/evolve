package main

import (
	"math"

	"evolve/example/tsp/internal/tsp"
)

type routeEvaluator struct {
	dists [][]float64
}

func newRouteEvaluator(cities []tsp.Point2D) *routeEvaluator {
	// Store precomputed distances in a 2D array.
	ncities := len(cities)

	dists := make([][]float64, ncities)
	for i := 0; i < ncities; i++ {
		dists[i] = make([]float64, ncities)
	}

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

// Fitness computes the perimeter of the polygon formed by the closed path.
func (e *routeEvaluator) Fitness(ind []int, pop [][]int) float64 {
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

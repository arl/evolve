package tsp

import (
	"math"
)

// SymmetricEvaluator evaluates candidates for the symmetric, unweighted
// traveling sales person. Symmetric means the distance between two cities is
// the same in each opposite direction.
type SymmetricEvaluator struct {
	// Distances holds the distances between every city pair.
	Distances [][]float64
}

// NewSymmetricEvaluator creates and initializes a symmetric TSP evaluator for
// the given list of citites. Euclidian distance between each city is
// precomputed, taking up len(cities)² space.
func NewSymmetricEvaluator(cities []Point2D) *SymmetricEvaluator {
	dists := make([][]float64, len(cities))
	for i := 0; i < len(cities); i++ {
		dists[i] = make([]float64, len(cities))
	}

	for i := range cities {
		for j := range cities {
			hypot := math.Hypot(cities[i].X-cities[j].X, cities[i].Y-cities[j].Y)
			dists[i][j] = hypot
			dists[j][i] = hypot
		}
	}

	return &SymmetricEvaluator{
		Distances: dists,
	}
}

// Fitness computes the perimeter of the polygon formed by the closed path
// passing through all all the cities in the order given by the candidate.
func (e *SymmetricEvaluator) Fitness(cand []int) float64 {
	var tot float64
	for i := 0; i < len(cand)-1; i++ {
		tot += e.Distances[cand[i]][cand[i+1]]
	}
	tot += e.Distances[0][len(cand)-1]

	return tot
}

func (e *SymmetricEvaluator) IsNatural() bool {
	// TSP optimizes for the shortest route.
	return false
}

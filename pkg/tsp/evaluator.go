package tsp

import (
	"math"

	"golang.org/x/exp/constraints"
)

// SymmetricTSPEvaluator evaluates candidates for the symmetrica and unweighted
// traveling sales person. Symmetric means the distance between two cities is
// the same in each opposite direction.
type SymmetricTSPEvaluator[T constraints.Integer] struct {
	// Distances holds the distances between every city pair.
	Distances [][]float64
}

// NewSymmetricEvaluator creates and initializes a symmetric TSP evaluator for
// the given list of citites. Euclidian distance between each city is
// precomputed, taking up len(cities)² space.
func NewSymmetricEvaluator[T constraints.Integer](cities []Point2D) *SymmetricTSPEvaluator[T] {
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

	return &SymmetricTSPEvaluator[T]{
		Distances: dists,
	}
}

// Fitness computes the perimeter of the polygon formed by the closed path
// passing through all all the cities in the order given by the candidate.
func (e *SymmetricTSPEvaluator[T]) Fitness(cand []T) float64 {
	var tot float64
	for i := 0; i < len(cand)-1; i++ {
		tot += e.Distances[cand[i]][cand[i+1]]
	}
	tot += e.Distances[0][len(cand)-1]

	return tot
}

func (e *SymmetricTSPEvaluator[T]) IsNatural() bool {
	// TSP optimizes for the shortest route.
	return false
}

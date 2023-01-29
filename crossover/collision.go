package crossover

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// The Collision crossover ...
type Collision[T constraints.Integer] struct {
	// Probability is the probability to apply the Collision crossover on a pair
	// of candidates.
	Probability generator.Float

	CalcDistance func(i, j int) float64
}

// paper: https://arxiv.org/pdf/1801.02335.pdf
func (op *Collision[T]) Apply(pop *evolve.Population[[]T], rng *rand.Rand) {
	// Shuffle candidates so that evolution is not biased by previous
	// operations.
	rand.Shuffle(pop.Len(), pop.Swap)

	for i := 0; i < pop.Len()-1; i += 2 {
		j := i + 1
		p1 := pop.Candidates[i]
		p2 := pop.Candidates[j]

		if rng.Float64() >= op.Probability.Next() {
			continue // Nothing to do.
		}

		if !pop.Evaluated[i] {
			// TODO(arl) this is too manual maybe? add some methods to Population for this.
			pop.Fitness[i] = pop.Evaluator.Fitness(p1)
			pop.Evaluated[i] = true
		}
		if !pop.Evaluated[j] {
			pop.Fitness[j] = pop.Evaluator.Fitness(p2)
			pop.Evaluated[j] = true
		}

		// Generate a random cut point from 2 to len(x1) -1
		cut := 2 + rng.Intn(len(p1)-3)

		// TODO(arl) cost need to be mapped to an int
		cost1 := int(pop.Fitness[i])
		cost2 := int(pop.Fitness[j])
		off1, off2 := op.collide(p1, p2, cost1, cost2, cut, rng)
		pop.Candidates[i] = off1
		pop.Candidates[j] = off2
	}
}

// collide performs the collision crossover
//
// This is based on the C# implementation of the Collision Crossover, courtesy
// of Prof. Ahmad Hassanat which kindly shared his unpublished code with me.
func (op *Collision[T]) collide(x1, x2 []T, cost1, cost2 int, cut int, rng *rand.Rand) (y1, y2 []T) {
	if len(x1) != len(x2) {
		panic("Collision cannot mate parents of different lengths")
	}

	tempCh1 := make([]T, 0, len(x1))
	tempCh2 := make([]T, 0, len(x1))

	// TODO(arl) uses a bitset (bitstring)
	visited1 := make([]bool, len(x1))
	visited2 := make([]bool, len(x1))

	// Collision

	tempCh1 = append(tempCh1, 0)
	visited1[0] = true
	tempCh2 = append(tempCh2, 0)
	visited2[0] = true

	var (
		v1, v2           float64
		m1, m2, v1p, v2p float64
	)

	// TODO(arl):
	//    compute cost based on the different in fitness of both candidates.

	v1 = float64(1 + rng.Intn(cost1-1))
	v2 = float64(1 + rng.Intn(cost2-1))
	v2 = -v2

	for j := 1; j < len(x1)-1; j++ {
		m1 = op.CalcDistance(int(x1[j-1]), int(x1[j])) + op.CalcDistance(int(x1[j]), int(x1[j+1]))
		m2 = op.CalcDistance(int(x2[j-1]), int(x2[j])) + op.CalcDistance(int(x2[j]), int(x2[j+1]))

		v1p = (v1 * (m1 - m2) / (m1 + m2)) + (v2 * 2.0 * m2 / (m1 + m2))
		v2p = (v1 * 2.0 * m1 / (m1 + m2)) - (v2 * (m1 - m2) / (m1 + m2))

		// Add to first offspring
		if v1p <= 0 {
			tempCh1 = append(tempCh1, x1[j])
			visited1[x1[j]] = true
		} else {
			tempCh1 = append(tempCh1, 0)
		}

		// Add to second offspring
		if v2p >= 0 {
			tempCh2 = append(tempCh2, x2[j])
			visited2[x2[j]] = true
		} else {
			tempCh2 = append(tempCh2, 0)
		}
	}

	// Do it for the last gene
	jj := len(x1) - 1

	m1 = op.CalcDistance(int(x1[jj-1]), int(x1[jj])) + op.CalcDistance(int(x1[jj]), int(x1[0]))
	m2 = op.CalcDistance(int(x2[jj-1]), int(x2[jj])) + op.CalcDistance(int(x2[jj]), int(x2[0]))

	v1p = (v1 * (m1 - m2) / (m1 + m2)) + (v2 * 2.0 * m2 / (m1 + m2))
	v2p = (v1 * 2.0 * m1 / (m1 + m2)) - (v2 * (m1 - m2) / (m1 + m2))

	// add to ch1

	if v1p <= 0 {
		// if the city refelected or stopped in parent1
		tempCh1 = append(tempCh1, x1[jj])
		visited1[x1[jj]] = true
	} else {
		tempCh1 = append(tempCh1, 0)
	}

	if v2p <= 0 {
		// if the city refelected or stopped in parent1
		tempCh2 = append(tempCh2, x2[jj])
		visited2[x2[jj]] = true
	} else {
		tempCh2 = append(tempCh2, 0)
	}

	// Fill the rest of genes the one which is not refelcted or stopped from the collision

	indx := slices.Index(tempCh1[1:], 0)
	indx++

	// while(indx){
	for j := 1; j < jj+1; j++ {
		if indx == 0 {
			j = jj + 1
			break
		}

		if !visited1[x2[j]] {
			tempCh1[indx] = x2[j] // fill the rest of cromosome1 from parent2

			visited1[x2[j]] = true
			indx = slices.Index(tempCh1[1:], 0)
			indx++
		} else {
			continue
		}

	}

	indx = slices.Index(tempCh2[1:], 0)
	indx++
	for j := 1; j < jj+1; j++ {

		if indx == 0 {
			j = jj + 1
			break
		}

		if !visited2[x1[j]] {
			tempCh2[indx] = x1[j] // fill the rest of cromosome2 from parent1
			visited2[x1[j]] = true

			indx = slices.Index(tempCh2[1:], 0)
			indx++
		} else {
			continue
		}
	}

	// TODO(arl) no need to recompute the fitness now

	// tempCh1.Cost = (ulong)CalcCostED(tempCh1, jj + 1);
	// tempCh2.Cost = (ulong)CalcCostED(tempCh2, jj + 1);

	return tempCh1, tempCh2 // output 2 children
}

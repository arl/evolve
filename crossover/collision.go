package crossover

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// The Collision crossover ...
// TODO(arl) document
// paper: https://arxiv.org/pdf/1801.02335.pdf
type Collision[T constraints.Integer] struct {
	// Probability is the probability to apply the Collision crossover on a pair
	// of candidates.
	Probability generator.Float

	EdgeWeight func(i, j int) float64
}

// seq returns a slice containing the sequence of consecutive numbers from 0 to n.
func seq[T constraints.Integer | constraints.Float](n int) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = T(i)
	}
	return s
}

func (op *Collision[T]) Apply(pop *evolve.Population[[]T], rng *rand.Rand) {
	indices := seq[int](pop.Len())
	rng.Shuffle(pop.Len(), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	for u := 0; u < pop.Len()-1; u += 2 {
		v := u + 1
		i, j := indices[u], indices[v]
		p1 := pop.Candidates[i]
		p2 := pop.Candidates[j]

		if rng.Float64() >= op.Probability.Next() {
			continue // Nothing to do.
		}

		if !pop.Evaluated[i] {
			// TODO(arl) this is too manual maybe? add some methods to Population for this.
			pop.Fitness[i] = pop.Evaluator.Fitness(p1)
		}
		if !pop.Evaluated[j] {
			pop.Fitness[j] = pop.Evaluator.Fitness(p2)
		}

		// TODO(arl) cost need to be mapped to an int
		cost1 := int(pop.Fitness[i])
		cost2 := int(pop.Fitness[j])
		off1, off2 := op.collide(p1, p2, cost1, cost2, rng)
		pop.Candidates[i] = off1
		pop.Candidates[j] = off2
	}
}

func (op *Collision[T]) collide(x1, x2 []T, cost1, cost2 int, rng *rand.Rand) (y1, y2 []T) {
	var zero T
	minusOne := zero - 1

	v1 := float64(1 + rng.Intn(cost1-1))
	v2 := -float64(1 + rng.Intn(cost2-1))

	off1, off2 := make([]T, len(x1)), make([]T, len(x2))

	vis1 := make([]bool, len(x1))
	vis2 := make([]bool, len(x2))

	for i := 0; i < len(x1); i++ {
		m1 := op.mass(x1, i)
		m2 := op.mass(x2, i)

		// we just want the sign of the velocity
		diff := m1 - m2
		v1p := (v1 * diff) + (v2 * 2.0 * m2)
		v2p := (v1 * 2.0 * m1) - (v2 * diff)

		if v1p <= 0 {
			// x1[i] is a good gene, keep it for off1.
			off1[i] = x1[i]
			vis1[x1[i]] = true
		} else {
			off1[i] = minusOne
		}

		if v2p <= 0 {
			// x2[i] is a good gene, keep it for off2.
			off2[i] = x2[i]
			vis2[x2[i]] = true
		} else {
			off2[i] = minusOne
		}
	}

	idx := slices.Index(off1, minusOne)
	for i := 0; i < len(off1); i++ {
		if idx == -1 {
			break
		}
		if !vis1[x2[i]] {
			off1[idx] = x2[i]
			vis1[x2[i]] = true
			idx = slices.Index(off1, minusOne)
		}
	}

	idx = slices.Index(off2, minusOne)
	for i := 0; i < len(off2); i++ {
		if idx == -1 {
			break
		}
		if !vis2[x1[i]] {
			off2[idx] = x1[i]
			vis2[x1[i]] = true
			idx = slices.Index(off2, minusOne)
		}
	}

	return off1, off2
}

// mass computes the mass of the city i in chromosome c.
func (op *Collision[T]) mass(c []T, i int) float64 {
	var prev, next int

	switch i {
	case 0:
		prev = len(c) - 1
		next = 1
	case len(c) - 1:
		prev = i - 1
		next = 0
	default:
		prev = i - 1
		next = i + 1
	}

	return op.EdgeWeight(int(c[prev]), int(c[i])) + op.EdgeWeight(int(c[i]), int(c[next]))
}

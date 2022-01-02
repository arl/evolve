package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestPMX(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New[[]int](PMX[int]{})
	xover.Points = generator.Const(2)
	xover.Probability = generator.Const(1.0)

	pop := [][]int{
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		[]int{3, 7, 5, 1, 6, 8, 2, 4},
	}

	// Perform multiple crossovers to check different crossover points.
	for i := 0; i < 50; i++ {
		pop = xover.Apply(pop, rng)

		for _, ind := range pop {
			for j := 1; j <= 8; j++ {
				assert.Containsf(t, ind, j, "offspring is missing element %d in slice ")
			}
		}
	}
}

func TestPMXDifferentLength(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New[[]int](PMX[int]{})
	xover.Points = generator.Const(2)

	pop := make([][]int, 2)
	pop[0] = []int{1, 2, 3, 4, 5, 6, 7, 8}
	pop[1] = []int{3, 7, 5, 1}

	assert.Panics(t, func() {
		xover.Apply(pop, rng)
	})
}

func TestPMX2CrossoverPoints(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New[[]int](PMX[int]{})
	xover.Points = generator.Const(3)

	pop := make([][]int, 2)
	pop[0] = []int{1, 2, 3, 4}
	pop[1] = []int{3, 7, 5, 1}

	assert.Panics(t, func() {
		xover.Apply(pop, rng)
	})
}

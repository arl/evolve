package operator

import (
	"math/rand"
	"testing"
)

// adjustInt mutates integers candidates by adding a fixed offset.
type adjustInt int

func (op adjustInt) Apply(cand []int, rng *rand.Rand) []int {
	result := make([]int, len(cand))
	for i, c := range cand {
		result[i] = c + int(op)
	}
	return result
}

func TestEvolutionPipeline(t *testing.T) {
	// Make sure that multiple operators in a pipeline are applied correctly
	// to the population and validate the cumulative effects.
	rng := rand.New(rand.NewSource(99))
	pop := make([]int, 0)
	for i := 10; i <= 100; i += 10 {
		pop = append(pop, i)
	}

	pipe := Pipeline[int]{adjustInt(1), adjustInt(3)}
	pop = pipe.Apply(pop, rng)

	// Net result should be each candidate increased by 4.
	var sum int
	for _, c := range pop {
		ic := c
		sum += ic
		if ic%10 != 4 {
			t.Error("candidate value should have increased by 4, got", ic)
		}
	}
	if sum != 590 {
		t.Error("want sum = 90, got", sum)
	}
}

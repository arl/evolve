package evolve

import (
	"math/rand"
	"testing"
)

func TestGeneratePopulation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	fac := FactoryFunc[int](
		func(rng *rand.Rand) int {
			return rng.Int()
		},
	)

	if pop := GeneratePopulation[int](fac, 10, rng); len(pop) != 10 {
		t.Errorf("len(pop) = %d, want %d", len(pop), 10)
	}
}

func occurrences(vals []int) map[int]int {
	m := make(map[int]int)
	for _, v := range vals {
		m[v]++
	}
	return m
}

func TestSeedPopulation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	// Factory that generates -1 (so we can check later whether a candidate is a
	// seed or not).
	fac := FactoryFunc[int](
		func(rng *rand.Rand) int {
			return -1
		},
	)

	t.Run("no seeds", func(t *testing.T) {
		var seeds []int
		pop := SeedPopulation[int](fac, 10, seeds, rng)
		if len(pop) != 10 {
			t.Errorf("len(pop) = %d, want %d", len(pop), 10)
		}
		m := occurrences(pop)
		if m[-1] != 10 {
			t.Errorf("got %d generated candidates (not seeded) in population, want %d", m[-1], 0)
		}
	})

	t.Run("2 seeds", func(t *testing.T) {
		seeds := []int{1, 2}
		pop := SeedPopulation[int](fac, 10, seeds, rng)
		if len(pop) != 10 {
			t.Errorf("len(pop) = %d, want %d", len(pop), 10)
		}
		m := occurrences(pop)
		if m[-1] != 8 {
			t.Errorf("got %d generated candidates (not seeded) in population, want %d", m[-1], 8)
		}
	})

	t.Run("only seeds", func(t *testing.T) {
		seeds := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		pop := SeedPopulation[int](fac, 10, seeds, rng)
		if len(pop) != 10 {
			t.Errorf("len(pop) = %d, want %d", len(pop), 10)
		}
		m := occurrences(pop)
		if m[-1] != 0 {
			t.Errorf("got %d generated candidates (not seeded) in population, want %d", m[-1], 0)
		}
	})

	t.Run("len(seeds)>n", func(t *testing.T) {
		seeds := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		pop := SeedPopulation[int](fac, 10, seeds, rng)
		if len(pop) != 10 {
			t.Errorf("len(pop) = %d, want %d", len(pop), 10)
		}
		m := occurrences(pop)
		if m[-1] != 0 {
			t.Errorf("got %d generated candidates (not seeded) in population, want %d", m[-1], 0)
		}
	})
}

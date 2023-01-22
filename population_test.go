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

	pop := GeneratePopulation[int](10, fac, nil, rng)
	if pop.Len() != 10 {
		t.Errorf("pop.Len() = %d, want %d", pop.Len(), 10)
	}
}

func TestSeedPopulation(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		seeds   []int
		wantGen int // number of generated candidates in the new population
	}{
		{
			name:    "n=0 and no seeds",
			n:       0,
			seeds:   nil,
			wantGen: 0,
		},
		{
			name:    "n=0 and empty seeds",
			n:       0,
			seeds:   []int{},
			wantGen: 0,
		},
		{
			name:    "n=0 and 2 seeds",
			n:       0,
			seeds:   []int{1, 2},
			wantGen: 0,
		},
		{
			name:    "n=10 and no seeds",
			n:       10,
			seeds:   []int{},
			wantGen: 10,
		},
		{
			name:    "n=10 and 3 seeds",
			n:       10,
			seeds:   []int{1, 2, 3},
			wantGen: 7,
		},
		{
			name:    "n=10 and 10 seeds",
			n:       10,
			seeds:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantGen: 0,
		},
		{
			name:    "n=5 and 10 seeds",
			n:       5,
			seeds:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantGen: 0,
		},
	}

	rng := rand.New(rand.NewSource(99))

	// Factory that generates -1 (so we can check later
	// whether a candidate is seeded or a generated).
	fac := FactoryFunc[int](func(rng *rand.Rand) int {
		return -1
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pop := SeedPopulation[int](tt.n, tt.seeds, fac, nil, rng)
			if pop.Len() != tt.n {
				t.Errorf("pop.Len() = %d, want %d", pop.Len(), 10)
			}

			generated := 0
			for _, cand := range pop.Candidates {
				if cand == -1 {
					generated++
				}
			}

			if generated != tt.wantGen {
				t.Errorf("got %d generated candidates in population, want %d", generated, tt.wantGen)
			}
		})
	}
}

package evolve_test

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

func sameStringPop(t *testing.T, a, b []string) {
	t.Helper()

	s1 := make([]string, 0)
	s2 := make([]string, 0)
	for _, a := range a {
		s1 = append(s1, string(a))
	}
	for _, b := range b {
		s2 = append(s2, string(b))
	}
	sort.Strings(s1)
	sort.Strings(s2)
	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("different strings\na = %v\nb = %v", s1, s2)
	}
}

func TestCrossoverProbabilityZero(t *testing.T) {
	org := []string{
		"abcde",
		"fghij",
		"klmno",
		"pqrst",
		"uvwxy",
	}

	xover := evolve.Crossover[string]{
		Mater:       nil,
		Probability: generator.Const(0.),
	}

	items := make([]string, len(org))
	copy(items, org)

	pop := evolve.NewPopulationOf(items, nil)

	rng := rand.New(rand.NewSource(99))
	xover.Apply(pop, rng)
	sameStringPop(t, pop.Candidates, org)
}

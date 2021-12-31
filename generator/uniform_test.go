package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"
)

func TestUniformInt(t *testing.T) {
	const (
		min, max = 150, 500
		rang     = max - min
	)

	const mean = (rang / 2) + min

	// For a uniformly distributed [0 n] range, standard deviation should be n/sqrt(12).
	stddev := float64(rang) / math.Sqrt(12)

	rng := rand.New(mt19937.New(0))
	g := NewUniformtInt[uint16](min, max, rng)

	checkUniformDistribution[uint16](t, g, mean, stddev)
}

func TestUniformAcrossZero(t *testing.T) {
	const (
		min, max = -150, 500
		rang     = max - min
	)

	const mean = (rang / 2) + min

	// For a uniformly distributed [0 n] range, standard deviation should be n/sqrt(12).
	stddev := float64(rang) / math.Sqrt(12)

	rng := rand.New(mt19937.New(0))
	g := NewUniformtInt[int16](min, max, rng)

	checkUniformDistribution[int16](t, g, mean, stddev)
}

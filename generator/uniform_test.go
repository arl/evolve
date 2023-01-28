package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"
	"golang.org/x/exp/constraints"
)

func TestUniformInt(t *testing.T) {
	testUniform[uint16](22, 122, t)
	testUniform[float64](100, 222, t)
	testUniform[uint32](22, 500, t)
	testUniform[byte](0, 255, t)
	testUniform[float64](100, 222, t)
	testUniform[float32](100, 222, t)
	testUniform[float32](0, math.MaxFloat32, t)
	testUniform[float64](0, math.MaxFloat32, t)
}

func testUniform[T constraints.Unsigned | constraints.Float](min, max T, t *testing.T) {
	diff := max - min
	mean := (diff / 2) + min

	// For a uniformly distributed [0 n] range, standard deviation should be
	// n/sqrt(12).
	stddev := float64(diff) / math.Sqrt(12)

	rng := rand.New(mt19937.New())
	g := Uniform(min, max, rng)

	checkUniformDistribution(t, g, float64(mean), stddev)
}

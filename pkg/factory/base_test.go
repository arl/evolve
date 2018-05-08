package factory

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intFactory struct{ BaseFactory }

func newIntFactory() *intFactory { return &intFactory{BaseFactory{intGenerator{}}} }

type intGenerator struct{}

func (intGenerator) GenerateCandidate(rng *rand.Rand) interface{} { return rng.Int() }

func TestAbstractCandidateFactoryPopulationCreation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	t.Run("generate whole population", func(t *testing.T) {
		cf := newIntFactory()
		pop := cf.GenPopulation(10, rng)
		assert.Len(t, pop, 10)
	})

	t.Run("seed initial population", func(t *testing.T) {
		cf := newIntFactory()

		// preseed 5 candidates over 10
		preseed := make([]interface{}, 5)
		for i := range preseed {
			preseed[i] = i
		}

		pop := cf.SeedPopulation(10, preseed, rng)
		assert.Len(t, pop, 10)
	})

	t.Run("too many seed candidates", func(t *testing.T) {
		cf := newIntFactory()

		// preseed 10 candidates
		preseed := make([]interface{}, 10)
		for i := range preseed {
			preseed[i] = i
		}

		// TODO: should return error instead of panicking
		assert.Panics(t, func() {
			cf.SeedPopulation(5, preseed, rng)
		})
	})
}

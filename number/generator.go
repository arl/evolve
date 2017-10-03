package number

import (
	"math"
	"math/rand"
	"time"
)

// IntegerGenerator is the interface implemented by objects providing sequence
// of integers.
type IntegerGenerator interface {

	// NextValue returns the next value from the generator.
	NextValue() int64
}

// ConstantIntegerGenerator is an integer generator that always returns the same
// value.
type ConstantIntegerGenerator struct {
	constant int64
}

// NewConstantIntegerGenerator returns an IntegerGenerator that always returns
// the provided constant.
func NewConstantIntegerGenerator(constant int64) IntegerGenerator {
	return ConstantIntegerGenerator{constant: constant}
}

// NextValue returns the next value from the generator.
func (g ConstantIntegerGenerator) NextValue() int64 {
	return g.constant
}

// Float64Generator is the interface implemented by objects providing sequence
// of float64
type Float64Generator interface {

	// NextValue returns the next value from the generator.
	NextValue() float64
}

// ConstantFloat64Generator is an integer generator that always returns the same
// value.
type ConstantFloat64Generator struct {
	constant float64
}

// NewConstantFloat64Generator returns an Float64Generator that always returns
// the provided constant.
func NewConstantFloat64Generator(constant float64) Float64Generator {
	return ConstantFloat64Generator{constant: constant}
}

// NextValue returns the next value from the generator.
func (g ConstantFloat64Generator) NextValue() float64 {
	return g.constant
}

// BoundedFloat64Generator produces float64 bounded in a given range
type BoundedFloat64Generator struct {
	min, max float64    // both min and max are exclusive
	rng      *rand.Rand // PRNG used to generate the values
}

// NewBoundedFloat64Generator returns a BoundedFloat64Generator that uses a
// default uniform PRNG.
func NewBoundedFloat64Generator(min, max float64) *BoundedFloat64Generator {
	g := &BoundedFloat64Generator{min: min, max: max}
	g.SetRNG(rand.New(rand.NewSource(time.Now().UnixNano())))
	return g
}

// SetRNG sets the pseudo random number generator used to produce the values.
func (g *BoundedFloat64Generator) SetRNG(rng *rand.Rand) {
	g.rng = rng
}

// NextValue returns the next value from the generator.
func (g *BoundedFloat64Generator) NextValue() float64 {
	return g.rng.Float64() + math.SmallestNonzeroFloat64
}

// ProbabilityGenerator is the interface implemented by objects providing sequence
// of probabilities (floating points numbers between 0 and 1).
type ProbabilityGenerator interface {

	// NextValue returns the next value from the generator.
	NextValue() Probability
}

// ConstantProbabilityGenerator is a probability generator that always returns
// the same value.
type ConstantProbabilityGenerator struct {
	constant Probability
}

// NewConstantProbabilityGenerator returns a ProbabilityGenerator that always
// returns the provided constant.
func NewConstantProbabilityGenerator(constant Probability) ProbabilityGenerator {
	return ConstantProbabilityGenerator{constant: constant}
}

// NextValue returns the next value from the generator.
func (g ConstantProbabilityGenerator) NextValue() Probability {
	return g.constant
}

package number

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

package number

// IntGenerator is the interface implemented by objects providing sequence of
// integers.
type IntegerGenerator interface {
	// NextValue returns the next value from the generator.
	NextValue() int64
}

// IntGenerator is the interface implemented by objects providing sequence of
// probabilities, floating points numbers between 0 and 1.
type ProbabilityGenerator interface {
	// NextValue returns the next value from the generator.
	NextValue() Probability
}

type ConstantIntegerGenerator struct {
	constant int64
}

func NewConstantIntegerGenerator(constant int64) ConstantIntegerGenerator {
	return ConstantIntegerGenerator{constant: constant}
}

func (g ConstantIntegerGenerator) NextValue() int64 {
	return g.constant
}

type ConstantProbabilityGenerator struct {
	constant Probability
}

func NewConstantProbabilityGenerator(constant Probability) ConstantProbabilityGenerator {
	return ConstantProbabilityGenerator{constant: constant}
}

func (g ConstantProbabilityGenerator) NextValue() Probability {
	return g.constant
}

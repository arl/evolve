package operators

import (
	"errors"
	"fmt"

	"github.com/aurelien-rainone/evolve/number"
)

// Option is the interface implemented by objects configuring a
// particular option of a genetic operator.
type Option interface {
	Apply(interface{}) error
}

// ConstantProbability configures a constant probability that a genetic operator
// applies.
func ConstantProbability(prob number.Probability) probabilityGeneratorOption {
	return probabilityGeneratorOption{
		gen: number.NewConstantProbabilityGenerator(prob),
	}
}

// VariableProbability configures, via a number.ProbabilityGenerator, the
// probability that a genetic operator applies to an individual.
func VariableProbability(gen number.ProbabilityGenerator) probabilityGeneratorOption {
	return probabilityGeneratorOption{
		gen: gen,
	}
}

type probabilityGeneratorOption struct {
	gen number.ProbabilityGenerator
	err error
}

func (opt probabilityGeneratorOption) Apply(ope interface{}) error {
	switch ope.(type) {
	case *AbstractCrossover:
		if opt.err == nil {
			crossover := ope.(*AbstractCrossover)
			crossover.crossoverProbability = opt.gen
		}
		return opt.err
	case *AbstractMutation:
		if opt.err == nil {
			mutation := ope.(*AbstractMutation)
			mutation.mutationProbability = opt.gen
		}
		return opt.err
	}
	return fmt.Errorf("can't apply option to object of type %T", ope)
}

// ConstantCrossoverPoints configures a constant number of crossover points.
//
// This option only applies to crossover operators.
func ConstantCrossoverPoints(points int64) integerGeneratorOption {
	var err error
	if points <= 0 {
		err = errors.New("number of crossover points must be positive")
	} else {
		err = nil
	}
	return integerGeneratorOption{
		gen: number.NewConstantIntegerGenerator(points),
		err: err,
	}
}

// VariableCrossoverPoints configures, via a number.IntegerGenerator, a
// crossover such as the number of crossover points varies.
//
// This option only applies to crossover operators.
func VariableCrossoverPoints(gen number.IntegerGenerator) integerGeneratorOption {
	return integerGeneratorOption{
		gen: gen,
	}
}

type integerGeneratorOption struct {
	gen number.IntegerGenerator
	err error
}

func (opt integerGeneratorOption) Apply(ope interface{}) error {
	switch ope.(type) {

	case *AbstractCrossover:
		if opt.err == nil {
			crossover := ope.(*AbstractCrossover)
			crossover.crossoverPoints = opt.gen
		}
		return opt.err

	case *AbstractMutation:
		mutation := ope.(*AbstractMutation)

		if opt.err == nil {

			switch mutation.Mutater.(type) {
			case *bitStringMutater:
				mutation.Mutater.(*bitStringMutater).mutationCount = opt.gen
			}
		}
		return opt.err
	}
	return fmt.Errorf("can't apply option to object of type %T", ope)
}

// ConstantMutationCount configures a constant number for the number of
// mutations in a candidate selected for mutation.
//
// This option only applies to some mutation operators.
func ConstantMutationCount(points int64) integerGeneratorOption {
	var err error
	if points <= 0 {
		err = errors.New("number of mutation count must be positive")
	} else {
		err = nil
	}
	return integerGeneratorOption{
		gen: number.NewConstantIntegerGenerator(points),
		err: err,
	}
}

// VariableMutationCount configures, via a number.IntegerGenerator, a
// mutation such as the number of mutations varies in a candidate selected for
// mutation.
//
// This option only applies to some mutation operators.
func VariableMutationCount(gen number.IntegerGenerator) integerGeneratorOption {
	return integerGeneratorOption{
		gen: gen,
	}
}

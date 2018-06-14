package evolve

// Evaluator calculates the fitness score of a given candidate of the
// appropriate type.
//
// Fitness evaluations may be executed concurrently and therefore any concurrent
// access to a shared state should be properly synchronized.
type Evaluator interface {

	// Fitness calculates a fitness score for the given candidate.
	//
	// Whether a higher score indicates a fitter candidate or not depends on
	// whether the fitness scores are natural. This method must always return a
	// value greater than or equal to zero. Framework behaviour is undefined for
	// negative fitness scores.
	//
	// cand is the candidate solution to calculate fitness for.
	//
	// pop is the entire population. This will include the specified candidate.
	// This is provided for fitness evaluators that evaluate individuals in the
	// context of the population that they are part of (e.g.  a program that
	// evolves game-playing strategies may wish to play each strategy against
	// each of the others). This parameter can be ignored by simple fitness
	// evaluators. When iterating over the population, a simple interface
	// equality check (==) can be used to identify which member of the
	// population is the specified candidate.
	//
	// Returns the fitness score for the specified candidate. Must always be a
	// non-negative value regardless of natural or non-natural evaluation is
	// being used.
	Fitness(cand interface{}, pop []interface{}) float64

	// IsNatural specifies whether this evaluator generates 'natural' fitness
	// scores or not.
	//
	// Natural fitness scores are those in which the fittest individual in a
	// population has the highest fitness value. In this case the algorithm is
	// attempting to maximise fitness scores.  There need not be a specified
	// maximum possible value.
	// In contrast, 'non-natural' fitness evaluation results in fitter
	// individuals being assigned lower scores than weaker individuals.  In the
	// case of non-natural fitness, the algorithm is attempting to minimise
	// fitness scores.
	//
	// An example of a situation in which non-natural fitness scores are
	// preferable is when the fitness corresponds to a cost and the algorithm is
	// attempting to minimise that cost.
	//
	// The terminology of 'natural' and 'non-natural' fitness scores is
	// introduced by the evolve Framework to describe the two types of fitness
	// scoring that exist within the framework. It does not correspond to either
	// standardised fitness or normalised fitness in the EA literature.
	// Standardised fitness evaluation generates non-natural
	// scores with a score of zero corresponding to the best possible fitness.
	// Normalised fitness evaluation is similar to standardised fitness but with
	// the scores adjusted to fall within the range 0 - 1.
	//
	// Returns true if a high fitness score means a fitter candidate or false if
	// a low fitness score means a fitter candidate.
	IsNatural() bool
}

// FitnessFunc is the type of function computing the fitness of a candidate
// solution.
type FitnessFunc func(interface{}, []interface{}) float64

type evaluatorFunc struct {
	f FitnessFunc
	n bool
}

func (e evaluatorFunc) Fitness(cand interface{}, pop []interface{}) float64 { return e.f(cand, pop) }
func (e evaluatorFunc) IsNatural() bool                                     { return e.n }

// EvaluatorFunc is an adapter to allow the use of ordinary functions as fitness
// evaluators. If f is a function with the appropriate signature, EvaluatorFunc
// returns an object satisfying the Evaluator interface, for which the Fitness
// method calls f and IsNatural returns natural.
func EvaluatorFunc(natural bool, f FitnessFunc) evaluatorFunc { // nolint: golint
	return evaluatorFunc{f: f, n: natural}
}

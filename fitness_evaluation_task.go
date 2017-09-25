package evolve

import (
	"fmt"

	"github.com/aurelien-rainone/evolve/framework"
)

// fitnessEvaluationTask is a task for performing parallel fitness evaluations.
type fitnessEvaluationTask struct {
	fitnessEvaluator framework.FitnessEvaluator
	candidate        framework.Candidate
	population       []framework.Candidate
}

// newFitnessEvaluationTask creates a task for performing fitness evaluations.
//
// - fitnessEvaluator is the fitness function used to determine candidate
// fitness.
// - population is the entire current population. This will include all of the
// candidates to evaluate along with any other individuals that are not being
// evaluated by this task.
func newFitnessEvaluationTask(fitnessEvaluator framework.FitnessEvaluator,
	population []framework.Candidate) *fitnessEvaluationTask {

	return &fitnessEvaluationTask{
		fitnessEvaluator: fitnessEvaluator,
		population:       population,
	}
}

func (t *fitnessEvaluationTask) compute(candidate framework.Candidate) *framework.EvaluatedCandidate {
	ec, err := framework.NewEvaluatedCandidate(candidate,
		t.fitnessEvaluator.Fitness(candidate, t.population))
	if err != nil {
		panic(fmt.Sprintf("Error during fitness computation of candidate %v: %v", candidate, err))
	}
	return ec
}

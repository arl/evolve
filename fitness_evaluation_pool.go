package evolve

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
)

// fitnessEvaluationPool is the class that actually runs the fitness evaluation
// tasks created by an EvolutionEngine.
//
// This responsibility is abstracted away from the evolution engine to permit
// the possibility of creating multiple instances across several machines, all
// fed by a single shared work queue.
type fitnessEvaluationPool struct {

	// Provide each worker instance with a unique name with which to prefix its threads.
	idSource *workerIDSource
}

// Creates a fitnessEvaluationPool that uses daemon threads.
func newFitnessEvaluationPool() *fitnessEvaluationPool {
	return &fitnessEvaluationPool{
		idSource: newWorkerIDSource("FitnessEvaluationWorker"),
	}
}

func (w *fitnessEvaluationPool) submit(task *fitnessEvaluationTask) framework.EvaluatedPopulation {
	// We want to limit the number of concurrent goroutine used in the
	// computation of population fitness
	var (
		maxConcurrency = 4                              // TODO: should be configurable
		throttle       = make(chan int, maxConcurrency) // used as a limiter
		wg             sync.WaitGroup
	)

	// create the slice to receive the evaluted candidates
	results := make(framework.EvaluatedPopulation, len(task.population))
	// and the slice to receive the result channels (as many as candidates)
	resultsChan := make([]chan *framework.EvaluatedCandidate, len(task.population))

	for i := 0; i < len(results); i++ {
		throttle <- 1
		resultsChan[i] = compute(task, task.population[i], &wg, throttle)
	}

	// wait for all the tasks to be completed
	wg.Wait()
	close(throttle)

	// read from the results channels
	for i, resultChan := range resultsChan {
		results[i] = <-resultChan
	}

	return results
}

func compute(
	task *fitnessEvaluationTask,
	candidate framework.Candidate,
	wg *sync.WaitGroup,
	throttle chan int) chan *framework.EvaluatedCandidate {

	// create the channel to  send the result
	result := make(chan *framework.EvaluatedCandidate, 1)
	wg.Add(1)

	go func(chan *framework.EvaluatedCandidate) {
		defer wg.Done()
		result <- task.compute(candidate)
		close(result)
		<-throttle
	}(result)

	return result
}

type workerIDSource struct {
	lock      *sync.Mutex
	startTime time.Time
	lastID    int32
	prefix    string
}

func newWorkerIDSource(prefix string) *workerIDSource {
	return &workerIDSource{
		lock:      &sync.Mutex{},
		startTime: time.Now(),
		lastID:    -1,
		prefix:    prefix,
	}
}

// nextID returs the next worker ID
//
// Safe for concurrent use by multiple goroutines.
func (s *workerIDSource) nextID() (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.lastID == math.MaxInt32 {
		hours := time.Since(s.startTime) / time.Hour
		return "", fmt.Errorf("32-bit ID source exhausted after %d", hours)
	}
	s.lastID++
	return fmt.Sprintf("%s%d", s.prefix, s.lastID), nil
}

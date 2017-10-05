package worker

import (
	"sync"
)

// A Worker represents a self-contained worker and its load of work, performed
// when Work is called. It should be self-contained as Worker objects are meant
// to execute their work concurrently with other workers.
type Worker interface {

	// Work performs the work and returns the result once done.
	Work() interface{}
}

// WorkWith is a convenience function implementing the Worker interface,
// allowing one to use a function as a Worker, even an anonymous one.
//
// Example:
//  var w Worker = WorkWith(func() interface{} {
//	    // anonymous function used as Worker
//      return "hello"
//  })
type WorkWith func() interface{}

// Work performs the work on the delegate function w
func (w WorkWith) Work() interface{} {
	// call delegate function
	return w()
}

// A Pool manages the concurrent execution of multiple Workers.
type Pool struct {
	maxConcurrency int
}

// NewPool creates a Pool that may use concurrent goroutines to perform the work
// of some Workers objects.
//
// maxConcurrency limits the number of concurrent goroutines, so that at any
// moment, there will be at most maxConcurrency goroutines. If maxConcurrency is
// set to 1 the workers will run synchronously.
func NewPool(maxConcurrency int) *Pool {
	return &Pool{
		maxConcurrency: maxConcurrency,
	}
}

// Submit indicates to the pool the workers the workers to be run concurrently.
// The workers will be started in any order and may run concurrently, depending
// on the maxConcurrency defined at pool creation.
//
// The returned slice will contain the results, indexed as the workers.
func (w *Pool) Submit(workers []Worker) (results []interface{}, err error) {
	var (
		throttle = make(chan int, w.maxConcurrency) // used as a limiter
		wg       sync.WaitGroup
	)

	// create the slice to store the results
	results = make([]interface{}, len(workers))
	for i := range workers {

		// increment the rate limiter
		throttle <- 1

		// increment the work group
		wg.Add(1)

		go func(idx int) {

			// ensure the waitgroup is decremented, as well as the throttle
			// channel, once the worker has eneded
			defer func() {
				wg.Done()
				<-throttle
			}()

			// perform the work
			results[idx] = workers[idx].Work()
		}(i)
	}

	// wait for all the tasks to be completed
	wg.Wait()
	close(throttle)
	return results, nil
}

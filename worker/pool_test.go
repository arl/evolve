package worker

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomSeed() int64 {
	return int64(time.Now().UnixNano())
}

const sleepTimeMs = 500

type waiter struct{}

func (w waiter) Work() interface{} {
	time.Sleep(sleepTimeMs * time.Millisecond)
	return struct{}{}
}

func TestConcurrentWorkerRatio(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestConcurrentWorkerRatio in short mode")
	}

	const epsilon = 0.05 // accepted ratio error (5%)

	var tests = []struct {
		count       int     // number of works/workers
		concurrency int     // max concurrent workers
		ratio       float64 // expected ratio
	}{
		{20, 10, 2},
		{20, 4, 5},
		{20, 2, 10},
		{20, 1, 20},
	}

	for _, tt := range tests {

		workers := make([]Worker, tt.count)
		for i := 0; i < tt.count; i++ {
			workers[i] = waiter{}
		}

		pool := NewPool(tt.concurrency)

		start := time.Now()
		_, err := pool.Submit(workers)
		assert.NoError(t, err)

		// returns the ratio between the actual elapsed time and the time spent
		// in one worker. The more concurrent the workers do their job, the less
		// the ratio is.
		ratio := float64(time.Since(start)/time.Millisecond) / float64(sleepTimeMs)
		//t.Logf("ratio for %v|%v => %v, expected: %v\n", tt.count, tt.concurrency, ratio, tt.ratio)
		assert.InEpsilon(t, ratio, tt.ratio, tt.ratio*epsilon)
	}
}

type dummyWorker struct {
	idx int
}

func (w dummyWorker) Work() interface{} {
	// sleep for a random number of millliseconds ([50, 300])
	ms := 50 + rand.Intn(250)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return w.idx
}

func TestConcurrentWorkerResults(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestConcurrentWorkerResults in short mode")
	}

	var tests = []struct {
		concurrency int // max concurrent workers
	}{
		{10},
		{4},
		{2},
		{1},
	}

	for _, tt := range tests {

		workers := make([]Worker, 20)
		for i := range workers {
			workers[i] = dummyWorker{idx: i}
		}
		pool := NewPool(tt.concurrency)

		results, err := pool.Submit(workers)
		assert.NoError(t, err)

		// check the results slice is indexed as the workers
		for i, result := range results {
			assert.Equal(t, result, workers[i].(dummyWorker).idx)
		}
	}
}

func TestConcurrentWorkerFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestConcurrentWorkerFunction in short mode")
	}
	workers := make([]Worker, 20)
	for i := range workers {
		// create a closure just for the need of the test
		func(idx int) {

			// the worker is an anonymous function
			workers[idx] = WorkWith(func() interface{} {
				return idx
			})

		}(i) // pass it to capture the value of i
	}
	pool := NewPool(2)

	results, err := pool.Submit(workers)
	assert.NoError(t, err)

	// check the results slice is indexed as the workers
	for i, result := range results {
		assert.Equal(t, result, i)
	}
}

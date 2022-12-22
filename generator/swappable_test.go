package generator

import (
	"math"
	"sync"
	"testing"
)

func TestSwappable(t *testing.T) {
	var v float64 = math.MaxFloat64
	s := NewSwappable(Const(v))
	if got := s.Next(); got != v {
		t.Errorf("got %v, want %v", got, v)
	}

	v = -12
	s.Swap(Const(v))
	if got := s.Next(); got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestAtomicSwappable(t *testing.T) {
	// To run with -test.race
	var v float64 = math.MaxFloat64
	s := NewAtomicSwappable(Const(v))
	if got := s.Next(); got != v {
		t.Errorf("got %v, want %v", got, v)
	}

	const numgs = 32
	var wg sync.WaitGroup
	wg.Add(numgs)
	for i := 0; i < numgs; i++ {
		i := i
		go func() {
			s.Swap(Const(float64(i)))
			if got := s.Next(); 0 > got || got >= numgs {
				t.Errorf("got=%v, want 0 > got >= %v", got, numgs)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

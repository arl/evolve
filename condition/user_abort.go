package condition

import (
	"sync/atomic"

	"github.com/arl/evolve"
)

// UserAbort is a termination condition which is satisfies once Abort has been
// called. UserAbort zero-value is a valid, non-satisfied termination condition.
type UserAbort[T any] uint32

// IsSatisfied reports whether or not Abort has been called.
// It is safe for concurrent use by multiple goroutines.
func (ua *UserAbort[T]) IsSatisfied(*evolve.PopulationStats[T]) bool {
	return atomic.LoadUint32((*uint32)(ua)) == 1
}

// Abort triggers the condition.
// It is safe for concurrent use by multiple goroutines.
func (ua *UserAbort[T]) Abort() {
	atomic.StoreUint32((*uint32)(ua), 1)
}

// Reset resets the abort condition to false so that it may be reused.
// It is safe for concurrent use by multiple goroutines.
func (ua *UserAbort[T]) Reset() {
	atomic.StoreUint32((*uint32)(ua), 0)
}

// String returns a string representation of this condition.
func (UserAbort[T]) String() string { return "Aborted" }

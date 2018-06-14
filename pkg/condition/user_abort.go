package condition

import (
	"sync"

	"github.com/arl/evolve"
)

// UserAbort is a condition satisfied when Abort has been called. It allows for
// user-initiated termination of an evolution algorithm.
type UserAbort struct {
	mutex   *sync.RWMutex
	aborted bool
}

// NewUserAbort creates a UserAbort condition.
func NewUserAbort() *UserAbort {
	return &UserAbort{
		mutex:   &sync.RWMutex{},
		aborted: false,
	}
}

// IsSatisfied reports whether or not Abort has been called.
func (ua *UserAbort) IsSatisfied(*evolve.PopulationData) bool {
	ua.mutex.RLock()
	defer ua.mutex.RUnlock()
	return ua.aborted
}

// Abort triggers the condition.
func (ua *UserAbort) Abort() {
	ua.mutex.Lock()
	ua.aborted = true
	ua.mutex.Unlock()
}

// Reset resets the abort condition to false so that it may be reused.
func (ua *UserAbort) Reset() {
	ua.mutex.Lock()
	ua.aborted = false
	ua.mutex.Unlock()
}

// String returns a string representation of this condition.
func (UserAbort) String() string { return "Abort called" }

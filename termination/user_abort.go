package termination

import (
	"sync"

	"github.com/aurelien-rainone/evolve/framework"
)

// UserAbort is an implementation of the TerminationCondition interface that
// allows for user-initiated termination of an evolutionary algorithm.
//
// This condition can be used, for instance, to provide a button on a GUI that
// terminates execution. The application should retain a reference to the
// instance after passing it to the evolution engine and should invoke the
// Abort() function to make the evolution terminate at the end of the current
// generation.
type UserAbort struct {
	mutex   *sync.RWMutex
	aborted bool
}

// NewUserAbort creates a UserAbort termination condition.
func NewUserAbort() *UserAbort {
	return &UserAbort{
		mutex:   &sync.RWMutex{},
		aborted: false,
	}
}

// ShouldTerminate reports whether or not evolution should finish at the current
// point.
//
// populationData is the information about the current state of evolution.  This
// may be used to determine whether evolution should continue or not.
func (ua *UserAbort) ShouldTerminate(populationData *framework.PopulationData) bool {
	return ua.IsAborted()
}

// Abort aborts any evolutionary algorithms that monitor this termination
// condition instance.
func (ua *UserAbort) Abort() {
	ua.mutex.Lock()
	ua.aborted = true
	ua.mutex.Unlock()
}

// IsAborted returns true if Abort has been invoked, false otherwise.
func (ua *UserAbort) IsAborted() bool {
	ua.mutex.RLock()
	defer ua.mutex.RUnlock()
	return ua.aborted
}

// Reset resets the abort condition to false so that it may be reused.
func (ua *UserAbort) Reset() {
	ua.mutex.Lock()
	ua.aborted = false
	ua.mutex.Unlock()
}

// String returns the termination condition representation as a string
func (ua *UserAbort) String() string {
	return "User aborted evolution"
}

package evolve

import "errors"

// ErrIllegalState is the error returned when trying a method of an engine has
// been called while its state doesn't allow that method call.
var ErrIllegalState = errors.New("illegal state")

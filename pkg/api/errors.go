package api

import "fmt"

// ErrIllegalState is return by an function (or method) when it has been called
// when the state of some argument (or the receiver) doesn't permit this
// function (or method) to be called.
type ErrIllegalState string

func (err ErrIllegalState) Error() string {
	return fmt.Sprintf("illegal state: %v", err)
}

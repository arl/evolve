package generator

import "constraints"

// A Generator generates sequences of values each of which is provided whenever Next is called.
type Generator[T constraints.Integer | constraints.Float] interface {
	Next() T
}

// Unsigned generates unsigned integer values.
type Unsigned[T constraints.Unsigned] interface {
	Next() T
}

// Signed generates signed integer values.
type Signed[T constraints.Signed] interface {
	Next() T
}

// Float generates float64 values.
type Float Generator[float64]

type constGen[T constraints.Integer | constraints.Float] struct {
	val T
}

func (c constGen[T]) Next() T { return c.val }

// Const returns a Generator that always returns val.
func Const[T constraints.Integer | constraints.Float](val T) Generator[T] {
	return constGen[T]{val: val}
}

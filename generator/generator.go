package generator

type UnsignedInteger interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

// UInt generates unsigned integer values.
type UInt[T UnsignedInteger] interface {
	Next() T
}

// Float generates float64 values.
type Float Generator[float64]

// A Generator generates sequences of values each of which is provided whenever Next is called.
type Generator[T any] interface {
	Next() T
}

type constGen[T any] struct {
	val T
}

func (c constGen[T]) Next() T { return c.val }

// Const returns a Generator that always returns val.
func Const[T any](val T) Generator[T] {
	return constGen[T]{val: val}
}

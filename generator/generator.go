package generator

// A Generator generates sequences of values each of which is provided whenever Next is called.
type Generator[T Number] interface {
	Next() T
}

type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

type Signed interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type Number interface {
	Unsigned | Signed | ~float32 | ~float64
}

// UInt generates unsigned integer values.
type UInt[T Unsigned] interface {
	Next() T
}

// Int generates signed integer values.
type Int[T Signed] interface {
	Next() T
}

// Float generates float64 values.
type Float Generator[float64]

type constGen[T Number] struct {
	val T
}

func (c constGen[T]) Next() T { return c.val }

// Const returns a Generator that always returns val.
func Const[T Number](val T) Generator[T] {
	return constGen[T]{val: val}
}

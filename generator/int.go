package generator

// An Int generates sequences of integers.
type Int interface {
	// Next returns the next number in the sequence.
	Next() int64
}

type ConstInt int64

// Next always returns i.
func (i ConstInt) Next() int64 {
	return int64(i)
}

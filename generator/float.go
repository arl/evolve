package generator

// An Float generates sequences of floating point numbers.
type Float interface {
	// Next returns the next number in the sequence.
	Next() float64
}

type ConstFloat64 float64

// Next always returns i.
func (f ConstFloat64) Next() float64 {
	return float64(f)
}

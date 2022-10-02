package generator

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

// Uniform returns a generator of random numbers which are uniformly distributed
// in the [min max) range. The range must be made of positive numbers only. In
// other words if T is a floating point number type, then min must be positive.
// Uniform panics if max <= min.
func Uniform[T constraints.Integer | constraints.Float](min, max T, rng *rand.Rand) Generator[T] {
	diff := max - min
	if diff <= 0 {
		panic("must have min < max")
	}
	if min < 0 {
		panic("must have min positive")
	}

	// TODO(generics) check status of generic type switches proposal
	// https://github.com/golang/go/issues/45380
	var t T
	switch any(t).(type) {
	case int, uint,
		int8, uint8,
		int16, uint16,
		int32, uint32,
		int64, uint64:
		idiff := int64(diff)
		return uniform[T](func() T {
			return min + T(rng.Int63n(idiff))
		})
	case float32:
		f32diff := float32(diff)
		return uniform[T](func() T {
			return min + T(rng.Float32())*T(f32diff)
		})
	case float64:
		return uniform[T](func() T {
			return min + T(rng.Float64())*T(diff)
		})
	}
	return nil
}

type uniform[T constraints.Integer | constraints.Float] func() T

func (u uniform[T]) Next() T { return u() }

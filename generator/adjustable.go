package generator

import (
	"math"
	"sync/atomic"
)

// AdjustableInt is an atomically adjustable int64 values generator.
type AdjustableInt struct {
	_ nocmp // disallow non-atomic comparison

	v int64
}

// NewAdjustableInt creates a new AdjustableInt.
func NewAdjustableInt(i int64) *AdjustableInt {
	return &AdjustableInt{v: i}
}

func (v *AdjustableInt) Set(newv int64) { atomic.StoreInt64(&v.v, newv) }
func (v *AdjustableInt) Next() int64    { return atomic.LoadInt64(&v.v) }

// AdjustableFloat is an atomically adjustable float64 values generator.
type AdjustableFloat struct {
	_ nocmp // disallow non-atomic comparison

	v uint64
}

var _float64_0 float64

// NewFloat64 creates a new Float64.
func NewAdjustableFloat(v float64) *AdjustableFloat {
	x := &AdjustableFloat{}
	if v != _float64_0 {
		x.Set(v)
	}
	return x
}

func (v *AdjustableFloat) Set(newv float64) {
	atomic.StoreUint64(&v.v, math.Float64bits(newv))
}

func (v *AdjustableFloat) Next() float64 {
	return math.Float64frombits(atomic.LoadUint64(&v.v))
}

/*

// Load atomically loads the wrapped float64.
func (x *Float64) Load() float64 {
	return math.Float64frombits(x.v.Load())
}

// Store atomically stores the passed float64.
func (x *Float64) Store(v float64) {
	x.v.Store(math.Float64bits(v))
}

/*
type AdjustableFloat float64

func (v *AdjustableFloat) Set(newv float64) {

	f := (float64)(*v)
	atomic.StoreUint64((*uint64)(), math.Float64bits(newv))
}

func (v *AdjustableFloat) Next() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(v)))
}*/

package generator

import (
	"math"
	"testing"
)

func TestAdjustableInt(t *testing.T) {
	tests := []int64{
		1,
		-1,
		0,
		-0,
		math.MaxInt64,
		math.MinInt64,
	}
	for _, want := range tests {
		var f AdjustableInt
		f.Set(want)

		if got := f.Next(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestAdjustableFloat(t *testing.T) {
	tests := []float64{
		1,
		-1,
		147,
		17,
		0,
		-0,
		math.MaxFloat64,
		-math.MaxFloat64,
	}

	var f AdjustableFloat

	for _, want := range tests {

		f.Set(want)
		if got := f.Next(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestAdjustableFloatNan(t *testing.T) {
	var f AdjustableFloat
	f.Set(math.NaN())
	if got := f.Next(); !math.IsNaN(got) {
		t.Errorf("got %v, want Nan", got)
	}

	f.Set(-math.NaN())
	if got := f.Next(); !math.IsNaN(got) {
		t.Errorf("got %v, want Nan", got)
	}
}

func TestAdjustableFloatInf(t *testing.T) {
	var f AdjustableFloat
	f.Set(math.Inf(1))
	if got := f.Next(); !math.IsInf(got, 1) {
		t.Errorf("got %v, want +Inf", got)
	}

	f.Set(math.Inf(-1))
	if got := f.Next(); !math.IsInf(got, -1) {
		t.Errorf("got %v, want -Inf", got)
	}
}

// func TestAdjustableFloatSetNext(t *testing.T) {
// 	const f = 147
// 	a := AdjustableFloat(f)
// 	if got := a.Next(); got != f {
// 		t.Errorf("got %v, want %v", got, f)
// 	}
// }

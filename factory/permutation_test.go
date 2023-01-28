package factory

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestPermutation(t *testing.T) {
	org := make([]int, 10)
	for i := 0; i < len(org); i++ {
		org[i] = i
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	f := Permutation[int](org)
	cand := f.New(rng)

	if len(cand) != len(org) {
		t.Errorf("permutated slice should have the same length")
	}

	if reflect.ValueOf(org).Pointer() == reflect.ValueOf(cand).Pointer() {
		t.Fatalf("new slice has the same backing array as the original")
	}

	if reflect.DeepEqual(org, cand) {
		t.Errorf("new slice should be a different permutation, got %+v", cand)
	}
}

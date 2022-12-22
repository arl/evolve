package factory

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPermutation(t *testing.T) {
	items := make([]int, 10)

	for i := 0; i < len(items); i++ {
		items[i] = i
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	f := Permutation[int](items)
	cpy := f.New(rng)

	assert.Len(t, cpy, len(items), "permutated slice should have the same length")

	same := reflect.ValueOf(items).Pointer() == reflect.ValueOf(cpy).Pointer()
	assert.False(t, same, "new slice has the same backing array as the original")

	assert.NotEqualValues(t, items, cpy, "new slice should have permuted values")
}

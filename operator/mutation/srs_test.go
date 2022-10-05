package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_srs(t *testing.T) {
	tests := []struct {
		org      []int
		xp1, xp2 int
		mut      []int // want
	}{
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//              |     |
			xp1: 2, xp2: 4,
			mut: []int{5, 6, 7, 3, 4, 2, 1},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//           |  |
			xp1: 1, xp2: 2,
			mut: []int{3, 4, 5, 6, 7, 2, 1},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//                       |  |
			xp1: 5, xp2: 6,
			mut: []int{7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		got := make([]int, len(tt.org))
		srs(tt.org, got, tt.xp1, tt.xp2)
		assert.Equal(t, tt.mut, got)
	}
}

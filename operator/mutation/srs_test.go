package mutation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_srs(t *testing.T) {
	tests := []struct {
		org      []int
		xp1, xp2 int
		want     []int // want
	}{
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//              |     |
			xp1: 2, xp2: 4,
			want: []int{5, 6, 7, 3, 4, 2, 1},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//           |  |
			xp1: 1, xp2: 2,
			want: []int{3, 4, 5, 6, 7, 2, 1},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6, 7},
			//                       |  |
			xp1: 5, xp2: 6,
			want: []int{7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		// Before calling srs, the slice to mutate is already a copy of the original slice.
		mut := make([]int, len(tt.org))
		copy(mut, tt.org)

		srs(tt.org, mut, tt.xp1, tt.xp2)
		if !cmp.Equal(tt.want, mut) {
			t.Errorf("srs(%+v, [res], %d, %d) = %+v, want %+v", tt.org, tt.xp1, tt.xp2, mut, tt.want)
		}
	}
}

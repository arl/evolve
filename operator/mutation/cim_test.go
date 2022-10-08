package mutation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_cim(t *testing.T) {
	tests := []struct {
		org  []int
		xp   int
		want []int // want
	}{
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//              |     |
			xp:   4,
			want: []int{4, 3, 2, 1, 6, 5},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//              |     |
			xp:   1,
			want: []int{1, 6, 5, 4, 3, 2},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//              |     |
			xp:   5,
			want: []int{5, 4, 3, 2, 1, 6},
		},
	}

	for _, tt := range tests {
		// Before calling cim, the slice to mutate is already a copy of the original slice.
		got := make([]int, len(tt.org))
		copy(got, tt.org)

		cim(got, tt.xp)
		if !cmp.Equal(tt.want, got) {
			t.Errorf("cim(%+v, %d) = %+v, want %+v", tt.org, tt.xp, got, tt.want)
		}
	}
}

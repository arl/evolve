package mutation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_cim(t *testing.T) {
	tests := []struct {
		org  []int
		cut  int
		want []int
	}{
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//                    |
			cut:  4,
			want: []int{4, 3, 2, 1, 6, 5},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//           |
			cut:  1,
			want: []int{1, 6, 5, 4, 3, 2},
		},
		{
			org: []int{1, 2, 3, 4, 5, 6},
			//                       |
			cut:  5,
			want: []int{5, 4, 3, 2, 1, 6},
		},
		{
			org: []int{1, 2},
			//         |
			cut:  0,
			want: []int{2, 1},
		},
		{
			org: []int{1, 2},
			//         |
			cut:  2,
			want: []int{2, 1},
		},
		{
			org: []int{1, 2},
			//           |
			cut:  1,
			want: []int{1, 2},
		},
		{
			org: []int{1},
			//           |
			cut:  1,
			want: []int{1},
		},
		{
			org: []int{1},
			//           |
			cut:  0,
			want: []int{1},
		},
	}

	for _, tt := range tests {
		// Before calling cim, the slice to mutate is already a copy of the original slice.
		got := make([]int, len(tt.org))
		copy(got, tt.org)

		cim(got, tt.cut)
		if !cmp.Equal(tt.want, got) {
			t.Errorf("cim(%+v, %d) = %+v, want %+v", tt.org, tt.cut, got, tt.want)
		}
	}
}

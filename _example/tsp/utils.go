package main

import "golang.org/x/exp/constraints"

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func worldBounds(cities []point) (maxw, maxh float64) {
	for _, c := range cities {
		maxw = max(maxw, c.X)
		maxh = max(maxw, c.Y)
	}
	return
}

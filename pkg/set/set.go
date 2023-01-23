// Package provides a set based on map.
package set

// Set is a set based on a map.s
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewOf creates and initializes a set of T.
func NewOf[T comparable]() Set[T] {
	return Set[T]{
		m: make(map[T]struct{}),
	}
}

// NewOf creates and initializes a set of T filled with the given elements.
func With[T comparable](elems ...T) Set[T] {
	s := NewOf[T]()
	for _, val := range elems {
		s.Insert(val)
	}
	return s
}

// Insert insert elem in the set.
func (s Set[T]) Insert(elem T) {
	s.m[elem] = struct{}{}
}

// Contains reports whether elem is a member of s.
func (s Set[T]) Contains(elem T) bool {
	_, ok := s.m[elem]
	return ok
}

// Delete deletes elemt from s.
func (s Set[T]) Delete(elem T) {
	delete(s.m, elem)
}

// Clear removes all elements from the set.
func (s Set[T]) Clear() {
	for e := range s.m {
		delete(s.m, e)
	}
}

// Len returns s cardinality.
func (s Set[T]) Len() int {
	return len(s.m)
}

// Each calls f for every element in s.
func (s Set[T]) Each(f func(elem T)) {
	for e := range s.m {
		f(e)
	}
}

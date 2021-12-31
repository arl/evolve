package generator

import (
	"sync/atomic"
)

// A Swappable embeds a generator that can be swapped with another. Not safe for concurrent use.
type Swappable[T any] struct {
	g Generator[T]
}

// NewSwappable returns a new Swappable generator, initialized with g.
func NewSwappable[T any](g Generator[T]) *Swappable[T] {
	return &Swappable[T]{g: g}
}

// Swap swaps the current generator with g.
func (m *Swappable[T]) Swap(g Generator[T]) {
	m.g = g
}

// Next calls Next on the current generator.
func (m *Swappable[T]) Next() T {
	return m.g.Next()
}

// An AtomicSwappable embeds a generator that can be swapped with another,
// atomically. This is the concurrent-safe counterpart of Swappable.
type AtomicSwappable[T any] struct {
	g atomic.Value
}

// NewAtomicSwappable returns a new AtomicSwappable generator, initialized with g.
func NewAtomicSwappable[T any](g Generator[T]) *AtomicSwappable[T] {
	var a AtomicSwappable[T]
	a.g.Store(g)
	return &a
}

// Swap atomically swaps the current generator with g.
func (a *AtomicSwappable[T]) Swap(g Generator[T]) {
	a.g.Swap(g)
}

// Next calls Next on the current generator.
func (v *AtomicSwappable[T]) Next() T {
	return v.g.Load().(Generator[T]).Next()
}

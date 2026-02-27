package stream

import "iter"

// ---------------------------------------------------------------------------
// Bridge: iter.Seq[T] → Stream[T]
// ---------------------------------------------------------------------------
// These functions connect Go 1.23+'s iterator protocol with Stream,
// enabling seamless interop with the standard library's slices.Values,
// maps.Keys, and other iter.Seq-producing functions.

// Collect creates a Stream[T] from an iter.Seq[T].
//
//	// Use with standard library iterators
//	keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()
func Collect[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}

// Collect2 creates a Stream[Pair[K,V]] from an iter.Seq2[K,V].
//
//	pairs := stream.Collect2(maps.All(myMap))
func Collect2[K, V any](seq iter.Seq2[K, V]) Stream[Pair[K, V]] {
	return Stream[Pair[K, V]]{seq: func(yield func(Pair[K, V]) bool) {
		for k, v := range seq {
			if !yield(Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}}
}

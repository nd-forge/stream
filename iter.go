package stream

import (
	"iter"
	"slices"
)

// ---------------------------------------------------------------------------
// Bridge: Stream[T] ↔ iter.Seq[T]
// ---------------------------------------------------------------------------
// These functions connect the eager Stream world with Go 1.23+'s
// iterator protocol, enabling seamless interop with the standard library's
// slices.Values, maps.Keys, and other iter.Seq-producing functions.

// Iter returns an iter.Seq[T] that yields all elements of the Stream.
//
//	s := stream.Of(1, 2, 3)
//	for v := range s.Iter() {
//	    fmt.Println(v)
//	}
func (s Stream[T]) Iter() iter.Seq[T] {
	return slices.Values(s.data)
}

// Iter2 returns an iter.Seq2[int, T] that yields index-element pairs.
//
//	for i, v := range stream.Of("a", "b", "c").Iter2() {
//	    fmt.Printf("%d: %s\n", i, v)
//	}
func (s Stream[T]) Iter2() iter.Seq2[int, T] {
	return slices.All(s.data)
}

// Collect creates a Stream[T] from an iter.Seq[T].
// This eagerly consumes the entire iterator into a slice.
//
//	// Use with standard library iterators
//	keys := stream.Collect(maps.Keys(myMap))
//	sorted := keys.Sort(cmp).ToSlice()
//
//	// Use with Pipeline (lazy → eager)
//	result := stream.Collect(pipeline.Seq())
func Collect[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{data: slices.Collect(seq)}
}

// Collect2 creates a Stream[Pair[K,V]] from an iter.Seq2[K,V].
//
//	pairs := stream.Collect2(maps.All(myMap))
func Collect2[K, V any](seq iter.Seq2[K, V]) Stream[Pair[K, V]] {
	var result []Pair[K, V]
	for k, v := range seq {
		result = append(result, Pair[K, V]{First: k, Second: v})
	}
	return Stream[Pair[K, V]]{data: result}
}

// Package stream provides a chainable collection operations library for Go.
//
// Stream[T] wraps an iter.Seq[T] and provides lazy evaluation with method
// chaining for filter, sort, take, skip, and other operations that preserve
// the element type. For type-changing operations (Map, FlatMap, Zip), use
// top-level functions.
//
// All operations are lazy by default. Operations that require full data
// (Sort, Reverse, Shuffle, TakeLast, Chunk, Partition) buffer internally
// and resume lazy iteration.
//
// Usage:
//
//	// Method chaining for same-type operations
//	result := stream.Of(1, 2, 3, 4, 5).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Sort(func(a, b int) int { return a - b }).
//	    Take(3).
//	    ToSlice()
//
//	// Top-level functions for type-changing operations
//	names := stream.Map(
//	    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
//	    func(u User) string { return u.Name },
//	).ToSlice()
package stream

import (
	"iter"
	"math/rand"
	"sort"
)

// Stream is a lazy evaluation wrapper around iter.Seq[T] that supports
// method chaining. All intermediate operations are deferred until a terminal
// operation is called. Operations that require full data (Sort, Reverse,
// Shuffle, TakeLast) buffer internally then resume lazy iteration.
//
// Stream is reusable: calling terminal operations multiple times produces
// the same result, as the underlying iter.Seq is re-executed each time.
type Stream[T any] struct {
	seq iter.Seq[T]
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// Of creates a new Stream from variadic arguments.
func Of[T any](items ...T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}}
}

// From creates a new Stream from an existing slice (copies the slice).
func From[T any](items []T) Stream[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return Of(copied...)
}

// Generate creates a Stream of n elements using a generator function.
func Generate[T any](n int, gen func(index int) T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(gen(i)) {
				return
			}
		}
	}}
}

// Range creates a Stream of integers from start (inclusive) to end (exclusive).
func Range(start, end int) Stream[int] {
	return Stream[int]{seq: func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}}
}

// ---------------------------------------------------------------------------
// Generators (infinite sequences)
// ---------------------------------------------------------------------------

// Repeat creates an infinite Stream that yields the same value.
// Must be combined with Take, TakeWhile, or Find to terminate.
//
//	zeros := stream.Repeat(0).Take(10).ToSlice()
func Repeat[T any](value T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for {
			if !yield(value) {
				return
			}
		}
	}}
}

// RepeatN creates a Stream that yields the same value n times.
//
//	dashes := stream.RepeatN("-", 20).ToSlice()
func RepeatN[T any](value T, n int) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(value) {
				return
			}
		}
	}}
}

// Iterate creates an infinite Stream: seed, fn(seed), fn(fn(seed)), ...
//
//	// Powers of 2: 1, 2, 4, 8, 16, ...
//	stream.Iterate(1, func(n int) int { return n * 2 }).Take(5)
func Iterate[T any](seed T, fn func(T) T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		v := seed
		for {
			if !yield(v) {
				return
			}
			v = fn(v)
		}
	}}
}

// Naturals returns an infinite Stream of natural numbers: 0, 1, 2, 3, ...
func Naturals() Stream[int] {
	return Iterate(0, func(n int) int { return n + 1 })
}

// ---------------------------------------------------------------------------
// Chainable methods (type-preserving: Stream[T] → Stream[T])
// ---------------------------------------------------------------------------

// Filter returns a Stream that yields only elements satisfying the predicate.
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			if predicate(v) {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// Reject returns a Stream that excludes elements satisfying the predicate.
func (s Stream[T]) Reject(predicate func(T) bool) Stream[T] {
	return s.Filter(func(v T) bool { return !predicate(v) })
}

// Sort buffers all elements, sorts them, and yields in sorted order.
// Note: This operation consumes all elements into memory, breaking pure laziness.
// However, subsequent operations in the chain remain lazy.
func (s Stream[T]) Sort(cmp func(a, b T) int) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		var buf []T
		for v := range seq {
			buf = append(buf, v)
		}
		sort.Slice(buf, func(i, j int) bool {
			return cmp(buf[i], buf[j]) < 0
		})
		for _, v := range buf {
			if !yield(v) {
				return
			}
		}
	}}
}

// Reverse buffers all elements and yields them in reverse order.
// Note: This operation consumes all elements into memory.
func (s Stream[T]) Reverse() Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		var buf []T
		for v := range seq {
			buf = append(buf, v)
		}
		for i := len(buf) - 1; i >= 0; i-- {
			if !yield(buf[i]) {
				return
			}
		}
	}}
}

// Take returns a Stream that yields only the first n elements.
// For infinite sequences, this is the primary way to limit output.
func (s Stream[T]) Take(n int) Stream[T] {
	if n <= 0 {
		return Stream[T]{seq: func(yield func(T) bool) {}}
	}
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		i := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			i++
			if i >= n {
				return
			}
		}
	}}
}

// TakeLast buffers all elements and returns only the last n elements.
// Note: This operation consumes all elements into memory.
func (s Stream[T]) TakeLast(n int) Stream[T] {
	if n <= 0 {
		return Stream[T]{seq: func(yield func(T) bool) {}}
	}
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		var buf []T
		for v := range seq {
			buf = append(buf, v)
		}
		start := len(buf) - n
		if start < 0 {
			start = 0
		}
		for _, v := range buf[start:] {
			if !yield(v) {
				return
			}
		}
	}}
}

// Skip returns a Stream that skips the first n elements.
func (s Stream[T]) Skip(n int) Stream[T] {
	if n <= 0 {
		return s
	}
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		i := 0
		for v := range seq {
			if i < n {
				i++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}}
}

// TakeWhile returns elements from the start as long as the predicate is true.
func (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			if !predicate(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}}
}

// DropWhile skips elements from the start while the predicate is true,
// then yields the rest.
func (s Stream[T]) DropWhile(predicate func(T) bool) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		dropping := true
		for v := range seq {
			if dropping {
				if predicate(v) {
					continue
				}
				dropping = false
			}
			if !yield(v) {
				return
			}
		}
	}}
}

// Distinct returns a Stream with duplicate elements removed.
// Uses the provided key function to determine equality.
// Note: Maintains a set of seen keys in memory.
func (s Stream[T]) Distinct(key func(T) string) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		seen := make(map[string]struct{})
		for v := range seq {
			k := key(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// Shuffle buffers all elements, randomizes their order, and yields them.
// Note: This operation consumes all elements into memory.
func (s Stream[T]) Shuffle() Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		var buf []T
		for v := range seq {
			buf = append(buf, v)
		}
		rand.Shuffle(len(buf), func(i, j int) {
			buf[i], buf[j] = buf[j], buf[i]
		})
		for _, v := range buf {
			if !yield(v) {
				return
			}
		}
	}}
}

// Peek executes a side-effect function for each element without modifying the Stream.
// Useful for debugging or logging within a lazy chain.
func (s Stream[T]) Peek(fn func(T)) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			fn(v)
			if !yield(v) {
				return
			}
		}
	}}
}

// Chain concatenates multiple Streams, yielding all elements from each in order.
//
//	combined := s1.Chain(s2, s3)
func (s Stream[T]) Chain(others ...Stream[T]) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			if !yield(v) {
				return
			}
		}
		for _, other := range others {
			for v := range other.seq {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// ---------------------------------------------------------------------------
// Terminal operations (consume the Stream)
// ---------------------------------------------------------------------------

// ToSlice collects all elements into a slice.
func (s Stream[T]) ToSlice() []T {
	var result []T
	for v := range s.seq {
		result = append(result, v)
	}
	if result == nil {
		return []T{}
	}
	return result
}

// Seq returns the underlying iter.Seq[T].
// Use this for interop with standard library functions like slices.Collect.
func (s Stream[T]) Seq() iter.Seq[T] {
	return s.seq
}

// ForEach executes a function for each element.
func (s Stream[T]) ForEach(fn func(T)) {
	for v := range s.seq {
		fn(v)
	}
}

// ForEachIndexed executes a function for each element with its index.
func (s Stream[T]) ForEachIndexed(fn func(int, T)) {
	i := 0
	for v := range s.seq {
		fn(i, v)
		i++
	}
}

// Reduce folds all elements into a single value of the same type.
// For reducing to a different type, use the top-level Reduce function.
func (s Stream[T]) Reduce(initial T, fn func(acc, item T) T) T {
	result := initial
	for v := range s.seq {
		result = fn(result, v)
	}
	return result
}

// First returns the first element and true, or zero value and false if empty.
// For infinite Streams, this returns immediately.
func (s Stream[T]) First() (T, bool) {
	for v := range s.seq {
		return v, true
	}
	var zero T
	return zero, false
}

// Last returns the last element and true, or zero value and false if empty.
// Warning: This consumes the entire Stream. Do not use on infinite sequences.
func (s Stream[T]) Last() (T, bool) {
	var last T
	found := false
	for v := range s.seq {
		last = v
		found = true
	}
	return last, found
}

// Find returns the first element matching the predicate.
// Short-circuits: stops iteration as soon as a match is found.
func (s Stream[T]) Find(predicate func(T) bool) (T, bool) {
	for v := range s.seq {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if any element satisfies the predicate.
// Short-circuits on first match.
func (s Stream[T]) Any(predicate func(T) bool) bool {
	for v := range s.seq {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements satisfy the predicate.
// Short-circuits on first non-match.
func (s Stream[T]) All(predicate func(T) bool) bool {
	for v := range s.seq {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// None returns true if no elements satisfy the predicate.
func (s Stream[T]) None(predicate func(T) bool) bool {
	return !s.Any(predicate)
}

// Count returns the total number of elements.
// Warning: Consumes the entire Stream. Do not use on infinite sequences.
func (s Stream[T]) Count() int {
	n := 0
	for range s.seq {
		n++
	}
	return n
}

// CountBy returns the number of elements satisfying the predicate.
// Warning: Consumes the entire Stream. Do not use on infinite sequences.
func (s Stream[T]) CountBy(predicate func(T) bool) int {
	n := 0
	for v := range s.seq {
		if predicate(v) {
			n++
		}
	}
	return n
}

// IsEmpty returns true if the Stream has no elements.
func (s Stream[T]) IsEmpty() bool {
	for range s.seq {
		return false
	}
	return true
}

// Contains returns true if any element satisfies the predicate.
func (s Stream[T]) Contains(predicate func(T) bool) bool {
	return s.Any(predicate)
}

// Partition collects all elements and splits them into two Streams:
// elements that match the predicate and those that don't.
// Note: This operation consumes all elements into memory.
func (s Stream[T]) Partition(predicate func(T) bool) (matched Stream[T], unmatched Stream[T]) {
	var yes, no []T
	for v := range s.seq {
		if predicate(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return From(yes), From(no)
}

// Chunk collects all elements and splits them into chunks of the specified size.
// Note: This operation consumes all elements into memory.
func (s Stream[T]) Chunk(size int) []Stream[T] {
	if size <= 0 {
		return nil
	}
	var buf []T
	for v := range s.seq {
		buf = append(buf, v)
	}
	var chunks []Stream[T]
	for i := 0; i < len(buf); i += size {
		end := i + size
		if end > len(buf) {
			end = len(buf)
		}
		chunk := make([]T, end-i)
		copy(chunk, buf[i:end])
		chunks = append(chunks, Of(chunk...))
	}
	return chunks
}

// MinBy returns the minimum element according to the comparison function.
// Warning: Consumes the entire Stream.
func (s Stream[T]) MinBy(less func(a, b T) bool) (T, bool) {
	var min T
	found := false
	for v := range s.seq {
		if !found || less(v, min) {
			min = v
			found = true
		}
	}
	return min, found
}

// MaxBy returns the maximum element according to the comparison function.
// Warning: Consumes the entire Stream.
func (s Stream[T]) MaxBy(less func(a, b T) bool) (T, bool) {
	var max T
	found := false
	for v := range s.seq {
		if !found || less(max, v) {
			max = v
			found = true
		}
	}
	return max, found
}

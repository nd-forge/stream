package stream

import (
	"iter"
	"sort"
)

// Pipeline is a lazy evaluation wrapper around iter.Seq[T] that supports
// method chaining. Unlike Stream which eagerly allocates intermediate slices,
// Pipeline defers all computation until a terminal operation is called.
//
// Pipeline is reusable: calling terminal operations multiple times produces
// the same result, as the underlying iter.Seq is re-executed each time.
//
// Usage:
//
//	// Lazy pipeline — no intermediate slices allocated
//	result := stream.Lazy(1, 2, 3, 4, 5).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Take(2).
//	    ToSlice()
//
//	// Convert Stream → Pipeline for lazy processing
//	result := stream.Of(items...).Lazy().Filter(pred).Take(10).ToSlice()
//
//	// Infinite sequences
//	evens := stream.Iterate(0, func(n int) int { return n + 2 }).Take(10).ToSlice()
type Pipeline[T any] struct {
	seq iter.Seq[T]
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// Lazy creates a Pipeline from variadic arguments.
func Lazy[T any](items ...T) Pipeline[T] {
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}}
}

// LazyFrom creates a Pipeline from an existing iter.Seq.
// Use this to wrap standard library iterators like slices.Values or maps.Keys.
//
//	pipe := stream.LazyFrom(slices.Values(mySlice))
//	pipe := stream.LazyFrom(maps.Keys(myMap))
func LazyFrom[T any](seq iter.Seq[T]) Pipeline[T] {
	return Pipeline[T]{seq: seq}
}

// Lazy converts a Stream into a Pipeline for lazy evaluation.
//
//	// Switch to lazy mode for efficient Take after heavy filter
//	result := stream.Of(items...).Lazy().Filter(pred).Take(5).ToSlice()
func (s Stream[T]) Lazy() Pipeline[T] {
	return LazyFrom[T](s.Iter())
}

// LazyRange creates a Pipeline of integers from start (inclusive) to end (exclusive).
// Unlike Range, this does not allocate a slice upfront.
func LazyRange(start, end int) Pipeline[int] {
	return Pipeline[int]{seq: func(yield func(int) bool) {
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

// Repeat creates an infinite Pipeline that yields the same value.
// Must be combined with Take, TakeWhile, or Find to terminate.
//
//	zeros := stream.Repeat(0).Take(10).ToSlice()
func Repeat[T any](value T) Pipeline[T] {
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for {
			if !yield(value) {
				return
			}
		}
	}}
}

// RepeatN creates a Pipeline that yields the same value n times.
//
//	dashes := stream.RepeatN("-", 20).ToSlice()
func RepeatN[T any](value T, n int) Pipeline[T] {
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(value) {
				return
			}
		}
	}}
}

// Iterate creates an infinite Pipeline: seed, fn(seed), fn(fn(seed)), ...
//
//	// Powers of 2: 1, 2, 4, 8, 16, ...
//	stream.Iterate(1, func(n int) int { return n * 2 }).Take(5)
//
//	// Fibonacci (using Pair):
//	stream.Iterate(
//	    stream.Pair[int,int]{First: 0, Second: 1},
//	    func(p stream.Pair[int,int]) stream.Pair[int,int] {
//	        return stream.Pair[int,int]{First: p.Second, Second: p.First + p.Second}
//	    },
//	).Take(10)
func Iterate[T any](seed T, fn func(T) T) Pipeline[T] {
	return Pipeline[T]{seq: func(yield func(T) bool) {
		v := seed
		for {
			if !yield(v) {
				return
			}
			v = fn(v)
		}
	}}
}

// Naturals returns an infinite Pipeline of natural numbers: 0, 1, 2, 3, ...
func Naturals() Pipeline[int] {
	return Iterate(0, func(n int) int { return n + 1 })
}

// ---------------------------------------------------------------------------
// Chainable methods (type-preserving: Pipeline[T] → Pipeline[T])
// ---------------------------------------------------------------------------

// Filter returns a Pipeline that yields only elements satisfying the predicate.
func (p Pipeline[T]) Filter(pred func(T) bool) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			if pred(v) {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// Reject returns a Pipeline that excludes elements satisfying the predicate.
func (p Pipeline[T]) Reject(pred func(T) bool) Pipeline[T] {
	return p.Filter(func(v T) bool { return !pred(v) })
}

// Take returns a Pipeline that yields only the first n elements.
// For infinite sequences, this is the primary way to limit output.
func (p Pipeline[T]) Take(n int) Pipeline[T] {
	if n <= 0 {
		return Pipeline[T]{seq: func(yield func(T) bool) {}}
	}
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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

// Skip returns a Pipeline that skips the first n elements.
func (p Pipeline[T]) Skip(n int) Pipeline[T] {
	if n <= 0 {
		return p
	}
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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
func (p Pipeline[T]) TakeWhile(pred func(T) bool) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			if !pred(v) {
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
func (p Pipeline[T]) DropWhile(pred func(T) bool) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
		dropping := true
		for v := range seq {
			if dropping {
				if pred(v) {
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

// Distinct returns a Pipeline with duplicate elements removed.
// Uses the provided key function to determine equality.
// Note: Maintains a set of seen keys in memory.
func (p Pipeline[T]) Distinct(key func(T) string) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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

// Peek executes a side-effect function for each element without modifying the Pipeline.
// Useful for debugging or logging within a lazy chain.
func (p Pipeline[T]) Peek(fn func(T)) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for v := range seq {
			fn(v)
			if !yield(v) {
				return
			}
		}
	}}
}

// Chain concatenates multiple Pipelines, yielding all elements from each in order.
//
//	combined := pipe1.Chain(pipe2, pipe3)
func (p Pipeline[T]) Chain(others ...Pipeline[T]) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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

// Sort buffers all elements, sorts them, and yields in sorted order.
// Note: This operation consumes all elements into memory, breaking pure laziness.
// However, subsequent operations in the chain remain lazy.
func (p Pipeline[T]) Sort(cmp func(a, b T) int) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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
func (p Pipeline[T]) Reverse() Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
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

// ---------------------------------------------------------------------------
// Terminal operations (consume the Pipeline)
// ---------------------------------------------------------------------------

// ToSlice collects all elements into a slice.
func (p Pipeline[T]) ToSlice() []T {
	var result []T
	for v := range p.seq {
		result = append(result, v)
	}
	if result == nil {
		return []T{}
	}
	return result
}

// ToStream collects all elements into an eager Stream.
// Useful when you need Stream's operations (e.g., TakeLast, Chunk, Shuffle).
//
//	result := pipeline.Filter(pred).ToStream().Shuffle().Take(5).ToSlice()
func (p Pipeline[T]) ToStream() Stream[T] {
	return Stream[T]{data: p.ToSlice()}
}

// Seq returns the underlying iter.Seq[T].
// Use this for interop with standard library functions like slices.Collect.
func (p Pipeline[T]) Seq() iter.Seq[T] {
	return p.seq
}

// ForEach executes a function for each element.
func (p Pipeline[T]) ForEach(fn func(T)) {
	for v := range p.seq {
		fn(v)
	}
}

// ForEachIndexed executes a function for each element with its index.
func (p Pipeline[T]) ForEachIndexed(fn func(int, T)) {
	i := 0
	for v := range p.seq {
		fn(i, v)
		i++
	}
}

// Reduce folds all elements into a single value of the same type.
// For reducing to a different type, use the top-level PipeReduce function.
func (p Pipeline[T]) Reduce(initial T, fn func(acc, item T) T) T {
	result := initial
	for v := range p.seq {
		result = fn(result, v)
	}
	return result
}

// First returns the first element and true, or zero value and false if empty.
// For infinite Pipelines, this returns immediately.
func (p Pipeline[T]) First() (T, bool) {
	for v := range p.seq {
		return v, true
	}
	var zero T
	return zero, false
}

// Last returns the last element and true, or zero value and false if empty.
// Warning: This consumes the entire Pipeline. Do not use on infinite sequences.
func (p Pipeline[T]) Last() (T, bool) {
	var last T
	found := false
	for v := range p.seq {
		last = v
		found = true
	}
	return last, found
}

// Find returns the first element matching the predicate.
// Short-circuits: stops iteration as soon as a match is found.
func (p Pipeline[T]) Find(pred func(T) bool) (T, bool) {
	for v := range p.seq {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if any element satisfies the predicate.
// Short-circuits on first match.
func (p Pipeline[T]) Any(pred func(T) bool) bool {
	for v := range p.seq {
		if pred(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements satisfy the predicate.
// Short-circuits on first non-match.
func (p Pipeline[T]) All(pred func(T) bool) bool {
	for v := range p.seq {
		if !pred(v) {
			return false
		}
	}
	return true
}

// None returns true if no elements satisfy the predicate.
func (p Pipeline[T]) None(pred func(T) bool) bool {
	return !p.Any(pred)
}

// Count returns the total number of elements.
// Warning: Consumes the entire Pipeline. Do not use on infinite sequences.
func (p Pipeline[T]) Count() int {
	n := 0
	for range p.seq {
		n++
	}
	return n
}

// CountBy returns the number of elements satisfying the predicate.
// Warning: Consumes the entire Pipeline. Do not use on infinite sequences.
func (p Pipeline[T]) CountBy(pred func(T) bool) int {
	n := 0
	for v := range p.seq {
		if pred(v) {
			n++
		}
	}
	return n
}

// IsEmpty returns true if the Pipeline has no elements.
func (p Pipeline[T]) IsEmpty() bool {
	for range p.seq {
		return false
	}
	return true
}

// Contains returns true if any element satisfies the predicate.
func (p Pipeline[T]) Contains(pred func(T) bool) bool {
	return p.Any(pred)
}

// MinBy returns the minimum element according to the comparison function.
// Warning: Consumes the entire Pipeline.
func (p Pipeline[T]) MinBy(less func(a, b T) bool) (T, bool) {
	var min T
	found := false
	for v := range p.seq {
		if !found || less(v, min) {
			min = v
			found = true
		}
	}
	return min, found
}

// MaxBy returns the maximum element according to the comparison function.
// Warning: Consumes the entire Pipeline.
func (p Pipeline[T]) MaxBy(less func(a, b T) bool) (T, bool) {
	var max T
	found := false
	for v := range p.seq {
		if !found || less(max, v) {
			max = v
			found = true
		}
	}
	return max, found
}

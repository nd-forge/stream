// Package stream provides a chainable collection operations library for Go.
//
// Stream[T] wraps a slice and provides method chaining for filter, sort,
// take, skip, and other operations that preserve the element type.
// For type-changing operations (Map, FlatMap, Zip), use top-level functions.
//
// Usage:
//
//	// Method chaining for same-type operations
//	result := stream.Of(1, 2, 3, 4, 5).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Sort(func(a, b int) int { return a - b }).
//	    ToSlice()
//
//	// Top-level functions for type-changing operations
//	names := stream.Map(
//	    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
//	    func(u User) string { return u.Name },
//	).ToSlice()
package stream

import (
	"math/rand"
	"sort"
)

// Stream is a lazy-evaluatable wrapper around a slice that supports method chaining.
type Stream[T any] struct {
	data []T
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// Of creates a new Stream from variadic arguments.
func Of[T any](items ...T) Stream[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return Stream[T]{data: copied}
}

// From creates a new Stream from an existing slice (copies the slice).
func From[T any](items []T) Stream[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return Stream[T]{data: copied}
}

// Generate creates a Stream of n elements using a generator function.
func Generate[T any](n int, gen func(index int) T) Stream[T] {
	data := make([]T, n)
	for i := 0; i < n; i++ {
		data[i] = gen(i)
	}
	return Stream[T]{data: data}
}

// Range creates a Stream of integers from start (inclusive) to end (exclusive).
func Range(start, end int) Stream[int] {
	if end <= start {
		return Stream[int]{data: nil}
	}
	data := make([]int, end-start)
	for i := range data {
		data[i] = start + i
	}
	return Stream[int]{data: data}
}

// ---------------------------------------------------------------------------
// Chainable methods (preserve type T → Stream[T])
// ---------------------------------------------------------------------------

// Filter returns a new Stream containing only elements that satisfy the predicate.
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	var result []T
	for _, v := range s.data {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return Stream[T]{data: result}
}

// Reject returns a new Stream excluding elements that satisfy the predicate.
func (s Stream[T]) Reject(predicate func(T) bool) Stream[T] {
	return s.Filter(func(v T) bool { return !predicate(v) })
}

// Sort returns a new sorted Stream using the provided comparison function.
// The cmp function should return negative if a < b, zero if a == b, positive if a > b.
func (s Stream[T]) Sort(cmp func(a, b T) int) Stream[T] {
	result := make([]T, len(s.data))
	copy(result, s.data)
	sort.Slice(result, func(i, j int) bool {
		return cmp(result[i], result[j]) < 0
	})
	return Stream[T]{data: result}
}

// Reverse returns a new Stream with elements in reverse order.
func (s Stream[T]) Reverse() Stream[T] {
	n := len(s.data)
	result := make([]T, n)
	for i, v := range s.data {
		result[n-1-i] = v
	}
	return Stream[T]{data: result}
}

// Take returns a new Stream with only the first n elements.
func (s Stream[T]) Take(n int) Stream[T] {
	if n >= len(s.data) {
		return s
	}
	if n <= 0 {
		return Stream[T]{data: nil}
	}
	result := make([]T, n)
	copy(result, s.data[:n])
	return Stream[T]{data: result}
}

// TakeLast returns a new Stream with only the last n elements.
func (s Stream[T]) TakeLast(n int) Stream[T] {
	if n >= len(s.data) {
		return s
	}
	if n <= 0 {
		return Stream[T]{data: nil}
	}
	result := make([]T, n)
	copy(result, s.data[len(s.data)-n:])
	return Stream[T]{data: result}
}

// Skip returns a new Stream with the first n elements removed.
func (s Stream[T]) Skip(n int) Stream[T] {
	if n >= len(s.data) {
		return Stream[T]{data: nil}
	}
	if n <= 0 {
		return s
	}
	result := make([]T, len(s.data)-n)
	copy(result, s.data[n:])
	return Stream[T]{data: result}
}

// TakeWhile returns elements from the start as long as the predicate is true.
func (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	var result []T
	for _, v := range s.data {
		if !predicate(v) {
			break
		}
		result = append(result, v)
	}
	return Stream[T]{data: result}
}

// DropWhile skips elements from the start as long as the predicate is true,
// then returns the rest.
func (s Stream[T]) DropWhile(predicate func(T) bool) Stream[T] {
	i := 0
	for ; i < len(s.data); i++ {
		if !predicate(s.data[i]) {
			break
		}
	}
	result := make([]T, len(s.data)-i)
	copy(result, s.data[i:])
	return Stream[T]{data: result}
}

// Distinct returns a new Stream with duplicate elements removed.
// Uses the provided key function to determine equality.
func (s Stream[T]) Distinct(key func(T) string) Stream[T] {
	seen := make(map[string]struct{})
	var result []T
	for _, v := range s.data {
		k := key(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, v)
		}
	}
	return Stream[T]{data: result}
}

// Shuffle returns a new Stream with elements in random order.
func (s Stream[T]) Shuffle() Stream[T] {
	result := make([]T, len(s.data))
	copy(result, s.data)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return Stream[T]{data: result}
}

// Peek executes a function for each element without modifying the Stream.
// Useful for debugging or logging within a chain.
func (s Stream[T]) Peek(fn func(T)) Stream[T] {
	for _, v := range s.data {
		fn(v)
	}
	return s
}

// ---------------------------------------------------------------------------
// Terminal operations (consume the Stream)
// ---------------------------------------------------------------------------

// ToSlice returns the underlying slice.
func (s Stream[T]) ToSlice() []T {
	if s.data == nil {
		return []T{}
	}
	return s.data
}

// ForEach executes a function for each element.
func (s Stream[T]) ForEach(fn func(T)) {
	for _, v := range s.data {
		fn(v)
	}
}

// ForEachIndexed executes a function for each element with its index.
func (s Stream[T]) ForEachIndexed(fn func(int, T)) {
	for i, v := range s.data {
		fn(i, v)
	}
}

// Reduce folds all elements into a single value of the same type.
// For reducing to a different type, use the top-level Reduce function.
func (s Stream[T]) Reduce(initial T, fn func(acc, item T) T) T {
	result := initial
	for _, v := range s.data {
		result = fn(result, v)
	}
	return result
}

// First returns the first element and true, or zero value and false if empty.
func (s Stream[T]) First() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[0], true
}

// Last returns the last element and true, or zero value and false if empty.
func (s Stream[T]) Last() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Find returns the first element matching the predicate.
func (s Stream[T]) Find(predicate func(T) bool) (T, bool) {
	for _, v := range s.data {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if any element satisfies the predicate.
func (s Stream[T]) Any(predicate func(T) bool) bool {
	for _, v := range s.data {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements satisfy the predicate.
func (s Stream[T]) All(predicate func(T) bool) bool {
	for _, v := range s.data {
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

// Count returns the number of elements in the Stream.
func (s Stream[T]) Count() int {
	return len(s.data)
}

// CountBy returns the number of elements that satisfy the predicate.
func (s Stream[T]) CountBy(predicate func(T) bool) int {
	count := 0
	for _, v := range s.data {
		if predicate(v) {
			count++
		}
	}
	return count
}

// IsEmpty returns true if the Stream has no elements.
func (s Stream[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// Contains returns true if the Stream contains an element matching the predicate.
func (s Stream[T]) Contains(predicate func(T) bool) bool {
	return s.Any(predicate)
}

// Partition splits the Stream into two: elements that match the predicate and those that don't.
func (s Stream[T]) Partition(predicate func(T) bool) (matched Stream[T], unmatched Stream[T]) {
	var yes, no []T
	for _, v := range s.data {
		if predicate(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return Stream[T]{data: yes}, Stream[T]{data: no}
}

// Chunk splits the Stream into chunks of the specified size.
func (s Stream[T]) Chunk(size int) []Stream[T] {
	if size <= 0 {
		return nil
	}
	var chunks []Stream[T]
	for i := 0; i < len(s.data); i += size {
		end := i + size
		if end > len(s.data) {
			end = len(s.data)
		}
		chunk := make([]T, end-i)
		copy(chunk, s.data[i:end])
		chunks = append(chunks, Stream[T]{data: chunk})
	}
	return chunks
}

// MinBy returns the minimum element according to the comparison function.
func (s Stream[T]) MinBy(less func(a, b T) bool) (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	min := s.data[0]
	for _, v := range s.data[1:] {
		if less(v, min) {
			min = v
		}
	}
	return min, true
}

// MaxBy returns the maximum element according to the comparison function.
func (s Stream[T]) MaxBy(less func(a, b T) bool) (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	max := s.data[0]
	for _, v := range s.data[1:] {
		if less(max, v) {
			max = v
		}
	}
	return max, true
}

package stream

import "iter"

// ---------------------------------------------------------------------------
// Top-level functions (type-changing operations: T → U)
// ---------------------------------------------------------------------------
// These cannot be methods because Go does not allow methods to introduce
// new type parameters. Use them to bridge between different Stream types.

// Map lazily transforms each element of type T into type U.
//
//	names := stream.Map(
//	    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
//	    func(u User) string { return u.Name },
//	).ToSlice()
func Map[T, U any](s Stream[T], fn func(T) U) Stream[U] {
	seq := s.seq
	return Stream[U]{seq: func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}}
}

// MapIndexed is like Map but also provides the index to the transform function.
func MapIndexed[T, U any](s Stream[T], fn func(int, T) U) Stream[U] {
	seq := s.seq
	return Stream[U]{seq: func(yield func(U) bool) {
		i := 0
		for v := range seq {
			if !yield(fn(i, v)) {
				return
			}
			i++
		}
	}}
}

// FlatMap lazily transforms each element into a slice and flattens the result.
//
//	allOrders := stream.FlatMap(
//	    stream.Of(users...),
//	    func(u User) []Order { return u.Orders },
//	)
func FlatMap[T, U any](s Stream[T], fn func(T) []U) Stream[U] {
	seq := s.seq
	return Stream[U]{seq: func(yield func(U) bool) {
		for v := range seq {
			for _, u := range fn(v) {
				if !yield(u) {
					return
				}
			}
		}
	}}
}

// Reduce folds all elements into a value of a different type.
//
//	total := stream.Reduce(orders, 0.0, func(acc float64, o Order) float64 {
//	    return acc + o.Amount
//	})
func Reduce[T, U any](s Stream[T], initial U, fn func(acc U, item T) U) U {
	result := initial
	for v := range s.seq {
		result = fn(result, v)
	}
	return result
}

// GroupBy groups elements by a key function and returns a map of key → slice.
//
//	bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })
//	// bySymbol["AUDUSD"] → []Trade
func GroupBy[T any, K comparable](s Stream[T], key func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for v := range s.seq {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

// Associate creates a map from Stream elements using a key-value function.
//
//	userMap := stream.Associate(users, func(u User) (int, string) {
//	    return u.ID, u.Name
//	})
func Associate[T any, K comparable, V any](s Stream[T], fn func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for v := range s.seq {
		k, val := fn(v)
		result[k] = val
	}
	return result
}

// Zip lazily combines two Streams into a Stream of pairs.
// Stops when either Stream is exhausted.
//
//	pairs := stream.Zip(names, scores)
func Zip[T, U any](s1 Stream[T], s2 Stream[U]) Stream[Pair[T, U]] {
	seq1, seq2 := s1.seq, s2.seq
	return Stream[Pair[T, U]]{seq: func(yield func(Pair[T, U]) bool) {
		next, stop := iter.Pull(seq2)
		defer stop()
		for v1 := range seq1 {
			v2, ok := next()
			if !ok {
				return
			}
			if !yield(Pair[T, U]{First: v1, Second: v2}) {
				return
			}
		}
	}}
}

// Pair holds two values of potentially different types.
type Pair[T, U any] struct {
	First  T
	Second U
}

// Flatten lazily flattens a Stream of slices into a flat Stream.
//
//	flat := stream.Flatten(stream.Of([]int{1, 2}, []int{3, 4}))
//	// yields 1, 2, 3, 4
func Flatten[T any](s Stream[[]T]) Stream[T] {
	seq := s.seq
	return Stream[T]{seq: func(yield func(T) bool) {
		for slice := range seq {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// ToMap collects a Stream of Pairs into a map.
func ToMap[K comparable, V any](s Stream[Pair[K, V]]) map[K]V {
	result := make(map[K]V)
	for pair := range s.seq {
		result[pair.First] = pair.Second
	}
	return result
}

// Enumerate wraps each element with its index as a Pair[int, T].
//
//	stream.Enumerate(stream.Of("a", "b", "c"))
//	// yields {0, "a"}, {1, "b"}, {2, "c"}
func Enumerate[T any](s Stream[T]) Stream[Pair[int, T]] {
	seq := s.seq
	return Stream[Pair[int, T]]{seq: func(yield func(Pair[int, T]) bool) {
		i := 0
		for v := range seq {
			if !yield(Pair[int, T]{First: i, Second: v}) {
				return
			}
			i++
		}
	}}
}

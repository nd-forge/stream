package stream

import "iter"

// ---------------------------------------------------------------------------
// Top-level functions for Pipeline (type-changing operations: T → U)
// ---------------------------------------------------------------------------
// These mirror the Stream transform functions but operate on Pipelines lazily.
// Prefixed with "Pipe" to distinguish from the eager Stream versions.

// PipeMap lazily transforms each element of type T into type U.
//
//	names := stream.PipeMap(
//	    stream.Lazy(users...).Filter(func(u User) bool { return u.IsActive }),
//	    func(u User) string { return u.Name },
//	).ToSlice()
func PipeMap[T, U any](p Pipeline[T], fn func(T) U) Pipeline[U] {
	seq := p.seq
	return Pipeline[U]{seq: func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}}
}

// PipeMapIndexed is like PipeMap but also provides the index.
func PipeMapIndexed[T, U any](p Pipeline[T], fn func(int, T) U) Pipeline[U] {
	seq := p.seq
	return Pipeline[U]{seq: func(yield func(U) bool) {
		i := 0
		for v := range seq {
			if !yield(fn(i, v)) {
				return
			}
			i++
		}
	}}
}

// PipeFlatMap lazily transforms each element into a slice and flattens the result.
//
//	allOrders := stream.PipeFlatMap(
//	    stream.Lazy(users...),
//	    func(u User) []Order { return u.Orders },
//	)
func PipeFlatMap[T, U any](p Pipeline[T], fn func(T) []U) Pipeline[U] {
	seq := p.seq
	return Pipeline[U]{seq: func(yield func(U) bool) {
		for v := range seq {
			for _, u := range fn(v) {
				if !yield(u) {
					return
				}
			}
		}
	}}
}

// PipeReduce folds all Pipeline elements into a value of a different type.
//
//	total := stream.PipeReduce(orders, 0.0, func(acc float64, o Order) float64 {
//	    return acc + o.Amount
//	})
func PipeReduce[T, U any](p Pipeline[T], initial U, fn func(acc U, item T) U) U {
	result := initial
	for v := range p.seq {
		result = fn(result, v)
	}
	return result
}

// PipeGroupBy groups Pipeline elements by a key function.
// This is a terminal operation that consumes the Pipeline.
//
//	groups := stream.PipeGroupBy(products, func(p Product) string { return p.Category })
func PipeGroupBy[T any, K comparable](p Pipeline[T], key func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for v := range p.seq {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

// PipeAssociate creates a map from Pipeline elements using a key-value function.
// This is a terminal operation that consumes the Pipeline.
func PipeAssociate[T any, K comparable, V any](p Pipeline[T], fn func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for v := range p.seq {
		k, val := fn(v)
		result[k] = val
	}
	return result
}

// PipeZip lazily combines two Pipelines into a Pipeline of pairs.
// Stops when either Pipeline is exhausted.
//
//	pairs := stream.PipeZip(names, scores)
func PipeZip[T, U any](p1 Pipeline[T], p2 Pipeline[U]) Pipeline[Pair[T, U]] {
	seq1, seq2 := p1.seq, p2.seq
	return Pipeline[Pair[T, U]]{seq: func(yield func(Pair[T, U]) bool) {
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

// PipeFlatten lazily flattens a Pipeline of slices into a flat Pipeline.
//
//	flat := stream.PipeFlatten(stream.Lazy([]int{1, 2}, []int{3, 4}))
//	// yields 1, 2, 3, 4
func PipeFlatten[T any](p Pipeline[[]T]) Pipeline[T] {
	seq := p.seq
	return Pipeline[T]{seq: func(yield func(T) bool) {
		for slice := range seq {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

// PipeToMap collects a Pipeline of Pairs into a map.
// This is a terminal operation that consumes the Pipeline.
func PipeToMap[K comparable, V any](p Pipeline[Pair[K, V]]) map[K]V {
	result := make(map[K]V)
	for pair := range p.seq {
		result[pair.First] = pair.Second
	}
	return result
}

// PipeEnumerate wraps each element with its index as a Pair[int, T].
//
//	stream.PipeEnumerate(stream.Lazy("a", "b", "c"))
//	// yields {0, "a"}, {1, "b"}, {2, "c"}
func PipeEnumerate[T any](p Pipeline[T]) Pipeline[Pair[int, T]] {
	seq := p.seq
	return Pipeline[Pair[int, T]]{seq: func(yield func(Pair[int, T]) bool) {
		i := 0
		for v := range seq {
			if !yield(Pair[int, T]{First: i, Second: v}) {
				return
			}
			i++
		}
	}}
}

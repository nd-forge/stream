package stream

// ---------------------------------------------------------------------------
// Top-level functions (type-changing operations: T → U)
// ---------------------------------------------------------------------------
// These cannot be methods because Go does not allow methods to introduce
// new type parameters. Use them to bridge between different Stream types.

// Map transforms each element of type T into type U and returns a new Stream[U].
//
//	names := stream.Map(
//	    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
//	    func(u User) string { return u.Name },
//	)
func Map[T, U any](s Stream[T], fn func(T) U) Stream[U] {
	result := make([]U, len(s.data))
	for i, v := range s.data {
		result[i] = fn(v)
	}
	return Stream[U]{data: result}
}

// MapIndexed is like Map but also provides the index to the transform function.
func MapIndexed[T, U any](s Stream[T], fn func(int, T) U) Stream[U] {
	result := make([]U, len(s.data))
	for i, v := range s.data {
		result[i] = fn(i, v)
	}
	return Stream[U]{data: result}
}

// FlatMap transforms each element into a slice and flattens the result.
//
//	allOrders := stream.FlatMap(
//	    stream.Of(users...),
//	    func(u User) []Order { return u.Orders },
//	)
func FlatMap[T, U any](s Stream[T], fn func(T) []U) Stream[U] {
	var result []U
	for _, v := range s.data {
		result = append(result, fn(v)...)
	}
	return Stream[U]{data: result}
}

// Reduce folds all elements into a value of a different type.
//
//	total := stream.Reduce(orders, 0.0, func(acc float64, o Order) float64 {
//	    return acc + o.Amount
//	})
func Reduce[T, U any](s Stream[T], initial U, fn func(acc U, item T) U) U {
	result := initial
	for _, v := range s.data {
		result = fn(result, v)
	}
	return result
}

// GroupBy groups elements by a key function and returns a map of key → Stream.
//
//	bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })
//	// bySymbol["AUDUSD"] → Stream[Trade]
func GroupBy[T any, K comparable](s Stream[T], key func(T) K) map[K]Stream[T] {
	groups := make(map[K][]T)
	for _, v := range s.data {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	result := make(map[K]Stream[T], len(groups))
	for k, v := range groups {
		result[k] = Stream[T]{data: v}
	}
	return result
}

// Associate creates a map from Stream elements using a key-value function.
//
//	userMap := stream.Associate(users, func(u User) (int, string) {
//	    return u.ID, u.Name
//	})
func Associate[T any, K comparable, V any](s Stream[T], fn func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(s.data))
	for _, v := range s.data {
		k, val := fn(v)
		result[k] = val
	}
	return result
}

// Zip combines two Streams into a Stream of pairs.
//
//	pairs := stream.Zip(dates, prices)
//	// Stream[Pair[time.Time, float64]]
func Zip[T, U any](s1 Stream[T], s2 Stream[U]) Stream[Pair[T, U]] {
	n := len(s1.data)
	if len(s2.data) < n {
		n = len(s2.data)
	}
	result := make([]Pair[T, U], n)
	for i := 0; i < n; i++ {
		result[i] = Pair[T, U]{First: s1.data[i], Second: s2.data[i]}
	}
	return Stream[Pair[T, U]]{data: result}
}

// Pair holds two values of potentially different types.
type Pair[T, U any] struct {
	First  T
	Second U
}

// Flatten converts a Stream of slices into a flat Stream.
//
//	flat := stream.Flatten(stream.Of([]int{1, 2}, []int{3, 4}))
//	// Stream[int]{1, 2, 3, 4}
func Flatten[T any](s Stream[[]T]) Stream[T] {
	var result []T
	for _, slice := range s.data {
		result = append(result, slice...)
	}
	return Stream[T]{data: result}
}

// ToMap converts a Stream of Pairs into a map.
func ToMap[K comparable, V any](s Stream[Pair[K, V]]) map[K]V {
	result := make(map[K]V, len(s.data))
	for _, p := range s.data {
		result[p.First] = p.Second
	}
	return result
}

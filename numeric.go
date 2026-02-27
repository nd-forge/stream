package stream

// Number is a constraint for numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum returns the sum of all elements in a numeric Stream.
func Sum[T Number](s Stream[T]) T {
	var total T
	for v := range s.seq {
		total += v
	}
	return total
}

// Avg returns the average of all elements in a numeric Stream.
func Avg[T Number](s Stream[T]) float64 {
	var total float64
	count := 0
	for v := range s.seq {
		total += float64(v)
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// Min returns the minimum element in a numeric Stream.
func Min[T Number](s Stream[T]) (T, bool) {
	var min T
	found := false
	for v := range s.seq {
		if !found || v < min {
			min = v
			found = true
		}
	}
	return min, found
}

// Max returns the maximum element in a numeric Stream.
func Max[T Number](s Stream[T]) (T, bool) {
	var max T
	found := false
	for v := range s.seq {
		if !found || v > max {
			max = v
			found = true
		}
	}
	return max, found
}

// SumBy extracts a numeric value from each element and returns the sum.
func SumBy[T any, N Number](s Stream[T], fn func(T) N) N {
	var total N
	for v := range s.seq {
		total += fn(v)
	}
	return total
}

// AvgBy extracts a numeric value from each element and returns the average.
func AvgBy[T any, N Number](s Stream[T], fn func(T) N) float64 {
	var total float64
	count := 0
	for v := range s.seq {
		total += float64(fn(v))
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

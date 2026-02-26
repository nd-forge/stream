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
	for _, v := range s.data {
		total += v
	}
	return total
}

// Avg returns the average of all elements in a numeric Stream.
func Avg[T Number](s Stream[T]) float64 {
	if len(s.data) == 0 {
		return 0
	}
	var total float64
	for _, v := range s.data {
		total += float64(v)
	}
	return total / float64(len(s.data))
}

// Min returns the minimum element in a numeric Stream.
func Min[T Number](s Stream[T]) (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	min := s.data[0]
	for _, v := range s.data[1:] {
		if v < min {
			min = v
		}
	}
	return min, true
}

// Max returns the maximum element in a numeric Stream.
func Max[T Number](s Stream[T]) (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	max := s.data[0]
	for _, v := range s.data[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

// SumBy extracts a numeric value from each element and returns the sum.
func SumBy[T any, N Number](s Stream[T], fn func(T) N) N {
	var total N
	for _, v := range s.data {
		total += fn(v)
	}
	return total
}

// AvgBy extracts a numeric value from each element and returns the average.
func AvgBy[T any, N Number](s Stream[T], fn func(T) N) float64 {
	if len(s.data) == 0 {
		return 0
	}
	var total float64
	for _, v := range s.data {
		total += float64(fn(v))
	}
	return total / float64(len(s.data))
}

package stream_test

import (
	"testing"

	"github.com/nd-forge/stream"
)

// ---------------------------------------------------------------------------
// Data setup
// ---------------------------------------------------------------------------

var benchData []int

func init() {
	benchData = make([]int, 10_000)
	for i := range benchData {
		benchData[i] = i
	}
}

// ---------------------------------------------------------------------------
// Filter benchmarks: Native vs Stream
// ---------------------------------------------------------------------------

func BenchmarkNativeFilter(b *testing.B) {
	for b.Loop() {
		var result []int
		for _, v := range benchData {
			if v%2 == 0 {
				result = append(result, v)
			}
		}
		_ = result
	}
}

func BenchmarkStreamFilter(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Filter + Take benchmarks: Shows lazy evaluation advantage
// ---------------------------------------------------------------------------

func BenchmarkNativeFilterTake(b *testing.B) {
	for b.Loop() {
		var result []int
		for _, v := range benchData {
			if v%2 == 0 {
				result = append(result, v)
				if len(result) >= 10 {
					break
				}
			}
		}
		_ = result
	}
}

func BenchmarkStreamFilterTake(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Filter(func(n int) bool { return n%2 == 0 }).Take(10).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Map + Filter benchmarks
// ---------------------------------------------------------------------------

func BenchmarkNativeMapFilter(b *testing.B) {
	for b.Loop() {
		var result []int
		for _, v := range benchData {
			doubled := v * 2
			if doubled < 100 {
				result = append(result, doubled)
			}
		}
		_ = result
	}
}

func BenchmarkStreamMapFilter(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = stream.Map(s, func(n int) int { return n * 2 }).
			Filter(func(n int) bool { return n < 100 }).
			ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Sort benchmarks
// ---------------------------------------------------------------------------

func BenchmarkNativeSort(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = 1000 - i
	}
	for b.Loop() {
		cp := make([]int, len(data))
		copy(cp, data)
		for i := 1; i < len(cp); i++ {
			for j := i; j > 0 && cp[j] < cp[j-1]; j-- {
				cp[j], cp[j-1] = cp[j-1], cp[j]
			}
		}
		_ = cp
	}
}

func BenchmarkStreamSort(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = 1000 - i
	}
	s := stream.From(data)
	for b.Loop() {
		_ = s.Sort(func(a, b int) int { return a - b }).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Chained operations benchmark
// ---------------------------------------------------------------------------

func BenchmarkNativeChained(b *testing.B) {
	for b.Loop() {
		// Filter -> Map -> Take 5
		var result []string
		for _, v := range benchData {
			if v%3 == 0 {
				s := string(rune('A' + v%26))
				result = append(result, s)
				if len(result) >= 5 {
					break
				}
			}
		}
		_ = result
	}
}

func BenchmarkStreamChained(b *testing.B) {
	s := stream.Map(
		stream.From(benchData).
			Filter(func(n int) bool { return n%3 == 0 }).
			Take(5),
		func(n int) string { return string(rune('A' + n%26)) },
	)
	for b.Loop() {
		_ = s.ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Reduce benchmarks
// ---------------------------------------------------------------------------

func BenchmarkNativeReduce(b *testing.B) {
	for b.Loop() {
		sum := 0
		for _, v := range benchData {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkStreamReduce(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Reduce(0, func(acc, v int) int { return acc + v })
	}
}

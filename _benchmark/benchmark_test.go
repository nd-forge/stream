package benchmark_test

import (
	"testing"

	"github.com/nd-forge/stream"
	"github.com/samber/lo"
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
// Filter benchmarks: Native vs lo vs Stream
// ---------------------------------------------------------------------------

func BenchmarkFilter_Native(b *testing.B) {
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

func BenchmarkFilter_Lo(b *testing.B) {
	for b.Loop() {
		_ = lo.Filter(benchData, func(v int, _ int) bool { return v%2 == 0 })
	}
}

func BenchmarkFilter_Stream(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Filter + Take benchmarks: Shows lazy evaluation advantage
// lo cannot short-circuit — it filters all then takes.
// ---------------------------------------------------------------------------

func BenchmarkFilterTake_Native(b *testing.B) {
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

func BenchmarkFilterTake_Lo(b *testing.B) {
	for b.Loop() {
		_ = lo.Subset(
			lo.Filter(benchData, func(v int, _ int) bool { return v%2 == 0 }),
			0, 10,
		)
	}
}

func BenchmarkFilterTake_Stream(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Filter(func(n int) bool { return n%2 == 0 }).Take(10).ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Map + Filter benchmarks
// ---------------------------------------------------------------------------

func BenchmarkMapFilter_Native(b *testing.B) {
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

func BenchmarkMapFilter_Lo(b *testing.B) {
	for b.Loop() {
		mapped := lo.Map(benchData, func(v int, _ int) int { return v * 2 })
		_ = lo.Filter(mapped, func(v int, _ int) bool { return v < 100 })
	}
}

func BenchmarkMapFilter_Stream(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = stream.Map(s, func(n int) int { return n * 2 }).
			Filter(func(n int) bool { return n < 100 }).
			ToSlice()
	}
}

// ---------------------------------------------------------------------------
// Reduce benchmarks
// ---------------------------------------------------------------------------

func BenchmarkReduce_Native(b *testing.B) {
	for b.Loop() {
		sum := 0
		for _, v := range benchData {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkReduce_Lo(b *testing.B) {
	for b.Loop() {
		_ = lo.Reduce(benchData, func(acc int, v int, _ int) int { return acc + v }, 0)
	}
}

func BenchmarkReduce_Stream(b *testing.B) {
	s := stream.From(benchData)
	for b.Loop() {
		_ = s.Reduce(0, func(acc, v int) int { return acc + v })
	}
}

// ---------------------------------------------------------------------------
// Chained operations benchmark: Filter -> Map -> Take 5
// ---------------------------------------------------------------------------

func BenchmarkChained_Native(b *testing.B) {
	for b.Loop() {
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

func BenchmarkChained_Lo(b *testing.B) {
	for b.Loop() {
		filtered := lo.Filter(benchData, func(v int, _ int) bool { return v%3 == 0 })
		mapped := lo.Map(filtered, func(v int, _ int) string { return string(rune('A' + v%26)) })
		_ = lo.Subset(mapped, 0, 5)
	}
}

func BenchmarkChained_Stream(b *testing.B) {
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

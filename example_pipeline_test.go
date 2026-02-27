package stream_test

import (
	"fmt"
	"strings"

	"github.com/nd-forge/stream"
)

// ---------------------------------------------------------------------------
// Pipeline constructors
// ---------------------------------------------------------------------------

func ExampleLazy() {
	result := stream.Lazy(1, 2, 3, 4, 5).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	fmt.Println(result)
	// Output: [2 4]
}

func ExampleLazyFrom() {
	seq := stream.Of(10, 20, 30).Iter()
	result := stream.LazyFrom(seq).
		Filter(func(n int) bool { return n > 10 }).
		ToSlice()
	fmt.Println(result)
	// Output: [20 30]
}

func ExampleStream_Lazy() {
	result := stream.Of(1, 2, 3, 4, 5).
		Lazy().
		Filter(func(n int) bool { return n > 3 }).
		ToSlice()
	fmt.Println(result)
	// Output: [4 5]
}

func ExampleLazyRange() {
	result := stream.LazyRange(0, 5).ToSlice()
	fmt.Println(result)
	// Output: [0 1 2 3 4]
}

// ---------------------------------------------------------------------------
// Generators (infinite sequences)
// ---------------------------------------------------------------------------

func ExampleRepeat() {
	result := stream.Repeat(42).Take(3).ToSlice()
	fmt.Println(result)
	// Output: [42 42 42]
}

func ExampleRepeatN() {
	result := stream.RepeatN("-", 4).ToSlice()
	fmt.Println(result)
	// Output: [- - - -]
}

func ExampleIterate() {
	// Powers of 2: 1, 2, 4, 8, 16
	result := stream.Iterate(1, func(n int) int { return n * 2 }).
		Take(5).ToSlice()
	fmt.Println(result)
	// Output: [1 2 4 8 16]
}

func ExampleNaturals() {
	result := stream.Naturals().Take(5).ToSlice()
	fmt.Println(result)
	// Output: [0 1 2 3 4]
}

// ---------------------------------------------------------------------------
// Pipeline chainable methods
// ---------------------------------------------------------------------------

func ExamplePipeline_Filter() {
	result := stream.Lazy(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	fmt.Println(result)
	// Output: [2 4 6]
}

func ExamplePipeline_Reject() {
	result := stream.Lazy(1, 2, 3, 4, 5).
		Reject(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 3 5]
}

func ExamplePipeline_Take() {
	result := stream.Lazy(10, 20, 30, 40, 50).Take(3).ToSlice()
	fmt.Println(result)
	// Output: [10 20 30]
}

func ExamplePipeline_Skip() {
	result := stream.Lazy(10, 20, 30, 40, 50).Skip(2).ToSlice()
	fmt.Println(result)
	// Output: [30 40 50]
}

func ExamplePipeline_TakeWhile() {
	result := stream.Lazy(1, 2, 3, 4, 5).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 2 3]
}

func ExamplePipeline_DropWhile() {
	result := stream.Lazy(1, 2, 3, 4, 5).
		DropWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	fmt.Println(result)
	// Output: [4 5]
}

func ExamplePipeline_Distinct() {
	result := stream.Lazy("a", "b", "a", "c", "b").
		Distinct(func(s string) string { return s }).
		ToSlice()
	fmt.Println(result)
	// Output: [a b c]
}

func ExamplePipeline_Sort() {
	result := stream.Lazy(3, 1, 4, 1, 5).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 1 3 4 5]
}

func ExamplePipeline_Reverse() {
	result := stream.Lazy(1, 2, 3).Reverse().ToSlice()
	fmt.Println(result)
	// Output: [3 2 1]
}

func ExamplePipeline_Chain() {
	p1 := stream.Lazy(1, 2, 3)
	p2 := stream.Lazy(4, 5, 6)
	result := p1.Chain(p2).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4 5 6]
}

func ExamplePipeline_Peek() {
	var log []int
	stream.Lazy(1, 2, 3).
		Peek(func(n int) { log = append(log, n) }).
		ToSlice()
	fmt.Println(log)
	// Output: [1 2 3]
}

// ---------------------------------------------------------------------------
// Pipeline terminal operations
// ---------------------------------------------------------------------------

func ExamplePipeline_First() {
	v, ok := stream.Lazy(10, 20, 30).First()
	fmt.Println(v, ok)
	// Output: 10 true
}

func ExamplePipeline_Last() {
	v, ok := stream.Lazy(10, 20, 30).Last()
	fmt.Println(v, ok)
	// Output: 30 true
}

func ExamplePipeline_Find() {
	v, ok := stream.Lazy(1, 2, 3, 4, 5).
		Find(func(n int) bool { return n > 3 })
	fmt.Println(v, ok)
	// Output: 4 true
}

func ExamplePipeline_Reduce() {
	sum := stream.Lazy(1, 2, 3, 4, 5).
		Reduce(0, func(acc, v int) int { return acc + v })
	fmt.Println(sum)
	// Output: 15
}

func ExamplePipeline_Count() {
	n := stream.Lazy(1, 2, 3, 4, 5).
		Filter(func(v int) bool { return v > 2 }).
		Count()
	fmt.Println(n)
	// Output: 3
}

func ExamplePipeline_ForEach() {
	stream.Lazy(1, 2, 3).ForEach(func(v int) {
		fmt.Printf("%d ", v)
	})
	fmt.Println()
	// Output: 1 2 3
}

func ExamplePipeline_ForEachIndexed() {
	stream.Lazy("a", "b", "c").ForEachIndexed(func(i int, s string) {
		fmt.Printf("%d:%s ", i, s)
	})
	fmt.Println()
	// Output: 0:a 1:b 2:c
}

func ExamplePipeline_Contains() {
	has := stream.Lazy(1, 2, 3).Contains(func(n int) bool { return n == 2 })
	fmt.Println(has)
	// Output: true
}

func ExamplePipeline_ToStream() {
	s := stream.Lazy(3, 1, 2).
		Filter(func(n int) bool { return n > 1 }).
		ToStream()
	fmt.Println(s.Count())
	// Output: 2
}

// ---------------------------------------------------------------------------
// Pipeline transform functions (type-changing)
// ---------------------------------------------------------------------------

func ExamplePipeMap() {
	result := stream.PipeMap(
		stream.Lazy(1, 2, 3),
		func(n int) string { return fmt.Sprintf("item_%d", n) },
	).ToSlice()
	fmt.Println(result)
	// Output: [item_1 item_2 item_3]
}

func ExamplePipeMapIndexed() {
	result := stream.PipeMapIndexed(
		stream.Lazy("a", "b", "c"),
		func(i int, s string) string { return fmt.Sprintf("%d:%s", i, s) },
	).ToSlice()
	fmt.Println(result)
	// Output: [0:a 1:b 2:c]
}

func ExamplePipeFlatMap() {
	result := stream.PipeFlatMap(
		stream.Lazy([]int{1, 2}, []int{3, 4}),
		func(s []int) []int { return s },
	).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4]
}

func ExamplePipeReduce() {
	sum := stream.PipeReduce(
		stream.Lazy(1, 2, 3, 4, 5),
		0,
		func(acc, v int) int { return acc + v },
	)
	fmt.Println(sum)
	// Output: 15
}

func ExamplePipeZip() {
	pairs := stream.PipeZip(
		stream.Lazy("a", "b", "c"),
		stream.Lazy(1, 2, 3),
	).ToSlice()
	for _, p := range pairs {
		fmt.Printf("%s=%d ", p.First, p.Second)
	}
	fmt.Println()
	// Output: a=1 b=2 c=3
}

func ExamplePipeFlatten() {
	result := stream.PipeFlatten(
		stream.Lazy([]int{1, 2}, []int{3, 4}, []int{5}),
	).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4 5]
}

func ExamplePipeEnumerate() {
	result := stream.PipeEnumerate(stream.Lazy("a", "b", "c")).ToSlice()
	for _, p := range result {
		fmt.Printf("%d:%s ", p.First, p.Second)
	}
	fmt.Println()
	// Output: 0:a 1:b 2:c
}

// ---------------------------------------------------------------------------
// Pipeline showcase examples
// ---------------------------------------------------------------------------

func Example_lazyFilterTake() {
	// Lazy pipeline: only evaluates elements as needed
	result := stream.Naturals().
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(5).
		ToSlice()
	fmt.Println(result)
	// Output: [0 2 4 6 8]
}

func Example_fibonacci() {
	// Fibonacci sequence using Iterate with Pair
	fib := stream.PipeMap(
		stream.Iterate(
			stream.Pair[int, int]{First: 0, Second: 1},
			func(p stream.Pair[int, int]) stream.Pair[int, int] {
				return stream.Pair[int, int]{First: p.Second, Second: p.First + p.Second}
			},
		).Take(10),
		func(p stream.Pair[int, int]) int { return p.First },
	).ToSlice()
	fmt.Println(fib)
	// Output: [0 1 1 2 3 5 8 13 21 34]
}

func Example_lazyToEager() {
	// Switch between lazy Pipeline and eager Stream
	result := stream.Naturals().
		Filter(func(n int) bool { return n%3 == 0 }).
		Take(6).
		ToStream().  // Pipeline -> Stream
		Reverse().   // Stream-only operation
		ToSlice()
	fmt.Println(result)
	// Output: [15 12 9 6 3 0]
}

func Example_pipelineTextProcessing() {
	result := stream.PipeMap(
		stream.Lazy("hello", "world", "hello", "go", "stream", "go").
			Distinct(func(s string) string { return s }).
			Sort(func(a, b string) int { return len(a) - len(b) }),
		strings.ToUpper,
	).ToSlice()
	fmt.Println(result)
	// Output: [GO HELLO WORLD STREAM]
}

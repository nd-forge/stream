package stream_test

import (
	"fmt"
	"maps"
	"strings"

	"github.com/nd-forge/stream"
)

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

func ExampleOf() {
	result := stream.Of(5, 2, 8, 1, 9).
		Filter(func(n int) bool { return n%2 == 0 }).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()
	fmt.Println(result)
	// Output: [2 8]
}

func ExampleFrom() {
	data := []string{"Go", "Rust", "Python"}
	result := stream.From(data).
		Filter(func(s string) bool { return len(s) <= 4 }).
		ToSlice()
	fmt.Println(result)
	// Output: [Go Rust]
}

func ExampleRange() {
	result := stream.Range(1, 6).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4 5]
}

func ExampleGenerate() {
	squares := stream.Generate(5, func(i int) int { return i * i }).ToSlice()
	fmt.Println(squares)
	// Output: [0 1 4 9 16]
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
// Chainable methods
// ---------------------------------------------------------------------------

func ExampleStream_Filter() {
	result := stream.Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	fmt.Println(result)
	// Output: [2 4 6]
}

func ExampleStream_Reject() {
	result := stream.Of(1, 2, 3, 4, 5).
		Reject(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 3 5]
}

func ExampleStream_Sort() {
	result := stream.Of(3, 1, 4, 1, 5).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 1 3 4 5]
}

func ExampleStream_Reverse() {
	result := stream.Of(1, 2, 3).Reverse().ToSlice()
	fmt.Println(result)
	// Output: [3 2 1]
}

func ExampleStream_Take() {
	result := stream.Of(10, 20, 30, 40, 50).Take(3).ToSlice()
	fmt.Println(result)
	// Output: [10 20 30]
}

func ExampleStream_TakeLast() {
	result := stream.Of(10, 20, 30, 40, 50).TakeLast(2).ToSlice()
	fmt.Println(result)
	// Output: [40 50]
}

func ExampleStream_Skip() {
	result := stream.Of(10, 20, 30, 40, 50).Skip(2).ToSlice()
	fmt.Println(result)
	// Output: [30 40 50]
}

func ExampleStream_TakeWhile() {
	result := stream.Of(1, 2, 3, 4, 5).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleStream_DropWhile() {
	result := stream.Of(1, 2, 3, 4, 5).
		DropWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	fmt.Println(result)
	// Output: [4 5]
}

func ExampleStream_Distinct() {
	result := stream.Of("a", "b", "a", "c", "b").
		Distinct(func(s string) string { return s }).
		ToSlice()
	fmt.Println(result)
	// Output: [a b c]
}

func ExampleStream_Peek() {
	var log []string
	stream.Of("hello", "world").
		Peek(func(s string) { log = append(log, "saw:"+s) }).
		ToSlice()
	fmt.Println(log)
	// Output: [saw:hello saw:world]
}

func ExampleStream_Shuffle() {
	s := stream.Of(1, 2, 3, 4, 5).Shuffle()
	// Shuffle randomizes order, so just check count
	fmt.Println(s.Count())
	// Output: 5
}

func ExampleStream_Chunk() {
	chunks := stream.Of(1, 2, 3, 4, 5).Chunk(2)
	for _, c := range chunks {
		fmt.Println(c.ToSlice())
	}
	// Output:
	// [1 2]
	// [3 4]
	// [5]
}

func ExampleStream_Partition() {
	evens, odds := stream.Of(1, 2, 3, 4, 5).
		Partition(func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens.ToSlice())
	fmt.Println("odds:", odds.ToSlice())
	// Output:
	// evens: [2 4]
	// odds: [1 3 5]
}

func ExampleStream_Chain() {
	s1 := stream.Of(1, 2, 3)
	s2 := stream.Of(4, 5, 6)
	result := s1.Chain(s2).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4 5 6]
}

// ---------------------------------------------------------------------------
// Terminal operations
// ---------------------------------------------------------------------------

func ExampleStream_ForEach() {
	stream.Of(1, 2, 3).ForEach(func(n int) {
		fmt.Printf("%d ", n)
	})
	fmt.Println()
	// Output: 1 2 3
}

func ExampleStream_ForEachIndexed() {
	stream.Of("a", "b", "c").ForEachIndexed(func(i int, s string) {
		fmt.Printf("%d:%s ", i, s)
	})
	fmt.Println()
	// Output: 0:a 1:b 2:c
}

func ExampleStream_Reduce() {
	sum := stream.Of(1, 2, 3, 4, 5).
		Reduce(0, func(acc, v int) int { return acc + v })
	fmt.Println(sum)
	// Output: 15
}

func ExampleStream_First() {
	v, ok := stream.Of(10, 20, 30).First()
	fmt.Println(v, ok)
	// Output: 10 true
}

func ExampleStream_Last() {
	v, ok := stream.Of(10, 20, 30).Last()
	fmt.Println(v, ok)
	// Output: 30 true
}

func ExampleStream_Find() {
	v, ok := stream.Of(1, 2, 3, 4, 5).
		Find(func(n int) bool { return n > 3 })
	fmt.Println(v, ok)
	// Output: 4 true
}

func ExampleStream_Count() {
	n := stream.Of(1, 2, 3).Count()
	fmt.Println(n)
	// Output: 3
}

func ExampleStream_CountBy() {
	n := stream.Of(1, 2, 3, 4, 5).
		CountBy(func(v int) bool { return v%2 == 0 })
	fmt.Println(n)
	// Output: 2
}

func ExampleStream_IsEmpty() {
	fmt.Println(stream.Of[int]().IsEmpty())
	fmt.Println(stream.Of(1).IsEmpty())
	// Output:
	// true
	// false
}

func ExampleStream_Contains() {
	has := stream.Of(1, 2, 3).Contains(func(n int) bool { return n == 2 })
	fmt.Println(has)
	// Output: true
}

func ExampleStream_Any() {
	fmt.Println(stream.Of(1, 2, 3).Any(func(n int) bool { return n > 2 }))
	// Output: true
}

func ExampleStream_All() {
	fmt.Println(stream.Of(2, 4, 6).All(func(n int) bool { return n%2 == 0 }))
	// Output: true
}

func ExampleStream_None() {
	fmt.Println(stream.Of(1, 2, 3).None(func(n int) bool { return n > 10 }))
	// Output: true
}

func ExampleStream_MinBy() {
	v, _ := stream.Of(3, 1, 4, 1, 5).
		MinBy(func(a, b int) bool { return a < b })
	fmt.Println(v)
	// Output: 1
}

func ExampleStream_MaxBy() {
	v, _ := stream.Of(3, 1, 4, 1, 5).
		MaxBy(func(a, b int) bool { return a < b })
	fmt.Println(v)
	// Output: 5
}

func ExampleStream_Seq() {
	var result []int
	for v := range stream.Of(1, 2, 3).Seq() {
		result = append(result, v)
	}
	fmt.Println(result)
	// Output: [1 2 3]
}

// ---------------------------------------------------------------------------
// Top-level transform functions
// ---------------------------------------------------------------------------

func ExampleMap() {
	result := stream.Map(
		stream.Of(1, 2, 3),
		func(n int) string { return fmt.Sprintf("item_%d", n) },
	).ToSlice()
	fmt.Println(result)
	// Output: [item_1 item_2 item_3]
}

func ExampleMapIndexed() {
	result := stream.MapIndexed(
		stream.Of("a", "b", "c"),
		func(i int, s string) string { return fmt.Sprintf("%d:%s", i, s) },
	).ToSlice()
	fmt.Println(result)
	// Output: [0:a 1:b 2:c]
}

func ExampleFlatMap() {
	result := stream.FlatMap(
		stream.Of([]int{1, 2}, []int{3, 4}),
		func(s []int) []int { return s },
	).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4]
}

func ExampleReduce() {
	sum := stream.Reduce(
		stream.Of(1, 2, 3, 4, 5),
		0,
		func(acc, v int) int { return acc + v },
	)
	fmt.Println(sum)
	// Output: 15
}

func ExampleZip() {
	pairs := stream.Zip(
		stream.Of("a", "b", "c"),
		stream.Of(1, 2, 3),
	).ToSlice()
	for _, p := range pairs {
		fmt.Printf("%s=%d ", p.First, p.Second)
	}
	fmt.Println()
	// Output: a=1 b=2 c=3
}

func ExampleFlatten() {
	result := stream.Flatten(
		stream.Of([]int{1, 2}, []int{3, 4}, []int{5}),
	).ToSlice()
	fmt.Println(result)
	// Output: [1 2 3 4 5]
}

func ExampleToMap() {
	m := stream.ToMap(stream.Zip(
		stream.Of("a", "b", "c"),
		stream.Of(1, 2, 3),
	))
	fmt.Println(m["a"], m["b"], m["c"])
	// Output: 1 2 3
}

func ExampleEnumerate() {
	result := stream.Enumerate(stream.Of("a", "b", "c")).ToSlice()
	for _, p := range result {
		fmt.Printf("%d:%s ", p.First, p.Second)
	}
	fmt.Println()
	// Output: 0:a 1:b 2:c
}

// ---------------------------------------------------------------------------
// Numeric functions
// ---------------------------------------------------------------------------

func ExampleSum() {
	fmt.Println(stream.Sum(stream.Of(1, 2, 3, 4, 5)))
	// Output: 15
}

func ExampleAvg() {
	fmt.Println(stream.Avg(stream.Of(10.0, 20.0, 30.0)))
	// Output: 20
}

func ExampleMin() {
	v, _ := stream.Min(stream.Of(3, 1, 4, 1, 5))
	fmt.Println(v)
	// Output: 1
}

func ExampleMax() {
	v, _ := stream.Max(stream.Of(3, 1, 4, 1, 5))
	fmt.Println(v)
	// Output: 5
}

// ---------------------------------------------------------------------------
// iter.Seq bridge
// ---------------------------------------------------------------------------

func ExampleCollect() {
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	s := stream.Collect(maps.Values(m)).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()
	fmt.Println(s)
	// Output: [1 2 3]
}

func ExampleCollect2() {
	m := map[string]int{"a": 1}
	pairs := stream.Collect2(maps.All(m)).ToSlice()
	fmt.Println(pairs[0].First, pairs[0].Second)
	// Output: a 1
}

// ---------------------------------------------------------------------------
// Showcase examples
// ---------------------------------------------------------------------------

func Example_chaining() {
	// Method chaining: filter even numbers, sort, take top 3
	result := stream.Of(9, 2, 7, 4, 1, 8, 3, 6, 5).
		Filter(func(n int) bool { return n%2 == 0 }).
		Sort(func(a, b int) int { return a - b }).
		Take(3).
		ToSlice()
	fmt.Println(result)
	// Output: [2 4 6]
}

func Example_textProcessing() {
	words := stream.Of("hello", "world", "hello", "go", "stream", "go").
		Distinct(func(s string) string { return s }).
		Sort(func(a, b string) int { return len(a) - len(b) })

	result := stream.Map(words, strings.ToUpper).ToSlice()
	fmt.Println(result)
	// Output: [GO HELLO WORLD STREAM]
}

func Example_lazyFilterTake() {
	// Lazy evaluation: only evaluates elements as needed
	result := stream.Naturals().
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(5).
		ToSlice()
	fmt.Println(result)
	// Output: [0 2 4 6 8]
}

func Example_fibonacci() {
	// Fibonacci sequence using Iterate with Pair
	fib := stream.Map(
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

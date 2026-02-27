package main

import (
	"fmt"
	"maps"
	"strings"

	"github.com/nd-forge/stream"
)

// --- Data types ---

type Product struct {
	Name     string
	Category string
	Price    float64
	InStock  bool
}

type User struct {
	Name     string
	Age      int
	IsActive bool
	Orders   []Order
}

type Order struct {
	Product  string
	Amount   float64
	Discount float64
}

func main() {
	// =====================================================================
	// Part 1: Stream (Eager Evaluation)
	// =====================================================================
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║           Part 1: Stream (Eager Evaluation)             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ===== Constructors =====
	fmt.Println("\n=== Constructors ===")

	nums := stream.Of(3, 1, 4, 1, 5, 9, 2, 6)
	fmt.Println("Of:", nums.ToSlice())

	r := stream.Range(1, 6)
	fmt.Println("Range(1,6):", r.ToSlice())

	squares := stream.Generate(5, func(i int) int { return i * i })
	fmt.Println("Generate squares:", squares.ToSlice())

	// ===== Filter, Sort, Take =====
	fmt.Println("\n=== Filter + Sort + Take ===")

	products := stream.Of(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
		Product{Name: "Sneakers", Category: "Clothing", Price: 120, InStock: false},
	)

	top3 := products.
		Filter(func(p Product) bool { return p.InStock }).
		Sort(func(a, b Product) int {
			if a.Price > b.Price {
				return -1
			}
			if a.Price < b.Price {
				return 1
			}
			return 0
		}).
		Take(3).
		ToSlice()

	fmt.Println("Top 3 in-stock products by price:")
	for _, p := range top3 {
		fmt.Printf("  %s: $%.0f\n", p.Name, p.Price)
	}

	// ===== Map (type-changing) =====
	fmt.Println("\n=== Map ===")

	categories := stream.Map(products, func(p Product) string {
		return p.Category
	}).Distinct(func(s string) string { return s }).ToSlice()

	fmt.Println("Unique categories:", categories)

	// ===== FlatMap =====
	fmt.Println("\n=== FlatMap ===")

	users := stream.Of(
		User{Name: "Alice", Age: 30, IsActive: true, Orders: []Order{
			{Product: "Laptop", Amount: 1200, Discount: 0.1},
			{Product: "Mouse", Amount: 50, Discount: 0},
		}},
		User{Name: "Bob", Age: 17, IsActive: true, Orders: []Order{
			{Product: "Keyboard", Amount: 80, Discount: 0.5},
		}},
		User{Name: "Charlie", Age: 25, IsActive: false, Orders: []Order{
			{Product: "Monitor", Amount: 500, Discount: 0.2},
		}},
	)

	allOrders := stream.FlatMap(
		users.Filter(func(u User) bool { return u.IsActive }),
		func(u User) []Order { return u.Orders },
	)
	fmt.Printf("Active user orders: %d\n", allOrders.Count())

	allOrders.ForEach(func(o Order) {
		fmt.Printf("  %s: $%.0f (%.0f%% off)\n", o.Product, o.Amount, o.Discount*100)
	})

	// ===== Reduce (type-changing) =====
	fmt.Println("\n=== Reduce ===")

	total := stream.Reduce(allOrders, 0.0, func(acc float64, o Order) float64 {
		return acc + o.Amount*(1-o.Discount)
	})
	fmt.Printf("Total after discount: $%.2f\n", total)

	// ===== GroupBy =====
	fmt.Println("\n=== GroupBy ===")

	byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })
	for category, group := range byCategory {
		totalPrice := stream.SumBy(group, func(p Product) float64 { return p.Price })
		fmt.Printf("  %s: total=$%.0f (%d products)\n", category, totalPrice, group.Count())
	}

	// ===== Associate =====
	fmt.Println("\n=== Associate ===")

	nameToAge := stream.Associate(users, func(u User) (string, int) {
		return u.Name, u.Age
	})
	fmt.Println("Name->Age:", nameToAge)

	// ===== Zip =====
	fmt.Println("\n=== Zip ===")

	names := stream.Of("Alice", "Bob", "Charlie")
	scores := stream.Of(85.0, 92.0, 78.0)

	stream.Zip(names, scores).ForEach(func(p stream.Pair[string, float64]) {
		fmt.Printf("  %s: %.0f\n", p.First, p.Second)
	})

	// ===== Partition =====
	fmt.Println("\n=== Partition ===")

	inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })
	fmt.Printf("In stock: %d, Out of stock: %d\n", inStock.Count(), outOfStock.Count())

	// ===== Chunk (batch processing) =====
	fmt.Println("\n=== Chunk ===")

	batches := stream.Range(0, 10).Chunk(3)
	for i, batch := range batches {
		fmt.Printf("  Batch %d: %v\n", i, batch.ToSlice())
	}

	// ===== Numeric operations =====
	fmt.Println("\n=== Numeric ===")

	values := stream.Of(10.0, 20.0, 30.0, 40.0, 50.0)
	fmt.Printf("Sum: %.0f\n", stream.Sum(values))
	fmt.Printf("Avg: %.1f\n", stream.Avg(values))

	min, _ := stream.Min(values)
	max, _ := stream.Max(values)
	fmt.Printf("Min: %.0f, Max: %.0f\n", min, max)

	avgPrice := stream.AvgBy(products, func(p Product) float64 { return p.Price })
	fmt.Printf("Avg product price: $%.2f\n", avgPrice)

	// ===== Text processing =====
	fmt.Println("\n=== Text Processing ===")

	words := stream.Of("hello", "world", "hello", "go", "stream", "go", "chain")

	result := stream.Map(
		words.
			Distinct(func(s string) string { return s }).
			Sort(func(a, b string) int { return len(a) - len(b) }),
		func(s string) string { return strings.ToUpper(s) },
	).ToSlice()

	fmt.Println("Unique words (sorted by length, uppercased):", result)

	// ===== Complex pipeline =====
	fmt.Println("\n=== Complex Pipeline ===")

	savings := stream.Reduce(
		stream.FlatMap(
			users.
				Filter(func(u User) bool { return u.IsActive }).
				Filter(func(u User) bool { return u.Age >= 20 }),
			func(u User) []Order { return u.Orders },
		).Filter(func(o Order) bool { return o.Discount > 0 }),
		0.0,
		func(acc float64, o Order) float64 {
			return acc + o.Amount*o.Discount
		},
	)
	fmt.Printf("Total savings from discounts: $%.2f\n", savings)

	// =====================================================================
	// Part 2: Pipeline (Lazy Evaluation)
	// =====================================================================
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          Part 2: Pipeline (Lazy Evaluation)             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ===== Lazy constructors =====
	fmt.Println("\n=== Lazy Constructors ===")

	lazy := stream.Lazy(1, 2, 3, 4, 5)
	fmt.Println("Lazy:", lazy.ToSlice())

	lazyRange := stream.LazyRange(0, 5)
	fmt.Println("LazyRange(0,5):", lazyRange.ToSlice())

	// ===== Infinite sequences =====
	fmt.Println("\n=== Infinite Sequences ===")

	// Natural numbers: 0, 1, 2, 3, ...
	naturals := stream.Naturals().Take(10).ToSlice()
	fmt.Println("Naturals (first 10):", naturals)

	// Powers of 2: 1, 2, 4, 8, 16, ...
	powers := stream.Iterate(1, func(n int) int { return n * 2 }).Take(10).ToSlice()
	fmt.Println("Powers of 2:", powers)

	// Repeat a value
	dashes := stream.Repeat("-").Take(5).ToSlice()
	fmt.Println("Repeat:", dashes)

	// RepeatN
	stars := stream.RepeatN("*", 3).ToSlice()
	fmt.Println("RepeatN:", stars)

	// ===== Fibonacci sequence =====
	fmt.Println("\n=== Fibonacci Sequence ===")

	fib := stream.PipeMap(
		stream.Iterate(
			stream.Pair[int, int]{First: 0, Second: 1},
			func(p stream.Pair[int, int]) stream.Pair[int, int] {
				return stream.Pair[int, int]{First: p.Second, Second: p.First + p.Second}
			},
		).Take(15),
		func(p stream.Pair[int, int]) int { return p.First },
	).ToSlice()

	fmt.Println("Fibonacci:", fib)

	// ===== Lazy evaluation advantage =====
	fmt.Println("\n=== Lazy Evaluation Advantage ===")

	// Pipeline only processes elements until Take is satisfied
	evaluated := 0
	result2 := stream.LazyRange(0, 1_000_000).
		Peek(func(int) { evaluated++ }).
		Filter(func(n int) bool { return n%1000 == 0 }).
		Take(3).
		ToSlice()

	fmt.Printf("Result: %v (evaluated only %d of 1,000,000 elements)\n", result2, evaluated)

	// ===== Pipeline Filter + Sort + Take =====
	fmt.Println("\n=== Pipeline: Filter + Sort + Take ===")

	lazyProducts := stream.Lazy(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
	)

	topLazy := lazyProducts.
		Filter(func(p Product) bool { return p.InStock }).
		Sort(func(a, b Product) int {
			if a.Price > b.Price {
				return -1
			}
			if a.Price < b.Price {
				return 1
			}
			return 0
		}).
		Take(2).
		ToSlice()

	fmt.Println("Top 2 in-stock (lazy):")
	for _, p := range topLazy {
		fmt.Printf("  %s: $%.0f\n", p.Name, p.Price)
	}

	// ===== Pipeline transform functions =====
	fmt.Println("\n=== Pipeline Transforms ===")

	// PipeMap: type-changing
	pipeNames := stream.PipeMap(
		stream.Lazy(
			User{Name: "Alice", Age: 30, IsActive: true},
			User{Name: "Bob", Age: 17, IsActive: true},
			User{Name: "Charlie", Age: 25, IsActive: false},
		).Filter(func(u User) bool { return u.IsActive }),
		func(u User) string { return u.Name },
	).ToSlice()
	fmt.Println("Active user names (PipeMap):", pipeNames)

	// PipeFlatMap
	pipeOrders := stream.PipeFlatMap(
		stream.Lazy(
			User{Name: "Alice", Orders: []Order{{Product: "A"}, {Product: "B"}}},
			User{Name: "Bob", Orders: []Order{{Product: "C"}}},
		),
		func(u User) []Order { return u.Orders },
	)
	fmt.Printf("All orders (PipeFlatMap): %d\n", pipeOrders.Count())

	// PipeZip
	pipeZip := stream.PipeZip(
		stream.Lazy("Alice", "Bob", "Charlie"),
		stream.Lazy(85, 92, 78),
	).ToSlice()
	fmt.Println("Zip:")
	for _, p := range pipeZip {
		fmt.Printf("  %s: %d\n", p.First, p.Second)
	}

	// PipeEnumerate
	enumerated := stream.PipeEnumerate(stream.Lazy("a", "b", "c")).ToSlice()
	fmt.Println("Enumerate:")
	for _, p := range enumerated {
		fmt.Printf("  [%d] %s\n", p.First, p.Second)
	}

	// PipeGroupBy
	groups := stream.PipeGroupBy(
		stream.Lazy(
			Product{Name: "Laptop", Category: "Electronics"},
			Product{Name: "T-Shirt", Category: "Clothing"},
			Product{Name: "Keyboard", Category: "Electronics"},
		),
		func(p Product) string { return p.Category },
	)
	fmt.Println("GroupBy (Pipeline):")
	for k, v := range groups {
		fmt.Printf("  %s: %d items\n", k, len(v))
	}

	// PipeReduce
	pipeTotal := stream.PipeReduce(
		stream.Lazy(
			Order{Amount: 100, Discount: 0.1},
			Order{Amount: 200, Discount: 0.2},
		),
		0.0,
		func(acc float64, o Order) float64 { return acc + o.Amount*(1-o.Discount) },
	)
	fmt.Printf("Total (PipeReduce): $%.2f\n", pipeTotal)

	// ===== Chain: concatenate pipelines =====
	fmt.Println("\n=== Pipeline Chain ===")

	p1 := stream.Lazy(1, 2, 3)
	p2 := stream.Lazy(4, 5, 6)
	chained := p1.Chain(p2).ToSlice()
	fmt.Println("Chain:", chained)

	// ===== Pipeline terminal operations =====
	fmt.Println("\n=== Pipeline Terminal Operations ===")

	pipe := stream.Lazy(10, 20, 30, 40, 50)

	first, _ := pipe.First()
	fmt.Println("First:", first)

	last, _ := pipe.Last()
	fmt.Println("Last:", last)

	found, _ := pipe.Find(func(n int) bool { return n > 25 })
	fmt.Println("Find(>25):", found)

	sum := pipe.Reduce(0, func(acc, v int) int { return acc + v })
	fmt.Println("Reduce(sum):", sum)

	fmt.Println("Any(>40):", pipe.Any(func(n int) bool { return n > 40 }))
	fmt.Println("All(>0):", pipe.All(func(n int) bool { return n > 0 }))
	fmt.Println("None(>100):", pipe.None(func(n int) bool { return n > 100 }))
	fmt.Println("Count:", pipe.Count())

	// =====================================================================
	// Part 3: Stream <-> Pipeline Bridge
	// =====================================================================
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║        Part 3: Stream <-> Pipeline Bridge               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ===== Stream -> Pipeline (lazy mode) =====
	fmt.Println("\n=== Stream -> Pipeline (.Lazy()) ===")

	eagerStream := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	lazyResult := eagerStream.
		Lazy(). // Switch to lazy evaluation
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(3).
		ToSlice()
	fmt.Println("Stream.Lazy().Filter().Take(3):", lazyResult)

	// ===== Pipeline -> Stream (eager mode) =====
	fmt.Println("\n=== Pipeline -> Stream (.ToStream()) ===")

	lazyPipe := stream.Naturals().
		Filter(func(n int) bool { return n%3 == 0 }).
		Take(6)

	// Convert to Stream to use Stream-only operations
	chunks := lazyPipe.ToStream().Chunk(3)
	fmt.Println("Pipeline.ToStream().Chunk(3):")
	for i, c := range chunks {
		fmt.Printf("  Chunk %d: %v\n", i, c.ToSlice())
	}

	// ===== iter.Seq bridge =====
	fmt.Println("\n=== iter.Seq Bridge ===")

	// Stream -> iter.Seq -> for range
	fmt.Print("Stream.Iter(): ")
	for v := range stream.Of(1, 2, 3).Iter() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// Stream -> iter.Seq2 -> for range (with index)
	fmt.Print("Stream.Iter2(): ")
	for i, v := range stream.Of("a", "b", "c").Iter2() {
		fmt.Printf("%d:%s ", i, v)
	}
	fmt.Println()

	// maps.Keys -> Collect -> Stream operations
	m := map[string]int{"Go": 1, "Rust": 2, "Python": 3}
	sorted := stream.Collect(maps.Keys(m)).
		Sort(func(a, b string) int { return strings.Compare(a, b) }).
		ToSlice()
	fmt.Println("Collect(maps.Keys):", sorted)

	// maps.All -> Collect2 -> Stream of Pairs
	pairs := stream.Collect2(maps.All(m)).ToSlice()
	fmt.Printf("Collect2(maps.All): %d pairs\n", len(pairs))

	// ===== Full interop: Stream -> iter.Seq -> Pipeline -> Stream =====
	fmt.Println("\n=== Full Interop ===")

	s := stream.Of(10, 20, 30, 40, 50) // Stream
	seq := s.Iter()                     // Stream -> iter.Seq
	p := stream.LazyFrom(seq)           // iter.Seq -> Pipeline
	finalStream := p.                   // Pipeline operations
					Filter(func(n int) bool { return n > 20 }).
					ToStream()  // Pipeline -> Stream
	finalSlice := finalStream.Reverse().ToSlice() // Stream operations
	fmt.Println("Stream->iter.Seq->Pipeline->Stream:", finalSlice)

	// ===== Practical example: Lazy filter, then eager Shuffle =====
	fmt.Println("\n=== Practical: Lazy + Eager Hybrid ===")

	// Use Pipeline for efficient filtering of large data, then Stream for Shuffle
	multiples := stream.Naturals().
		Filter(func(n int) bool { return n%7 == 0 }).
		Skip(1).   // skip 0
		Take(10).  // first 10 multiples of 7
		ToStream() // switch to Stream for Shuffle

	fmt.Println("Multiples of 7 (shuffled):", multiples.Shuffle().ToSlice())
	fmt.Println("Multiples of 7 (reversed):", multiples.Reverse().Take(5).ToSlice())
}

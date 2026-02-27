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
	// Part 1: Basic Operations
	// =====================================================================
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║              Part 1: Basic Operations                    ║")
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

	categories := stream.Map(
		stream.Of(
			Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
			Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
			Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
			Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
			Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
			Product{Name: "Sneakers", Category: "Clothing", Price: 120, InStock: false},
		), func(p Product) string {
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

	// Re-create because Count consumed it
	allOrders2 := stream.FlatMap(
		stream.Of(
			User{Name: "Alice", Age: 30, IsActive: true, Orders: []Order{
				{Product: "Laptop", Amount: 1200, Discount: 0.1},
				{Product: "Mouse", Amount: 50, Discount: 0},
			}},
			User{Name: "Bob", Age: 17, IsActive: true, Orders: []Order{
				{Product: "Keyboard", Amount: 80, Discount: 0.5},
			}},
		).Filter(func(u User) bool { return u.IsActive }),
		func(u User) []Order { return u.Orders },
	)

	allOrders2.ForEach(func(o Order) {
		fmt.Printf("  %s: $%.0f (%.0f%% off)\n", o.Product, o.Amount, o.Discount*100)
	})

	// ===== Reduce (type-changing) =====
	fmt.Println("\n=== Reduce ===")

	total := stream.Reduce(
		stream.FlatMap(
			stream.Of(
				User{Name: "Alice", Age: 30, IsActive: true, Orders: []Order{
					{Product: "Laptop", Amount: 1200, Discount: 0.1},
					{Product: "Mouse", Amount: 50, Discount: 0},
				}},
				User{Name: "Bob", Age: 17, IsActive: true, Orders: []Order{
					{Product: "Keyboard", Amount: 80, Discount: 0.5},
				}},
			).Filter(func(u User) bool { return u.IsActive }),
			func(u User) []Order { return u.Orders },
		),
		0.0,
		func(acc float64, o Order) float64 {
			return acc + o.Amount*(1-o.Discount)
		},
	)
	fmt.Printf("Total after discount: $%.2f\n", total)

	// ===== GroupBy =====
	fmt.Println("\n=== GroupBy ===")

	byCategory := stream.GroupBy(
		stream.Of(
			Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
			Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
			Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
			Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
			Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
			Product{Name: "Sneakers", Category: "Clothing", Price: 120, InStock: false},
		),
		func(p Product) string { return p.Category },
	)
	for category, group := range byCategory {
		totalPrice := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
		fmt.Printf("  %s: total=$%.0f (%d products)\n", category, totalPrice, len(group))
	}

	// ===== Associate =====
	fmt.Println("\n=== Associate ===")

	nameToAge := stream.Associate(
		stream.Of(
			User{Name: "Alice", Age: 30},
			User{Name: "Bob", Age: 25},
			User{Name: "Charlie", Age: 25},
		),
		func(u User) (string, int) { return u.Name, u.Age },
	)
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

	inStock, outOfStock := stream.Of(
		Product{Name: "Laptop", Price: 1200, InStock: true},
		Product{Name: "Headphones", Price: 150, InStock: false},
		Product{Name: "Keyboard", Price: 75, InStock: true},
		Product{Name: "Sneakers", Price: 120, InStock: false},
	).Partition(func(p Product) bool { return p.InStock })
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
	fmt.Printf("Avg: %.1f\n", stream.Avg(stream.Of(10.0, 20.0, 30.0, 40.0, 50.0)))

	min, _ := stream.Min(stream.Of(10.0, 20.0, 30.0, 40.0, 50.0))
	max, _ := stream.Max(stream.Of(10.0, 20.0, 30.0, 40.0, 50.0))
	fmt.Printf("Min: %.0f, Max: %.0f\n", min, max)

	avgPrice := stream.AvgBy(
		stream.Of(
			Product{Price: 1200},
			Product{Price: 150},
			Product{Price: 75},
			Product{Price: 80},
			Product{Price: 25},
			Product{Price: 120},
		),
		func(p Product) float64 { return p.Price },
	)
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

	// =====================================================================
	// Part 2: Infinite Sequences & Lazy Evaluation
	// =====================================================================
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║      Part 2: Infinite Sequences & Lazy Evaluation       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

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

	fib := stream.Map(
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

	// Stream only processes elements until Take is satisfied
	evaluated := 0
	result2 := stream.Range(0, 1_000_000).
		Peek(func(int) { evaluated++ }).
		Filter(func(n int) bool { return n%1000 == 0 }).
		Take(3).
		ToSlice()

	fmt.Printf("Result: %v (evaluated only %d of 1,000,000 elements)\n", result2, evaluated)

	// ===== Filter + Sort + Take =====
	fmt.Println("\n=== Filter + Sort + Take ===")

	topProducts := stream.Of(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
	).
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

	fmt.Println("Top 2 in-stock:")
	for _, p := range topProducts {
		fmt.Printf("  %s: $%.0f\n", p.Name, p.Price)
	}

	// ===== Chain: concatenate streams =====
	fmt.Println("\n=== Chain ===")

	p1 := stream.Of(1, 2, 3)
	p2 := stream.Of(4, 5, 6)
	chained := p1.Chain(p2).ToSlice()
	fmt.Println("Chain:", chained)

	// =====================================================================
	// Part 3: iter.Seq Bridge
	// =====================================================================
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║              Part 3: iter.Seq Bridge                     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ===== Stream.Seq() -> for range =====
	fmt.Println("\n=== Seq() ===")

	fmt.Print("Seq(): ")
	for v := range stream.Of(1, 2, 3).Seq() {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// ===== Enumerate =====
	fmt.Println("\n=== Enumerate ===")

	enumerated := stream.Enumerate(stream.Of("a", "b", "c")).ToSlice()
	for _, p := range enumerated {
		fmt.Printf("  [%d] %s\n", p.First, p.Second)
	}

	// maps.Keys -> Collect -> Stream operations
	fmt.Println("\n=== Collect from maps ===")

	m := map[string]int{"Go": 1, "Rust": 2, "Python": 3}
	sorted := stream.Collect(maps.Keys(m)).
		Sort(func(a, b string) int { return strings.Compare(a, b) }).
		ToSlice()
	fmt.Println("Collect(maps.Keys):", sorted)

	// maps.All -> Collect2 -> Stream of Pairs
	pairs := stream.Collect2(maps.All(m)).ToSlice()
	fmt.Printf("Collect2(maps.All): %d pairs\n", len(pairs))

	// ===== Practical example: Shuffle from infinite =====
	fmt.Println("\n=== Practical: Infinite + Shuffle ===")

	multiples := stream.Naturals().
		Filter(func(n int) bool { return n%7 == 0 }).
		Skip(1).  // skip 0
		Take(10). // first 10 multiples of 7
		Shuffle()

	fmt.Println("Multiples of 7 (shuffled):", multiples.ToSlice())
	fmt.Println("Multiples of 7 (reversed):", stream.Naturals().
		Filter(func(n int) bool { return n%7 == 0 }).
		Skip(1).
		Take(10).
		Reverse().
		Take(5).
		ToSlice())
}

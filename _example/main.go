package main

import (
	"fmt"
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
	// ===== Constructors =====
	fmt.Println("=== Constructors ===")

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

	// ===== Chaining: complex pipeline =====
	fmt.Println("\n=== Complex Pipeline ===")

	// Active adult users -> all orders -> discount > 0 -> total savings
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
}

package main

import (
	"fmt"
	"strings"

	"github.com/nd-forge/stream"
)

// --- Data types ---

type Trade struct {
	Symbol string
	Amount float64
	PnL    float64
	IsOpen bool
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

	trades := stream.Of(
		Trade{Symbol: "AUDUSD", Amount: 10000, PnL: 150, IsOpen: false},
		Trade{Symbol: "NZDUSD", Amount: 5000, PnL: -80, IsOpen: false},
		Trade{Symbol: "AUDUSD", Amount: 20000, PnL: 300, IsOpen: true},
		Trade{Symbol: "EURUSD", Amount: 15000, PnL: -200, IsOpen: false},
		Trade{Symbol: "AUDUSD", Amount: 8000, PnL: 50, IsOpen: false},
		Trade{Symbol: "NZDUSD", Amount: 12000, PnL: 100, IsOpen: false},
	)

	top3 := trades.
		Filter(func(t Trade) bool { return !t.IsOpen }).
		Filter(func(t Trade) bool { return t.PnL > 0 }).
		Sort(func(a, b Trade) int {
			if a.PnL > b.PnL {
				return -1
			}
			if a.PnL < b.PnL {
				return 1
			}
			return 0
		}).
		Take(3).
		ToSlice()

	fmt.Println("Top 3 profitable closed trades:")
	for _, t := range top3 {
		fmt.Printf("  %s: PnL=%.0f\n", t.Symbol, t.PnL)
	}

	// ===== Map (type-changing) =====
	fmt.Println("\n=== Map ===")

	symbols := stream.Map(trades, func(t Trade) string {
		return t.Symbol
	}).Distinct(func(s string) string { return s }).ToSlice()

	fmt.Println("Unique symbols:", symbols)

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

	bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })
	for symbol, group := range bySymbol {
		totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
		fmt.Printf("  %s: PnL=%.0f (%d trades)\n", symbol, totalPnL, group.Count())
	}

	// ===== Associate =====
	fmt.Println("\n=== Associate ===")

	nameToAge := stream.Associate(users, func(u User) (string, int) {
		return u.Name, u.Age
	})
	fmt.Println("Name→Age:", nameToAge)

	// ===== Zip =====
	fmt.Println("\n=== Zip ===")

	dates := stream.Of("2024-01", "2024-02", "2024-03")
	prices := stream.Of(1.05, 1.06, 1.04)

	stream.Zip(dates, prices).ForEach(func(p stream.Pair[string, float64]) {
		fmt.Printf("  %s: %.2f\n", p.First, p.Second)
	})

	// ===== Partition =====
	fmt.Println("\n=== Partition ===")

	profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
	fmt.Printf("Profit: %d trades, Loss: %d trades\n", profit.Count(), loss.Count())

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

	avgPnL := stream.AvgBy(trades, func(t Trade) float64 { return t.PnL })
	fmt.Printf("Avg PnL: %.2f\n", avgPnL)

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

	// Active adult users → all orders → discount > 0 → total savings
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

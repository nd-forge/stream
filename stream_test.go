package stream_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nd-forge/stream"
)

// ---------------------------------------------------------------------------
// Test types
// ---------------------------------------------------------------------------

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
	UserID   int
	Product  string
	Amount   float64
	Discount float64
}

// ---------------------------------------------------------------------------
// Constructor tests
// ---------------------------------------------------------------------------

func TestOf(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5)
	if s.Count() != 5 {
		t.Errorf("expected 5, got %d", s.Count())
	}
}

func TestFrom(t *testing.T) {
	original := []int{1, 2, 3}
	s := stream.From(original)
	original[0] = 999 // mutating original should not affect stream
	if s.ToSlice()[0] != 1 {
		t.Error("From should copy the slice")
	}
}

func TestRange(t *testing.T) {
	s := stream.Range(0, 5)
	expected := []int{0, 1, 2, 3, 4}
	result := s.ToSlice()
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Range: expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestGenerate(t *testing.T) {
	s := stream.Generate(5, func(i int) int { return i * i })
	expected := []int{0, 1, 4, 9, 16}
	result := s.ToSlice()
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Generate: expected %d, got %d", expected[i], v)
		}
	}
}

// ---------------------------------------------------------------------------
// Chainable method tests
// ---------------------------------------------------------------------------

func TestFilter(t *testing.T) {
	result := stream.Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()

	if len(result) != 3 || result[0] != 2 || result[1] != 4 || result[2] != 6 {
		t.Errorf("Filter: unexpected result %v", result)
	}
}

func TestReject(t *testing.T) {
	result := stream.Of(1, 2, 3, 4, 5).
		Reject(func(n int) bool { return n%2 == 0 }).
		ToSlice()

	if len(result) != 3 || result[0] != 1 {
		t.Errorf("Reject: unexpected result %v", result)
	}
}

func TestSort(t *testing.T) {
	result := stream.Of(3, 1, 4, 1, 5, 9).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()

	expected := []int{1, 1, 3, 4, 5, 9}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Sort: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestReverse(t *testing.T) {
	result := stream.Of(1, 2, 3).Reverse().ToSlice()
	if result[0] != 3 || result[1] != 2 || result[2] != 1 {
		t.Errorf("Reverse: unexpected result %v", result)
	}
}

func TestTakeAndSkip(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5)

	top3 := s.Take(3).ToSlice()
	if len(top3) != 3 || top3[2] != 3 {
		t.Errorf("Take: unexpected %v", top3)
	}

	rest := s.Skip(3).ToSlice()
	if len(rest) != 2 || rest[0] != 4 {
		t.Errorf("Skip: unexpected %v", rest)
	}
}

func TestTakeLast(t *testing.T) {
	result := stream.Of(1, 2, 3, 4, 5).TakeLast(2).ToSlice()
	if len(result) != 2 || result[0] != 4 || result[1] != 5 {
		t.Errorf("TakeLast: unexpected %v", result)
	}
}

func TestTakeWhileAndDropWhile(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5, 1, 2)

	taken := s.TakeWhile(func(n int) bool { return n < 4 }).ToSlice()
	if len(taken) != 3 {
		t.Errorf("TakeWhile: expected 3 elements, got %v", taken)
	}

	dropped := s.DropWhile(func(n int) bool { return n < 4 }).ToSlice()
	if len(dropped) != 4 || dropped[0] != 4 {
		t.Errorf("DropWhile: unexpected %v", dropped)
	}
}

func TestDistinct(t *testing.T) {
	result := stream.Of("a", "b", "a", "c", "b").
		Distinct(func(s string) string { return s }).
		ToSlice()

	if len(result) != 3 {
		t.Errorf("Distinct: expected 3 unique, got %v", result)
	}
}

func TestPartition(t *testing.T) {
	inStock, outOfStock := stream.Of(
		Product{Name: "Laptop", Price: 1200, InStock: true},
		Product{Name: "Headphones", Price: 150, InStock: false},
		Product{Name: "Keyboard", Price: 75, InStock: true},
	).Partition(func(p Product) bool { return p.InStock })

	if inStock.Count() != 2 || outOfStock.Count() != 1 {
		t.Errorf("Partition: unexpected split %d/%d", inStock.Count(), outOfStock.Count())
	}
}

func TestChunk(t *testing.T) {
	chunks := stream.Of(1, 2, 3, 4, 5).Chunk(2)
	if len(chunks) != 3 {
		t.Errorf("Chunk: expected 3 chunks, got %d", len(chunks))
	}
	if chunks[2].Count() != 1 {
		t.Errorf("Chunk: last chunk should have 1 element")
	}
}

// ---------------------------------------------------------------------------
// Terminal operation tests
// ---------------------------------------------------------------------------

func TestFirstAndLast(t *testing.T) {
	s := stream.Of(10, 20, 30)

	first, ok := s.First()
	if !ok || first != 10 {
		t.Errorf("First: expected 10, got %d", first)
	}

	last, ok := s.Last()
	if !ok || last != 30 {
		t.Errorf("Last: expected 30, got %d", last)
	}

	_, ok = stream.Of[int]().First()
	if ok {
		t.Error("First on empty should return false")
	}
}

func TestFind(t *testing.T) {
	p, ok := stream.Of(
		Product{Name: "Laptop", Price: 1200},
		Product{Name: "Headphones", Price: 150},
	).Find(func(p Product) bool { return p.Name == "Headphones" })

	if !ok || p.Price != 150 {
		t.Error("Find: failed to find Headphones")
	}
}

func TestAnyAllNone(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5)

	if !s.Any(func(n int) bool { return n > 4 }) {
		t.Error("Any: should find element > 4")
	}
	if !s.All(func(n int) bool { return n > 0 }) {
		t.Error("All: all should be > 0")
	}
	if !s.None(func(n int) bool { return n > 10 }) {
		t.Error("None: no element > 10")
	}
}

func TestMinByMaxBy(t *testing.T) {
	products := stream.Of(
		Product{Name: "Laptop", Price: 1200},
		Product{Name: "Headphones", Price: 150},
		Product{Name: "Keyboard", Price: 75},
	)

	most, _ := products.MaxBy(func(a, b Product) bool { return a.Price < b.Price })
	if most.Name != "Laptop" {
		t.Errorf("MaxBy: expected Laptop, got %s", most.Name)
	}

	least, _ := products.MinBy(func(a, b Product) bool { return a.Price < b.Price })
	if least.Name != "Keyboard" {
		t.Errorf("MinBy: expected Keyboard, got %s", least.Name)
	}
}

// ---------------------------------------------------------------------------
// Transform function tests (type-changing)
// ---------------------------------------------------------------------------

func TestMap(t *testing.T) {
	result := stream.Map(
		stream.Of(1, 2, 3),
		func(n int) string { return fmt.Sprintf("item_%d", n) },
	).ToSlice()

	if result[0] != "item_1" || result[2] != "item_3" {
		t.Errorf("Map: unexpected %v", result)
	}
}

func TestFlatMap(t *testing.T) {
	users := stream.Of(
		User{Name: "Alice", Orders: []Order{{Product: "A"}, {Product: "B"}}},
		User{Name: "Bob", Orders: []Order{{Product: "C"}}},
	)

	orders := stream.FlatMap(users, func(u User) []Order { return u.Orders })
	if orders.Count() != 3 {
		t.Errorf("FlatMap: expected 3 orders, got %d", orders.Count())
	}
}

func TestGroupBy(t *testing.T) {
	products := stream.Of(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75},
	)

	groups := stream.GroupBy(products, func(p Product) string { return p.Category })

	if groups["Electronics"].Count() != 2 {
		t.Errorf("GroupBy: expected 2 Electronics products")
	}
	if groups["Clothing"].Count() != 1 {
		t.Errorf("GroupBy: expected 1 Clothing product")
	}
}

func TestZip(t *testing.T) {
	names := stream.Of("Alice", "Bob", "Charlie")
	scores := stream.Of(85.0, 92.0, 78.0)

	pairs := stream.Zip(names, scores).ToSlice()
	if len(pairs) != 3 || pairs[0].First != "Alice" || pairs[0].Second != 85.0 {
		t.Errorf("Zip: unexpected %v", pairs)
	}
}

func TestAssociate(t *testing.T) {
	users := stream.Of(
		User{Name: "Alice", Age: 30},
		User{Name: "Bob", Age: 25},
	)

	nameToAge := stream.Associate(users, func(u User) (string, int) {
		return u.Name, u.Age
	})

	if nameToAge["Alice"] != 30 || nameToAge["Bob"] != 25 {
		t.Errorf("Associate: unexpected %v", nameToAge)
	}
}

func TestReduce(t *testing.T) {
	total := stream.Reduce(
		stream.Of(
			Order{Amount: 100, Discount: 0.1},
			Order{Amount: 200, Discount: 0.2},
			Order{Amount: 300, Discount: 0.0},
		),
		0.0,
		func(acc float64, o Order) float64 {
			return acc + o.Amount*(1-o.Discount)
		},
	)

	expected := 90.0 + 160.0 + 300.0
	if total != expected {
		t.Errorf("Reduce: expected %.1f, got %.1f", expected, total)
	}
}

// ---------------------------------------------------------------------------
// Numeric function tests
// ---------------------------------------------------------------------------

func TestSum(t *testing.T) {
	total := stream.Sum(stream.Of(1, 2, 3, 4, 5))
	if total != 15 {
		t.Errorf("Sum: expected 15, got %d", total)
	}
}

func TestAvg(t *testing.T) {
	avg := stream.Avg(stream.Of(10.0, 20.0, 30.0))
	if avg != 20.0 {
		t.Errorf("Avg: expected 20.0, got %f", avg)
	}
}

func TestSumByAvgBy(t *testing.T) {
	products := stream.Of(
		Product{Price: 1200},
		Product{Price: 150},
		Product{Price: 75},
	)

	total := stream.SumBy(products, func(p Product) float64 { return p.Price })
	if total != 1425 {
		t.Errorf("SumBy: expected 1425, got %f", total)
	}

	avg := stream.AvgBy(products, func(p Product) float64 { return p.Price })
	expected := 1425.0 / 3.0
	if avg < expected-0.01 || avg > expected+0.01 {
		t.Errorf("AvgBy: expected ~%.2f, got %.2f", expected, avg)
	}
}

func TestMinMax(t *testing.T) {
	s := stream.Of(3.0, 1.0, 4.0, 1.0, 5.0)

	min, _ := stream.Min(s)
	if min != 1.0 {
		t.Errorf("Min: expected 1.0, got %f", min)
	}

	max, _ := stream.Max(s)
	if max != 5.0 {
		t.Errorf("Max: expected 5.0, got %f", max)
	}
}

// ---------------------------------------------------------------------------
// Chaining integration tests — practical examples
// ---------------------------------------------------------------------------

func TestChaining_ProductAnalysis(t *testing.T) {
	products := stream.Of(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
		Product{Name: "Sneakers", Category: "Clothing", Price: 120, InStock: false},
	)

	// 在庫ありの商品を価格降順でトップ3
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

	if len(top3) != 3 || top3[0].Name != "Laptop" {
		t.Errorf("Chaining: unexpected top3 %v", top3)
	}

	// カテゴリ別の価格合計
	byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })
	for category, group := range byCategory {
		totalPrice := stream.SumBy(group, func(p Product) float64 { return p.Price })
		t.Logf("%s: total = %.2f (%d products)", category, totalPrice, group.Count())
	}
}

func TestChaining_UserOrders(t *testing.T) {
	users := stream.Of(
		User{Name: "Alice", Age: 30, IsActive: true, Orders: []Order{
			{Amount: 100, Discount: 0.1},
			{Amount: 200, Discount: 0},
		}},
		User{Name: "Bob", Age: 17, IsActive: true, Orders: []Order{
			{Amount: 50, Discount: 0.5},
		}},
		User{Name: "Charlie", Age: 25, IsActive: false, Orders: []Order{
			{Amount: 300, Discount: 0.2},
		}},
	)

	// アクティブな成人ユーザーの全注文を取得し、割引後の合計を計算
	total := stream.Reduce(
		stream.FlatMap(
			users.
				Filter(func(u User) bool { return u.IsActive }).
				Filter(func(u User) bool { return u.Age >= 20 }),
			func(u User) []Order { return u.Orders },
		),
		0.0,
		func(acc float64, o Order) float64 {
			return acc + o.Amount*(1-o.Discount)
		},
	)

	// Alice: 100*0.9 + 200*1.0 = 290
	if total != 290 {
		t.Errorf("UserOrders: expected 290, got %.1f", total)
	}
}

func TestChaining_TextProcessing(t *testing.T) {
	words := stream.Of("hello", "world", "hello", "go", "stream", "go", "chain")

	// ユニーク単語を大文字にして、長さ順ソート
	result := stream.Map(
		words.
			Distinct(func(s string) string { return s }).
			Sort(func(a, b string) int { return len(a) - len(b) }),
		func(s string) string { return strings.ToUpper(s) },
	).ToSlice()

	t.Logf("Processed words: %v", result)

	if result[0] != "GO" {
		t.Errorf("TextProcessing: expected GO first, got %s", result[0])
	}
}

func TestChaining_BatchProcessing(t *testing.T) {
	// 100件のアイテムを20件ずつバッチ処理
	items := stream.Range(0, 100)
	batches := items.Chunk(20)

	if len(batches) != 5 {
		t.Errorf("Batch: expected 5 batches, got %d", len(batches))
	}

	for i, batch := range batches {
		if batch.Count() != 20 {
			t.Errorf("Batch %d: expected 20 items, got %d", i, batch.Count())
		}
	}
}

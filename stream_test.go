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
	result := stream.Range(0, 5).ToSlice()
	expected := []int{0, 1, 2, 3, 4}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Range: expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestGenerate(t *testing.T) {
	result := stream.Generate(5, func(i int) int { return i * i }).ToSlice()
	expected := []int{0, 1, 4, 9, 16}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Generate: expected %d, got %d", expected[i], v)
		}
	}
}

// ---------------------------------------------------------------------------
// Generator tests (infinite sequences)
// ---------------------------------------------------------------------------

func TestRepeat(t *testing.T) {
	result := stream.Repeat(42).Take(5).ToSlice()
	if len(result) != 5 {
		t.Errorf("Repeat: expected 5 elements, got %d", len(result))
	}
	for _, v := range result {
		if v != 42 {
			t.Errorf("Repeat: expected 42, got %d", v)
		}
	}
}

func TestRepeatN(t *testing.T) {
	result := stream.RepeatN("x", 3).ToSlice()
	if len(result) != 3 || result[0] != "x" {
		t.Errorf("RepeatN: unexpected %v", result)
	}
}

func TestIterate(t *testing.T) {
	// Powers of 2: 1, 2, 4, 8, 16
	result := stream.Iterate(1, func(n int) int { return n * 2 }).Take(5).ToSlice()
	expected := []int{1, 2, 4, 8, 16}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Iterate: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestNaturals(t *testing.T) {
	result := stream.Naturals().Take(5).ToSlice()
	expected := []int{0, 1, 2, 3, 4}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Naturals: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestInfiniteFilterTake(t *testing.T) {
	// First 5 even naturals: 0, 2, 4, 6, 8
	result := stream.Naturals().
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(5).
		ToSlice()
	expected := []int{0, 2, 4, 6, 8}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("InfiniteFilterTake: expected %d at %d, got %d", expected[i], i, v)
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

func TestTakeZero(t *testing.T) {
	result := stream.Of(1, 2, 3).Take(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("Take(0): expected empty, got %v", result)
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

func TestShuffle(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5).Shuffle()
	result := s.ToSlice()
	if len(result) != 5 {
		t.Errorf("Shuffle: expected 5 elements, got %d", len(result))
	}
}

func TestChain(t *testing.T) {
	s1 := stream.Of(1, 2, 3)
	s2 := stream.Of(4, 5, 6)
	result := s1.Chain(s2).ToSlice()
	if len(result) != 6 || result[3] != 4 {
		t.Errorf("Chain: unexpected %v", result)
	}
}

func TestSortThenTake(t *testing.T) {
	// Sort buffers all, but Take after Sort is still lazy
	result := stream.Of(5, 3, 1, 4, 2).
		Sort(func(a, b int) int { return a - b }).
		Take(3).
		ToSlice()
	expected := []int{1, 2, 3}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("SortThenTake: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestPeek(t *testing.T) {
	var peeked []int
	result := stream.Of(1, 2, 3).
		Peek(func(n int) { peeked = append(peeked, n) }).
		Filter(func(n int) bool { return n > 1 }).
		ToSlice()

	if len(peeked) != 3 {
		t.Errorf("Peek: expected 3 peeks, got %d", len(peeked))
	}
	if len(result) != 2 {
		t.Errorf("Peek: filter should produce 2, got %v", result)
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

func TestCount(t *testing.T) {
	n := stream.Of(1, 2, 3, 4, 5).
		Filter(func(v int) bool { return v > 2 }).
		Count()
	if n != 3 {
		t.Errorf("Count: expected 3, got %d", n)
	}
}

func TestCountBy(t *testing.T) {
	n := stream.Of(1, 2, 3, 4, 5).
		CountBy(func(v int) bool { return v%2 == 0 })
	if n != 2 {
		t.Errorf("CountBy: expected 2, got %d", n)
	}
}

func TestIsEmpty(t *testing.T) {
	if !stream.Of[int]().IsEmpty() {
		t.Error("IsEmpty: empty should return true")
	}
	if stream.Of(1).IsEmpty() {
		t.Error("IsEmpty: non-empty should return false")
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

	least, _ := stream.Of(
		Product{Name: "Laptop", Price: 1200},
		Product{Name: "Headphones", Price: 150},
		Product{Name: "Keyboard", Price: 75},
	).MinBy(func(a, b Product) bool { return a.Price < b.Price })
	if least.Name != "Keyboard" {
		t.Errorf("MinBy: expected Keyboard, got %s", least.Name)
	}
}

func TestForEach(t *testing.T) {
	var result []int
	stream.Of(1, 2, 3).ForEach(func(v int) { result = append(result, v*2) })
	if len(result) != 3 || result[0] != 2 || result[2] != 6 {
		t.Errorf("ForEach: unexpected %v", result)
	}
}

func TestForEachIndexed(t *testing.T) {
	var indices []int
	stream.Of("a", "b", "c").ForEachIndexed(func(i int, _ string) {
		indices = append(indices, i)
	})
	if len(indices) != 3 || indices[2] != 2 {
		t.Errorf("ForEachIndexed: unexpected %v", indices)
	}
}

func TestReduce(t *testing.T) {
	sum := stream.Of(1, 2, 3, 4, 5).
		Reduce(0, func(acc, v int) int { return acc + v })
	if sum != 15 {
		t.Errorf("Reduce: expected 15, got %d", sum)
	}
}

func TestSeq(t *testing.T) {
	s := stream.Of(1, 2, 3)
	var result []int
	for v := range s.Seq() {
		result = append(result, v)
	}
	if len(result) != 3 {
		t.Errorf("Seq: unexpected %v", result)
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

func TestMapIndexed(t *testing.T) {
	result := stream.MapIndexed(
		stream.Of("a", "b", "c"),
		func(i int, s string) string { return fmt.Sprintf("%d:%s", i, s) },
	).ToSlice()

	if result[0] != "0:a" || result[2] != "2:c" {
		t.Errorf("MapIndexed: unexpected %v", result)
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

	if len(groups["Electronics"]) != 2 {
		t.Errorf("GroupBy: expected 2 Electronics products")
	}
	if len(groups["Clothing"]) != 1 {
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

func TestZipUnevenLengths(t *testing.T) {
	s1 := stream.Of(1, 2, 3, 4, 5)
	s2 := stream.Of("a", "b")

	pairs := stream.Zip(s1, s2).ToSlice()
	if len(pairs) != 2 {
		t.Errorf("Zip uneven: expected 2, got %d", len(pairs))
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

func TestReduceTypeChanging(t *testing.T) {
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

func TestFlatten(t *testing.T) {
	result := stream.Flatten(
		stream.Of([]int{1, 2}, []int{3, 4}, []int{5}),
	).ToSlice()
	if len(result) != 5 || result[0] != 1 || result[4] != 5 {
		t.Errorf("Flatten: unexpected %v", result)
	}
}

func TestToMap(t *testing.T) {
	m := stream.ToMap(stream.Zip(
		stream.Of("a", "b", "c"),
		stream.Of(1, 2, 3),
	))
	if m["a"] != 1 || m["c"] != 3 {
		t.Errorf("ToMap: unexpected %v", m)
	}
}

func TestEnumerate(t *testing.T) {
	result := stream.Enumerate(stream.Of("a", "b", "c")).ToSlice()
	if len(result) != 3 || result[0].First != 0 || result[0].Second != "a" {
		t.Errorf("Enumerate: unexpected %v", result)
	}
	if result[2].First != 2 || result[2].Second != "c" {
		t.Errorf("Enumerate: unexpected last %v", result[2])
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

	products2 := stream.Of(
		Product{Price: 1200},
		Product{Price: 150},
		Product{Price: 75},
	)
	avg := stream.AvgBy(products2, func(p Product) float64 { return p.Price })
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

	s2 := stream.Of(3.0, 1.0, 4.0, 1.0, 5.0)
	max, _ := stream.Max(s2)
	if max != 5.0 {
		t.Errorf("Max: expected 5.0, got %f", max)
	}
}

// ---------------------------------------------------------------------------
// Lazy evaluation proof tests
// ---------------------------------------------------------------------------

func TestShortCircuit_Take(t *testing.T) {
	// Prove that Take stops iteration early
	evaluated := 0
	result := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
		Peek(func(int) { evaluated++ }).
		Take(3).
		ToSlice()

	if len(result) != 3 {
		t.Errorf("ShortCircuit.Take: expected 3 elements, got %d", len(result))
	}
	if evaluated != 3 {
		t.Errorf("ShortCircuit.Take: expected 3 evaluations, got %d (lazy broken!)", evaluated)
	}
}

func TestShortCircuit_Find(t *testing.T) {
	evaluated := 0
	v, ok := stream.Of(1, 2, 3, 4, 5).
		Peek(func(int) { evaluated++ }).
		Find(func(n int) bool { return n == 3 })

	if !ok || v != 3 {
		t.Errorf("ShortCircuit.Find: expected 3, got %d", v)
	}
	if evaluated != 3 {
		t.Errorf("ShortCircuit.Find: expected 3 evaluations, got %d", evaluated)
	}
}

func TestShortCircuit_Any(t *testing.T) {
	evaluated := 0
	found := stream.Of(1, 2, 3, 4, 5).
		Peek(func(int) { evaluated++ }).
		Any(func(n int) bool { return n == 2 })

	if !found {
		t.Error("ShortCircuit.Any: should find 2")
	}
	if evaluated != 2 {
		t.Errorf("ShortCircuit.Any: expected 2 evaluations, got %d", evaluated)
	}
}

func TestNoAllocation_FilterTake(t *testing.T) {
	evaluated := 0
	result := stream.Range(0, 1_000_000).
		Peek(func(int) { evaluated++ }).
		Filter(func(n int) bool { return n%1000 == 0 }).
		Take(3).
		ToSlice()

	if len(result) != 3 || result[0] != 0 || result[1] != 1000 || result[2] != 2000 {
		t.Errorf("NoAllocation: unexpected %v", result)
	}
	// Should evaluate at most 2001 elements (0..2000), not 1_000_000
	if evaluated > 2001 {
		t.Errorf("NoAllocation: evaluated %d elements (should be ~2001)", evaluated)
	}
}

// ---------------------------------------------------------------------------
// Reusability test
// ---------------------------------------------------------------------------

func TestReusable(t *testing.T) {
	s := stream.Of(1, 2, 3, 4, 5).Filter(func(n int) bool { return n > 2 })

	// Calling ToSlice() multiple times should return same result
	r1 := s.ToSlice()
	r2 := s.ToSlice()

	if len(r1) != len(r2) {
		t.Errorf("Reusable: lengths differ %d vs %d", len(r1), len(r2))
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			t.Errorf("Reusable: differ at %d: %d vs %d", i, r1[i], r2[i])
		}
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
	products2 := stream.Of(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
		Product{Name: "Sneakers", Category: "Clothing", Price: 120, InStock: false},
	)
	byCategory := stream.GroupBy(products2, func(p Product) string { return p.Category })
	for category, group := range byCategory {
		totalPrice := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
		t.Logf("%s: total = %.2f (%d products)", category, totalPrice, len(group))
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

func TestIntegration_InfiniteSequence(t *testing.T) {
	// Fibonacci sequence using Pair
	fib := stream.Map(
		stream.Iterate(
			stream.Pair[int, int]{First: 0, Second: 1},
			func(p stream.Pair[int, int]) stream.Pair[int, int] {
				return stream.Pair[int, int]{First: p.Second, Second: p.First + p.Second}
			},
		).Take(10),
		func(p stream.Pair[int, int]) int { return p.First },
	).ToSlice()

	expected := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	for i, v := range fib {
		if v != expected[i] {
			t.Errorf("Fibonacci: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestIntegration_ChunkFromInfinite(t *testing.T) {
	// Use lazy Stream for filtering, then Chunk (buffers internally)
	result := stream.Naturals().
		Filter(func(n int) bool { return n%3 == 0 }).
		Take(12).
		Chunk(4)

	if len(result) != 3 {
		t.Errorf("ChunkFromInfinite: expected 3 chunks, got %d", len(result))
	}
	if result[0].Count() != 4 {
		t.Errorf("ChunkFromInfinite: first chunk should have 4 elements, got %d", result[0].Count())
	}
}

func TestIntegration_CollectFromIter(t *testing.T) {
	// Show interop: iter.Seq → Stream
	s := stream.Of(10, 20, 30, 40, 50)
	seq := s.Seq()
	result := stream.Collect(seq).
		Filter(func(n int) bool { return n > 20 }).
		ToSlice()

	if len(result) != 3 || result[0] != 30 {
		t.Errorf("CollectFromIter: unexpected %v", result)
	}
}

// ---------------------------------------------------------------------------
// Early termination (yield returns false) coverage tests
// ---------------------------------------------------------------------------

func TestGenerate_EarlyBreak(t *testing.T) {
	v, ok := stream.Generate(100, func(i int) int { return i }).First()
	if !ok || v != 0 {
		t.Errorf("Generate early break: expected 0, got %d", v)
	}
}

func TestRepeatN_EarlyBreak(t *testing.T) {
	result := stream.RepeatN("x", 100).Take(1).ToSlice()
	if len(result) != 1 || result[0] != "x" {
		t.Errorf("RepeatN early break: unexpected %v", result)
	}
}

func TestReverse_EarlyBreak(t *testing.T) {
	v, ok := stream.Of(1, 2, 3).Reverse().First()
	if !ok || v != 3 {
		t.Errorf("Reverse early break: expected 3, got %d", v)
	}
}

func TestReverse_Empty(t *testing.T) {
	result := stream.Of[int]().Reverse().ToSlice()
	if len(result) != 0 {
		t.Errorf("Reverse empty: expected empty, got %v", result)
	}
}

func TestTake_EarlyBreak(t *testing.T) {
	// Take(5) but consumer only takes 1
	v, ok := stream.Of(1, 2, 3, 4, 5).Take(5).First()
	if !ok || v != 1 {
		t.Errorf("Take early break: expected 1, got %d", v)
	}
}

func TestTakeNegative(t *testing.T) {
	result := stream.Of(1, 2, 3).Take(-1).ToSlice()
	if len(result) != 0 {
		t.Errorf("Take(-1): expected empty, got %v", result)
	}
}

func TestTakeLast_Zero(t *testing.T) {
	result := stream.Of(1, 2, 3).TakeLast(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("TakeLast(0): expected empty, got %v", result)
	}
}

func TestTakeLast_MoreThanLength(t *testing.T) {
	result := stream.Of(1, 2).TakeLast(10).ToSlice()
	if len(result) != 2 || result[0] != 1 || result[1] != 2 {
		t.Errorf("TakeLast(>len): expected [1 2], got %v", result)
	}
}

func TestTakeLast_EarlyBreak(t *testing.T) {
	v, ok := stream.Of(1, 2, 3, 4, 5).TakeLast(3).First()
	if !ok || v != 3 {
		t.Errorf("TakeLast early break: expected 3, got %d", v)
	}
}

func TestTakeLast_Negative(t *testing.T) {
	result := stream.Of(1, 2, 3).TakeLast(-1).ToSlice()
	if len(result) != 0 {
		t.Errorf("TakeLast(-1): expected empty, got %v", result)
	}
}

func TestSkip_Zero(t *testing.T) {
	result := stream.Of(1, 2, 3).Skip(0).ToSlice()
	if len(result) != 3 {
		t.Errorf("Skip(0): expected 3 elements, got %v", result)
	}
}

func TestSkip_Negative(t *testing.T) {
	result := stream.Of(1, 2, 3).Skip(-1).ToSlice()
	if len(result) != 3 {
		t.Errorf("Skip(-1): expected 3 elements, got %v", result)
	}
}

func TestSkip_EarlyBreak(t *testing.T) {
	v, ok := stream.Of(1, 2, 3, 4, 5).Skip(2).First()
	if !ok || v != 3 {
		t.Errorf("Skip early break: expected 3, got %d", v)
	}
}

func TestTakeWhile_EarlyBreak(t *testing.T) {
	// All elements pass predicate but consumer stops early
	v, ok := stream.Of(1, 2, 3, 4, 5).
		TakeWhile(func(n int) bool { return n < 10 }).
		First()
	if !ok || v != 1 {
		t.Errorf("TakeWhile early break: expected 1, got %d", v)
	}
}

func TestDropWhile_EarlyBreak(t *testing.T) {
	v, ok := stream.Of(1, 2, 3, 4, 5).
		DropWhile(func(n int) bool { return n < 3 }).
		First()
	if !ok || v != 3 {
		t.Errorf("DropWhile early break: expected 3, got %d", v)
	}
}

func TestDistinct_EarlyBreak(t *testing.T) {
	v, ok := stream.Of("a", "b", "c").
		Distinct(func(s string) string { return s }).
		First()
	if !ok || v != "a" {
		t.Errorf("Distinct early break: expected 'a', got %s", v)
	}
}

func TestShuffle_EarlyBreak(t *testing.T) {
	v, ok := stream.Of(1, 2, 3, 4, 5).Shuffle().First()
	if !ok {
		t.Error("Shuffle early break: expected a value")
	}
	_ = v
}

func TestChain_EarlyBreakFirst(t *testing.T) {
	// Break during the first stream
	v, ok := stream.Of(1, 2, 3).Chain(stream.Of(4, 5, 6)).First()
	if !ok || v != 1 {
		t.Errorf("Chain early break first: expected 1, got %d", v)
	}
}

func TestChain_EarlyBreakSecond(t *testing.T) {
	// Break during the second stream
	result := stream.Of(1).Chain(stream.Of(2, 3, 4)).Take(2).ToSlice()
	if len(result) != 2 || result[0] != 1 || result[1] != 2 {
		t.Errorf("Chain early break second: expected [1 2], got %v", result)
	}
}

func TestFind_NotFound(t *testing.T) {
	_, ok := stream.Of(1, 2, 3).Find(func(n int) bool { return n > 10 })
	if ok {
		t.Error("Find: should not find element > 10")
	}
}

func TestAll_Empty(t *testing.T) {
	result := stream.Of[int]().All(func(n int) bool { return n > 0 })
	if !result {
		t.Error("All on empty: should return true")
	}
}

func TestAll_SomeFail(t *testing.T) {
	result := stream.Of(1, 2, 3, 4, 5).All(func(n int) bool { return n < 3 })
	if result {
		t.Error("All: should return false when some fail")
	}
}

func TestAvg_Empty(t *testing.T) {
	avg := stream.Avg(stream.Of[int]())
	if avg != 0 {
		t.Errorf("Avg empty: expected 0, got %f", avg)
	}
}

func TestAvgBy_Empty(t *testing.T) {
	avg := stream.AvgBy(stream.Of[Product](), func(p Product) float64 { return p.Price })
	if avg != 0 {
		t.Errorf("AvgBy empty: expected 0, got %f", avg)
	}
}

func TestChunk_Zero(t *testing.T) {
	result := stream.Of(1, 2, 3).Chunk(0)
	if result != nil {
		t.Errorf("Chunk(0): expected nil, got %v", result)
	}
}

func TestChunk_Negative(t *testing.T) {
	result := stream.Of(1, 2, 3).Chunk(-1)
	if result != nil {
		t.Errorf("Chunk(-1): expected nil, got %v", result)
	}
}

func TestMap_EarlyBreak(t *testing.T) {
	v, ok := stream.Map(
		stream.Of(1, 2, 3),
		func(n int) string { return fmt.Sprintf("%d", n) },
	).First()
	if !ok || v != "1" {
		t.Errorf("Map early break: expected '1', got %s", v)
	}
}

func TestMapIndexed_EarlyBreak(t *testing.T) {
	v, ok := stream.MapIndexed(
		stream.Of("a", "b", "c"),
		func(i int, s string) string { return fmt.Sprintf("%d:%s", i, s) },
	).First()
	if !ok || v != "0:a" {
		t.Errorf("MapIndexed early break: expected '0:a', got %s", v)
	}
}

func TestFlatMap_EarlyBreak(t *testing.T) {
	v, ok := stream.FlatMap(
		stream.Of(
			User{Name: "Alice", Orders: []Order{{Product: "A"}, {Product: "B"}}},
			User{Name: "Bob", Orders: []Order{{Product: "C"}}},
		),
		func(u User) []Order { return u.Orders },
	).First()
	if !ok || v.Product != "A" {
		t.Errorf("FlatMap early break: expected product A, got %s", v.Product)
	}
}

func TestFlatMap_EarlyBreakMidSlice(t *testing.T) {
	// Break in the middle of a flattened slice
	result := stream.FlatMap(
		stream.Of([]int{1, 2, 3}, []int{4, 5, 6}),
		func(s []int) []int { return s },
	).Take(2).ToSlice()
	if len(result) != 2 || result[0] != 1 || result[1] != 2 {
		t.Errorf("FlatMap mid-slice break: expected [1 2], got %v", result)
	}
}

func TestZip_EarlyBreak(t *testing.T) {
	v, ok := stream.Zip(
		stream.Of(1, 2, 3),
		stream.Of("a", "b", "c"),
	).First()
	if !ok || v.First != 1 || v.Second != "a" {
		t.Errorf("Zip early break: unexpected %v", v)
	}
}

func TestFlatten_EarlyBreak(t *testing.T) {
	v, ok := stream.Flatten(
		stream.Of([]int{1, 2}, []int{3, 4}),
	).First()
	if !ok || v != 1 {
		t.Errorf("Flatten early break: expected 1, got %d", v)
	}
}

func TestFlatten_EarlyBreakMidSlice(t *testing.T) {
	result := stream.Flatten(
		stream.Of([]int{1, 2, 3}, []int{4, 5}),
	).Take(2).ToSlice()
	if len(result) != 2 || result[0] != 1 || result[1] != 2 {
		t.Errorf("Flatten mid-slice break: expected [1 2], got %v", result)
	}
}

func TestEnumerate_EarlyBreak(t *testing.T) {
	v, ok := stream.Enumerate(stream.Of("a", "b", "c")).First()
	if !ok || v.First != 0 || v.Second != "a" {
		t.Errorf("Enumerate early break: unexpected %v", v)
	}
}

func TestCollect2_EarlyBreak(t *testing.T) {
	seq2 := func(yield func(string, int) bool) {
		if !yield("a", 1) {
			return
		}
		if !yield("b", 2) {
			return
		}
		yield("c", 3)
	}
	v, ok := stream.Collect2(seq2).First()
	if !ok || v.First != "a" || v.Second != 1 {
		t.Errorf("Collect2 early break: unexpected %v", v)
	}
}

func TestLast_Empty(t *testing.T) {
	_, ok := stream.Of[int]().Last()
	if ok {
		t.Error("Last on empty: should return false")
	}
}

func TestMinBy_Empty(t *testing.T) {
	_, ok := stream.Of[int]().MinBy(func(a, b int) bool { return a < b })
	if ok {
		t.Error("MinBy on empty: should return false")
	}
}

func TestMaxBy_Empty(t *testing.T) {
	_, ok := stream.Of[int]().MaxBy(func(a, b int) bool { return a < b })
	if ok {
		t.Error("MaxBy on empty: should return false")
	}
}

func TestMin_Empty(t *testing.T) {
	_, ok := stream.Min(stream.Of[int]())
	if ok {
		t.Error("Min on empty: should return false")
	}
}

func TestMax_Empty(t *testing.T) {
	_, ok := stream.Max(stream.Of[int]())
	if ok {
		t.Error("Max on empty: should return false")
	}
}

func TestContains(t *testing.T) {
	if !stream.Of(1, 2, 3).Contains(func(n int) bool { return n == 2 }) {
		t.Error("Contains: should find 2")
	}
	if stream.Of(1, 2, 3).Contains(func(n int) bool { return n == 5 }) {
		t.Error("Contains: should not find 5")
	}
}

func TestToSlice_Empty(t *testing.T) {
	result := stream.Of[int]().ToSlice()
	if result == nil || len(result) != 0 {
		t.Errorf("ToSlice empty: expected non-nil empty slice, got %v", result)
	}
}

func TestSkip_AllElements(t *testing.T) {
	result := stream.Of(1, 2, 3).Skip(5).ToSlice()
	if len(result) != 0 {
		t.Errorf("Skip all: expected empty, got %v", result)
	}
}

func TestTakeWhile_NoneMatch(t *testing.T) {
	result := stream.Of(5, 6, 7).TakeWhile(func(n int) bool { return n < 5 }).ToSlice()
	if len(result) != 0 {
		t.Errorf("TakeWhile none: expected empty, got %v", result)
	}
}

func TestDropWhile_AllMatch(t *testing.T) {
	result := stream.Of(1, 2, 3).DropWhile(func(n int) bool { return n < 10 }).ToSlice()
	if len(result) != 0 {
		t.Errorf("DropWhile all: expected empty, got %v", result)
	}
}

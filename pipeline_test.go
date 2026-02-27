package stream_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nd-forge/stream"
)

// ---------------------------------------------------------------------------
// Pipeline constructor tests
// ---------------------------------------------------------------------------

func TestPipeline_Lazy(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5).ToSlice()
	if len(result) != 5 || result[0] != 1 || result[4] != 5 {
		t.Errorf("Lazy: unexpected %v", result)
	}
}

func TestPipeline_LazyFrom(t *testing.T) {
	s := stream.Of(1, 2, 3)
	p := stream.LazyFrom(s.Iter())
	result := p.Filter(func(n int) bool { return n > 1 }).ToSlice()
	if len(result) != 2 || result[0] != 2 {
		t.Errorf("LazyFrom: unexpected %v", result)
	}
}

func TestPipeline_StreamToLazy(t *testing.T) {
	result := stream.Of(1, 2, 3, 4, 5).
		Lazy().
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(1).
		ToSlice()
	if len(result) != 1 || result[0] != 2 {
		t.Errorf("Stream.Lazy: unexpected %v", result)
	}
}

func TestPipeline_LazyRange(t *testing.T) {
	result := stream.LazyRange(0, 5).ToSlice()
	expected := []int{0, 1, 2, 3, 4}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("LazyRange: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

// ---------------------------------------------------------------------------
// Generator tests (infinite sequences)
// ---------------------------------------------------------------------------

func TestPipeline_Repeat(t *testing.T) {
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

func TestPipeline_RepeatN(t *testing.T) {
	result := stream.RepeatN("x", 3).ToSlice()
	if len(result) != 3 || result[0] != "x" {
		t.Errorf("RepeatN: unexpected %v", result)
	}
}

func TestPipeline_Iterate(t *testing.T) {
	// Powers of 2: 1, 2, 4, 8, 16
	result := stream.Iterate(1, func(n int) int { return n * 2 }).Take(5).ToSlice()
	expected := []int{1, 2, 4, 8, 16}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Iterate: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestPipeline_Naturals(t *testing.T) {
	result := stream.Naturals().Take(5).ToSlice()
	expected := []int{0, 1, 2, 3, 4}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Naturals: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestPipeline_InfiniteFilterTake(t *testing.T) {
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

func TestPipeline_Filter(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	if len(result) != 3 || result[0] != 2 || result[1] != 4 || result[2] != 6 {
		t.Errorf("Pipeline.Filter: unexpected %v", result)
	}
}

func TestPipeline_Reject(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5).
		Reject(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	if len(result) != 3 || result[0] != 1 {
		t.Errorf("Pipeline.Reject: unexpected %v", result)
	}
}

func TestPipeline_Take(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5).Take(3).ToSlice()
	if len(result) != 3 || result[2] != 3 {
		t.Errorf("Pipeline.Take: unexpected %v", result)
	}
}

func TestPipeline_TakeZero(t *testing.T) {
	result := stream.Lazy(1, 2, 3).Take(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("Pipeline.Take(0): expected empty, got %v", result)
	}
}

func TestPipeline_Skip(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5).Skip(3).ToSlice()
	if len(result) != 2 || result[0] != 4 {
		t.Errorf("Pipeline.Skip: unexpected %v", result)
	}
}

func TestPipeline_TakeWhile(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5, 1, 2).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	if len(result) != 3 || result[2] != 3 {
		t.Errorf("Pipeline.TakeWhile: unexpected %v", result)
	}
}

func TestPipeline_DropWhile(t *testing.T) {
	result := stream.Lazy(1, 2, 3, 4, 5, 1, 2).
		DropWhile(func(n int) bool { return n < 4 }).
		ToSlice()
	if len(result) != 4 || result[0] != 4 {
		t.Errorf("Pipeline.DropWhile: unexpected %v", result)
	}
}

func TestPipeline_Distinct(t *testing.T) {
	result := stream.Lazy("a", "b", "a", "c", "b").
		Distinct(func(s string) string { return s }).
		ToSlice()
	if len(result) != 3 {
		t.Errorf("Pipeline.Distinct: expected 3 unique, got %v", result)
	}
}

func TestPipeline_Chain(t *testing.T) {
	p1 := stream.Lazy(1, 2, 3)
	p2 := stream.Lazy(4, 5, 6)
	result := p1.Chain(p2).ToSlice()
	if len(result) != 6 || result[3] != 4 {
		t.Errorf("Pipeline.Chain: unexpected %v", result)
	}
}

func TestPipeline_Sort(t *testing.T) {
	result := stream.Lazy(3, 1, 4, 1, 5, 9).
		Sort(func(a, b int) int { return a - b }).
		ToSlice()
	expected := []int{1, 1, 3, 4, 5, 9}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Pipeline.Sort: expected %d at %d, got %d", expected[i], i, v)
		}
	}
}

func TestPipeline_SortThenTake(t *testing.T) {
	// Sort buffers all, but Take after Sort is still lazy
	result := stream.Lazy(5, 3, 1, 4, 2).
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

func TestPipeline_Reverse(t *testing.T) {
	result := stream.Lazy(1, 2, 3).Reverse().ToSlice()
	if result[0] != 3 || result[2] != 1 {
		t.Errorf("Pipeline.Reverse: unexpected %v", result)
	}
}

func TestPipeline_Peek(t *testing.T) {
	var peeked []int
	result := stream.Lazy(1, 2, 3).
		Peek(func(n int) { peeked = append(peeked, n) }).
		Filter(func(n int) bool { return n > 1 }).
		ToSlice()

	if len(peeked) != 3 {
		t.Errorf("Pipeline.Peek: expected 3 peeks, got %d", len(peeked))
	}
	if len(result) != 2 {
		t.Errorf("Pipeline.Peek: filter should produce 2, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Terminal operation tests
// ---------------------------------------------------------------------------

func TestPipeline_ToStream(t *testing.T) {
	// Pipeline → Stream: can use Stream-only operations like Chunk
	s := stream.Lazy(1, 2, 3, 4, 5).
		Filter(func(n int) bool { return n > 2 }).
		ToStream()

	chunks := s.Chunk(2)
	if len(chunks) != 2 {
		t.Errorf("ToStream.Chunk: expected 2 chunks, got %d", len(chunks))
	}
}

func TestPipeline_Seq(t *testing.T) {
	p := stream.Lazy(1, 2, 3)
	var result []int
	for v := range p.Seq() {
		result = append(result, v)
	}
	if len(result) != 3 {
		t.Errorf("Seq: unexpected %v", result)
	}
}

func TestPipeline_FirstAndLast(t *testing.T) {
	p := stream.Lazy(10, 20, 30)

	first, ok := p.First()
	if !ok || first != 10 {
		t.Errorf("Pipeline.First: expected 10, got %d", first)
	}

	last, ok := p.Last()
	if !ok || last != 30 {
		t.Errorf("Pipeline.Last: expected 30, got %d", last)
	}

	_, ok = stream.Lazy[int]().First()
	if ok {
		t.Error("Pipeline.First on empty should return false")
	}
}

func TestPipeline_Find(t *testing.T) {
	v, ok := stream.Lazy(1, 2, 3, 4, 5).
		Find(func(n int) bool { return n > 3 })
	if !ok || v != 4 {
		t.Errorf("Pipeline.Find: expected 4, got %d", v)
	}
}

func TestPipeline_AnyAllNone(t *testing.T) {
	p := stream.Lazy(1, 2, 3, 4, 5)

	if !p.Any(func(n int) bool { return n > 4 }) {
		t.Error("Pipeline.Any: should find element > 4")
	}
	if !p.All(func(n int) bool { return n > 0 }) {
		t.Error("Pipeline.All: all should be > 0")
	}
	if !p.None(func(n int) bool { return n > 10 }) {
		t.Error("Pipeline.None: no element > 10")
	}
}

func TestPipeline_Count(t *testing.T) {
	n := stream.Lazy(1, 2, 3, 4, 5).
		Filter(func(v int) bool { return v > 2 }).
		Count()
	if n != 3 {
		t.Errorf("Pipeline.Count: expected 3, got %d", n)
	}
}

func TestPipeline_CountBy(t *testing.T) {
	n := stream.Lazy(1, 2, 3, 4, 5).
		CountBy(func(v int) bool { return v%2 == 0 })
	if n != 2 {
		t.Errorf("Pipeline.CountBy: expected 2, got %d", n)
	}
}

func TestPipeline_IsEmpty(t *testing.T) {
	if !stream.Lazy[int]().IsEmpty() {
		t.Error("Pipeline.IsEmpty: empty should return true")
	}
	if stream.Lazy(1).IsEmpty() {
		t.Error("Pipeline.IsEmpty: non-empty should return false")
	}
}

func TestPipeline_MinByMaxBy(t *testing.T) {
	type item struct {
		name  string
		value int
	}
	p := stream.Lazy(
		item{"a", 30}, item{"b", 10}, item{"c", 20},
	)

	min, ok := p.MinBy(func(a, b item) bool { return a.value < b.value })
	if !ok || min.name != "b" {
		t.Errorf("Pipeline.MinBy: expected b, got %s", min.name)
	}

	max, ok := p.MaxBy(func(a, b item) bool { return a.value < b.value })
	if !ok || max.name != "a" {
		t.Errorf("Pipeline.MaxBy: expected a, got %s", max.name)
	}
}

func TestPipeline_Reduce(t *testing.T) {
	sum := stream.Lazy(1, 2, 3, 4, 5).
		Reduce(0, func(acc, v int) int { return acc + v })
	if sum != 15 {
		t.Errorf("Pipeline.Reduce: expected 15, got %d", sum)
	}
}

func TestPipeline_ForEach(t *testing.T) {
	var result []int
	stream.Lazy(1, 2, 3).ForEach(func(v int) { result = append(result, v*2) })
	if len(result) != 3 || result[0] != 2 || result[2] != 6 {
		t.Errorf("Pipeline.ForEach: unexpected %v", result)
	}
}

func TestPipeline_ForEachIndexed(t *testing.T) {
	var indices []int
	stream.Lazy("a", "b", "c").ForEachIndexed(func(i int, _ string) {
		indices = append(indices, i)
	})
	if len(indices) != 3 || indices[2] != 2 {
		t.Errorf("Pipeline.ForEachIndexed: unexpected %v", indices)
	}
}

// ---------------------------------------------------------------------------
// Pipeline transform function tests (type-changing)
// ---------------------------------------------------------------------------

func TestPipeMap(t *testing.T) {
	result := stream.PipeMap(
		stream.Lazy(1, 2, 3),
		func(n int) string { return fmt.Sprintf("item_%d", n) },
	).ToSlice()

	if result[0] != "item_1" || result[2] != "item_3" {
		t.Errorf("PipeMap: unexpected %v", result)
	}
}

func TestPipeMapIndexed(t *testing.T) {
	result := stream.PipeMapIndexed(
		stream.Lazy("a", "b", "c"),
		func(i int, s string) string { return fmt.Sprintf("%d:%s", i, s) },
	).ToSlice()

	if result[0] != "0:a" || result[2] != "2:c" {
		t.Errorf("PipeMapIndexed: unexpected %v", result)
	}
}

func TestPipeFlatMap(t *testing.T) {
	users := stream.Lazy(
		User{Name: "Alice", Orders: []Order{{Product: "A"}, {Product: "B"}}},
		User{Name: "Bob", Orders: []Order{{Product: "C"}}},
	)

	orders := stream.PipeFlatMap(users, func(u User) []Order { return u.Orders })
	if orders.Count() != 3 {
		t.Errorf("PipeFlatMap: expected 3, got %d", orders.Count())
	}
}

func TestPipeReduce(t *testing.T) {
	total := stream.PipeReduce(
		stream.Lazy(
			Order{Amount: 100, Discount: 0.1},
			Order{Amount: 200, Discount: 0.2},
		),
		0.0,
		func(acc float64, o Order) float64 { return acc + o.Amount*(1-o.Discount) },
	)
	// 90 + 160 = 250
	if total != 250 {
		t.Errorf("PipeReduce: expected 250, got %.1f", total)
	}
}

func TestPipeGroupBy(t *testing.T) {
	products := stream.Lazy(
		Product{Name: "Laptop", Category: "Electronics"},
		Product{Name: "T-Shirt", Category: "Clothing"},
		Product{Name: "Keyboard", Category: "Electronics"},
	)

	groups := stream.PipeGroupBy(products, func(p Product) string { return p.Category })
	if len(groups["Electronics"]) != 2 {
		t.Errorf("PipeGroupBy: expected 2 Electronics, got %d", len(groups["Electronics"]))
	}
}

func TestPipeAssociate(t *testing.T) {
	m := stream.PipeAssociate(
		stream.Lazy(User{Name: "Alice", Age: 30}, User{Name: "Bob", Age: 25}),
		func(u User) (string, int) { return u.Name, u.Age },
	)
	if m["Alice"] != 30 || m["Bob"] != 25 {
		t.Errorf("PipeAssociate: unexpected %v", m)
	}
}

func TestPipeZip(t *testing.T) {
	names := stream.Lazy("Alice", "Bob", "Charlie")
	scores := stream.Lazy(85.0, 92.0, 78.0)

	pairs := stream.PipeZip(names, scores).ToSlice()
	if len(pairs) != 3 || pairs[0].First != "Alice" || pairs[0].Second != 85.0 {
		t.Errorf("PipeZip: unexpected %v", pairs)
	}
}

func TestPipeZip_UnevenLengths(t *testing.T) {
	p1 := stream.Lazy(1, 2, 3, 4, 5)
	p2 := stream.Lazy("a", "b")

	pairs := stream.PipeZip(p1, p2).ToSlice()
	if len(pairs) != 2 {
		t.Errorf("PipeZip uneven: expected 2, got %d", len(pairs))
	}
}

func TestPipeFlatten(t *testing.T) {
	result := stream.PipeFlatten(
		stream.Lazy([]int{1, 2}, []int{3, 4}, []int{5}),
	).ToSlice()
	if len(result) != 5 || result[0] != 1 || result[4] != 5 {
		t.Errorf("PipeFlatten: unexpected %v", result)
	}
}

func TestPipeToMap(t *testing.T) {
	m := stream.PipeToMap(stream.PipeZip(
		stream.Lazy("a", "b", "c"),
		stream.Lazy(1, 2, 3),
	))
	if m["a"] != 1 || m["c"] != 3 {
		t.Errorf("PipeToMap: unexpected %v", m)
	}
}

func TestPipeEnumerate(t *testing.T) {
	result := stream.PipeEnumerate(stream.Lazy("a", "b", "c")).ToSlice()
	if len(result) != 3 || result[0].First != 0 || result[0].Second != "a" {
		t.Errorf("PipeEnumerate: unexpected %v", result)
	}
	if result[2].First != 2 || result[2].Second != "c" {
		t.Errorf("PipeEnumerate: unexpected last %v", result[2])
	}
}

// ---------------------------------------------------------------------------
// Lazy evaluation proof tests
// ---------------------------------------------------------------------------

func TestPipeline_ShortCircuit_Take(t *testing.T) {
	// Prove that Take stops iteration early
	evaluated := 0
	result := stream.Lazy(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
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

func TestPipeline_ShortCircuit_Find(t *testing.T) {
	evaluated := 0
	v, ok := stream.Lazy(1, 2, 3, 4, 5).
		Peek(func(int) { evaluated++ }).
		Find(func(n int) bool { return n == 3 })

	if !ok || v != 3 {
		t.Errorf("ShortCircuit.Find: expected 3, got %d", v)
	}
	if evaluated != 3 {
		t.Errorf("ShortCircuit.Find: expected 3 evaluations, got %d", evaluated)
	}
}

func TestPipeline_ShortCircuit_Any(t *testing.T) {
	evaluated := 0
	found := stream.Lazy(1, 2, 3, 4, 5).
		Peek(func(int) { evaluated++ }).
		Any(func(n int) bool { return n == 2 })

	if !found {
		t.Error("ShortCircuit.Any: should find 2")
	}
	if evaluated != 2 {
		t.Errorf("ShortCircuit.Any: expected 2 evaluations, got %d", evaluated)
	}
}

func TestPipeline_NoAllocation_FilterTake(t *testing.T) {
	// Pipeline with Filter + Take doesn't create intermediate slices
	// We prove this by applying to a large range — if eager, it would
	// allocate a huge intermediate slice. With lazy eval, only 3 elements
	// are ever processed.
	evaluated := 0
	result := stream.LazyRange(0, 1_000_000).
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

func TestPipeline_Reusable(t *testing.T) {
	p := stream.Lazy(1, 2, 3, 4, 5).Filter(func(n int) bool { return n > 2 })

	// Calling ToSlice() multiple times should return same result
	r1 := p.ToSlice()
	r2 := p.ToSlice()

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
// Pipeline integration tests — practical examples
// ---------------------------------------------------------------------------

func TestPipeline_Integration_ProductAnalysis(t *testing.T) {
	products := stream.Lazy(
		Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
		Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
		Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
		Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
		Product{Name: "Keyboard", Category: "Electronics", Price: 75, InStock: true},
	)

	// Lazy filter → Sort (buffers) → Take (lazy again)
	top2 := products.
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

	if len(top2) != 2 || top2[0].Name != "Laptop" || top2[1].Name != "Jeans" {
		t.Errorf("Integration.ProductAnalysis: unexpected %v", top2)
	}
}

func TestPipeline_Integration_UserOrders(t *testing.T) {
	users := stream.Lazy(
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

	total := stream.PipeReduce(
		stream.PipeFlatMap(
			users.
				Filter(func(u User) bool { return u.IsActive }).
				Filter(func(u User) bool { return u.Age >= 20 }),
			func(u User) []Order { return u.Orders },
		),
		0.0,
		func(acc float64, o Order) float64 { return acc + o.Amount*(1-o.Discount) },
	)

	if total != 290 {
		t.Errorf("Integration.UserOrders: expected 290, got %.1f", total)
	}
}

func TestPipeline_Integration_TextProcessing(t *testing.T) {
	result := stream.PipeMap(
		stream.Lazy("hello", "world", "hello", "go", "stream", "go", "chain").
			Distinct(func(s string) string { return s }).
			Sort(func(a, b string) int { return len(a) - len(b) }),
		func(s string) string { return strings.ToUpper(s) },
	).ToSlice()

	if result[0] != "GO" {
		t.Errorf("Integration.TextProcessing: expected GO first, got %s", result[0])
	}
}

func TestPipeline_Integration_InfiniteSequence(t *testing.T) {
	// Fibonacci sequence using Pair
	fib := stream.PipeMap(
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

func TestPipeline_Integration_LazyToEagerBridge(t *testing.T) {
	// Use lazy Pipeline for filtering, then switch to eager Stream for Shuffle + Chunk
	result := stream.Naturals().
		Filter(func(n int) bool { return n%3 == 0 }).
		Take(12).
		ToStream().  // Pipeline → Stream
		Chunk(4)     // Stream-only operation

	if len(result) != 3 {
		t.Errorf("LazyToEager: expected 3 chunks, got %d", len(result))
	}
	if result[0].Count() != 4 {
		t.Errorf("LazyToEager: first chunk should have 4 elements, got %d", result[0].Count())
	}
}

func TestPipeline_Integration_CollectFromIter(t *testing.T) {
	// Show full interop: Stream → iter.Seq → Pipeline → Stream
	s := stream.Of(10, 20, 30, 40, 50)
	seq := s.Iter()                                                  // Stream → iter.Seq
	p := stream.LazyFrom(seq)                                        // iter.Seq → Pipeline
	result := p.Filter(func(n int) bool { return n > 20 }).ToStream() // Pipeline → Stream

	if result.Count() != 3 || result.ToSlice()[0] != 30 {
		t.Errorf("CollectFromIter: unexpected %v", result.ToSlice())
	}
}

# stream

[![Go Reference](https://pkg.go.dev/badge/github.com/nd-forge/stream.svg)](https://pkg.go.dev/github.com/nd-forge/stream)
[![CI](https://github.com/nd-forge/stream/actions/workflows/check-code.yml/badge.svg)](https://github.com/nd-forge/stream/actions/workflows/check-code.yml)
[![coverage](https://img.shields.io/badge/coverage-100.0%25-brightgreen)](https://github.com/nd-forge/stream/actions/workflows/check-code.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nd-forge/stream)](https://goreportcard.com/report/github.com/nd-forge/stream)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.23-blue)](https://go.dev/)
[![Go Playground](https://img.shields.io/badge/Go-Playground-00ADD8?logo=go&logoColor=white)](https://go.dev/play/p/QaQ_rdqYn1Y)

**English** | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

A Go generic stream processing library. Chainable collection operations for filter, map, sort, group, and more — with **lazy evaluation** by default. [Try it online!](https://go.dev/play/p/QaQ_rdqYn1Y)

All operations are lazy. Operations that require full data (Sort, Reverse, Shuffle, TakeLast, Chunk, Partition) buffer internally and resume lazy iteration automatically.

## Requirements

- Go 1.23 or later (uses `iter.Seq[T]` and range-over-function)

## Install

```bash
go get github.com/nd-forge/stream
```

## Quick Start

```go
import "github.com/nd-forge/stream"

// Method chaining for same-type operations
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// Top-level functions for type-changing operations
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## Design

### Lazy by Default

All operations build a lazy pipeline internally using `iter.Seq[T]`. No intermediate slices are allocated until a terminal operation (`ToSlice`, `ForEach`, `Reduce`, etc.) is called.

Operations that inherently need all data — `Sort`, `Reverse`, `Shuffle`, `TakeLast`, `Chunk`, `Partition` — buffer internally, then resume lazy iteration for subsequent operations.

### Type Parameters

Go does not allow methods to introduce new type parameters. This library separates:

| Kind | Implementation | Signature |
|---|---|---|
| Type-preserving (Filter, Sort, Take...) | **Methods** — chainable | `Stream[T] → Stream[T]` |
| Type-changing (Map, FlatMap, GroupBy...) | **Top-level functions** | `Stream[T] → Stream[U]` |

## API

### Constructors

| Function | Description |
|---|---|
| `Of[T](items ...T)` | Create from variadic args |
| `From[T](items []T)` | Create from slice (copies) |
| `Range(start, end)` | Create integer sequence `[start, end)` |
| `Generate[T](n, fn)` | Create n elements with generator |

### Generators (Infinite Sequences)

| Function | Description |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | Infinite repetition of value |
| `RepeatN(value, n)` | Repeat value n times |

### Chainable Methods

Operations that return `Stream[T]` and can be chained.

| Method | Description |
|---|---|
| `Filter(predicate)` | Keep elements matching predicate |
| `Reject(predicate)` | Remove elements matching predicate |
| `Sort(cmp)` | Sort by comparison function |
| `Reverse()` | Reverse order |
| `Take(n)` / `TakeLast(n)` | First / last n elements |
| `Skip(n)` | Remove first n elements |
| `TakeWhile(pred)` / `DropWhile(pred)` | Take / skip from start while true |
| `Distinct(key)` | Remove duplicates by key |
| `Shuffle()` | Random order |
| `Peek(fn)` | Execute side effect without modifying |
| `Chain(others...)` | Concatenate multiple streams |

> `Sort`, `Reverse`, `Shuffle`, `TakeLast` buffer all elements internally.

### Terminal Operations

| Method | Returns |
|---|---|
| `ToSlice()` | `[]T` |
| `First()` / `Last()` | `(T, bool)` |
| `Find(predicate)` | `(T, bool)` |
| `Reduce(initial, fn)` | `T` |
| `Any(pred)` / `All(pred)` / `None(pred)` | `bool` |
| `Count()` / `CountBy(pred)` | `int` |
| `IsEmpty()` | `bool` |
| `Contains(predicate)` | `bool` |
| `MinBy(less)` / `MaxBy(less)` | `(T, bool)` |
| `Partition(pred)` | `(Stream[T], Stream[T])` |
| `Chunk(size)` | `[]Stream[T]` |
| `ForEach(fn)` / `ForEachIndexed(fn)` | — |
| `Seq()` | `iter.Seq[T]` |

### Transform Functions

Top-level functions for type-changing operations.

| Function | Description |
|---|---|
| `Map(s, fn)` | Transform `T → U` |
| `MapIndexed(s, fn)` | Transform with index |
| `FlatMap(s, fn)` | Transform and flatten `T → []U` |
| `Reduce(s, initial, fn)` | Fold into different type `T → U` |
| `GroupBy(s, key)` | Group by key `→ map[K][]T` |
| `Associate(s, fn)` | Build map `→ map[K]V` |
| `Zip(s1, s2)` | Pair two streams `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | Flatten `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | Convert `Stream[Pair[K,V]] → map[K]V` |
| `Enumerate(s)` | Add index `→ Stream[Pair[int,T]]` |

### Numeric Functions

Specialized operations for numeric streams (`int`, `float64`, etc.).

| Function | Description |
|---|---|
| `Sum(s)` / `Avg(s)` | Sum / average |
| `Min(s)` / `Max(s)` | Minimum / maximum |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | Sum / average of extracted values |

### iter.Seq Bridge

| Function | Description |
|---|---|
| `Seq()` | `Stream[T]` → `iter.Seq[T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## Examples

Types used in examples below:

```go
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
```

### Filter, Sort, Take [[playground]](https://go.dev/play/p/W5fo1cfb_VA)

```go
products := stream.Of(
    Product{Name: "Laptop", Category: "Electronics", Price: 1200, InStock: true},
    Product{Name: "T-Shirt", Category: "Clothing", Price: 25, InStock: true},
    Product{Name: "Headphones", Category: "Electronics", Price: 150, InStock: false},
    Product{Name: "Jeans", Category: "Clothing", Price: 80, InStock: true},
)

// In-stock products, sorted by price descending, top 3
top3 := products.
    Filter(func(p Product) bool { return p.InStock }).
    Sort(func(a, b Product) int {
        if a.Price > b.Price { return -1 }
        if a.Price < b.Price { return 1 }
        return 0
    }).
    Take(3).
    ToSlice()
```

### Map and FlatMap [[playground]](https://go.dev/play/p/GyrjAXMYnZD)

```go
// Extract product names
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// Flatten nested orders from users
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy and Aggregate [[playground]](https://go.dev/play/p/nzL3i-4Kgj4)

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f count=%d\n", category, total, len(group))
}
```

### Partition and Chunk [[playground]](https://go.dev/play/p/c9KhEbkeKat)

```go
// Split by condition
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// Batch processing
batches := stream.From(items).Chunk(100)
for _, batch := range batches {
    api.Send(batch.ToSlice())
}
```

### Zip [[playground]](https://go.dev/play/p/QUdJ_GonDTa)

```go
names := stream.Of("Alice", "Bob", "Charlie")
scores := stream.Of(85.0, 92.0, 78.0)

pairs := stream.Zip(names, scores).ToSlice()
// [{Alice 85}, {Bob 92}, {Charlie 78}]
```

### Infinite Sequences [[playground]](https://go.dev/play/p/9prnIg-NjtF)

```go
// First 5 even natural numbers
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// Fibonacci sequence
fib := stream.Map(
    stream.Iterate(
        stream.Pair[int, int]{First: 0, Second: 1},
        func(p stream.Pair[int, int]) stream.Pair[int, int] {
            return stream.Pair[int, int]{First: p.Second, Second: p.First + p.Second}
        },
    ).Take(10),
    func(p stream.Pair[int, int]) int { return p.First },
).ToSlice()
// [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]

// Lazy evaluation: processes only 2,001 of 1,000,000 elements
result := stream.Range(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### iter.Seq Bridge [[playground]](https://go.dev/play/p/i-rNvVdYut5)

```go
// Standard library interop
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range support
for v := range stream.Of(1, 2, 3).Seq() {
    fmt.Println(v)
}
```

## Benchmark

10,000 `int` elements — Apple M1. Compared with [samber/lo](https://github.com/samber/lo).

### Filter + Take (lazy evaluation advantage)

```
Benchmark                    ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         111        248          5
Stream                         346        528         15   ← 3.1x native
lo                          25,000     81,920          1   ← 72x slower than Stream
```

Stream's lazy evaluation short-circuits — only processes elements until `Take` is satisfied. lo must filter **all 10,000 elements** first, then take.

### Chained: Filter → Map → Take 5

```
Benchmark                    ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         171        260          9
Stream                         384        544         19   ← 2.2x native
lo                          64,600    152,601      3,336   ← 168x slower than Stream
```

### Full scan (no early termination)

```
Benchmark                    ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native Filter              19,400    128,249         16
lo     Filter              25,600     81,920          1
Stream Filter              42,200    128,409         22

Native Reduce               3,300          0          0
lo     Reduce               9,700          0          0
Stream Reduce              23,400         64          2
```

For full scans without early termination, lo is faster than Stream due to lower abstraction overhead. Stream's advantage shines in pipelines with `Take`, `First`, `Find`, or other short-circuiting operations.

> Run `cd _benchmark && go test -bench=. -benchmem ./...` to reproduce.

## License

MIT

# stream

**English** | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

A Go generic stream processing library. Chainable collection operations for filter, map, sort, group, and more.

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

### Transform Functions

Top-level functions for type-changing operations.

| Function | Description |
|---|---|
| `Map(s, fn)` | Transform `T → U` |
| `MapIndexed(s, fn)` | Transform with index |
| `FlatMap(s, fn)` | Transform and flatten `T → []U` |
| `Reduce(s, initial, fn)` | Fold into different type `T → U` |
| `GroupBy(s, key)` | Group by key `→ map[K]Stream[T]` |
| `Associate(s, fn)` | Build map `→ map[K]V` |
| `Zip(s1, s2)` | Pair two streams `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | Flatten `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | Convert `Stream[Pair[K,V]] → map[K]V` |

### Numeric Functions

Specialized operations for numeric streams (`int`, `float64`, etc.).

| Function | Description |
|---|---|
| `Sum(s)` / `Avg(s)` | Sum / average |
| `Min(s)` / `Max(s)` | Minimum / maximum |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | Sum / average of extracted values |

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

### Filter, Sort, Take

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

### Map and FlatMap

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

### GroupBy and Aggregate

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(group, func(p Product) float64 { return p.Price })
    avg := stream.AvgBy(group, func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f avg=$%.2f count=%d\n", category, total, avg, group.Count())
}
```

### Partition and Chunk

```go
// Split by condition
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// Batch processing
batches := stream.From(items).Chunk(100)
for _, batch := range batches {
    api.Send(batch.ToSlice())
}
```

### Zip

```go
names := stream.Of("Alice", "Bob", "Charlie")
scores := stream.Of(85.0, 92.0, 78.0)

pairs := stream.Zip(names, scores).ToSlice()
// [{Alice 85}, {Bob 92}, {Charlie 78}]
```

## License

MIT

# stream

[English](README.md) | [日本語](README_ja.md) | **中文** | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

Go 泛型流处理库。支持链式调用的集合操作：filter、map、sort、group 等 — 默认采用**惰性求值**。

所有操作都是惰性的。需要完整数据的操作（Sort、Reverse、Shuffle、TakeLast、Chunk、Partition）会在内部缓冲，然后自动恢复惰性迭代。

## 安装

```bash
go get github.com/nd-forge/stream
```

## 快速开始

```go
import "github.com/nd-forge/stream"

// 方法链式调用（保持类型不变的操作）
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// 顶层函数（类型转换操作）
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## 设计理念

### 默认惰性求值

所有操作在内部使用 `iter.Seq[T]` 构建惰性管道。在调用终端操作（`ToSlice`、`ForEach`、`Reduce` 等）之前，不会分配任何中间切片。

本质上需要全部数据的操作 — `Sort`、`Reverse`、`Shuffle`、`TakeLast`、`Chunk`、`Partition` — 会在内部缓冲，然后对后续操作恢复惰性迭代。

### 类型参数

Go 不允许方法引入新的类型参数。本库做了以下分离：

| 类别 | 实现方式 | 签名 |
|---|---|---|
| 保持类型的操作 (Filter, Sort, Take...) | **方法** — 可链式调用 | `Stream[T] → Stream[T]` |
| 类型转换的操作 (Map, FlatMap, GroupBy...) | **顶层函数** | `Stream[T] → Stream[U]` |

## API

### 构造器

| 函数 | 说明 |
|---|---|
| `Of[T](items ...T)` | 从可变参数创建 |
| `From[T](items []T)` | 从切片创建（复制） |
| `Range(start, end)` | 创建整数序列 `[start, end)` |
| `Generate[T](n, fn)` | 用生成函数创建 n 个元素 |

### 生成器（无限序列）

| 函数 | 说明 |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | 无限重复值 |
| `RepeatN(value, n)` | 重复值 n 次 |

### 可链式方法

返回 `Stream[T]`，支持链式调用。

| 方法 | 说明 |
|---|---|
| `Filter(predicate)` | 保留匹配的元素 |
| `Reject(predicate)` | 排除匹配的元素 |
| `Sort(cmp)` | 按比较函数排序 |
| `Reverse()` | 反转顺序 |
| `Take(n)` / `TakeLast(n)` | 前 / 后 n 个元素 |
| `Skip(n)` | 跳过前 n 个元素 |
| `TakeWhile(pred)` / `DropWhile(pred)` | 条件式获取 / 跳过 |
| `Distinct(key)` | 按键去重 |
| `Shuffle()` | 随机排序 |
| `Peek(fn)` | 执行副作用（不修改） |
| `Chain(others...)` | 连接多个流 |

> `Sort`、`Reverse`、`Shuffle`、`TakeLast` 会在内部缓冲所有元素。

### 终端操作

| 方法 | 返回值 |
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

### 转换函数

用于类型转换的顶层函数。

| 函数 | 说明 |
|---|---|
| `Map(s, fn)` | 转换 `T → U` |
| `MapIndexed(s, fn)` | 带索引转换 |
| `FlatMap(s, fn)` | 转换并展平 `T → []U` |
| `Reduce(s, initial, fn)` | 折叠为不同类型 `T → U` |
| `GroupBy(s, key)` | 按键分组 `→ map[K][]T` |
| `Associate(s, fn)` | 构建映射 `→ map[K]V` |
| `Zip(s1, s2)` | 配对两个流 `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 展平 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 转换 `Stream[Pair[K,V]] → map[K]V` |
| `Enumerate(s)` | 添加索引 `→ Stream[Pair[int,T]]` |

### 数值函数

数值流的专用操作（`int`、`float64` 等）。

| 函数 | 说明 |
|---|---|
| `Sum(s)` / `Avg(s)` | 求和 / 平均 |
| `Min(s)` / `Max(s)` | 最小 / 最大 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 提取值的求和 / 平均 |

### iter.Seq 桥接

| 函数 | 说明 |
|---|---|
| `Seq()` | `Stream[T]` → `iter.Seq[T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## 示例

以下示例中使用的类型：

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

// 有库存的商品，按价格降序，取前3个
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

### Map 和 FlatMap

```go
// 提取商品名称
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// 展平用户的嵌套订单
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy 和聚合

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f count=%d\n", category, total, len(group))
}
```

### Partition 和 Chunk

```go
// 按条件分割
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// 批量处理
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

### 无限序列

```go
// 前5个偶数自然数
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// 斐波那契数列
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

// 惰性求值：100万元素中只处理2,001个
result := stream.Range(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### iter.Seq 桥接

```go
// 标准库互操作
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range 支持
for v := range stream.Of(1, 2, 3).Seq() {
    fmt.Println(v)
}
```

## 基准测试

10,000 个 `int` 元素，`Filter(偶数)` 然后 `Take(10)` — Apple M1：

```
基准测试               ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
StreamFilterTake         315      464       13   ← 惰性：原生的2.5倍
```

惰性求值会短路 — 只处理元素直到 `Take` 满足为止。

全量扫描（无提前终止）：

```
基准测试               ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> 运行 `go test -bench=. -benchmem ./...` 来复现。

## License

MIT

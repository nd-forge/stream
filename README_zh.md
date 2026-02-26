# stream

[English](README.md) | [日本語](README_ja.md) | **中文** | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

Go 泛型流处理库。支持链式调用的集合操作：filter、map、sort、group 等。

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

### 转换函数

用于类型转换的顶层函数。

| 函数 | 说明 |
|---|---|
| `Map(s, fn)` | 转换 `T → U` |
| `MapIndexed(s, fn)` | 带索引转换 |
| `FlatMap(s, fn)` | 转换并展平 `T → []U` |
| `Reduce(s, initial, fn)` | 折叠为不同类型 `T → U` |
| `GroupBy(s, key)` | 按键分组 `→ map[K]Stream[T]` |
| `Associate(s, fn)` | 构建映射 `→ map[K]V` |
| `Zip(s1, s2)` | 配对两个流 `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 展平 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 转换 `Stream[Pair[K,V]] → map[K]V` |

### 数值函数

数值流的专用操作（`int`、`float64` 等）。

| 函数 | 说明 |
|---|---|
| `Sum(s)` / `Avg(s)` | 求和 / 平均 |
| `Min(s)` / `Max(s)` | 最小 / 最大 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 提取值的求和 / 平均 |

## 示例

### Filter, Sort, Take

```go
top3 := trades.
    Filter(func(t Trade) bool { return !t.IsOpen }).
    Filter(func(t Trade) bool { return t.PnL > 0 }).
    Sort(func(a, b Trade) int {
        if a.PnL > b.PnL { return -1 }
        if a.PnL < b.PnL { return 1 }
        return 0
    }).
    Take(3).
    ToSlice()
```

### Map 和 FlatMap

```go
// 提取名称
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()

// 展平嵌套订单
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy 和聚合

```go
bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })

for symbol, group := range bySymbol {
    totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
    avgPnL := stream.AvgBy(group, func(t Trade) float64 { return t.PnL })
    fmt.Printf("%s: total=%.2f avg=%.2f count=%d\n", symbol, totalPnL, avgPnL, group.Count())
}
```

### Partition 和 Chunk

```go
// 按条件分割
profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
winRate := float64(profit.Count()) / float64(trades.Count()) * 100

// 批量处理
batches := stream.From(items).Chunk(100)
for _, batch := range batches {
    api.Send(batch.ToSlice())
}
```

### Zip

```go
pairs := stream.Zip(
    stream.Of("2024-01", "2024-02", "2024-03"),
    stream.Of(1.05, 1.06, 1.04),
).ToSlice()
// [{2024-01 1.05}, {2024-02 1.06}, {2024-03 1.04}]
```

## License

MIT

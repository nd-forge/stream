# stream

[English](README.md) | **日本語** | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

Go ジェネリクス対応のストリーム処理ライブラリ。filter, map, sort, group などのコレクション操作をメソッドチェーンで記述できます。**即時評価** (`Stream`) と **遅延評価** (`Pipeline`) の両方をサポート。

## インストール

```bash
go get github.com/nd-forge/stream
```

## クイックスタート

```go
import "github.com/nd-forge/stream"

// メソッドチェーンで型を保持する操作
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// トップレベル関数で型変換する操作
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## 設計方針

### 即時評価 vs 遅延評価

| 型 | 評価方式 | 最適な用途 |
|---|---|---|
| `Stream[T]` | **即時評価** — 中間スライスを生成 | 小〜中規模データ、ランダムアクセス (`Shuffle`, `Chunk`, `TakeLast`) |
| `Pipeline[T]` | **遅延評価** — 中間割り当てゼロ | 大規模データ、無限シーケンス、早期終了 (`Filter+Take`) |

自由に切替可能: `stream.Lazy()` / `pipeline.ToStream()`

### 型パラメータ

Go ではメソッドに新しい型パラメータを追加できません。本ライブラリでは以下のように分離しています:

| 種類 | 実装方法 | シグネチャ |
|---|---|---|
| 型を保持する操作 (Filter, Sort, Take...) | **メソッド** — チェーン可能 | `Stream[T] → Stream[T]` |
| 型を変換する操作 (Map, FlatMap, GroupBy...) | **トップレベル関数** | `Stream[T] → Stream[U]` |

## API

### コンストラクタ

| 関数 | 説明 |
|---|---|
| `Of[T](items ...T)` | 可変長引数から作成 |
| `From[T](items []T)` | スライスから作成（コピー） |
| `Range(start, end)` | 整数の連番 `[start, end)` |
| `Generate[T](n, fn)` | 生成関数で n 要素を作成 |

### チェーン可能メソッド

`Stream[T]` を返し、チェーンできる操作。

| メソッド | 説明 |
|---|---|
| `Filter(predicate)` | 条件に合う要素のみ保持 |
| `Reject(predicate)` | 条件に合う要素を除外 |
| `Sort(cmp)` | 比較関数でソート |
| `Reverse()` | 逆順 |
| `Take(n)` / `TakeLast(n)` | 先頭 / 末尾 n 件 |
| `Skip(n)` | 先頭 n 件をスキップ |
| `TakeWhile(pred)` / `DropWhile(pred)` | 条件ベースの取得 / スキップ |
| `Distinct(key)` | キーによる重複排除 |
| `Shuffle()` | ランダム並替 |
| `Peek(fn)` | 副作用を実行（変更なし） |

### 終端操作

| メソッド | 戻り値 |
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

### トップレベル関数（型変換）

| 関数 | 説明 |
|---|---|
| `Map(s, fn)` | 各要素を変換 `T → U` |
| `MapIndexed(s, fn)` | インデックス付き変換 |
| `FlatMap(s, fn)` | 変換して平坦化 `T → []U` |
| `Reduce(s, initial, fn)` | 異なる型への畳み込み `T → U` |
| `GroupBy(s, key)` | キーでグルーピング `→ map[K]Stream[T]` |
| `Associate(s, fn)` | マップ化 `→ map[K]V` |
| `Zip(s1, s2)` | 2つの Stream をペアに `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 平坦化 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 変換 `Stream[Pair[K,V]] → map[K]V` |

### 数値操作

数値型の Stream 向け特殊操作（`int`, `float64` など）。

| 関数 | 説明 |
|---|---|
| `Sum(s)` / `Avg(s)` | 合計 / 平均 |
| `Min(s)` / `Max(s)` | 最小 / 最大 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 抽出値の合計 / 平均 |

### Pipeline（遅延評価）

#### コンストラクタ

| 関数 | 説明 |
|---|---|
| `Lazy[T](items ...T)` | 可変長引数から遅延パイプライン作成 |
| `LazyFrom[T](seq)` | `iter.Seq[T]` から作成 |
| `LazyRange(start, end)` | 遅延整数列 `[start, end)` |
| `stream.Lazy()` | `Stream[T]` → `Pipeline[T]` 変換 |

#### ジェネレータ（無限シーケンス）

| 関数 | 説明 |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | 値の無限繰り返し |
| `RepeatN(value, n)` | 値を n 回繰り返し |

#### チェーン可能メソッド

Stream と同じ API: `Filter`, `Reject`, `Sort`, `Reverse`, `Take`, `Skip`, `TakeWhile`, `DropWhile`, `Distinct`, `Peek`, `Chain`

#### 終端操作

Stream と同じ: `ToSlice`, `First`, `Last`, `Find`, `Reduce`, `Any`, `All`, `None`, `Count`, `CountBy`, `IsEmpty`, `Contains`, `MinBy`, `MaxBy`, `ForEach`, `ForEachIndexed`

追加: `ToStream()`（即時 Stream に変換）, `Seq()`（`iter.Seq[T]` を取得）

#### トップレベル関数（型変換）

| 関数 | 説明 |
|---|---|
| `PipeMap(p, fn)` | 変換 `T → U` |
| `PipeMapIndexed(p, fn)` | インデックス付き変換 |
| `PipeFlatMap(p, fn)` | 変換して平坦化 `T → []U` |
| `PipeReduce(p, initial, fn)` | 異なる型への畳み込み `T → U` |
| `PipeGroupBy(p, key)` | キーでグルーピング `→ map[K][]T` |
| `PipeAssociate(p, fn)` | マップ化 `→ map[K]V` |
| `PipeZip(p1, p2)` | 2つのパイプラインをペアに `→ Pipeline[Pair[T,U]]` |
| `PipeFlatten(p)` | 平坦化 `Pipeline[[]T] → Pipeline[T]` |
| `PipeToMap(p)` | 変換 `Pipeline[Pair[K,V]] → map[K]V` |
| `PipeEnumerate(p)` | インデックス付加 `→ Pipeline[Pair[int,T]]` |

### iter.Seq ブリッジ

| 関数 | 説明 |
|---|---|
| `stream.Iter()` | `Stream[T]` → `iter.Seq[T]` |
| `stream.Iter2()` | `Stream[T]` → `iter.Seq2[int, T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## 使用例

以下の例で使用する型:

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

// 在庫ありの商品を価格降順で上位3件
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

### Map と FlatMap

```go
// 商品名の抽出
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// ユーザーのネストされた注文を平坦化
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy と集計

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(group, func(p Product) float64 { return p.Price })
    avg := stream.AvgBy(group, func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f avg=$%.2f count=%d\n", category, total, avg, group.Count())
}
```

### Partition と Chunk

```go
// 条件で分割
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// バッチ処理
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

### Pipeline（遅延評価）

```go
// 無限シーケンス: 最初の5つの偶数
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// フィボナッチ数列
fib := stream.PipeMap(
    stream.Iterate(
        stream.Pair[int, int]{First: 0, Second: 1},
        func(p stream.Pair[int, int]) stream.Pair[int, int] {
            return stream.Pair[int, int]{First: p.Second, Second: p.First + p.Second}
        },
    ).Take(10),
    func(p stream.Pair[int, int]) int { return p.First },
).ToSlice()
// [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]

// 遅延評価: 100万要素のうち2,001個だけ処理
result := stream.LazyRange(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### 即時評価と遅延評価の切替

```go
// Stream → Pipeline（遅延処理へ）
result := stream.Of(items...).Lazy().Filter(pred).Take(10).ToSlice()

// Pipeline → Stream（即時評価専用操作へ）
chunks := stream.Naturals().Take(12).ToStream().Chunk(4)
```

### iter.Seq ブリッジ

```go
// 標準ライブラリとの連携
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range で使用
for v := range stream.Of(1, 2, 3).Iter() {
    fmt.Println(v)
}
```

## ベンチマーク

10,000 個の `int` 要素で `Filter(偶数)` → `Take(10)` — Apple M1:

```
ベンチマーク              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
PipelineFilterTake       315      464       13   ← 遅延: ネイティブの2.5倍
StreamFilterTake      30,831  128,329       17   ← 即時: 全件走査
```

`Filter+Take` で Pipeline は Stream より**約100倍高速**。遅延評価により `Take` が満たされた時点で処理を打ち切るため、Stream が全スライスをフィルタするのに対し必要最小限の要素のみ処理します。

全件走査（早期終了なし）:

```
ベンチマーク              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          30,529  128,249       16
PipelineFilter        42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> `go test -bench=. -benchmem ./...` で再現可能。

## License

MIT

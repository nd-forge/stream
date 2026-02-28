# stream

[English](README.md) | **日本語** | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

Go ジェネリクス対応のストリーム処理ライブラリ。filter, map, sort, group などのコレクション操作をメソッドチェーンで記述できます。**デフォルトで遅延評価**。

全操作は遅延評価です。全データが必要な操作（Sort, Reverse, Shuffle, TakeLast, Chunk, Partition）は内部でバッファリングし、後続の操作は自動的に遅延評価に戻ります。

## 要件

- Go 1.23 以上（`iter.Seq[T]` および range-over-function を使用）

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

### デフォルトで遅延評価

全操作は内部的に `iter.Seq[T]` を使った遅延パイプラインを構築します。終端操作（`ToSlice`, `ForEach`, `Reduce` 等）が呼ばれるまで中間スライスは一切割り当てられません。

全データが必要な操作 — `Sort`, `Reverse`, `Shuffle`, `TakeLast`, `Chunk`, `Partition` — は内部でバッファリングし、後続の操作では再び遅延評価に戻ります。

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

### ジェネレータ（無限シーケンス）

| 関数 | 説明 |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | 値の無限繰り返し |
| `RepeatN(value, n)` | 値を n 回繰り返し |

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
| `Chain(others...)` | 複数ストリームを連結 |

> `Sort`, `Reverse`, `Shuffle`, `TakeLast` は内部で全要素をバッファリングします。

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
| `Seq()` | `iter.Seq[T]` |

### トップレベル関数（型変換）

| 関数 | 説明 |
|---|---|
| `Map(s, fn)` | 各要素を変換 `T → U` |
| `MapIndexed(s, fn)` | インデックス付き変換 |
| `FlatMap(s, fn)` | 変換して平坦化 `T → []U` |
| `Reduce(s, initial, fn)` | 異なる型への畳み込み `T → U` |
| `GroupBy(s, key)` | キーでグルーピング `→ map[K][]T` |
| `Associate(s, fn)` | マップ化 `→ map[K]V` |
| `Zip(s1, s2)` | 2つの Stream をペアに `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 平坦化 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 変換 `Stream[Pair[K,V]] → map[K]V` |
| `Enumerate(s)` | インデックス付加 `→ Stream[Pair[int,T]]` |

### 数値操作

数値型の Stream 向け特殊操作（`int`, `float64` など）。

| 関数 | 説明 |
|---|---|
| `Sum(s)` / `Avg(s)` | 合計 / 平均 |
| `Min(s)` / `Max(s)` | 最小 / 最大 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 抽出値の合計 / 平均 |

### iter.Seq ブリッジ

| 関数 | 説明 |
|---|---|
| `Seq()` | `Stream[T]` → `iter.Seq[T]` |
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
    total := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f count=%d\n", category, total, len(group))
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

### 無限シーケンス

```go
// 最初の5つの偶数
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// フィボナッチ数列
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

// 遅延評価: 100万要素のうち2,001個だけ処理
result := stream.Range(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### iter.Seq ブリッジ

```go
// 標準ライブラリとの連携
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range で使用
for v := range stream.Of(1, 2, 3).Seq() {
    fmt.Println(v)
}
```

## ベンチマーク

10,000 個の `int` 要素 — Apple M1。[samber/lo](https://github.com/samber/lo) との比較。

### Filter + Take（遅延評価の優位性）

```
ベンチマーク                  ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         111        248          5
Stream                         346        528         15   ← ネイティブの3.1倍
lo                          25,000     81,920          1   ← Stream より72倍遅い
```

Stream の遅延評価は `Take` が満たされた時点で処理を打ち切ります。lo は**全 10,000 要素**をフィルタした後に切り出すため、短絡できません。

### チェーン: Filter → Map → Take 5

```
ベンチマーク                  ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         171        260          9
Stream                         384        544         19   ← ネイティブの2.2倍
lo                          64,600    152,601      3,336   ← Stream より168倍遅い
```

### 全件走査（早期終了なし）

```
ベンチマーク                  ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native Filter              19,400    128,249         16
lo     Filter              25,600     81,920          1
Stream Filter              42,200    128,409         22

Native Reduce               3,300          0          0
lo     Reduce               9,700          0          0
Stream Reduce              23,400         64          2
```

早期終了なしの全件走査では、lo は抽象化オーバーヘッドが低いため Stream より高速です。Stream の優位性は `Take`、`First`、`Find` などの短絡操作を含むパイプラインで発揮されます。

> `cd _benchmark && go test -bench=. -benchmem ./...` で再現可能。

## License

MIT

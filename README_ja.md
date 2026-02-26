# stream

[English](README.md) | **日本語** | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

Go ジェネリクス対応のストリーム処理ライブラリ。filter, map, sort, group などのコレクション操作をメソッドチェーンで記述できます。

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

## 使用例

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

### Map と FlatMap

```go
// 名前の抽出
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()

// ネストされた注文を平坦化
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy と集計

```go
bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })

for symbol, group := range bySymbol {
    totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
    avgPnL := stream.AvgBy(group, func(t Trade) float64 { return t.PnL })
    fmt.Printf("%s: total=%.2f avg=%.2f count=%d\n", symbol, totalPnL, avgPnL, group.Count())
}
```

### Partition と Chunk

```go
// 条件で分割
profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
winRate := float64(profit.Count()) / float64(trades.Count()) * 100

// バッチ処理
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

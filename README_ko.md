# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | **한국어** | [Español](README_es.md) | [Português](README_pt.md)

Go 제네릭 스트림 처리 라이브러리. filter, map, sort, group 등의 컬렉션 연산을 메서드 체이닝으로 작성할 수 있습니다.

## 설치

```bash
go get github.com/nd-forge/stream
```

## 빠른 시작

```go
import "github.com/nd-forge/stream"

// 메서드 체이닝 (타입을 유지하는 연산)
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// 최상위 함수 (타입 변환 연산)
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## 설계 방침

Go에서는 메서드에 새로운 타입 파라미터를 추가할 수 없습니다. 본 라이브러리는 다음과 같이 분리합니다:

| 종류 | 구현 방식 | 시그니처 |
|---|---|---|
| 타입 유지 연산 (Filter, Sort, Take...) | **메서드** — 체이닝 가능 | `Stream[T] → Stream[T]` |
| 타입 변환 연산 (Map, FlatMap, GroupBy...) | **최상위 함수** | `Stream[T] → Stream[U]` |

## API

### 생성자

| 함수 | 설명 |
|---|---|
| `Of[T](items ...T)` | 가변 인자로 생성 |
| `From[T](items []T)` | 슬라이스로 생성 (복사) |
| `Range(start, end)` | 정수 시퀀스 `[start, end)` |
| `Generate[T](n, fn)` | 생성 함수로 n개 요소 생성 |

### 체이닝 가능 메서드

`Stream[T]`을 반환하며 체이닝할 수 있는 연산.

| 메서드 | 설명 |
|---|---|
| `Filter(predicate)` | 조건에 맞는 요소만 유지 |
| `Reject(predicate)` | 조건에 맞는 요소 제외 |
| `Sort(cmp)` | 비교 함수로 정렬 |
| `Reverse()` | 역순 |
| `Take(n)` / `TakeLast(n)` | 앞 / 뒤 n개 요소 |
| `Skip(n)` | 앞 n개 건너뛰기 |
| `TakeWhile(pred)` / `DropWhile(pred)` | 조건 기반 가져오기 / 건너뛰기 |
| `Distinct(key)` | 키 기반 중복 제거 |
| `Shuffle()` | 랜덤 순서 |
| `Peek(fn)` | 부수 효과 실행 (수정 없음) |

### 종단 연산

| 메서드 | 반환값 |
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

### 변환 함수

타입 변환을 위한 최상위 함수.

| 함수 | 설명 |
|---|---|
| `Map(s, fn)` | 변환 `T → U` |
| `MapIndexed(s, fn)` | 인덱스 포함 변환 |
| `FlatMap(s, fn)` | 변환 후 평탄화 `T → []U` |
| `Reduce(s, initial, fn)` | 다른 타입으로 폴딩 `T → U` |
| `GroupBy(s, key)` | 키로 그룹화 `→ map[K]Stream[T]` |
| `Associate(s, fn)` | 맵 구축 `→ map[K]V` |
| `Zip(s1, s2)` | 두 스트림 페어링 `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 평탄화 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 변환 `Stream[Pair[K,V]] → map[K]V` |

### 수치 함수

수치 스트림 전용 연산 (`int`, `float64` 등).

| 함수 | 설명 |
|---|---|
| `Sum(s)` / `Avg(s)` | 합계 / 평균 |
| `Min(s)` / `Max(s)` | 최솟값 / 최댓값 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 추출값의 합계 / 평균 |

## 예제

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

### Map과 FlatMap

```go
// 이름 추출
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()

// 중첩된 주문 평탄화
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy와 집계

```go
bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })

for symbol, group := range bySymbol {
    totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
    avgPnL := stream.AvgBy(group, func(t Trade) float64 { return t.PnL })
    fmt.Printf("%s: total=%.2f avg=%.2f count=%d\n", symbol, totalPnL, avgPnL, group.Count())
}
```

### Partition과 Chunk

```go
// 조건으로 분할
profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
winRate := float64(profit.Count()) / float64(trades.Count()) * 100

// 배치 처리
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

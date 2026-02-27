# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | **한국어** | [Español](README_es.md) | [Português](README_pt.md)

Go 제네릭 스트림 처리 라이브러리. filter, map, sort, group 등의 컬렉션 연산을 메서드 체이닝으로 작성할 수 있습니다 — **즉시 평가** (`Stream`)와 **지연 평가** (`Pipeline`) 모두 지원.

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

### 즉시 평가 vs 지연 평가

| 타입 | 평가 방식 | 최적 용도 |
|---|---|---|
| `Stream[T]` | **즉시 평가** — 중간 슬라이스 생성 | 소/중규모 데이터, 랜덤 액세스 (`Shuffle`, `Chunk`, `TakeLast`) |
| `Pipeline[T]` | **지연 평가** — 중간 할당 제로 | 대규모 데이터, 무한 시퀀스, 조기 종료 (`Filter+Take`) |

자유롭게 전환: `stream.Lazy()` / `pipeline.ToStream()`

### 타입 파라미터

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

### Pipeline (지연 평가)

#### 생성자

| 함수 | 설명 |
|---|---|
| `Lazy[T](items ...T)` | 가변 인자에서 지연 파이프라인 생성 |
| `LazyFrom[T](seq)` | `iter.Seq[T]`에서 생성 |
| `LazyRange(start, end)` | 지연 정수 시퀀스 `[start, end)` |
| `stream.Lazy()` | `Stream[T]` → `Pipeline[T]` 변환 |

#### 생성기 (무한 시퀀스)

| 함수 | 설명 |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | 값의 무한 반복 |
| `RepeatN(value, n)` | 값을 n번 반복 |

#### 체이닝 가능 메서드

Stream과 동일한 API: `Filter`, `Reject`, `Sort`, `Reverse`, `Take`, `Skip`, `TakeWhile`, `DropWhile`, `Distinct`, `Peek`, `Chain`

#### 종단 연산

Stream과 동일: `ToSlice`, `First`, `Last`, `Find`, `Reduce`, `Any`, `All`, `None`, `Count`, `CountBy`, `IsEmpty`, `Contains`, `MinBy`, `MaxBy`, `ForEach`, `ForEachIndexed`

추가: `ToStream()` (즉시 Stream으로 변환), `Seq()` (기본 `iter.Seq[T]` 가져오기)

#### 변환 함수 (타입 변환)

| 함수 | 설명 |
|---|---|
| `PipeMap(p, fn)` | 변환 `T → U` |
| `PipeMapIndexed(p, fn)` | 인덱스 포함 변환 |
| `PipeFlatMap(p, fn)` | 변환 후 평탄화 `T → []U` |
| `PipeReduce(p, initial, fn)` | 다른 타입으로 폴딩 `T → U` |
| `PipeGroupBy(p, key)` | 키로 그룹화 `→ map[K][]T` |
| `PipeAssociate(p, fn)` | 맵 구축 `→ map[K]V` |
| `PipeZip(p1, p2)` | 두 파이프라인 페어링 `→ Pipeline[Pair[T,U]]` |
| `PipeFlatten(p)` | 평탄화 `Pipeline[[]T] → Pipeline[T]` |
| `PipeToMap(p)` | 변환 `Pipeline[Pair[K,V]] → map[K]V` |
| `PipeEnumerate(p)` | 인덱스 추가 `→ Pipeline[Pair[int,T]]` |

### iter.Seq 브릿지

| 함수 | 설명 |
|---|---|
| `stream.Iter()` | `Stream[T]` → `iter.Seq[T]` |
| `stream.Iter2()` | `Stream[T]` → `iter.Seq2[int, T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## 예제

아래 예제에서 사용하는 타입:

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

// 재고 있는 상품을 가격 내림차순으로 상위 3개
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

### Map과 FlatMap

```go
// 상품명 추출
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// 사용자의 중첩된 주문 평탄화
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy와 집계

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(group, func(p Product) float64 { return p.Price })
    avg := stream.AvgBy(group, func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f avg=$%.2f count=%d\n", category, total, avg, group.Count())
}
```

### Partition과 Chunk

```go
// 조건으로 분할
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// 배치 처리
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

### Pipeline (지연 평가)

```go
// 무한 시퀀스: 처음 5개의 짝수
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// 피보나치 수열
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

// 지연 평가: 100만 요소 중 2,001개만 처리
result := stream.LazyRange(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### 즉시 평가와 지연 평가 전환

```go
// Stream → Pipeline (지연 처리로)
result := stream.Of(items...).Lazy().Filter(pred).Take(10).ToSlice()

// Pipeline → Stream (즉시 전용 연산으로)
chunks := stream.Naturals().Take(12).ToStream().Chunk(4)
```

### iter.Seq 브릿지

```go
// 표준 라이브러리 연동
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range 지원
for v := range stream.Of(1, 2, 3).Iter() {
    fmt.Println(v)
}
```

## 벤치마크

10,000개 `int` 요소에서 `Filter(짝수)` → `Take(10)` — Apple M1:

```
벤치마크                ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
PipelineFilterTake       315      464       13   ← 지연: 네이티브의 2.5배
StreamFilterTake      30,831  128,329       17   ← 즉시: 전체 스캔
```

`Filter+Take`에서 Pipeline은 Stream보다 **약 100배 빠릅니다**. 지연 평가로 `Take`가 충족되면 즉시 중단하지만, Stream은 전체 슬라이스를 먼저 필터링합니다.

전체 스캔 (조기 종료 없음):

```
벤치마크                ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          30,529  128,249       16
PipelineFilter        42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> `go test -bench=. -benchmem ./...`로 재현 가능.

## License

MIT

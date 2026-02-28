# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | **한국어** | [Español](README_es.md) | [Português](README_pt.md)

Go 제네릭 스트림 처리 라이브러리. filter, map, sort, group 등의 컬렉션 연산을 메서드 체이닝으로 작성할 수 있습니다 — 기본적으로 **지연 평가**를 지원합니다. [온라인에서 사용해보기!](https://go.dev/play/p/QaQ_rdqYn1Y)

모든 연산은 지연 평가됩니다. 전체 데이터가 필요한 연산(Sort, Reverse, Shuffle, TakeLast, Chunk, Partition)은 내부적으로 버퍼링한 후 자동으로 지연 이터레이션을 재개합니다.

## 요구 사항

- Go 1.23 이상 (`iter.Seq[T]` 및 range-over-function 사용)

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

### 기본 지연 평가

모든 연산은 `iter.Seq[T]`를 사용하여 내부적으로 지연 파이프라인을 구축합니다. 종단 연산(`ToSlice`, `ForEach`, `Reduce` 등)이 호출될 때까지 중간 슬라이스가 할당되지 않습니다.

본질적으로 전체 데이터가 필요한 연산 — `Sort`, `Reverse`, `Shuffle`, `TakeLast`, `Chunk`, `Partition` — 은 내부적으로 버퍼링한 후, 이후 연산에 대해 지연 이터레이션을 재개합니다.

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

### 생성기 (무한 시퀀스)

| 함수 | 설명 |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | 값의 무한 반복 |
| `RepeatN(value, n)` | 값을 n번 반복 |

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
| `Chain(others...)` | 여러 스트림 연결 |

> `Sort`, `Reverse`, `Shuffle`, `TakeLast`는 내부적으로 모든 요소를 버퍼링합니다.

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
| `Seq()` | `iter.Seq[T]` |

### 변환 함수

타입 변환을 위한 최상위 함수.

| 함수 | 설명 |
|---|---|
| `Map(s, fn)` | 변환 `T → U` |
| `MapIndexed(s, fn)` | 인덱스 포함 변환 |
| `FlatMap(s, fn)` | 변환 후 평탄화 `T → []U` |
| `Reduce(s, initial, fn)` | 다른 타입으로 폴딩 `T → U` |
| `GroupBy(s, key)` | 키로 그룹화 `→ map[K][]T` |
| `Associate(s, fn)` | 맵 구축 `→ map[K]V` |
| `Zip(s1, s2)` | 두 스트림 페어링 `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | 평탄화 `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | 변환 `Stream[Pair[K,V]] → map[K]V` |
| `Enumerate(s)` | 인덱스 추가 `→ Stream[Pair[int,T]]` |

### 수치 함수

수치 스트림 전용 연산 (`int`, `float64` 등).

| 함수 | 설명 |
|---|---|
| `Sum(s)` / `Avg(s)` | 합계 / 평균 |
| `Min(s)` / `Max(s)` | 최솟값 / 최댓값 |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | 추출값의 합계 / 평균 |

### iter.Seq 브릿지

| 함수 | 설명 |
|---|---|
| `Seq()` | `Stream[T]` → `iter.Seq[T]` |
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

### Filter, Sort, Take [[playground]](https://go.dev/play/p/W5fo1cfb_VA)

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

### Map과 FlatMap [[playground]](https://go.dev/play/p/GyrjAXMYnZD)

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

### GroupBy와 집계 [[playground]](https://go.dev/play/p/nzL3i-4Kgj4)

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f count=%d\n", category, total, len(group))
}
```

### Partition과 Chunk [[playground]](https://go.dev/play/p/c9KhEbkeKat)

```go
// 조건으로 분할
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// 배치 처리
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

### 무한 시퀀스 [[playground]](https://go.dev/play/p/9prnIg-NjtF)

```go
// 처음 5개의 짝수 자연수
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// 피보나치 수열
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

// 지연 평가: 100만 요소 중 2,001개만 처리
result := stream.Range(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### iter.Seq 브릿지 [[playground]](https://go.dev/play/p/i-rNvVdYut5)

```go
// 표준 라이브러리 연동
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// for-range 지원
for v := range stream.Of(1, 2, 3).Seq() {
    fmt.Println(v)
}
```

## 벤치마크

10,000개 `int` 요소 — Apple M1. [samber/lo](https://github.com/samber/lo)와 비교.

### Filter + Take (지연 평가 이점)

```
벤치마크                      ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         111        248          5
Stream                         346        528         15   ← 네이티브의 3.1배
lo                          25,000     81,920          1   ← Stream보다 72배 느림
```

Stream의 지연 평가는 `Take`가 충족되면 즉시 쇼트서킷합니다. lo는 **전체 10,000개 요소**를 먼저 필터링해야 하므로 쇼트서킷이 불가능합니다.

### 체이닝: Filter → Map → Take 5

```
벤치마크                      ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native                         171        260          9
Stream                         384        544         19   ← 네이티브의 2.2배
lo                          64,600    152,601      3,336   ← Stream보다 168배 느림
```

### 전체 스캔 (조기 종료 없음)

```
벤치마크                      ns/op       B/op    allocs/op
────────────────────────────────────────────────────────────
Native Filter              19,400    128,249         16
lo     Filter              25,600     81,920          1
Stream Filter              42,200    128,409         22

Native Reduce               3,300          0          0
lo     Reduce               9,700          0          0
Stream Reduce              23,400         64          2
```

조기 종료 없는 전체 스캔에서는 lo가 추상화 오버헤드가 낮아 Stream보다 빠릅니다. Stream의 이점은 `Take`, `First`, `Find` 등 쇼트서킷 연산이 포함된 파이프라인에서 발휘됩니다.

> `cd _benchmark && go test -bench=. -benchmem ./...`로 재현 가능.

## License

MIT

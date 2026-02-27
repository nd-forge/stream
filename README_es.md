# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | **Español** | [Português](README_pt.md)

Una biblioteca genérica de procesamiento de streams para Go. Operaciones encadenables sobre colecciones: filter, map, sort, group y más — con evaluación **eager** (`Stream`) y **lazy** (`Pipeline`).

## Instalación

```bash
go get github.com/nd-forge/stream
```

## Inicio rápido

```go
import "github.com/nd-forge/stream"

// Encadenamiento de métodos (operaciones que preservan el tipo)
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// Funciones de nivel superior (operaciones que cambian el tipo)
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## Diseño

### Eager vs Lazy

| Tipo | Evaluación | Ideal para |
|---|---|---|
| `Stream[T]` | **Eager** — asigna slices intermedios | Datos pequeños/medianos, acceso aleatorio (`Shuffle`, `Chunk`, `TakeLast`) |
| `Pipeline[T]` | **Lazy** — cero asignación intermedia | Datos grandes, secuencias infinitas, terminación temprana (`Filter+Take`) |

Cambie libremente: `stream.Lazy()` / `pipeline.ToStream()`

### Parámetros de tipo

Go no permite que los métodos introduzcan nuevos parámetros de tipo. Esta biblioteca los separa así:

| Tipo | Implementación | Firma |
|---|---|---|
| Preservan el tipo (Filter, Sort, Take...) | **Métodos** — encadenables | `Stream[T] → Stream[T]` |
| Cambian el tipo (Map, FlatMap, GroupBy...) | **Funciones de nivel superior** | `Stream[T] → Stream[U]` |

## API

### Constructores

| Función | Descripción |
|---|---|
| `Of[T](items ...T)` | Crear desde argumentos variádicos |
| `From[T](items []T)` | Crear desde slice (copia) |
| `Range(start, end)` | Crear secuencia de enteros `[start, end)` |
| `Generate[T](n, fn)` | Crear n elementos con generador |

### Métodos encadenables

Operaciones que devuelven `Stream[T]` y pueden encadenarse.

| Método | Descripción |
|---|---|
| `Filter(predicate)` | Mantener elementos que coincidan |
| `Reject(predicate)` | Excluir elementos que coincidan |
| `Sort(cmp)` | Ordenar por función de comparación |
| `Reverse()` | Invertir orden |
| `Take(n)` / `TakeLast(n)` | Primeros / últimos n elementos |
| `Skip(n)` | Omitir primeros n elementos |
| `TakeWhile(pred)` / `DropWhile(pred)` | Tomar / omitir mientras sea verdadero |
| `Distinct(key)` | Eliminar duplicados por clave |
| `Shuffle()` | Orden aleatorio |
| `Peek(fn)` | Ejecutar efecto secundario (sin modificar) |

### Operaciones terminales

| Método | Retorno |
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

### Funciones de transformación

Funciones de nivel superior para operaciones que cambian el tipo.

| Función | Descripción |
|---|---|
| `Map(s, fn)` | Transformar `T → U` |
| `MapIndexed(s, fn)` | Transformar con índice |
| `FlatMap(s, fn)` | Transformar y aplanar `T → []U` |
| `Reduce(s, initial, fn)` | Plegar a tipo diferente `T → U` |
| `GroupBy(s, key)` | Agrupar por clave `→ map[K]Stream[T]` |
| `Associate(s, fn)` | Construir mapa `→ map[K]V` |
| `Zip(s1, s2)` | Emparejar dos streams `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | Aplanar `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | Convertir `Stream[Pair[K,V]] → map[K]V` |

### Funciones numéricas

Operaciones especializadas para streams numéricos (`int`, `float64`, etc.).

| Función | Descripción |
|---|---|
| `Sum(s)` / `Avg(s)` | Suma / promedio |
| `Min(s)` / `Max(s)` | Mínimo / máximo |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | Suma / promedio de valores extraídos |

### Pipeline (Evaluación Lazy)

#### Constructores

| Función | Descripción |
|---|---|
| `Lazy[T](items ...T)` | Crear pipeline lazy desde args |
| `LazyFrom[T](seq)` | Crear desde `iter.Seq[T]` |
| `LazyRange(start, end)` | Secuencia lazy de enteros `[start, end)` |
| `stream.Lazy()` | Convertir `Stream[T]` a `Pipeline[T]` |

#### Generadores (Secuencias Infinitas)

| Función | Descripción |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | Repetición infinita de valor |
| `RepeatN(value, n)` | Repetir valor n veces |

#### Métodos encadenables

Misma API que Stream: `Filter`, `Reject`, `Sort`, `Reverse`, `Take`, `Skip`, `TakeWhile`, `DropWhile`, `Distinct`, `Peek`, `Chain`

#### Operaciones terminales

Igual que Stream: `ToSlice`, `First`, `Last`, `Find`, `Reduce`, `Any`, `All`, `None`, `Count`, `CountBy`, `IsEmpty`, `Contains`, `MinBy`, `MaxBy`, `ForEach`, `ForEachIndexed`

Adicional: `ToStream()` (convertir a Stream eager), `Seq()` (obtener `iter.Seq[T]` subyacente)

#### Funciones de transformación (cambio de tipo)

| Función | Descripción |
|---|---|
| `PipeMap(p, fn)` | Transformar `T → U` |
| `PipeMapIndexed(p, fn)` | Transformar con índice |
| `PipeFlatMap(p, fn)` | Transformar y aplanar `T → []U` |
| `PipeReduce(p, initial, fn)` | Plegar a tipo diferente `T → U` |
| `PipeGroupBy(p, key)` | Agrupar por clave `→ map[K][]T` |
| `PipeAssociate(p, fn)` | Construir mapa `→ map[K]V` |
| `PipeZip(p1, p2)` | Emparejar pipelines `→ Pipeline[Pair[T,U]]` |
| `PipeFlatten(p)` | Aplanar `Pipeline[[]T] → Pipeline[T]` |
| `PipeToMap(p)` | Convertir `Pipeline[Pair[K,V]] → map[K]V` |
| `PipeEnumerate(p)` | Agregar índice `→ Pipeline[Pair[int,T]]` |

### Puente iter.Seq

| Función | Descripción |
|---|---|
| `stream.Iter()` | `Stream[T]` → `iter.Seq[T]` |
| `stream.Iter2()` | `Stream[T]` → `iter.Seq2[int, T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## Ejemplos

Tipos utilizados en los ejemplos:

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

// Productos en stock, ordenados por precio descendente, top 3
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

### Map y FlatMap

```go
// Extraer nombres de productos
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// Aplanar pedidos anidados de usuarios
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy y agregación

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(group, func(p Product) float64 { return p.Price })
    avg := stream.AvgBy(group, func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f avg=$%.2f count=%d\n", category, total, avg, group.Count())
}
```

### Partition y Chunk

```go
// Dividir por condición
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// Procesamiento por lotes
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

### Pipeline (Lazy)

```go
// Secuencia infinita: primeros 5 números pares
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// Secuencia de Fibonacci
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

// Evaluación lazy: procesa solo 2,001 de 1,000,000 elementos
result := stream.LazyRange(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### Cambiar entre Eager y Lazy

```go
// Stream → Pipeline (para procesamiento lazy)
result := stream.Of(items...).Lazy().Filter(pred).Take(10).ToSlice()

// Pipeline → Stream (para operaciones solo-eager)
chunks := stream.Naturals().Take(12).ToStream().Chunk(4)
```

### Puente iter.Seq

```go
// Interoperabilidad con biblioteca estándar
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// Soporte for-range
for v := range stream.Of(1, 2, 3).Iter() {
    fmt.Println(v)
}
```

## Benchmark

10,000 elementos `int`, `Filter(par)` luego `Take(10)` — Apple M1:

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
PipelineFilterTake       315      464       13   ← lazy: 2.5x nativo
StreamFilterTake      30,831  128,329       17   ← eager: escanea todo
```

Pipeline con `Filter+Take` es **~100x más rápido** que Stream porque la evaluación lazy cortocircuita — solo procesa elementos hasta satisfacer `Take`.

Escaneo completo (sin terminación temprana):

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          30,529  128,249       16
PipelineFilter        42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> Ejecute `go test -bench=. -benchmem ./...` para reproducir.

## License

MIT

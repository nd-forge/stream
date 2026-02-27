# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | **Español** | [Português](README_pt.md)

Una biblioteca generica de procesamiento de streams para Go. Operaciones encadenables sobre colecciones: filter, map, sort, group y mas — con **evaluacion lazy** por defecto.

Todas las operaciones son lazy. Las operaciones que requieren todos los datos (Sort, Reverse, Shuffle, TakeLast, Chunk, Partition) almacenan internamente en buffer y reanudan la iteracion lazy automaticamente.

## Instalacion

```bash
go get github.com/nd-forge/stream
```

## Inicio rapido

```go
import "github.com/nd-forge/stream"

// Encadenamiento de metodos para operaciones del mismo tipo
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// Funciones de nivel superior para operaciones que cambian el tipo
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## Diseno

### Lazy por defecto

Todas las operaciones construyen un pipeline lazy internamente usando `iter.Seq[T]`. No se asignan slices intermedios hasta que se invoca una operacion terminal (`ToSlice`, `ForEach`, `Reduce`, etc.).

Las operaciones que inherentemente necesitan todos los datos — `Sort`, `Reverse`, `Shuffle`, `TakeLast`, `Chunk`, `Partition` — almacenan internamente en buffer y luego reanudan la iteracion lazy para las operaciones siguientes.

### Parametros de tipo

Go no permite que los metodos introduzcan nuevos parametros de tipo. Esta biblioteca los separa asi:

| Tipo | Implementacion | Firma |
|---|---|---|
| Preservan el tipo (Filter, Sort, Take...) | **Metodos** — encadenables | `Stream[T] -> Stream[T]` |
| Cambian el tipo (Map, FlatMap, GroupBy...) | **Funciones de nivel superior** | `Stream[T] -> Stream[U]` |

## API

### Constructores

| Funcion | Descripcion |
|---|---|
| `Of[T](items ...T)` | Crear desde argumentos variadicos |
| `From[T](items []T)` | Crear desde slice (copia) |
| `Range(start, end)` | Crear secuencia de enteros `[start, end)` |
| `Generate[T](n, fn)` | Crear n elementos con generador |

### Generadores (Secuencias infinitas)

| Funcion | Descripcion |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | Repeticion infinita de valor |
| `RepeatN(value, n)` | Repetir valor n veces |

### Metodos encadenables

Operaciones que devuelven `Stream[T]` y pueden encadenarse.

| Metodo | Descripcion |
|---|---|
| `Filter(predicate)` | Mantener elementos que coincidan con el predicado |
| `Reject(predicate)` | Excluir elementos que coincidan con el predicado |
| `Sort(cmp)` | Ordenar por funcion de comparacion |
| `Reverse()` | Invertir orden |
| `Take(n)` / `TakeLast(n)` | Primeros / ultimos n elementos |
| `Skip(n)` | Omitir primeros n elementos |
| `TakeWhile(pred)` / `DropWhile(pred)` | Tomar / omitir desde el inicio mientras sea verdadero |
| `Distinct(key)` | Eliminar duplicados por clave |
| `Shuffle()` | Orden aleatorio |
| `Peek(fn)` | Ejecutar efecto secundario sin modificar |
| `Chain(others...)` | Concatenar multiples streams |

> `Sort`, `Reverse`, `Shuffle`, `TakeLast` almacenan todos los elementos internamente en buffer.

### Operaciones terminales

| Metodo | Retorno |
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

### Funciones de transformacion

Funciones de nivel superior para operaciones que cambian el tipo.

| Funcion | Descripcion |
|---|---|
| `Map(s, fn)` | Transformar `T -> U` |
| `MapIndexed(s, fn)` | Transformar con indice |
| `FlatMap(s, fn)` | Transformar y aplanar `T -> []U` |
| `Reduce(s, initial, fn)` | Plegar a tipo diferente `T -> U` |
| `GroupBy(s, key)` | Agrupar por clave `-> map[K][]T` |
| `Associate(s, fn)` | Construir mapa `-> map[K]V` |
| `Zip(s1, s2)` | Emparejar dos streams `-> Stream[Pair[T,U]]` |
| `Flatten(s)` | Aplanar `Stream[[]T] -> Stream[T]` |
| `ToMap(s)` | Convertir `Stream[Pair[K,V]] -> map[K]V` |
| `Enumerate(s)` | Agregar indice `-> Stream[Pair[int,T]]` |

### Funciones numericas

Operaciones especializadas para streams numericos (`int`, `float64`, etc.).

| Funcion | Descripcion |
|---|---|
| `Sum(s)` / `Avg(s)` | Suma / promedio |
| `Min(s)` / `Max(s)` | Minimo / maximo |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | Suma / promedio de valores extraidos |

### Puente iter.Seq

| Funcion | Descripcion |
|---|---|
| `Seq()` | `Stream[T]` -> `iter.Seq[T]` |
| `Collect(seq)` | `iter.Seq[T]` -> `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` -> `Stream[Pair[K,V]]` |

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

### GroupBy y agregacion

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(stream.Of(group...), func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f count=%d\n", category, total, len(group))
}
```

### Partition y Chunk

```go
// Dividir por condicion
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

### Secuencias infinitas

```go
// Primeros 5 numeros pares naturales
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// Secuencia de Fibonacci
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

// Evaluacion lazy: procesa solo 2,001 de 1,000,000 elementos
result := stream.Range(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### Puente iter.Seq

```go
// Interoperabilidad con biblioteca estandar
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// Soporte for-range
for v := range stream.Of(1, 2, 3).Seq() {
    fmt.Println(v)
}
```

## Benchmark

10,000 elementos `int`, `Filter(par)` luego `Take(10)` — Apple M1:

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
StreamFilterTake         315      464       13   <- lazy: 2.5x nativo
```

La evaluacion lazy cortocircuita — solo procesa elementos hasta que `Take` se satisface.

Escaneo completo (sin terminacion temprana):

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> Ejecute `go test -bench=. -benchmem ./...` para reproducir.

## License

MIT

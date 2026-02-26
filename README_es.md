# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | **Español** | [Português](README_pt.md)

Una biblioteca genérica de procesamiento de streams para Go. Operaciones encadenables sobre colecciones: filter, map, sort, group y más.

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

## Ejemplos

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

### Map y FlatMap

```go
// Extraer nombres
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()

// Aplanar pedidos anidados
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy y agregación

```go
bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })

for symbol, group := range bySymbol {
    totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
    avgPnL := stream.AvgBy(group, func(t Trade) float64 { return t.PnL })
    fmt.Printf("%s: total=%.2f avg=%.2f count=%d\n", symbol, totalPnL, avgPnL, group.Count())
}
```

### Partition y Chunk

```go
// Dividir por condición
profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
winRate := float64(profit.Count()) / float64(trades.Count()) * 100

// Procesamiento por lotes
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

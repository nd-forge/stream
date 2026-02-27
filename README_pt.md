# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | **Português**

Uma biblioteca genérica de processamento de streams para Go. Operações encadeáveis sobre coleções: filter, map, sort, group e mais — com avaliação **eager** (`Stream`) e **lazy** (`Pipeline`).

## Instalação

```bash
go get github.com/nd-forge/stream
```

## Início rápido

```go
import "github.com/nd-forge/stream"

// Encadeamento de métodos (operações que preservam o tipo)
result := stream.Of(5, 2, 8, 1, 9, 4, 7, 3, 6).
    Filter(func(n int) bool { return n%2 == 0 }).
    Sort(func(a, b int) int { return b - a }).
    Take(3).
    ToSlice()
// [8, 6, 4]

// Funções de nível superior (operações que mudam o tipo)
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()
```

## Design

### Eager vs Lazy

| Tipo | Avaliação | Ideal para |
|---|---|---|
| `Stream[T]` | **Eager** — aloca slices intermediários | Dados pequenos/médios, acesso aleatório (`Shuffle`, `Chunk`, `TakeLast`) |
| `Pipeline[T]` | **Lazy** — zero alocação intermediária | Dados grandes, sequências infinitas, terminação antecipada (`Filter+Take`) |

Alterne livremente: `stream.Lazy()` / `pipeline.ToStream()`

### Parâmetros de tipo

Go não permite que métodos introduzam novos parâmetros de tipo. Esta biblioteca separa:

| Tipo | Implementação | Assinatura |
|---|---|---|
| Preservam o tipo (Filter, Sort, Take...) | **Métodos** — encadeáveis | `Stream[T] → Stream[T]` |
| Mudam o tipo (Map, FlatMap, GroupBy...) | **Funções de nível superior** | `Stream[T] → Stream[U]` |

## API

### Construtores

| Função | Descrição |
|---|---|
| `Of[T](items ...T)` | Criar a partir de argumentos variádicos |
| `From[T](items []T)` | Criar a partir de slice (cópia) |
| `Range(start, end)` | Criar sequência de inteiros `[start, end)` |
| `Generate[T](n, fn)` | Criar n elementos com gerador |

### Métodos encadeáveis

Operações que retornam `Stream[T]` e podem ser encadeadas.

| Método | Descrição |
|---|---|
| `Filter(predicate)` | Manter elementos que correspondam |
| `Reject(predicate)` | Excluir elementos que correspondam |
| `Sort(cmp)` | Ordenar por função de comparação |
| `Reverse()` | Inverter ordem |
| `Take(n)` / `TakeLast(n)` | Primeiros / últimos n elementos |
| `Skip(n)` | Pular primeiros n elementos |
| `TakeWhile(pred)` / `DropWhile(pred)` | Pegar / pular enquanto verdadeiro |
| `Distinct(key)` | Remover duplicatas por chave |
| `Shuffle()` | Ordem aleatória |
| `Peek(fn)` | Executar efeito colateral (sem modificar) |

### Operações terminais

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

### Funções de transformação

Funções de nível superior para operações que mudam o tipo.

| Função | Descrição |
|---|---|
| `Map(s, fn)` | Transformar `T → U` |
| `MapIndexed(s, fn)` | Transformar com índice |
| `FlatMap(s, fn)` | Transformar e achatar `T → []U` |
| `Reduce(s, initial, fn)` | Dobrar para tipo diferente `T → U` |
| `GroupBy(s, key)` | Agrupar por chave `→ map[K]Stream[T]` |
| `Associate(s, fn)` | Construir mapa `→ map[K]V` |
| `Zip(s1, s2)` | Parear dois streams `→ Stream[Pair[T,U]]` |
| `Flatten(s)` | Achatar `Stream[[]T] → Stream[T]` |
| `ToMap(s)` | Converter `Stream[Pair[K,V]] → map[K]V` |

### Funções numéricas

Operações especializadas para streams numéricos (`int`, `float64`, etc.).

| Função | Descrição |
|---|---|
| `Sum(s)` / `Avg(s)` | Soma / média |
| `Min(s)` / `Max(s)` | Mínimo / máximo |
| `SumBy(s, fn)` / `AvgBy(s, fn)` | Soma / média de valores extraídos |

### Pipeline (Avaliação Lazy)

#### Construtores

| Função | Descrição |
|---|---|
| `Lazy[T](items ...T)` | Criar pipeline lazy de args |
| `LazyFrom[T](seq)` | Criar de `iter.Seq[T]` |
| `LazyRange(start, end)` | Sequência lazy de inteiros `[start, end)` |
| `stream.Lazy()` | Converter `Stream[T]` para `Pipeline[T]` |

#### Geradores (Sequências Infinitas)

| Função | Descrição |
|---|---|
| `Naturals()` | 0, 1, 2, 3, ... |
| `Iterate(seed, fn)` | seed, fn(seed), fn(fn(seed)), ... |
| `Repeat(value)` | Repetição infinita de valor |
| `RepeatN(value, n)` | Repetir valor n vezes |

#### Métodos encadeáveis

Mesma API do Stream: `Filter`, `Reject`, `Sort`, `Reverse`, `Take`, `Skip`, `TakeWhile`, `DropWhile`, `Distinct`, `Peek`, `Chain`

#### Operações terminais

Igual ao Stream: `ToSlice`, `First`, `Last`, `Find`, `Reduce`, `Any`, `All`, `None`, `Count`, `CountBy`, `IsEmpty`, `Contains`, `MinBy`, `MaxBy`, `ForEach`, `ForEachIndexed`

Adicional: `ToStream()` (converter para Stream eager), `Seq()` (obter `iter.Seq[T]` subjacente)

#### Funções de transformação (mudança de tipo)

| Função | Descrição |
|---|---|
| `PipeMap(p, fn)` | Transformar `T → U` |
| `PipeMapIndexed(p, fn)` | Transformar com índice |
| `PipeFlatMap(p, fn)` | Transformar e achatar `T → []U` |
| `PipeReduce(p, initial, fn)` | Dobrar para tipo diferente `T → U` |
| `PipeGroupBy(p, key)` | Agrupar por chave `→ map[K][]T` |
| `PipeAssociate(p, fn)` | Construir mapa `→ map[K]V` |
| `PipeZip(p1, p2)` | Parear pipelines `→ Pipeline[Pair[T,U]]` |
| `PipeFlatten(p)` | Achatar `Pipeline[[]T] → Pipeline[T]` |
| `PipeToMap(p)` | Converter `Pipeline[Pair[K,V]] → map[K]V` |
| `PipeEnumerate(p)` | Adicionar índice `→ Pipeline[Pair[int,T]]` |

### Ponte iter.Seq

| Função | Descrição |
|---|---|
| `stream.Iter()` | `Stream[T]` → `iter.Seq[T]` |
| `stream.Iter2()` | `Stream[T]` → `iter.Seq2[int, T]` |
| `Collect(seq)` | `iter.Seq[T]` → `Stream[T]` |
| `Collect2(seq)` | `iter.Seq2[K,V]` → `Stream[Pair[K,V]]` |

## Exemplos

Tipos utilizados nos exemplos:

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

// Produtos em estoque, ordenados por preço decrescente, top 3
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

### Map e FlatMap

```go
// Extrair nomes de produtos
names := stream.Map(
    products.Filter(func(p Product) bool { return p.InStock }),
    func(p Product) string { return p.Name },
).ToSlice()

// Achatar pedidos aninhados de usuários
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy e agregação

```go
byCategory := stream.GroupBy(products, func(p Product) string { return p.Category })

for category, group := range byCategory {
    total := stream.SumBy(group, func(p Product) float64 { return p.Price })
    avg := stream.AvgBy(group, func(p Product) float64 { return p.Price })
    fmt.Printf("%s: total=$%.2f avg=$%.2f count=%d\n", category, total, avg, group.Count())
}
```

### Partition e Chunk

```go
// Dividir por condição
inStock, outOfStock := products.Partition(func(p Product) bool { return p.InStock })

// Processamento em lotes
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
// Sequência infinita: primeiros 5 números pares
evens := stream.Naturals().
    Filter(func(n int) bool { return n%2 == 0 }).
    Take(5).
    ToSlice()
// [0, 2, 4, 6, 8]

// Sequência de Fibonacci
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

// Avaliação lazy: processa apenas 2.001 de 1.000.000 elementos
result := stream.LazyRange(0, 1_000_000).
    Filter(func(n int) bool { return n%1000 == 0 }).
    Take(3).
    ToSlice()
// [0, 1000, 2000]
```

### Alternar entre Eager e Lazy

```go
// Stream → Pipeline (para processamento lazy)
result := stream.Of(items...).Lazy().Filter(pred).Take(10).ToSlice()

// Pipeline → Stream (para operações somente-eager)
chunks := stream.Naturals().Take(12).ToStream().Chunk(4)
```

### Ponte iter.Seq

```go
// Interoperabilidade com biblioteca padrão
keys := stream.Collect(maps.Keys(myMap)).Sort(cmp).ToSlice()

// Suporte for-range
for v := range stream.Of(1, 2, 3).Iter() {
    fmt.Println(v)
}
```

## Benchmark

10.000 elementos `int`, `Filter(par)` depois `Take(10)` — Apple M1:

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilterTake         124      248        5
PipelineFilterTake       315      464       13   ← lazy: 2.5x nativo
StreamFilterTake      30,831  128,329       17   ← eager: escaneia tudo
```

Pipeline com `Filter+Take` é **~100x mais rápido** que Stream porque a avaliação lazy faz curto-circuito — processa elementos apenas até satisfazer `Take`.

Escaneamento completo (sem terminação antecipada):

```
Benchmark              ns/op     B/op    allocs/op
─────────────────────────────────────────────────────
NativeFilter          18,746  128,249       16
StreamFilter          30,529  128,249       16
PipelineFilter        42,359  128,377       21
NativeReduce           3,245        0        0
StreamReduce           9,740        0        0
```

> Execute `go test -bench=. -benchmem ./...` para reproduzir.

## License

MIT

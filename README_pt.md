# stream

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | **Português**

Uma biblioteca genérica de processamento de streams para Go. Operações encadeáveis sobre coleções: filter, map, sort, group e mais.

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

## Exemplos

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

### Map e FlatMap

```go
// Extrair nomes
names := stream.Map(
    stream.Of(users...).Filter(func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Name },
).ToSlice()

// Achatar pedidos aninhados
allOrders := stream.FlatMap(
    stream.Of(users...),
    func(u User) []Order { return u.Orders },
)
```

### GroupBy e agregação

```go
bySymbol := stream.GroupBy(trades, func(t Trade) string { return t.Symbol })

for symbol, group := range bySymbol {
    totalPnL := stream.SumBy(group, func(t Trade) float64 { return t.PnL })
    avgPnL := stream.AvgBy(group, func(t Trade) float64 { return t.PnL })
    fmt.Printf("%s: total=%.2f avg=%.2f count=%d\n", symbol, totalPnL, avgPnL, group.Count())
}
```

### Partition e Chunk

```go
// Dividir por condição
profit, loss := trades.Partition(func(t Trade) bool { return t.PnL > 0 })
winRate := float64(profit.Count()) / float64(trades.Count()) * 100

// Processamento em lotes
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

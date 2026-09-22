# 04 — Funções Utilitárias, Iteradores e Boas Práticas

## 🧰 O trio clássico: Filter, Map e Reduce

```go
// Filtrar: mantém só os que passam na condição
func Filtrar[T any](lista []T, manter func(T) bool) []T {
	var resultado []T
	for _, v := range lista {
		if manter(v) {
			resultado = append(resultado, v)
		}
	}
	return resultado
}

// Mapear: transforma cada elemento
func Mapear[T, U any](lista []T, f func(T) U) []U {
	resultado := make([]U, len(lista))
	for i, v := range lista {
		resultado[i] = f(v)
	}
	return resultado
}

// Reduzir: combina tudo num único valor
func Reduzir[T, A any](lista []T, inicial A, f func(A, T) A) A {
	acc := inicial
	for _, v := range lista {
		acc = f(acc, v)
	}
	return acc
}
```

Juntando tudo:

```go
type Produto struct {
	Nome    string
	Preco   float64
	Estoque int
}

produtos := []Produto{
	{"Caneta", 2.5, 100},
	{"Caderno", 25, 0},
	{"Mochila", 120, 5},
}

emEstoque := Filtrar(produtos, func(p Produto) bool { return p.Estoque > 0 })
nomes := Mapear(emEstoque, func(p Produto) string { return p.Nome })
valorTotal := Reduzir(emEstoque, 0.0, func(total float64, p Produto) float64 {
	return total + p.Preco*float64(p.Estoque)
})

fmt.Println(nomes)      // [Caneta Mochila]
fmt.Println(valorTotal) // 850
```

> 💡 A biblioteca padrão **não tem** `Map`/`Filter`/`Reduce` para slices (é uma escolha de design: um laço `for` costuma ser igualmente claro). Mas escrever os seus é um ótimo exercício, e em alguns projetos eles ajudam.

---

## 📦 A biblioteca padrão já é genérica

Desde o Go 1.21, os pacotes `slices` e `maps` usam generics. Você já usou (módulo 06):

```go
slices.Sort(nums)
slices.Contains(nomes, "Ana")
slices.Index(lista, alvo)
slices.Max(precos)
maps.Clone(m)
```

A assinatura de `slices.Contains`, por exemplo:

```go
func Contains[S ~[]E, E comparable](s S, v E) bool
```

Repare no `S ~[]E`: aceita **qualquer tipo cujo tipo base seja um slice de `E`**, até tipos próprios como `type Lista []string`.

Outras funções genéricas úteis:

```go
slices.IndexFunc(produtos, func(p Produto) bool { return p.Nome == "Mochila" }) // 2
slices.ContainsFunc(produtos, func(p Produto) bool { return p.Estoque == 0 })   // true
slices.SortFunc(produtos, func(a, b Produto) int { return cmp.Compare(a.Preco, b.Preco) })
slices.MaxFunc(produtos, func(a, b Produto) int { return cmp.Compare(a.Preco, b.Preco) })
```

---

## 🔁 Iteradores: `range` sobre funções (Go 1.23+)

Desde o Go 1.23, o `for range` também funciona com **funções iteradoras**. O pacote `iter` define os tipos:

```go
type Seq[V any]     func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)
```

A ideia: a função chama `yield(valor)` para cada elemento. Se `yield` retornar `false` (o laço deu `break`), ela para.

### Criando um iterador

```go
import "iter"

func Contar(ate int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; i <= ate; i++ {
			if !yield(i) {
				return // quem está no for deu break
			}
		}
	}
}

func main() {
	for n := range Contar(5) {
		fmt.Print(n, " ")
	}
	// 1 2 3 4 5
}
```

### Iterador para um tipo genérico

Nossa `Pilha[T]` pode ser percorrida com `range`:

```go
func (p *Pilha[T]) Todos() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(p.itens) - 1; i >= 0; i-- { // do topo para a base
			if !yield(p.itens[i]) {
				return
			}
		}
	}
}

var p Pilha[string]
p.Empilhar("a")
p.Empilhar("b")
p.Empilhar("c")

for v := range p.Todos() {
	fmt.Print(v, " ") // c b a
}
```

### Iteradores genéricos que transformam iteradores

```go
func Filtrado[T any](seq iter.Seq[T], manter func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if manter(v) && !yield(v) {
				return
			}
		}
	}
}

pares := Filtrado(Contar(10), func(n int) bool { return n%2 == 0 })
for n := range pares {
	fmt.Print(n, " ") // 2 4 6 8 10
}
```

### Iteradores na biblioteca padrão

```go
nums := []int{3, 1, 2}
for v := range slices.Values(nums) { ... }        // só os valores
for i, v := range slices.All(nums) { ... }        // índice e valor
for k := range maps.Keys(m) { ... }
for _, k := range slices.Sorted(maps.Keys(m)) { ... } // chaves ordenadas
coletado := slices.Collect(Contar(3))              // [1 2 3]
```

> Vantagem dos iteradores: os valores são produzidos **sob demanda**, sem criar slices intermediários. Dá até pra ter sequências **infinitas**.

---

## ✅ Boas práticas: quando usar (e não usar) generics

### ✅ Use generics para:

1. **Estruturas de dados** genéricas: pilha, fila, árvore, cache, conjunto
2. **Funções sobre slices, maps e channels** que não dependem do tipo dos elementos: `Contains`, `Filter`, `Keys`
3. Quando você perceber que está **copiando e colando** a mesma função mudando só o tipo

### ❌ Evite generics quando:

1. **Uma interface resolve**: se você só chama métodos, prefira a interface

   ```go
   // ❌ desnecessário
   func Salvar[T io.Writer](w T, dados []byte)

   // ✅ mais simples
   func Salvar(w io.Writer, dados []byte)
   ```

2. **A implementação muda para cada tipo** (aí você acaba com `switch any(v).(type)`)
3. **Só existe um tipo** usando o código, e não há previsão de outro
4. Deixa o código **mais difícil de ler** do que duplicar 10 linhas

> 🧠 *"Escreva código, não tipos."* Comece com funções concretas. Generalize quando a **repetição aparecer**, não antes.

### Sobre desempenho

O compilador do Go gera código **compartilhado** entre tipos com "formato" parecido na memória (*GC shape stenciling*). Na prática, generics costumam ser tão rápidos quanto código escrito à mão, e normalmente **mais rápidos que `any`** (que exige type assertions). Mas, como sempre: **meça antes de otimizar**.

---

## 🧾 Resumão do módulo

| Conceito | Exemplo |
|---|---|
| Função genérica | `func Maior[T cmp.Ordered](s []T) T` |
| Inferência | `Maior([]int{1, 2})` |
| Instanciação explícita | `Zero[int]()` |
| Valor zero de T | `var zero T` |
| `any` / `comparable` / `cmp.Ordered` | restrições prontas |
| Restrição própria | `type Numero interface { ~int \| ~float64 }` |
| `~T` | inclui tipos derivados de T |
| Tipo genérico | `type Pilha[T any] struct { itens []T }` |
| Método de tipo genérico | `func (p *Pilha[T]) Empilhar(v T)` |
| Iterador | `func Todos() iter.Seq[T]` + `for v := range` |

---

## ✍️ Exercícios

1. Implemente `Filtrar`, `Mapear` e `Reduzir` e use-os para: pegar os números pares de 1 a 20, elevar ao quadrado e somar tudo.
2. Escreva `AgruparPor[T any, K comparable](lista []T, chave func(T) K) map[K][]T` e agrupe palavras pela primeira letra.
3. Escreva `Unicos[T comparable](lista []T) []T` que remova duplicados **mantendo a ordem**.
4. Crie um iterador `Fibonacci() iter.Seq[int]` **infinito** e imprima os 15 primeiros com `break`.
5. Adicione um método `Todos() iter.Seq[T]` na sua `Fila[T]` e percorra com `range`.
6. Escreva `Mapeado[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U]` e combine com `Filtrado`.
7. Reescreva uma função que recebe `[]any` do módulo 10 usando generics. Ficou melhor? Em que casos **não** ficaria?

---

⬅️ Anterior: [Tipos genéricos](03-tipos-genericos.md) · 🏠 [Voltar ao roteiro](../README.md)

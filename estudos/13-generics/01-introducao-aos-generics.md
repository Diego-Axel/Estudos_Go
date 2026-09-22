# 01 — Introdução aos Generics

**Generics** (Go 1.18+) permitem escrever funções e tipos que funcionam com **vários tipos**, mantendo a **verificação de tipos** do compilador.

---

## 😩 O problema

Queremos uma função que retorne o maior valor de um slice. Sem generics, há duas saídas ruins:

### Opção 1: duplicar o código

```go
func MaiorInt(nums []int) int {
	maior := nums[0]
	for _, n := range nums[1:] {
		if n > maior {
			maior = n
		}
	}
	return maior
}

func MaiorFloat(nums []float64) float64 {
	maior := nums[0]
	for _, n := range nums[1:] { // ← exatamente o mesmo código...
		if n > maior {
			maior = n
		}
	}
	return maior
}

// ...e MaiorString, MaiorInt64, MaiorUint... 😩
```

### Opção 2: usar `any`

```go
func Maior(nums []any) any {
	// ❌ não dá pra usar ">" em any
	// ❌ precisa de type switch para cada tipo
	// ❌ quem chama recebe any e precisa de type assertion
}
```

Perde-se a **segurança de tipos**, e o código fica feio.

---

## ✨ A solução: parâmetros de tipo

```go
import "cmp"

func Maior[T cmp.Ordered](nums []T) T {
	maior := nums[0]
	for _, n := range nums[1:] {
		if n > maior {
			maior = n
		}
	}
	return maior
}

func main() {
	fmt.Println(Maior([]int{3, 9, 2}))            // 9
	fmt.Println(Maior([]float64{1.5, 0.2}))       // 1.5
	fmt.Println(Maior([]string{"pera", "uva"}))   // uva
}
```

**Uma** função, **vários** tipos, e o compilador continua checando tudo. 🎉

---

## 🔍 Anatomia

```go
func Maior[T cmp.Ordered](nums []T) T
          │ │             │        │
          │ │             │        └─ retorna um T
          │ │             └────────── recebe um slice de T
          │ └──────────────────────── restrição: T precisa ser "ordenável" (<, >)
          └────────────────────────── parâmetro de tipo chamado T
```

- Os **parâmetros de tipo** ficam entre **colchetes** `[ ]`, antes dos parâmetros normais
- Cada um tem uma **restrição** (*constraint*): o que o tipo precisa suportar
- Por convenção, usam-se letras maiúsculas curtas: `T`, `K`, `V`, `E`

---

## 🧠 Inferência de tipos

Na maioria das vezes, o Go **descobre o `T`** pelos argumentos:

```go
Maior([]int{1, 2, 3}) // T = int (inferido)
```

Mas você também pode ser **explícito** (a *instanciação*):

```go
Maior[int]([]int{1, 2, 3})
Maior[float64]([]float64{1, 2})
```

Quando é **obrigatório** ser explícito? Quando o tipo não aparece nos argumentos:

```go
func Zero[T any]() T {
	var z T
	return z
}

// Zero()       // ❌ erro: cannot infer T
Zero[int]()     // ✅ 0
Zero[string]()  // ✅ ""
```

---

## 🧮 Vários parâmetros de tipo

```go
func Mapear[T, U any](lista []T, f func(T) U) []U {
	resultado := make([]U, 0, len(lista))
	for _, v := range lista {
		resultado = append(resultado, f(v))
	}
	return resultado
}

func main() {
	nums := []int{1, 2, 3}

	dobros := Mapear(nums, func(n int) int { return n * 2 })
	textos := Mapear(nums, func(n int) string { return fmt.Sprintf("#%d", n) })

	fmt.Println(dobros) // [2 4 6]
	fmt.Println(textos) // [#1 #2 #3]
}
```

Aqui `T` é o tipo de **entrada** e `U` o de **saída**, ambos inferidos.

### Com tipos de restrições diferentes

```go
func ChavesDe[K comparable, V any](m map[K]V) []K {
	chaves := make([]K, 0, len(m))
	for k := range m {
		chaves = append(chaves, k)
	}
	return chaves
}

idades := map[string]int{"Ana": 30, "Bia": 25}
fmt.Println(len(ChavesDe(idades))) // 2
```

> `K` precisa ser `comparable` porque **chaves de map** precisam suportar `==`.

---

## 🕳️ O valor zero de `T`

Às vezes é preciso retornar "nada" de um tipo genérico. Use `var zero T`:

```go
func Primeiro[T any](lista []T) (T, bool) {
	if len(lista) == 0 {
		var zero T
		return zero, false // 0, "", nil, struct vazia... depende do T
	}
	return lista[0], true
}

v, ok := Primeiro([]string{})
fmt.Printf("%q %v\n", v, ok) // "" false
```

---

## 🆚 Generics vs. interfaces

Os dois permitem "código que funciona com vários tipos", mas de formas diferentes:

| | Interface | Generics |
|---|---|---|
| O que abstrai | **comportamento** (métodos) | **o tipo** em si |
| Tipo decidido | em **tempo de execução** | em **tempo de compilação** |
| Retorno | normalmente a interface | o **tipo concreto** (`T`) |
| Operadores (`+`, `<`) | ❌ não dá | ✅ com a restrição certa |
| Exemplo clássico | `io.Writer`, `error` | `slices.Sort`, `Maior[T]` |

```go
// Interface: "qualquer coisa que saiba escrever"
func Salvar(w io.Writer, dados []byte)

// Generics: "um slice de qualquer tipo ordenável, e devolvo o MESMO tipo"
func Maior[T cmp.Ordered](nums []T) T
```

> 🧠 Regra prática: se você está escrevendo o **mesmo código** para tipos diferentes, use **generics**. Se tipos diferentes têm **comportamentos diferentes** com a mesma "cara", use **interfaces**.

---

## ✍️ Exercícios

1. Escreva `Menor[T cmp.Ordered](nums []T) T` e teste com `int`, `float64` e `string`.
2. Escreva `Contem[T comparable](lista []T, alvo T) bool`.
3. Escreva `Inverter[T any](lista []T) []T` que retorne um **novo** slice invertido.
4. Escreva `Ultimo[T any](lista []T) (T, bool)`, retornando o valor zero e `false` se estiver vazio.
5. Escreva `Mapear` como no exemplo e use para transformar `[]string` em `[]int` (o tamanho de cada palavra).
6. Por que `Zero()` não compila, mas `Zero[int]()` sim?

---

🏠 [Módulo 13](README.md) · ➡️ Próximo: [Restrições (constraints)](02-restricoes.md)

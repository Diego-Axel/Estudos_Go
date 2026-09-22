# 02 — Restrições (Constraints)

A **restrição** diz **o que o tipo `T` precisa suportar**. Ela determina o que você pode fazer com valores do tipo `T` dentro da função.

---

## 🌌 `any`: qualquer tipo

Aceita tudo, mas você **não pode fazer quase nada** com o valor (nem `==`, nem `+`, nem `<`):

```go
func Imprimir[T any](v T) {
	fmt.Println(v) // ✅ passar adiante, guardar, retornar
}

func Igual[T any](a, b T) bool {
	return a == b // ❌ erro: incomparable types in type set
}
```

---

## ⚖️ `comparable`: suporta `==` e `!=`

```go
func Igual[T comparable](a, b T) bool {
	return a == b // ✅
}

func Indice[T comparable](lista []T, alvo T) int {
	for i, v := range lista {
		if v == alvo {
			return i
		}
	}
	return -1
}

fmt.Println(Indice([]string{"a", "b", "c"}, "c")) // 2
```

Tipos comparáveis: números, strings, bool, ponteiros, channels, arrays e structs de campos comparáveis, interfaces.
**Não** são: slices, maps, funções.

> Chaves de map genérico sempre precisam de `comparable`: `map[K]V` com `K comparable`.

---

## 📶 `cmp.Ordered`: suporta `<`, `<=`, `>`, `>=` (Go 1.21+)

```go
import "cmp"

func Limitar[T cmp.Ordered](v, minimo, maximo T) T {
	if v < minimo {
		return minimo
	}
	if v > maximo {
		return maximo
	}
	return v
}

fmt.Println(Limitar(15, 0, 10))    // 10
fmt.Println(Limitar(-2.5, 0, 1.0)) // 0
```

`cmp.Ordered` inclui todos os inteiros, floats e `string`.

O pacote `cmp` também tem funções úteis:

```go
cmp.Compare(1, 2)  // -1 (menor), 0 (igual), 1 (maior)
cmp.Less(1, 2)     // true
cmp.Or("", "padrão") // "padrão" (primeiro valor não-zero) Go 1.22+
```

E existem as funções nativas `min` e `max` (Go 1.21+):

```go
fmt.Println(min(3, 1, 2))      // 1
fmt.Println(max(2.5, 1.0))     // 2.5
fmt.Println(min("b", "a"))     // a
```

---

## 🧩 Criando suas próprias restrições

Restrições são **interfaces**. Além de métodos, elas podem listar **conjuntos de tipos** com `|`:

```go
type Numero interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Somar[T Numero](nums []T) T {
	var total T
	for _, n := range nums {
		total += n // ✅ todos os tipos da lista suportam +
	}
	return total
}

fmt.Println(Somar([]int{1, 2, 3}))       // 6
fmt.Println(Somar([]float64{1.5, 2.5}))  // 4
```

O compilador permite um operador **só se todos** os tipos do conjunto o suportarem.

### Restrição direto na assinatura

Para casos simples, dá pra escrever inline:

```go
func Dobro[T int | float64](v T) T {
	return v * 2
}
```

---

## 〰️ O til `~`: incluindo tipos derivados

Lembra dos tipos próprios, como `type Celsius float64`? Eles **não** entram em `float64`:

```go
type Numero interface {
	int | float64
}

type Celsius float64

Somar([]Celsius{36.5, 37}) // ❌ erro: Celsius does not satisfy Numero
```

Com **`~`**, a restrição aceita **qualquer tipo cujo tipo base** seja aquele:

```go
type Numero interface {
	~int | ~float64 // int, float64 E qualquer tipo baseado neles
}

Somar([]Celsius{36.5, 37}) // ✅ 73.5
```

> ✅ Quase sempre você quer o `~`. O `cmp.Ordered`, por exemplo, é definido com `~int | ~int8 | ... | ~string`.

---

## 🔧 Restrições com métodos

Uma restrição pode exigir **métodos**, como uma interface comum:

```go
type Descritivel interface {
	Descricao() string
}

func Listar[T Descritivel](itens []T) {
	for i, item := range itens {
		fmt.Printf("%d. %s\n", i+1, item.Descricao())
	}
}
```

E pode **combinar** métodos + conjunto de tipos:

```go
type NumeroFormatavel interface {
	~int | ~float64
	String() string
}
```

---

## ⚠️ Restrições com conjuntos de tipos **não** são tipos comuns

Interfaces com `|` ou `~` **só** podem ser usadas como restrição:

```go
var n Numero // ❌ erro: cannot use type Numero outside a type constraint: interface contains type constraints
```

Interfaces **só com métodos** (como `Descritivel`) podem ser usadas nos dois lugares.

---

## 🔀 Type switch dentro de uma função genérica

Não dá pra fazer `switch` direto em `T`, mas dá pra converter para `any` antes:

```go
func Descrever[T any](v T) string {
	switch x := any(v).(type) {
	case int:
		return fmt.Sprintf("inteiro %d", x)
	case string:
		return fmt.Sprintf("texto %q", x)
	default:
		return fmt.Sprintf("outro: %v", x)
	}
}
```

> ⚠️ Se você precisa disso com frequência, talvez generics **não** sejam a ferramenta certa. Generics brilham quando o código é **igual** para todos os tipos.

---

## 📋 Resumo das restrições

| Restrição | Permite | Exemplos de tipos |
|---|---|---|
| `any` | guardar, passar, retornar | qualquer um |
| `comparable` | `==`, `!=`, chave de map | int, string, structs simples |
| `cmp.Ordered` | `<`, `>`, `<=`, `>=`, `==`, `+` (numéricos e strings) | int, float64, string |
| `int \| float64` | operadores comuns aos tipos | só esses dois |
| `~int \| ~float64` | idem | esses + tipos baseados neles |
| interface com métodos | chamar os métodos | quem implementa |

---

## ✍️ Exercícios

1. Escreva `Media[T Numero](nums []T) float64`.
2. Crie `type Reais float64` e mostre que ela **não** funciona com `int | float64`, mas funciona com `~int | ~float64`.
3. Escreva `Remover[T comparable](lista []T, alvo T) []T` que remova todas as ocorrências.
4. Escreva `ContarSe[T any](lista []T, cond func(T) bool) int`.
5. Crie a restrição `Identificavel` (método `ID() string`) e `BuscarPorID[T Identificavel](itens []T, id string) (T, bool)`.
6. Tente declarar `var x cmp.Ordered` e explique o erro.

---

⬅️ Anterior: [Introdução aos generics](01-introducao-aos-generics.md) · ➡️ Próximo: [Tipos genéricos](03-tipos-genericos.md)

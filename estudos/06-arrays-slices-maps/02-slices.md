# 02 — Slices

Um **slice** é uma **"janela" dinâmica** sobre um array. Diferente do array, ele **pode crescer**, e é a estrutura de lista mais usada em Go.

---

## 🧠 Como um slice funciona por dentro

Um slice é uma pequena estrutura com **3 campos**:

```
slice
┌──────────┬─────┬─────┐
│ ponteiro │ len │ cap │
└────┬─────┴─────┴─────┘
     │
     ▼
   array de verdade (na memória)
   ┌───┬───┬───┬───┬───┐
   │ 1 │ 2 │ 3 │   │   │
   └───┴───┴───┴───┴───┘
     ◄── len=3 ─►
     ◄──── cap=5 ─────►
```

- **ponteiro:** onde os dados começam no array de baixo
- **`len` (tamanho):** quantos elementos o slice **tem**
- **`cap` (capacidade):** quantos elementos **cabem** a partir do início, antes de precisar de um array novo

---

## 📝 Criando slices

```go
// 1. Literal (sem tamanho entre os colchetes!)
frutas := []string{"maçã", "banana", "uva"}

// 2. make(tipo, len)
nums := make([]int, 3)        // [0 0 0], len=3, cap=3

// 3. make(tipo, len, cap)
buf := make([]int, 0, 10)     // [], len=0, cap=10

// 4. Slice nil (declarado sem valor)
var vazio []int               // nil, len=0, cap=0

// 5. Fatiando um array
arr := [5]int{10, 20, 30, 40, 50}
s := arr[1:4]                 // [20 30 40]
```

> 🔍 A diferença: `[3]int{...}` é **array**; `[]int{...}` é **slice**.

```go
fmt.Println(len(frutas), cap(frutas)) // 3 3
fmt.Println(len(buf), cap(buf))       // 0 10
```

---

## ➕ `append`: adicionando elementos

```go
var nums []int
nums = append(nums, 1)
nums = append(nums, 2, 3, 4)
fmt.Println(nums) // [1 2 3 4]

outros := []int{5, 6}
nums = append(nums, outros...) // junta dois slices
fmt.Println(nums) // [1 2 3 4 5 6]
```

> ⚠️ **Sempre** reatribua o resultado: `nums = append(nums, x)`. O `append` pode devolver um slice apontando para um **array novo**.

### Como o slice cresce

Quando o `len` vai passar do `cap`, o `append` **cria um array maior**, copia tudo para ele e retorna um slice apontando para o novo array:

```go
var s []int
for i := range 10 {
	s = append(s, i)
	fmt.Printf("len=%-2d cap=%d\n", len(s), cap(s))
}
// len=1  cap=1
// len=2  cap=2
// len=3  cap=4
// len=4  cap=4
// len=5  cap=8
// ...
// len=9  cap=16
```

Para slices pequenos a capacidade costuma **dobrar**; para grandes, cresce numa proporção menor. (Os valores exatos são detalhe interno do runtime.)

### ⚡ Dica de desempenho: pré-alocar

Se você **sabe** quantos elementos vão entrar, use `make` com capacidade e evite realocações:

```go
quadrados := make([]int, 0, 1000) // len 0, cap 1000
for i := range 1000 {
	quadrados = append(quadrados, i*i) // nunca realoca
}
```

---

## 🔪 Fatiamento (*slicing*): `s[inicio:fim]`

Pega do índice `inicio` **até `fim - 1`** (o `fim` **não** entra):

```go
s := []int{0, 10, 20, 30, 40, 50}

fmt.Println(s[1:4]) // [10 20 30]
fmt.Println(s[:3])  // [0 10 20]     (do começo)
fmt.Println(s[3:])  // [30 40 50]    (até o fim)
fmt.Println(s[:])   // [0 10 20 30 40 50] (tudo)
```

| Expressão | `len` | `cap` |
|---|---|---|
| `s[a:b]` | `b - a` | `cap(s) - a` |
| `s[a:b:c]` | `b - a` | `c - a` |

```go
s := make([]int, 6, 10)
t := s[2:4]
fmt.Println(len(t), cap(t)) // 2 8
```

A terceira forma (`s[a:b:c]`, *full slice expression*) **limita a capacidade**. A utilidade disso aparece no próximo arquivo.

---

## 🕳️ `nil` slice vs. slice vazio

```go
var a []int          // nil
b := []int{}         // vazio, mas não nil
c := make([]int, 0)  // vazio, mas não nil

fmt.Println(a == nil, b == nil, c == nil) // true false false
fmt.Println(len(a), len(b), len(c))       // 0 0 0
```

Na prática, **se comportam igual**: `len` é 0, `range` não itera, `append` funciona.

```go
var lista []string
lista = append(lista, "funciona!") // ✅ append em nil é ok
```

> ✅ **Idioma Go:** verifique vazio com `len(s) == 0`, **não** com `s == nil`.
>
> Uma diferença aparece no **JSON**: `nil` vira `null` e `[]int{}` vira `[]`.

---

## 🔁 Percorrendo

```go
frutas := []string{"maçã", "banana", "uva"}

for i, f := range frutas {
	fmt.Println(i, f)
}
```

---

## 🧮 Slice de slices (matriz dinâmica)

```go
matriz := [][]int{
	{1, 2, 3},
	{4, 5},
	{6},
}
fmt.Println(matriz[1][0]) // 4 (as linhas podem ter tamanhos diferentes!)

// Criando uma matriz linhas × colunas
linhas, colunas := 3, 4
grade := make([][]int, linhas)
for i := range grade {
	grade[i] = make([]int, colunas)
}
grade[1][2] = 7
fmt.Println(grade) // [[0 0 0 0] [0 0 7 0] [0 0 0 0]]
```

---

## ✍️ Exercícios

1. Crie um slice vazio e adicione os números de 1 a 20 com `append`. Imprima `len` e `cap` a cada volta.
2. Dado `s := []int{10, 20, 30, 40, 50, 60, 70}`, imprima: os 3 primeiros, os 3 últimos e os do meio (índices 2 a 4).
3. Leia números do teclado até o usuário digitar `0`, guardando num slice. No final, mostre soma, média, maior e menor.
4. Qual o `len` e o `cap` de `s[2:5]` se `s := make([]int, 8, 12)`? E de `s[2:5:6]`?
5. Crie uma função `filtrarPares(nums []int) []int` que retorne um **novo** slice só com os pares (pré-aloque com `make(..., 0, len(nums))`).
6. Crie uma matriz 5×5 onde cada posição vale `linha * coluna` e imprima como tabela.

---

⬅️ Anterior: [Arrays](01-arrays.md) · ➡️ Próximo: [Slices: avançado](03-slices-avancado.md)

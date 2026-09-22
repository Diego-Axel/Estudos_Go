# 03 — Funções como Valor e Closures

Em Go, funções são **cidadãs de primeira classe**: podem ser guardadas em variáveis, passadas como argumento e retornadas por outras funções.

---

## 📦 Função em variável

```go
func soma(a, b int) int { return a + b }

func main() {
	operacao := soma          // sem parênteses: é a função, não a chamada
	fmt.Println(operacao(2, 3)) // 5

	fmt.Printf("%T\n", operacao) // func(int, int) int
}
```

O **tipo** de uma função é a sua **assinatura**: `func(int, int) int`.

---

## 🕶️ Funções anônimas

Funções **sem nome**, definidas na hora:

```go
dobro := func(n int) int {
	return n * 2
}
fmt.Println(dobro(5)) // 10
```

Executando na hora (IIFE: *immediately invoked function expression*):

```go
func() {
	fmt.Println("executei na hora!")
}() // ← os () no final chamam a função

resultado := func(a, b int) int { return a * b }(3, 4)
fmt.Println(resultado) // 12
```

> Muito usado com `defer` e `go` (goroutines).

---

## 🔀 Funções como parâmetro (*higher-order functions*)

```go
func aplicar(nums []int, f func(int) int) []int {
	resultado := make([]int, len(nums))
	for i, n := range nums {
		resultado[i] = f(n)
	}
	return resultado
}

func main() {
	nums := []int{1, 2, 3, 4}

	dobrados := aplicar(nums, func(n int) int { return n * 2 })
	quadrados := aplicar(nums, func(n int) int { return n * n })

	fmt.Println(dobrados)  // [2 4 6 8]
	fmt.Println(quadrados) // [1 4 9 16]
}
```

Exemplo real da biblioteca padrão, ordenar com critério próprio:

```go
import "sort"

nomes := []string{"Caio", "Ana", "Bernardo"}
sort.Slice(nomes, func(i, j int) bool {
	return len(nomes[i]) < len(nomes[j]) // ordena pelo tamanho
})
fmt.Println(nomes) // [Ana Caio Bernardo]
```

---

## 🏷️ Criando um tipo de função

Deixa assinaturas longas mais legíveis:

```go
type Operacao func(float64, float64) float64

func calcular(a, b float64, op Operacao) float64 {
	return op(a, b)
}

func main() {
	soma := func(a, b float64) float64 { return a + b }
	fmt.Println(calcular(10, 5, soma)) // 15
}
```

### Map de funções: uma calculadora sem `switch`

```go
operacoes := map[string]Operacao{
	"+": func(a, b float64) float64 { return a + b },
	"-": func(a, b float64) float64 { return a - b },
	"*": func(a, b float64) float64 { return a * b },
	"/": func(a, b float64) float64 { return a / b },
}

if op, ok := operacoes["*"]; ok {
	fmt.Println(op(6, 7)) // 42
}
```

---

## ↩️ Funções que retornam funções

```go
func multiplicador(fator int) func(int) int {
	return func(n int) int {
		return n * fator
	}
}

func main() {
	dobro := multiplicador(2)
	triplo := multiplicador(3)

	fmt.Println(dobro(10))  // 20
	fmt.Println(triplo(10)) // 30
}
```

---

## 🔒 Closures

Uma **closure** é uma função que **"captura" variáveis** do escopo onde foi criada, e **continua tendo acesso a elas** mesmo depois que esse escopo terminou.

```go
func contador() func() int {
	count := 0 // essa variável "sobrevive" dentro da closure
	return func() int {
		count++
		return count
	}
}

func main() {
	c1 := contador()
	fmt.Println(c1()) // 1
	fmt.Println(c1()) // 2
	fmt.Println(c1()) // 3

	c2 := contador() // novo contador, com seu próprio count
	fmt.Println(c2()) // 1
}
```

Cada chamada de `contador()` cria um **`count` novo**. `c1` e `c2` não interferem um no outro.

### A closure captura a **variável**, não o valor

```go
x := 10
mostrar := func() { fmt.Println(x) }

x = 20
mostrar() // 20 ← enxerga o valor ATUAL de x
```

E pode alterá-la:

```go
total := 0
somar := func(n int) { total += n }

somar(5)
somar(10)
fmt.Println(total) // 15
```

### Uso prático: gerador de IDs

```go
func geradorDeID(prefixo string) func() string {
	id := 0
	return func() string {
		id++
		return fmt.Sprintf("%s-%03d", prefixo, id)
	}
}

func main() {
	novoLivro := geradorDeID("LIV")
	fmt.Println(novoLivro()) // LIV-001
	fmt.Println(novoLivro()) // LIV-002
}
```

---

## ⚠️ Comparando funções

Funções **só podem ser comparadas com `nil`**:

```go
var f func()
fmt.Println(f == nil) // true

// f()  // 💥 panic: chamar função nil
```

---

## ✍️ Exercícios

1. Crie `filtrar(nums []int, cond func(int) bool) []int` e use para pegar só os pares e depois só os maiores que 10.
2. Crie `reduzir(nums []int, inicial int, f func(int, int) int) int` e use para calcular a soma e o produto de um slice.
3. Crie `somador(inicial int) func(int) int` que acumule os valores: `s := somador(10); s(5) → 15; s(3) → 18`.
4. Faça uma calculadora usando um `map[string]func(float64, float64) float64`.
5. Crie uma closure `fibonacci() func() int` que, a cada chamada, retorne o próximo número da sequência.
6. Ordene um slice de nomes **em ordem alfabética decrescente** com `sort.Slice`.

---

⬅️ Anterior: [Funções variádicas](02-variadicas.md) · ➡️ Próximo: [Recursão](04-recursao.md)

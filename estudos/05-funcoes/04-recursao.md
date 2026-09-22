# 04 — Recursão

## 🔁 O que é?

Uma função **recursiva** é uma função que **chama a si mesma**. Toda recursão precisa de duas partes:

1. **Caso base:** a condição de parada (sem ela, a recursão nunca termina)
2. **Caso recursivo:** a chamada a si mesma com um problema **menor**

---

## 🧮 Fatorial

`5! = 5 × 4 × 3 × 2 × 1 = 120`

Repare que `5! = 5 × 4!`, e `4! = 4 × 3!`...

```go
func fatorial(n int) int {
	if n <= 1 { // caso base
		return 1
	}
	return n * fatorial(n-1) // caso recursivo
}

func main() {
	fmt.Println(fatorial(5)) // 120
}
```

### Como funciona por dentro (pilha de chamadas)

```
fatorial(5)
= 5 * fatorial(4)
= 5 * (4 * fatorial(3))
= 5 * (4 * (3 * fatorial(2)))
= 5 * (4 * (3 * (2 * fatorial(1))))
= 5 * (4 * (3 * (2 * 1)))       ← caso base, agora "volta"
= 120
```

---

## 🐚 Fibonacci

`0, 1, 1, 2, 3, 5, 8, 13, 21...` → cada número é a soma dos dois anteriores.

```go
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

fmt.Println(fib(10)) // 55
```

### ⚠️ Problema: recálculo exponencial

`fib(5)` calcula `fib(3)` duas vezes, `fib(2)` três vezes... `fib(45)` já demora **segundos**.

**Solução 1: memoização** (guardar resultados já calculados):

```go
func fibMemo(n int, memo map[int]int) int {
	if n < 2 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
	return memo[n]
}

fmt.Println(fibMemo(90, map[int]int{})) // instantâneo
```

**Solução 2: versão iterativa** (a mais eficiente):

```go
func fibIter(n int) int {
	a, b := 0, 1
	for range n {
		a, b = b, a+b
	}
	return a
}
```

---

## 📚 Mais exemplos clássicos

### Soma dos dígitos

```go
func somaDigitos(n int) int {
	if n < 10 {
		return n
	}
	return n%10 + somaDigitos(n/10)
}

fmt.Println(somaDigitos(1234)) // 10
```

### Potência

```go
func potencia(base, exp int) int {
	if exp == 0 {
		return 1
	}
	return base * potencia(base, exp-1)
}

fmt.Println(potencia(2, 10)) // 1024
```

### Inverter uma string

```go
func inverter(s string) string {
	r := []rune(s)
	if len(r) <= 1 {
		return s
	}
	return inverter(string(r[1:])) + string(r[0])
}

fmt.Println(inverter("Golang")) // gnaloG
```

### Busca binária (em slice ordenado)

```go
func buscaBinaria(nums []int, alvo, inicio, fim int) int {
	if inicio > fim {
		return -1 // não encontrado
	}
	meio := (inicio + fim) / 2
	switch {
	case nums[meio] == alvo:
		return meio
	case nums[meio] < alvo:
		return buscaBinaria(nums, alvo, meio+1, fim)
	default:
		return buscaBinaria(nums, alvo, inicio, meio-1)
	}
}

nums := []int{1, 3, 5, 7, 9, 11}
fmt.Println(buscaBinaria(nums, 7, 0, len(nums)-1)) // 3
```

---

## 🧨 Estouro de pilha (*stack overflow*)

Sem caso base (ou com um caso base que nunca é atingido), a recursão cresce até estourar a memória:

```go
func infinita(n int) int {
	return infinita(n + 1)
}
// runtime: goroutine stack exceeds 1000000000-byte limit
// fatal error: stack overflow
```

> Em Go, a pilha das goroutines **cresce dinamicamente**, então dá pra ir bem fundo (milhões de chamadas). Mas Go **não** faz otimização de *tail call*: recursões muito profundas continuam gastando memória.

---

## 🆚 Recursão ou laço?

| Recursão | Laço (`for`) |
|---|---|
| ✅ Natural para **estruturas recursivas**: árvores, pastas, JSON aninhado | ✅ Mais rápido e econômico em memória |
| ✅ Código mais curto para divisão e conquista (quicksort, mergesort) | ✅ Sem risco de estouro de pilha |
| ❌ Pode recalcular coisas (fib) | ❌ Às vezes fica mais verboso |

> 🧠 **Em Go, o idioma é preferir laços**, e usar recursão quando o problema é naturalmente recursivo (ex: percorrer pastas com subpastas).

---

## ✍️ Exercícios

1. Escreva `somaAte(n int) int` recursiva que some de 1 até `n`.
2. Escreva `contarRegressivo(n int)` que imprima de `n` até 1 e depois "Fogo! 🚀".
3. Escreva `ehPalindromo(s string) bool` recursiva (ex: `"arara"` → true).
4. Escreva `mdc(a, b int) int` usando o **algoritmo de Euclides**: `mdc(a, 0) = a`; `mdc(a, b) = mdc(b, a % b)`.
5. Escreva `maiorRecursivo(nums []int) int` que encontre o maior valor de um slice sem usar laço.
6. Compare o tempo de `fib(40)` recursivo com `fibIter(40)` (use `time.Now()` e `time.Since()`).
7. Desafio: resolva a **Torre de Hanói** para 3 discos, imprimindo cada movimento.

---

⬅️ Anterior: [Funções como valor e closures](03-funcoes-como-valor-e-closures.md) · ➡️ Próximo: [defer](05-defer.md)

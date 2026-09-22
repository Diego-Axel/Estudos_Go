# 01 — Operadores Aritméticos

## ➕ Os operadores

| Operador | Nome | Exemplo | Resultado |
|---|---|---|---|
| `+` | soma | `7 + 3` | `10` |
| `-` | subtração | `7 - 3` | `4` |
| `*` | multiplicação | `7 * 3` | `21` |
| `/` | divisão | `7 / 3` | `2` ⚠️ |
| `%` | resto (módulo) | `7 % 3` | `1` |

```go
a, b := 7, 3
fmt.Println(a+b, a-b, a*b, a/b, a%b) // 10 4 21 2 1
```

> ⚠️ Os dois lados precisam ser do **mesmo tipo**. `int + float64` não compila (veja [Conversão de tipos](../02-fundamentos/04-conversao-de-tipos.md)).

---

## ➗ Divisão: inteira vs. decimal

Se **os dois** operandos são inteiros, o resultado é **inteiro** (a parte decimal é **cortada**):

```go
fmt.Println(7 / 2)     // 3
fmt.Println(7.0 / 2)   // 3.5  (constante 7.0 é float)
fmt.Println(-7 / 2)    // -3   (corta em direção ao zero)
```

Com variáveis, converta antes:

```go
x, y := 7, 2
fmt.Println(x / y)                   // 3
fmt.Println(float64(x) / float64(y)) // 3.5
```

---

## 🔁 Resto `%`

- Só funciona com **inteiros** (para floats use `math.Mod`).
- O sinal do resultado segue o **primeiro** operando.

```go
fmt.Println(7 % 3)   // 1
fmt.Println(-7 % 3)  // -1
fmt.Println(7 % -3)  // 1

fmt.Println(math.Mod(7.5, 2)) // 1.5
```

### Usos clássicos

```go
n := 10
ehPar := n%2 == 0          // true  (par ou ímpar)
ultimoDigito := 1234 % 10  // 4
multiploDe5 := n%5 == 0    // true
```

---

## 💥 Divisão por zero

| Situação | O que acontece |
|---|---|
| Inteiro ÷ **constante** 0 | ❌ erro de **compilação**: `division by zero` |
| Inteiro ÷ **variável** que vale 0 | 💥 **panic** em execução: `integer divide by zero` |
| Float ÷ variável que vale 0 | `+Inf`, `-Inf` ou `NaN` (sem panic) |

```go
zero := 0
// fmt.Println(10 / 0)  // ❌ não compila
// fmt.Println(10 / zero) // 💥 panic: runtime error: integer divide by zero

fz := 0.0
fmt.Println(1 / fz)  // +Inf
fmt.Println(-1 / fz) // -Inf
fmt.Println(fz / fz) // NaN
```

> ✅ Sempre verifique o divisor antes de dividir: `if divisor != 0 { ... }`

---

## ⬆️⬇️ Incremento e decremento: `++` e `--`

```go
i := 5
i++ // i = 6
i-- // i = 5
```

Em Go eles são **comandos** (*statements*), **não expressões**. Isso significa:

```go
i++          // ✅
// ++i       // ❌ não existe forma prefixada
// x := i++  // ❌ não pode usar como valor
// fmt.Println(i++) // ❌
```

> 🧠 Isso elimina de vez aquelas pegadinhas do C/Java tipo `a = i++ + ++i`.

---

## ➖ Operadores unários `+` e `-`

```go
x := 5
fmt.Println(-x)  // -5
fmt.Println(+x)  // 5
```

---

## 🔤 `+` com strings = concatenação

```go
nome := "Diego"
msg := "Olá, " + nome + "!"
fmt.Println(msg) // Olá, Diego!

// msg := "Idade: " + 25 // ❌ erro: não mistura string com int
msg2 := "Idade: " + strconv.Itoa(25) // ✅
```

> 💡 Para juntar **muitas** strings num laço, use `strings.Builder` (mais eficiente). Veremos no módulo de Strings.

---

## 🧮 E a potência?

Go **não tem** operador de potência (`**` ou `^`). Use `math.Pow`:

```go
import "math"

fmt.Println(math.Pow(2, 10)) // 1024 (float64)
fmt.Println(math.Sqrt(16))   // 4   (raiz quadrada)
```

> ⚠️ `2 ^ 10` em Go **não é potência**! É o operador XOR bit a bit e dá `8`. Veja em [Operadores bit a bit](03-bit-a-bit.md).

Para potência de 2 com inteiros, dá pra usar deslocamento: `1 << 10` = `1024`.

---

## ✍️ Exercícios

1. Leia dois números inteiros e mostre soma, subtração, multiplicação, divisão inteira, divisão decimal e resto.
2. Faça um programa que diga se um número é **par ou ímpar**.
3. Dado um total de segundos (ex: `3725`), mostre quantas **horas, minutos e segundos** ele tem (use `/` e `%`). Resposta esperada: `1h 2min 5s`.
4. Qual o resultado de `-10 % 3`? E de `-10 / 3`? Explique.
5. Calcule a área de um círculo de raio `5` (`π × r²`) usando `math.Pi` e `math.Pow`.
6. O que acontece com `x := 10; y := 0; fmt.Println(x / y)`? E se `x` e `y` forem `float64`?

---

🏠 [Módulo 03](README.md) · ➡️ Próximo: [Comparação e lógicos](02-comparacao-e-logicos.md)

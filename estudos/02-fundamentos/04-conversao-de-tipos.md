# 04 — Conversão de Tipos

## 🚫 Go **não** converte sozinho

Em muitas linguagens, `int + float` "simplesmente funciona". Em Go, **não**. Toda conversão entre tipos diferentes precisa ser **explícita**.

```go
a := 10      // int
b := 2.5     // float64

// c := a + b // ❌ erro: mismatched types int and float64
c := float64(a) + b // ✅ 12.5
```

Isso vale até entre inteiros de tamanhos diferentes:

```go
var x int32 = 5
var y int64 = 10
// z := x + y      // ❌ erro
z := int64(x) + y  // ✅
```

---

## 🔄 Sintaxe: `Tipo(valor)`

```go
i := 42
f := float64(i)   // 42.0
u := uint(f)      // 42
```

### Float → Int: **trunca** (corta a parte decimal, não arredonda)

```go
f := 9.99
fmt.Println(int(f))  // 9
fmt.Println(int(-9.99)) // -9
```

Para arredondar, use o pacote `math`:

```go
import "math"

fmt.Println(math.Round(9.5))  // 10
fmt.Println(math.Floor(9.99)) // 9  (para baixo)
fmt.Println(math.Ceil(9.01))  // 10 (para cima)
fmt.Println(int(math.Round(9.5))) // 10 como int
```

### ⚠️ Inteiro grande → inteiro pequeno: **perde dados**

```go
var grande int = 300
pequeno := int8(grande)
fmt.Println(pequeno) // 44 😱 (300 - 256)
```

O compilador **não reclama** se o valor está numa variável. Cuidado!

---

## 🔤 Números ↔ Strings: use `strconv`!

### ❌ A armadilha do `string(numero)`

```go
n := 65
s := string(n)
fmt.Println(s) // "A"  ← NÃO é "65"!
```

`string(int)` interpreta o número como **código de caractere Unicode**. O `go vet` até avisa sobre isso.

### ✅ O jeito certo: pacote `strconv`

```go
import "strconv"

// int -> string
s := strconv.Itoa(65)        // "65"   (Itoa = Integer to ASCII)

// string -> int
n, err := strconv.Atoi("123") // 123   (Atoi = ASCII to Integer)
if err != nil {
	fmt.Println("não é um número válido:", err)
}

n2, err := strconv.Atoi("abc")
fmt.Println(n2, err) // 0 strconv.Atoi: parsing "abc": invalid syntax
```

> 📌 Converter **de** string **pode falhar** (o usuário pode digitar "abc"), por isso as funções retornam um `error`. Sempre trate!

### Outras conversões com `strconv`

```go
// string -> float64
f, err := strconv.ParseFloat("3.14", 64)

// string -> bool  (aceita "true", "false", "1", "0", "T", "F"...)
b, err := strconv.ParseBool("true")

// string -> int com base e tamanho específicos
i, err := strconv.ParseInt("ff", 16, 64) // 255 (hexadecimal)

// float -> string
s := strconv.FormatFloat(3.14159, 'f', 2, 64) // "3.14"

// bool -> string
s2 := strconv.FormatBool(true) // "true"

// int -> string em outra base
s3 := strconv.FormatInt(255, 2) // "11111111"
```

> O seu `main.go` da livraria já usa isso: `strconv.ParseFloat` para o preço e `strconv.Atoi` para o estoque. 😉

### Alternativa rápida: `fmt.Sprint`

Converte **qualquer coisa** para string (mais lento que `strconv`, mas prático):

```go
s := fmt.Sprint(42)          // "42"
s2 := fmt.Sprintf("%.2f", 3.14159) // "3.14"
```

---

## 🔠 String ↔ bytes e runes

```go
s := "Olá"

bytes := []byte(s)      // [79 108 195 161]
runes := []rune(s)      // [79 108 225]

fmt.Println(len(bytes)) // 4
fmt.Println(len(runes)) // 3 ← quantidade real de caracteres

volta := string(runes)  // "Olá"
```

---

## 🏷️ Conversão entre tipos personalizados

Tipos com a **mesma base** podem ser convertidos entre si:

```go
type Celsius float64
type Fahrenheit float64

func CParaF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

fmt.Println(CParaF(100)) // 212
```

---

## 📋 Resumo

| De → Para | Como |
|---|---|
| `int` → `float64` | `float64(i)` |
| `float64` → `int` | `int(f)` (trunca) ou `int(math.Round(f))` |
| `int` → `string` | `strconv.Itoa(i)` ✅ · ~~`string(i)`~~ ❌ |
| `string` → `int` | `strconv.Atoi(s)` (retorna erro) |
| `string` → `float64` | `strconv.ParseFloat(s, 64)` |
| `string` → `bool` | `strconv.ParseBool(s)` |
| qualquer → `string` | `fmt.Sprint(v)` |
| `string` → bytes/runes | `[]byte(s)` / `[]rune(s)` |

---

## ✍️ Exercícios

1. Declare `a := 7` e `b := 2`. Calcule a divisão **com casas decimais** (`3.5`).
2. Converta `7.8` para `int`. O resultado é `7` ou `8`? Como obter `8`?
3. O que `string(97)` retorna? E `strconv.Itoa(97)`?
4. Faça um programa que leia a idade do usuário como texto, converta para `int` com `strconv.Atoi` e, **se houver erro**, mostre "idade inválida".
5. Converta a string `"101010"` (binário) para inteiro usando `strconv.ParseInt`.
6. Crie os tipos `Reais` e `Dolares` e uma função que converta reais para dólares (use uma cotação fixa).

---

⬅️ Anterior: [Constantes e iota](03-constantes-e-iota.md) · 🏠 [Voltar ao roteiro](../README.md)

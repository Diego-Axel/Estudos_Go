# 02 — Tipos de Dados

Go tem os **tipos básicos** abaixo. Os tipos compostos (arrays, slices, maps, structs...) terão módulos próprios.

```
bool

string

int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr

byte  // apelido para uint8
rune  // apelido para int32 (representa um caractere Unicode)

float32 float64

complex64 complex128
```

---

## ✅ Booleano: `bool`

Só dois valores: `true` ou `false`.

```go
ativo := true
maiorDeIdade := idade >= 18
```

> ⚠️ Em Go, `bool` **não** é número. `if 1 { }` é erro. Não existe "truthy/falsy" como em JavaScript/Python.

---

## 🔢 Inteiros

### Com sinal (aceitam negativos)

| Tipo | Tamanho | Faixa |
|---|---|---|
| `int8` | 8 bits | -128 a 127 |
| `int16` | 16 bits | -32.768 a 32.767 |
| `int32` | 32 bits | ~ -2,1 bilhões a 2,1 bilhões |
| `int64` | 64 bits | ~ -9,2 × 10¹⁸ a 9,2 × 10¹⁸ |
| `int` | 32 **ou** 64 bits | depende da plataforma (64 bits nos PCs atuais) |

### Sem sinal (só zero e positivos)

| Tipo | Tamanho | Faixa |
|---|---|---|
| `uint8` / `byte` | 8 bits | 0 a 255 |
| `uint16` | 16 bits | 0 a 65.535 |
| `uint32` | 32 bits | 0 a ~4,2 bilhões |
| `uint64` | 64 bits | 0 a ~1,8 × 10¹⁹ |
| `uint` | 32 ou 64 bits | depende da plataforma |
| `uintptr` | tamanho de ponteiro | uso de baixo nível |

> 💡 **Regra prática:** use **`int`** para quase tudo. Use os tamanhos específicos só quando houver motivo (protocolos, arquivos binários, economia de memória em grandes volumes).

### Formas de escrever números inteiros

```go
decimal := 255
binario := 0b11111111
octal   := 0o377
hexa    := 0xFF
grande  := 1_000_000 // underscore só para facilitar leitura

fmt.Println(decimal, binario, octal, hexa, grande) // 255 255 255 255 1000000
```

### ⚠️ Overflow

Passar do limite **"dá a volta"**, sem aviso em tempo de execução:

```go
var x int8 = 127
x++
fmt.Println(x) // -128 😱
```

Mas se for uma **constante**, o compilador pega:

```go
var y int8 = 200 // ❌ erro: cannot use 200 (untyped int constant) as int8 value (overflows)
```

### Limites no pacote `math`

```go
import "math"

fmt.Println(math.MaxInt8, math.MinInt8) // 127 -128
fmt.Println(math.MaxInt64)              // 9223372036854775807
fmt.Println(uint64(math.MaxUint64))     // 18446744073709551615
```

---

## 🌊 Ponto flutuante

| Tipo | Precisão |
|---|---|
| `float32` | ~7 dígitos |
| `float64` | ~15 dígitos ⭐ **(padrão)** |

```go
pi := 3.14159         // float64 (inferência sempre dá float64)
var f32 float32 = 1.5
cientifico := 1.5e3   // 1500
```

### ⚠️ Imprecisão (acontece em qualquer linguagem)

```go
fmt.Println(0.1 + 0.2) // 0.30000000000000004
```

> 💰 **Nunca use float para dinheiro!** Guarde em centavos com `int64` (ex: R$ 10,50 → `1050`).

### Valores especiais

```go
fmt.Println(math.Inf(1))   // +Inf
fmt.Println(math.NaN())    // NaN
fmt.Println(math.MaxFloat64)
```

---

## 🌀 Números complexos

Pouco usados no dia a dia, mas existem:

```go
c := 3 + 4i               // complex128
fmt.Println(real(c))      // 3
fmt.Println(imag(c))      // 4
```

---

## 🔤 `string`

Sequência de **bytes** (normalmente texto em UTF-8), **imutável**.

```go
s := "Olá, Go!"          // aspas duplas: aceita escapes (\n, \t, \")
raw := `C:\pasta\arquivo
várias linhas, sem escapes` // crase: string "crua" (raw string)
```

### Escapes comuns

| Escape | Significado |
|---|---|
| `\n` | nova linha |
| `\t` | tab |
| `\"` | aspas |
| `\\` | barra invertida |
| `\u00e9` | caractere Unicode (é) |

### Operações básicas

```go
s := "Go"
fmt.Println(len(s))       // 2 (quantidade de BYTES)
fmt.Println(s + "lang")   // Golang (concatenação)
fmt.Println(s[0])         // 71 (é um byte! o código de 'G')
fmt.Println(string(s[0])) // G

// s[0] = 'g' // ❌ erro: strings são imutáveis
```

> ⚠️ `len` conta **bytes**, não letras: `len("ã")` é `2`. Para contar caracteres use `utf8.RuneCountInString`. Isso vai ser aprofundado no módulo **Strings, Runes e Bytes**.

---

## 🅰️ `byte` e `rune`

- **`byte`** = `uint8` → representa **um byte** (ex: dados binários, caracteres ASCII)
- **`rune`** = `int32` → representa **um caractere Unicode** (code point)

Caracteres usam **aspas simples**:

```go
var b byte = 'A'
var r rune = 'ã'

fmt.Println(b, r)                 // 65 227
fmt.Printf("%c %c\n", b, r)       // A ã
fmt.Printf("%T %T\n", b, r)       // uint8 int32
```

Percorrendo uma string com `range`, você recebe **runes**:

```go
for i, r := range "Olá" {
	fmt.Printf("%d: %c\n", i, r)
}
// 0: O
// 1: l
// 2: á
```

---

## 🏷️ Criando seus próprios tipos

Você pode criar um **tipo novo** baseado em outro:

```go
type Celsius float64
type Fahrenheit float64

var temp Celsius = 36.5
var f Fahrenheit = 97.7

// temp = f // ❌ erro: são tipos DIFERENTES, mesmo ambos sendo float64
temp = Celsius(f) // ✅ conversão explícita
```

Isso ajuda o compilador a evitar erros de lógica (misturar unidades, IDs de coisas diferentes, etc.).

### Apelido de tipo (*alias*)

Com `=`, não cria tipo novo, só outro **nome** para o mesmo tipo:

```go
type MeuInt = int // alias: MeuInt e int são exatamente o mesmo tipo
```

É assim que `byte` e `rune` são definidos na linguagem.

---

## 🔎 Descobrindo o tipo

```go
fmt.Printf("%T\n", 42)        // int
fmt.Printf("%T\n", 3.14)      // float64
fmt.Printf("%T\n", "oi")      // string
fmt.Printf("%T\n", 'x')       // int32
fmt.Printf("%T\n", true)      // bool
fmt.Printf("%T\n", 2+3i)      // complex128
```

---

## ✍️ Exercícios

1. Qual o maior valor de um `uint8`? E de um `int8`? Imprima usando o pacote `math`.
2. Crie `var x uint8 = 255` e faça `x++`. Qual o resultado? Por quê?
3. Escreva o número `42` em decimal, binário, octal e hexadecimal e imprima todos.
4. Imprima `len("Go")` e `len("Ação")`. Por que os resultados não batem com a quantidade de letras?
5. Crie os tipos `Metros` e `Pes` (ambos `float64`) e tente somar uma variável de cada. Depois corrija com conversão.
6. Por que não devemos usar `float64` para valores monetários? Como guardar R$ 19,99?

---

⬅️ Anterior: [Variáveis](01-variaveis.md) · ➡️ Próximo: [Constantes e iota](03-constantes-e-iota.md)

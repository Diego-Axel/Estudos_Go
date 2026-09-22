# 01 — UTF-8, Bytes e Runes

Pra dominar strings em Go, você precisa entender **uma coisa**:

> 🧠 **Uma string em Go é uma sequência de BYTES (somente leitura), e não de caracteres.**
> Por convenção, esses bytes estão em **UTF-8**.

---

## 🔤 O que é UTF-8?

**Unicode** dá um número (*code point*) para cada caractere do mundo: `A` = 65, `ã` = 227, `€` = 8364, `😀` = 128512...

**UTF-8** é a forma de guardar esses números em bytes, usando **de 1 a 4 bytes** por caractere:

| Caractere | Code point | Bytes em UTF-8 | Qtd. |
|---|---|---|---|
| `A` | U+0041 | `41` | 1 |
| `ã` | U+00E3 | `C3 A3` | 2 |
| `€` | U+20AC | `E2 82 AC` | 3 |
| `😀` | U+1F600 | `F0 9F 98 80` | 4 |

Os caracteres ASCII (letras sem acento, números, símbolos básicos) usam **1 byte**, então textos em inglês ficam compactos. Acentos já ocupam **2**.

> 💡 Curiosidade: Rob Pike e Ken Thompson, dois dos criadores de Go, também **inventaram o UTF-8** (em 1992).

---

## 📏 `len` conta bytes, não letras

```go
import "unicode/utf8"

s := "Olá"
fmt.Println(len(s))                    // 4  ← bytes (O, l, e 2 bytes do á)
fmt.Println(utf8.RuneCountInString(s)) // 3  ← caracteres

fmt.Println(len("😀"))                 // 4
fmt.Println(len([]rune("😀")))         // 1
```

---

## 🧱 `byte` e `rune`

| Tipo | É apelido de | Representa |
|---|---|---|
| `byte` | `uint8` | **um byte** (0 a 255) |
| `rune` | `int32` | **um code point Unicode** (um "caractere") |

```go
var b byte = 'A'
var r rune = 'ã'

fmt.Println(b, r)          // 65 227
fmt.Printf("%c %c\n", b, r) // A ã
fmt.Printf("%U\n", r)       // U+00E3
```

> Literais entre **aspas simples** (`'a'`) são runes. Entre **aspas duplas** (`"a"`) são strings.

---

## 🎯 Indexar uma string retorna um **byte**

```go
s := "Olá"
fmt.Println(s[0])         // 79  (byte do 'O')
fmt.Println(s[2])         // 195 (primeiro byte do 'á')
fmt.Printf("%c\n", s[2])  // Ã   ← lixo! pegou só metade do 'á'
```

Fatiar também trabalha com **bytes**. Cuidado para não cortar um caractere ao meio:

```go
s := "ação"
fmt.Println(s[0:2]) // "a\xc3" → imprime "a�" (cortou o ç ao meio)
fmt.Println(s[0:3]) // "aç"
```

---

## 🔁 Percorrendo: bytes vs. runes

### `for` clássico → **bytes**

```go
s := "Olá"
for i := 0; i < len(s); i++ {
	fmt.Printf("%d:%x ", i, s[i])
}
// 0:4f 1:6c 2:c3 3:a1
```

### `for range` → **runes** ⭐

```go
for i, r := range "Olá" {
	fmt.Printf("%d:%c ", i, r)
}
// 0:O 1:l 2:á
```

O `i` é o índice do **byte** onde o caractere começa (por isso pode "pular" números).

> ✅ **Regra:** para processar **texto**, use `range` ou converta para `[]rune`.

---

## 🔄 Conversões

```go
s := "Olá"

b := []byte(s)   // [79 108 195 161]
r := []rune(s)   // [79 108 225]

fmt.Println(string(b)) // Olá
fmt.Println(string(r)) // Olá

fmt.Println(string(r[2]))  // á   ✅ (índice de caractere)
fmt.Println(len(r))        // 3
```

> ⚠️ `[]byte(s)` e `[]rune(s)` **copiam** os dados (strings são imutáveis, então a cópia é necessária).

### Invertendo uma string (do jeito certo)

```go
func inverter(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

fmt.Println(inverter("ação")) // oãça
```

Se você invertesse os **bytes**, os acentos virariam lixo.

---

## 🔒 Strings são imutáveis

```go
s := "gato"
// s[0] = 'p' // ❌ erro: cannot assign to s[0] (neither addressable nor a map index expression)

// Para "alterar", crie uma nova string:
b := []byte(s)
b[0] = 'p'
s = string(b)
fmt.Println(s) // pato
```

Por serem imutáveis, strings podem ser **compartilhadas** com segurança e fatiar (`s[2:5]`) é barato: não copia os bytes.

---

## 🧰 Pacote `unicode/utf8`

```go
import "unicode/utf8"

s := "Olá, 世界"

fmt.Println(utf8.RuneCountInString(s)) // 7
fmt.Println(utf8.RuneLen('世'))        // 3 (bytes que ocupa)
fmt.Println(utf8.ValidString(s))       // true

// Decodifica a primeira rune e diz quantos bytes ela ocupa
r, tamanho := utf8.DecodeRuneInString("ãbc")
fmt.Printf("%c %d\n", r, tamanho)      // ã 2
```

### Bytes inválidos

Se a string tiver bytes que **não** formam UTF-8 válido, o `range` devolve o caractere de substituição `�` (U+FFFD):

```go
invalida := "a\xffb"
for _, r := range invalida {
	fmt.Printf("%c ", r)
}
// a � b
fmt.Println(utf8.ValidString(invalida)) // false
```

---

## 📝 Literais de string

```go
normal := "Linha 1\nLinha 2\t(tab) \"aspas\""

crua := `Não processa \n nem \t
pode ter várias linhas
e "aspas" à vontade`

unicode := "é \U0001F600" // é 😀
hexa := "\x41\x42"             // AB
```

> 💡 Strings cruas (entre crases) são ótimas para **regex**, **JSON**, **SQL** e **caminhos do Windows**.

---

## ✍️ Exercícios

1. Imprima `len` e `utf8.RuneCountInString` de: `"Go"`, `"Golang é top"`, `"Coração"`, `"日本"`, `"👍🏽"`. Explique as diferenças.
2. Percorra `"Pão de Queijo"` com `for` clássico e depois com `range`, mostrando índice e valor. Compare.
3. Escreva `primeiraLetra(s string) string` que retorne o primeiro **caractere** (funcione com `"Ágata"`).
4. Escreva `ehPalindromo(s string) bool` que funcione com acentos (`"ovo"`, `"açaça"`).
5. Mostre os bytes em hexadecimal de `"ç"` e o code point com `%U`.
6. Por que `s[0] = 'x'` não compila? Como trocar a primeira letra de uma string?
7. Escreva `truncar(s string, n int) string` que retorne os `n` primeiros **caracteres** (sem cortar acentos ao meio).

---

🏠 [Módulo 07](README.md) · ➡️ Próximo: [Pacote strings](02-pacote-strings.md)

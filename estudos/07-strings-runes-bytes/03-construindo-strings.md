# 03 — Construindo Strings com Eficiência

## 🐢 O problema de concatenar com `+` em laço

Strings são **imutáveis**. Cada `+` cria uma **string nova**, copiando tudo que já existia:

```go
resultado := ""
for i := range 10000 {
	resultado += strconv.Itoa(i) + "," // copia TUDO de novo a cada volta 😬
}
```

Com 10.000 voltas, são milhões de bytes copiados à toa. Para **poucas** concatenações, `+` é ótimo. Para **laços**, use `strings.Builder`.

---

## 🏗️ `strings.Builder` ⭐

Ele guarda os pedaços num buffer interno que cresce conforme precisa, e só monta a string final **uma vez**:

```go
import "strings"

var sb strings.Builder // valor zero já está pronto pra usar

for i := range 5 {
	sb.WriteString("item")
	sb.WriteString(strconv.Itoa(i))
	sb.WriteByte(' ')
}

resultado := sb.String()
fmt.Println(resultado) // item0 item1 item2 item3 item4
```

### Métodos

```go
var sb strings.Builder

sb.WriteString("Olá")   // adiciona string
sb.WriteByte(',')       // adiciona um byte
sb.WriteRune('✓')       // adiciona uma rune
sb.Write([]byte(" ok")) // adiciona []byte

fmt.Println(sb.String()) // Olá,✓ ok
fmt.Println(sb.Len())    // tamanho em bytes

sb.Reset()               // esvazia
```

### Com formatação: `fmt.Fprintf`

O `Builder` pode ser usado como **destino** do `fmt`:

```go
var sb strings.Builder
for i, nome := range []string{"Ana", "Bia"} {
	fmt.Fprintf(&sb, "%d. %s\n", i+1, nome)
}
fmt.Print(sb.String())
// 1. Ana
// 2. Bia
```

### Pré-alocando

Se tem noção do tamanho final, use `Grow`:

```go
var sb strings.Builder
sb.Grow(1024) // reserva ~1KB de uma vez
```

> ⚠️ **Não copie** um `Builder` depois de usá-lo (passe `*strings.Builder` para funções). Copiar causa panic.

---

## 🧮 Comparando as formas

| Forma | Quando usar |
|---|---|
| `a + b` | poucas concatenações, código simples |
| `fmt.Sprintf("%s-%d", a, n)` | quando precisa **formatar** (números, padding...) |
| `strings.Join(lista, ",")` | juntar um slice com separador ⭐ |
| `strings.Builder` | montar texto em **laço** / muitas partes |
| `bytes.Buffer` | quando também precisa **ler** do buffer ou trabalhar com `[]byte` |

```go
// Todas produzem "a,b,c"
s1 := "a" + "," + "b" + "," + "c"
s2 := fmt.Sprintf("%s,%s,%s", "a", "b", "c")
s3 := strings.Join([]string{"a", "b", "c"}, ",")
```

---

## 🧺 `bytes.Buffer`

Parecido com o `Builder`, mas mais genérico: dá pra **escrever e ler**:

```go
import "bytes"

var buf bytes.Buffer
buf.WriteString("Olá ")
buf.WriteString("mundo")

fmt.Println(buf.String()) // Olá mundo
```

> Para **só montar** strings, o `strings.Builder` é o recomendado (é um pouco mais eficiente no `String()`).

---

## 📦 Pacote `bytes`: o "irmão" do `strings`

Quase toda função de `strings` tem uma versão em `bytes` que trabalha com `[]byte`:

```go
import "bytes"

dados := []byte("Go é legal")
fmt.Println(bytes.Contains(dados, []byte("legal"))) // true
fmt.Println(string(bytes.ToUpper(dados)))           // GO É LEGAL
```

Útil ao lidar com arquivos, rede e JSON, que normalmente trabalham com `[]byte`, **sem** precisar converter para string.

---

## 📖 Lendo strings como fluxo: `strings.NewReader`

Transforma uma string em algo que pode ser **lido** (um `io.Reader`). Útil para testar código que lê de arquivos ou do teclado:

```go
import (
	"bufio"
	"strings"
)

entrada := strings.NewReader("linha 1\nlinha 2\nlinha 3")
scanner := bufio.NewScanner(entrada)

for scanner.Scan() {
	fmt.Println(">", scanner.Text())
}
// > linha 1
// > linha 2
// > linha 3
```

> É o mesmo `bufio.Scanner` que seu `main.go` usa com `os.Stdin`. Troque `os.Stdin` por `strings.NewReader(...)` e dá pra testar o menu sem digitar nada!

---

## 🔢 Números ↔ strings (revisão)

```go
import "strconv"

strconv.Itoa(42)                     // "42"
strconv.Atoi("42")                   // 42, nil
strconv.FormatFloat(3.14159, 'f', 2, 64) // "3.14"
strconv.ParseFloat("3.14", 64)       // 3.14, nil
strconv.Quote("olá\n")               // "\"olá\\n\"" (com aspas e escapes)
```

Veja mais em [Conversão de tipos](../02-fundamentos/04-conversao-de-tipos.md).

---

## ✍️ Exercícios

1. Monte com `strings.Builder` uma string com os números de 1 a 100 separados por vírgula (sem vírgula no final!).
2. Faça o mesmo do exercício 1 com `strings.Join` (dica: crie um `[]string` antes). Qual ficou mais simples?
3. Escreva `repetirComSeparador(s, sep string, n int) string`: `("ab", "-", 3)` → `"ab-ab-ab"`.
4. Gere uma tabela de tabuada (1 a 10) com `fmt.Fprintf` escrevendo num `strings.Builder` e imprima no final.
5. Use `strings.NewReader` + `bufio.Scanner` para ler um texto de várias linhas e contar quantas linhas **não estão vazias**.
6. Escreva uma função que receba `[]string` e retorne uma lista HTML: `<ul><li>a</li><li>b</li></ul>`.

---

⬅️ Anterior: [Pacote strings](02-pacote-strings.md) · ➡️ Próximo: [Unicode e formatação](04-unicode-e-formatacao.md)

# 04 — Pacote `unicode` e Formatação com `fmt`

## 🔣 Pacote `unicode`: classificando caracteres

Trabalha com **runes** e entende acentos e outros alfabetos:

```go
import "unicode"

unicode.IsLetter('ã')  // true
unicode.IsLetter('7')  // false
unicode.IsDigit('7')   // true
unicode.IsNumber('½')  // true  (IsNumber é mais amplo que IsDigit)
unicode.IsSpace('\t')  // true  (espaço, tab, quebra de linha...)
unicode.IsPunct('!')   // true
unicode.IsUpper('Ç')   // true
unicode.IsLower('ç')   // true

unicode.ToUpper('ã')   // 'Ã'
unicode.ToLower('Ç')   // 'ç'
```

### Exemplo: analisando uma senha

```go
func analisarSenha(senha string) {
	var letras, digitos, maiusculas, especiais int

	for _, r := range senha {
		switch {
		case unicode.IsDigit(r):
			digitos++
		case unicode.IsLetter(r):
			letras++
			if unicode.IsUpper(r) {
				maiusculas++
			}
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			especiais++
		}
	}

	fmt.Printf("letras=%d dígitos=%d maiúsculas=%d especiais=%d\n",
		letras, digitos, maiusculas, especiais)
}

analisarSenha("Senh@123") // letras=4 dígitos=3 maiúsculas=1 especiais=1
```

### Exemplo: capitalizar a primeira letra

```go
func capitalizar(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

fmt.Println(capitalizar("ágata")) // Ágata
```

---

## 🖨️ Formatação com `fmt`: guia completo

Você já conhece os verbos básicos (`%v`, `%d`, `%s`...) do [primeiro programa](../01-introducao/03-primeiro-programa.md). Agora vamos ao **controle fino**.

### Estrutura de um verbo

```
%[flags][largura][.precisão]verbo
```

### Largura e alinhamento

```go
fmt.Printf("|%5d|\n", 42)     // |   42|  ← largura 5, alinhado à direita
fmt.Printf("|%-5d|\n", 42)    // |42   |  ← "-" alinha à esquerda
fmt.Printf("|%05d|\n", 42)    // |00042|  ← "0" preenche com zeros

fmt.Printf("|%10s|\n", "Go")  // |        Go|
fmt.Printf("|%-10s|\n", "Go") // |Go        |
```

> A largura de `%s` conta **caracteres (runes)**, não bytes, então acentos alinham certinho.

### Precisão

```go
fmt.Printf("%.2f\n", 3.14159)   // 3.14      (casas decimais)
fmt.Printf("%8.2f|\n", 3.14159) //     3.14| (largura 8 + 2 casas)
fmt.Printf("%.3s\n", "Golang")  // Gol       (em string: corta!)
```

### Flags

| Flag | Efeito | Exemplo | Saída |
|---|---|---|---|
| `-` | alinha à esquerda | `%-6d\|` | `42    \|` |
| `0` | preenche com zeros | `%06.2f` | `003.14` |
| `+` | sempre mostra o sinal | `%+d` | `+42` |
| ` ` (espaço) | espaço no lugar do `+` / separa bytes em `% x` | `% x` em `"Go"` | `47 6f` |
| `#` | forma alternativa | `%#x` | `0x2a` |

### Verbos para números

```go
n := 255
fmt.Printf("%d %b %o %x %X %#x\n", n, n, n, n, n, n)
// 255 11111111 377 ff FF 0xff

f := 1234567.891
fmt.Printf("%f | %.1f | %e | %g\n", f, f, f, f)
// 1234567.891000 | 1234567.9 | 1.234568e+06 | 1.234567891e+06
```

### Verbos para strings e runes

```go
s := "Olá"
fmt.Printf("%s\n", s)   // Olá
fmt.Printf("%q\n", s)   // "Olá"
fmt.Printf("%x\n", s)   // 4f6cc3a1      (bytes em hexa)
fmt.Printf("% x\n", s)  // 4f 6c c3 a1

r := 'ã'
fmt.Printf("%c %q %U %d\n", r, r, r, r) // ã 'ã' U+00E3 227
```

### Verbos genéricos

```go
type Livro struct {
	Titulo string
	Paginas int
}
l := Livro{"Go", 300}

fmt.Printf("%v\n", l)   // {Go 300}
fmt.Printf("%+v\n", l)  // {Titulo:Go Paginas:300}
fmt.Printf("%#v\n", l)  // main.Livro{Titulo:"Go", Paginas:300}
fmt.Printf("%T\n", l)   // main.Livro
```

### Reutilizando argumentos: `%[n]`

```go
fmt.Printf("%[1]d em binário é %[1]b e em hexa é %[1]x\n", 10)
// 10 em binário é 1010 e em hexa é a
```

### Largura dinâmica: `*`

```go
largura := 8
fmt.Printf("|%*d|\n", largura, 42) // |      42|
```

---

## 📊 Montando uma tabela alinhada

```go
produtos := []struct {
	Nome  string
	Preco float64
	Qtd   int
}{
	{"Caneta", 2.5, 100},
	{"Caderno universitário", 24.9, 12},
	{"Lápis", 1.2, 250},
}

fmt.Printf("%-22s %10s %5s\n", "PRODUTO", "PREÇO", "QTD")
fmt.Println(strings.Repeat("-", 39))
for _, p := range produtos {
	fmt.Printf("%-22s %10.2f %5d\n", p.Nome, p.Preco, p.Qtd)
}
```

Saída:

```
PRODUTO                     PREÇO   QTD
---------------------------------------
Caneta                       2.50   100
Caderno universitário       24.90    12
Lápis                        1.20   250
```

> Para tabelas mais complexas, veja o pacote `text/tabwriter`.

---

## 🧾 Família `Print`

| Função | Escreve em | Retorna |
|---|---|---|
| `Print`, `Println`, `Printf` | saída padrão (terminal) | — |
| `Sprint`, `Sprintln`, `Sprintf` | — | **string** |
| `Fprint`, `Fprintln`, `Fprintf` | qualquer `io.Writer` (arquivo, `os.Stderr`, `Builder`...) | — |
| `Errorf` | — | **error** formatado |

```go
fmt.Fprintln(os.Stderr, "erro: arquivo não encontrado") // saída de erro
err := fmt.Errorf("usuário %q não encontrado", "ana")
```

---

## 🧾 Resumão do módulo

- String = **bytes** imutáveis em UTF-8; `len` conta **bytes**
- `rune` = caractere Unicode; use `range` ou `[]rune` para processar texto
- Pacote **`strings`**: buscar, dividir, juntar, substituir, limpar
- **`strings.Builder`** para montar strings em laço
- Pacote **`unicode`** para classificar caracteres
- **`fmt`**: largura, precisão e flags para formatar a saída

---

## ✍️ Exercícios

1. Conte quantas **letras, dígitos, espaços e pontuações** tem a frase `"Go 1.25 chegou! É rápido, simples e 100% legal."`.
2. Valide uma senha: mínimo 8 caracteres, ao menos uma maiúscula, uma minúscula, um dígito e um caractere especial.
3. Imprima os números de 1 a 16 em uma tabela com colunas: decimal, binário (8 dígitos com zeros), octal e hexadecimal.
4. Monte um "recibo" alinhado com nome do produto, quantidade, preço unitário e subtotal, com o **total** no final.
5. Use `%[1]` para imprimir um número em decimal, hexa e binário passando o número **uma única vez**.
6. Escreva `centralizar(s string, largura int) string` que centralize um texto preenchendo com espaços.
7. **Desafio:** cifra de César: desloque cada **letra** 3 posições (`a→d`, `z→c`), mantendo maiúsculas/minúsculas e sem alterar outros caracteres.

---

⬅️ Anterior: [Construindo strings](03-construindo-strings.md) · 🏠 [Voltar ao roteiro](../README.md)

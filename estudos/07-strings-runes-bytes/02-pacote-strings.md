# 02 — O Pacote `strings`

O pacote [`strings`](https://pkg.go.dev/strings) tem praticamente tudo que você precisa para manipular texto. Todas as funções **retornam uma string nova** (a original nunca muda).

```go
import "strings"
```

---

## 🔍 Buscar

```go
s := "Go é simples e Go é rápido"

strings.Contains(s, "simples")     // true
strings.ContainsAny(s, "xyz")      // false (algum desses caracteres?)
strings.ContainsRune(s, 'é')       // true

strings.HasPrefix(s, "Go")         // true (começa com)
strings.HasSuffix(s, "rápido")     // true (termina com)

strings.Index(s, "Go")             // 0   (primeira ocorrência, em BYTES)
strings.LastIndex(s, "Go")         // 16  (última ocorrência)
strings.Index(s, "Rust")           // -1  (não encontrado)

strings.Count(s, "Go")             // 2
```

> ⚠️ Os índices retornados são posições em **bytes**, não em caracteres.

---

## ✂️ Dividir e juntar

### `Split`

```go
csv := "ana,bia,caio"
partes := strings.Split(csv, ",")
fmt.Println(partes)      // [ana bia caio]
fmt.Println(len(partes)) // 3
```

Pegadinhas do `Split`:

```go
strings.Split("a,,b", ",")  // [a  b]  ← 3 itens, o do meio é ""
strings.Split("", ",")      // []      ← parece vazio, mas tem 1 item: ""
strings.Split("abc", "")    // [a b c] ← separador vazio: divide em caracteres

strings.SplitN("a,b,c,d", ",", 2) // [a b,c,d] ← no máximo 2 partes
```

### `Fields`: divide por espaços (qualquer quantidade)

```go
frase := "  Go   é    muito   bom  "
fmt.Println(strings.Fields(frase))       // [Go é muito bom]
fmt.Println(len(strings.Fields(frase)))  // 4 → contar palavras!
```

### `Join`

```go
nomes := []string{"Ana", "Bia", "Caio"}
fmt.Println(strings.Join(nomes, ", ")) // Ana, Bia, Caio
```

### `Cut` (Go 1.18+) ⭐: divide em **duas** partes

```go
antes, depois, achou := strings.Cut("chave=valor", "=")
fmt.Println(antes, depois, achou) // chave valor true

usuario, dominio, _ := strings.Cut("diego@email.com", "@")
fmt.Println(usuario, dominio)     // diego email.com
```

> Muito mais limpo que `Index` + fatiamento para ler `chave=valor`, e-mails, cabeçalhos HTTP...

---

## 🔁 Substituir

```go
s := "banana"

strings.Replace(s, "a", "o", 2)  // "bonona" (só as 2 primeiras)
strings.Replace(s, "a", "o", -1) // "bonono" (todas; -1 = sem limite)
strings.ReplaceAll(s, "a", "o")  // "bonono" (mesmo que -1)
```

Várias substituições de uma vez com `Replacer`:

```go
r := strings.NewReplacer("ã", "a", "ç", "c", "é", "e")
fmt.Println(r.Replace("ação é"))  // acao e
```

---

## 🔠 Maiúsculas e minúsculas

```go
strings.ToUpper("ação")           // AÇÃO  (funciona com acentos)
strings.ToLower("GoLang")         // golang

strings.EqualFold("GoLang", "GOLANG") // true (compara ignorando caixa)
```

> ⚠️ `strings.Title` está **obsoleta**. Para "Primeira Letra Maiúscula", use o pacote `golang.org/x/text/cases` ou faça na mão com runes (exercício abaixo).

---

## 🧽 Remover espaços e caracteres das pontas

```go
strings.TrimSpace("  \t olá \n ")       // "olá"

strings.TrimPrefix("R$ 10,00", "R$ ")  // "10,00"
strings.TrimSuffix("foto.png", ".png") // "foto"

strings.Trim("***oi***", "*")          // "oi"
strings.TrimLeft("000123", "0")        // "123"
strings.TrimRight("oi!!!", "!")        // "oi"
```

### ⚠️ `TrimLeft` ≠ `TrimPrefix`

`Trim`, `TrimLeft` e `TrimRight` recebem um **conjunto de caracteres** (*cutset*), não uma palavra:

```go
strings.TrimPrefix("abacate", "ab") // "acate"  ← remove o prefixo "ab" uma vez
strings.TrimLeft("abacate", "ab")   // "cate"   ← remove qualquer 'a' ou 'b' do início!
```

---

## 🔂 Outros úteis

```go
strings.Repeat("ab", 3)        // "ababab"
strings.Repeat("-", 20)        // "--------------------"

strings.Compare("a", "b")      // -1 (a < b), 0 (igual), 1 (a > b)
                               // (normalmente use ==, <, > direto)

// Map: aplica uma função a cada rune
semVogais := strings.Map(func(r rune) rune {
	if strings.ContainsRune("aeiouAEIOU", r) {
		return -1 // -1 remove o caractere
	}
	return r
}, "Golang é legal")
fmt.Println(semVogais)         // Glng é lgl
```

---

## 📋 Resumo

| Precisa... | Use |
|---|---|
| saber se contém | `Contains`, `HasPrefix`, `HasSuffix` |
| achar posição | `Index`, `LastIndex` |
| contar | `Count` |
| dividir | `Split`, `Fields`, `Cut`, `SplitN` |
| juntar | `Join` |
| substituir | `ReplaceAll`, `Replace`, `NewReplacer` |
| mudar caixa | `ToUpper`, `ToLower` |
| comparar sem caixa | `EqualFold` |
| limpar pontas | `TrimSpace`, `TrimPrefix`, `TrimSuffix`, `Trim` |
| repetir | `Repeat` |
| transformar cada caractere | `Map` |

---

## ✍️ Exercícios

1. Leia uma frase e mostre: número de palavras, se contém "Go" (ignorando maiúsculas) e a frase toda em maiúsculas.
2. Dado `"nome=Diego;idade=25;cidade=Natal"`, monte um `map[string]string` usando `Split` e `Cut`.
3. Escreva `capitalizar(s string) string` que deixe a primeira letra de **cada palavra** maiúscula: `"olá mundo go"` → `"Olá Mundo Go"`.
4. Escreva `slug(s string) string`: `"  Aprendendo Go é Legal "` → `"aprendendo-go-e-legal"` (use `TrimSpace`, `ToLower`, `NewReplacer`, `Fields` e `Join`).
5. Verifique se um e-mail é "válido": contém exatamente um `@`, e a parte depois do `@` contém um `.`.
6. Conte quantas vezes cada palavra aparece em um texto (combine `Fields`, `ToLower` e um map).
7. Explique a diferença de resultado entre `strings.TrimLeft("xxyxz", "xy")` e `strings.TrimPrefix("xxyxz", "xy")`.

---

⬅️ Anterior: [UTF-8, bytes e runes](01-utf8-bytes-e-runes.md) · ➡️ Próximo: [Construindo strings](03-construindo-strings.md)

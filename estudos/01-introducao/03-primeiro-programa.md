# 03 — Primeiro Programa

## 🚀 Criando o projeto

```bash
mkdir ola-mundo
cd ola-mundo
go mod init ola-mundo
```

O `go mod init` cria o arquivo **`go.mod`**, que identifica o projeto (módulo):

```go
module ola-mundo

go 1.25
```

---

## 👋 Olá, Mundo!

Crie o arquivo `main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("Olá, Mundo!")
}
```

Rode:

```bash
go run main.go
# ou, dentro da pasta do módulo:
go run .
```

Saída:

```
Olá, Mundo!
```

---

## 🔍 Linha por linha

```go
package main
```
Todo arquivo Go pertence a um **pacote**. O pacote **`main`** é especial: ele indica que isso é um **programa executável** (e não uma biblioteca).

```go
import "fmt"
```
Importa o pacote **`fmt`** (de *format*), da biblioteca padrão, usado para imprimir e ler texto.

```go
func main() {
```
A função **`main`** é o **ponto de entrada**: é onde o programa começa. Não recebe argumentos e não retorna nada.

```go
	fmt.Println("Olá, Mundo!")
```
Chama a função `Println` do pacote `fmt`. Repare que ela começa com **letra maiúscula** — em Go, isso significa que ela é **exportada** (pública). Letra minúscula = privada ao pacote.

---

## 📦 Importando vários pacotes

```go
import (
	"fmt"
	"math"
	"strings"
)
```

> ⚠️ **Import não usado é ERRO de compilação** em Go! O mesmo vale para **variável local não usada**. Isso mantém o código limpo.

```go
import (
	"fmt"
	"os" // ❌ erro: "os" imported and not used
)
```

---

## 🖨️ Imprimindo coisas com `fmt`

```go
fmt.Print("sem quebra de linha")
fmt.Println("com quebra de linha", 42, true) // separa com espaço
fmt.Printf("Nome: %s, Idade: %d\n", "Diego", 25) // formatado
```

### Verbos de formatação mais usados (`Printf`)

| Verbo | Uso | Exemplo | Saída |
|---|---|---|---|
| `%v` | valor padrão (qualquer tipo) | `Printf("%v", 10)` | `10` |
| `%+v` | struct com nomes dos campos | `Printf("%+v", p)` | `{Nome:Ana}` |
| `%T` | **tipo** do valor | `Printf("%T", 3.14)` | `float64` |
| `%d` | inteiro | `Printf("%d", 42)` | `42` |
| `%f` | float | `Printf("%.2f", 3.14159)` | `3.14` |
| `%s` | string | `Printf("%s", "go")` | `go` |
| `%q` | string com aspas | `Printf("%q", "go")` | `"go"` |
| `%t` | booleano | `Printf("%t", true)` | `true` |
| `%c` | caractere (rune) | `Printf("%c", 'A')` | `A` |
| `%x` | hexadecimal | `Printf("%x", 255)` | `ff` |
| `%b` | binário | `Printf("%b", 5)` | `101` |
| `%p` | ponteiro (endereço) | `Printf("%p", &x)` | `0xc000012345` |
| `%%` | o próprio `%` | `Printf("100%%")` | `100%` |

### `Sprintf` — formata e **retorna** a string (não imprime)

```go
msg := fmt.Sprintf("Tenho %d anos", 25)
fmt.Println(msg) // Tenho 25 anos
```

---

## ⌨️ Lendo do teclado

Forma simples com `fmt.Scan`:

```go
package main

import "fmt"

func main() {
	var nome string
	fmt.Print("Qual seu nome? ")
	fmt.Scan(&nome) // o & passa o ENDEREÇO da variável (veremos em Ponteiros)
	fmt.Println("Olá,", nome)
}
```

> ⚠️ `fmt.Scan` para no primeiro espaço. Para ler **a linha inteira**, use `bufio.Scanner` (é o que o seu `main.go` da livraria faz!):

```go
scanner := bufio.NewScanner(os.Stdin)
scanner.Scan()
linha := scanner.Text()
```

---

## 🎨 Formatação automática

Go tem **um estilo oficial único**. Não existe briga de "tabs vs espaços": é **tab**, e ponto.

```bash
go fmt ./...
```

Também: a chave `{` **tem que** ficar na mesma linha:

```go
func main() {   // ✅ certo
}

func main()
{               // ❌ erro de compilação
}
```

E **não precisa de `;`** no fim das linhas (o compilador coloca sozinho).

---

## 💬 Comentários

```go
// Comentário de uma linha

/*
   Comentário de
   várias linhas
*/

// Soma retorna a soma de a e b.   <- comentário de documentação:
func Soma(a, b int) int {         //    fica logo acima e começa com o nome
	return a + b
}
```

---

## ✍️ Exercícios

1. Crie um módulo `exercicio01` e faça um programa que imprima seu nome, idade e cidade, **cada um em uma linha**.
2. Use `Printf` para imprimir: `O valor de PI é aproximadamente 3.14` (use `%.2f` com `3.14159`).
3. Use `%T` para descobrir o tipo de: `42`, `3.14`, `"go"`, `true`, `'A'`.
4. Faça um programa que pergunte o nome e o ano de nascimento do usuário e imprima quantos anos ele tem (aproximadamente).
5. Adicione um `import "os"` sem usar e veja o erro. Depois remova.

---

⬅️ Anterior: [Instalação e ambiente](02-instalacao-e-ambiente.md) · ➡️ Próximo módulo: [Fundamentos](../02-fundamentos/README.md)

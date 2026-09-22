# 01 — Pacotes

## 📦 O que é um pacote?

Um **pacote** é uma **pasta** com arquivos `.go` que declaram o **mesmo nome de pacote** na primeira linha. É a unidade de organização e reutilização de código em Go.

```
calculadora/
├── soma.go         → package calculadora
├── multiplicacao.go → package calculadora
└── util.go         → package calculadora
```

Regras:
- **Uma pasta = um pacote**. Todos os `.go` da pasta (exceto testes, às vezes) têm o mesmo `package`.
- Arquivos do mesmo pacote **enxergam tudo** uns dos outros, até o que é privado, como se fossem um arquivo só.
- Convenção: o nome do pacote é **igual ao nome da pasta**.

---

## 🚀 O pacote `main`

O pacote **`main`** é especial: ele gera um **executável**. Precisa ter a função `func main()`.

Qualquer outro nome de pacote gera uma **biblioteca**, que outros pacotes importam.

| `package main` | `package qualquercoisa` |
|---|---|
| vira executável (`go build`) | vira biblioteca |
| precisa de `func main()` | não tem `main` |
| não pode ser importado | é importado por outros |

---

## 🔓 Visibilidade: exportado ou não

A regra que você já conhece, agora com o motivo: ela funciona **entre pacotes**.

```go
package calculadora

func Somar(a, b int) int { // ✅ Exportado: começa com maiúscula
	return a + b
}

func validar(n int) bool { // 🔒 Não exportado: só dentro de "calculadora"
	return n >= 0
}

const Pi = 3.14159 // ✅ exportada
var versao = "1.0" // 🔒 privada
```

Vale para **tudo**: funções, tipos, variáveis, constantes, campos de struct e métodos.

```go
package main

import "meuprojeto/calculadora"

func main() {
	calculadora.Somar(1, 2)   // ✅
	// calculadora.validar(1) // ❌ erro: cannot refer to unexported name calculadora.validar
}
```

> Não existe `public`, `private` ou `protected` em Go. **A primeira letra decide tudo.**

---

## 📥 Importando pacotes

```go
import "fmt"            // biblioteca padrão
import "math/rand/v2"   // subpacote da biblioteca padrão

import (
	"fmt"
	"strings"

	"github.com/google/uuid"          // pacote externo (de terceiros)
	"meuprojeto/internal/estoque"     // pacote do seu próprio módulo
)
```

> 📌 O `goimports` (e o VS Code) organizam os imports automaticamente: padrão primeiro, depois os externos, separados por uma linha em branco.

O **caminho** do import é o endereço do pacote. O **nome** usado no código é o que está no `package` (normalmente o último pedaço do caminho):

```go
import "math/rand/v2" // caminho
rand.IntN(10)         // nome do pacote: rand
```

### Alias (apelido)

Útil quando dois pacotes têm o mesmo nome, ou o nome é ruim:

```go
import (
	crand "crypto/rand"
	mrand "math/rand/v2"
)

crand.Read(buf)
mrand.IntN(10)
```

### Import "em branco" `_`

Importa o pacote **só pelos efeitos colaterais** (roda o `init()` dele), sem usar nada diretamente:

```go
import (
	"database/sql"
	_ "github.com/lib/pq" // registra o driver do PostgreSQL no database/sql
)

import _ "image/png" // registra o decodificador de PNG
```

### Import com ponto `.` (evite!)

```go
import . "fmt"

Println("sem o fmt. na frente") // funciona, mas confunde: de onde veio Println?
```

> ⚠️ Só se vê isso em alguns testes. No código normal, **não use**.

---

## ⚙️ A função `init()`

Um pacote pode ter funções `init()`, que rodam **automaticamente**, **antes** do `main`, sem ninguém chamar:

```go
package config

import "os"

var Ambiente string

func init() {
	Ambiente = os.Getenv("AMBIENTE")
	if Ambiente == "" {
		Ambiente = "desenvolvimento"
	}
}
```

### Ordem de inicialização

```
1. Pacotes importados são inicializados primeiro (recursivamente)
2. Variáveis de pacote do pacote atual
3. Funções init() do pacote atual (na ordem em que aparecem)
4. main()
```

```go
package main

import "fmt"

var x = inicializarX()

func inicializarX() int {
	fmt.Println("1. variável de pacote")
	return 10
}

func init() {
	fmt.Println("2. init()")
}

func main() {
	fmt.Println("3. main()")
}
// 1. variável de pacote
// 2. init()
// 3. main()
```

Detalhes:
- Pode haver **vários** `init()` no mesmo pacote (até no mesmo arquivo)
- `init()` não recebe nem retorna nada, e **não pode** ser chamada manualmente
- Cada pacote é inicializado **uma única vez**, mesmo importado por vários

> ⚠️ **Use `init()` com moderação.** Código "mágico" que roda escondido é difícil de testar e entender. Prefira funções explícitas (`config.Carregar()`) chamadas no `main`.

---

## 🏷️ Como nomear pacotes

✅ **Bons nomes:**
- **curtos**, **minúsculos**, **uma palavra**: `http`, `json`, `estoque`, `usuario`
- sem `_` nem letras maiúsculas: ~~`meu_pacote`~~, ~~`MeuPacote`~~
- substantivos que dizem **o que o pacote fornece**

❌ **Evite:**
- `util`, `utils`, `common`, `helpers`, `misc`: não dizem nada e viram "gaveta de bagunça"
- repetir o nome do pacote nos identificadores ("*stuttering*"):

```go
estoque.EstoqueProduto // ❌ repetitivo
estoque.Produto        // ✅

http.HTTPServer        // ❌
http.Server            // ✅ (é assim que a biblioteca padrão faz)
```

> 🧠 Lembre que o nome do pacote **sempre aparece** junto: `strings.Builder`, `json.Marshal`. Pense em como fica **na chamada**.

---

## 🔁 Imports cíclicos são proibidos

Se `a` importa `b`, então `b` **não pode** importar `a` (nem indiretamente):

```
pedido → cliente → pedido   ❌ import cycle not allowed
```

Soluções comuns:
- mover o código compartilhado para um **terceiro pacote** que os dois importam
- usar uma **interface** no pacote que "consome" (módulo 10)
- repensar: talvez os dois devessem ser **um pacote só**

---

## ✍️ Exercícios

1. Crie um pacote `calculadora` com `Somar`, `Subtrair` e uma função privada `validar`. Use no `main` e tente chamar a privada.
2. Crie dois arquivos no mesmo pacote e mostre que um enxerga as funções **privadas** do outro.
3. Crie um programa com uma variável de pacote, dois `init()` e o `main`, cada um imprimindo algo. Confira a ordem.
4. Importe `crypto/rand` e `math/rand/v2` no mesmo arquivo usando aliases.
5. Renomeie estes identificadores para evitar repetição: `usuario.UsuarioService`, `config.ConfigLoad`, `pedido.NovoPedido`.
6. Desenhe (no papel) dois pacotes que se importam mutuamente e proponha uma forma de quebrar o ciclo.

---

🏠 [Módulo 12](README.md) · ➡️ Próximo: [Módulos e go.mod](02-modulos-e-go-mod.md)

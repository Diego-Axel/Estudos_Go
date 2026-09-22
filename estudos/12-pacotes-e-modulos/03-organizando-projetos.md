# 03 — Organizando Projetos

## 🌱 Comece simples

Não existe **uma** estrutura oficial obrigatória. A recomendação do time do Go é: **comece pequeno e cresça conforme a necessidade**.

### Nível 1: tudo num pacote só

Para programas pequenos, scripts e estudos:

```
livraria/
├── go.mod
├── main.go
├── produto.go
└── livraria.go
```

Todos os arquivos são `package main`. Simples e perfeitamente válido. (O seu `main.go` da livraria cabe aqui.)

### Nível 2: separando em pacotes

Quando o código cresce e surgem **responsabilidades** diferentes:

```
livraria/
├── go.mod
├── main.go                 ← só monta as peças e inicia
├── produto/
│   ├── produto.go          ← tipo Produto e regras
│   └── produto_test.go
└── armazenamento/
    └── memoria.go          ← onde os produtos ficam guardados
```

### Nível 3: projeto maior com vários executáveis

```
livraria/
├── go.mod
├── cmd/
│   ├── api/
│   │   └── main.go         ← executável 1: servidor HTTP
│   └── cli/
│       └── main.go         ← executável 2: ferramenta de terminal
├── internal/
│   ├── produto/
│   ├── pedido/
│   └── armazenamento/
└── README.md
```

```bash
go run ./cmd/api
go build -o livraria-cli ./cmd/cli
```

---

## 🔒 A pasta `internal/` ⭐

Pacotes dentro de `internal/` só podem ser importados por código que está **dentro da pasta-mãe** do `internal`. **O compilador garante isso.**

```
livraria/
├── internal/
│   └── preco/          ← importável só por livraria/...
└── cmd/api/main.go     ← ✅ pode importar livraria/internal/preco
```

```go
// num OUTRO módulo:
import "github.com/Diego-Axel/livraria/internal/preco"
// ❌ erro: use of internal package ... not allowed
```

Por que isso é ótimo? Você pode **mudar à vontade** o que está em `internal/` sem quebrar ninguém de fora. É a sua "área privada" no nível de módulo.

> 🧠 **Dica:** na dúvida, coloque em `internal/`. É fácil tornar público depois; o contrário quebra quem já depende.

---

## 📂 Pastas comuns (convenções da comunidade)

| Pasta | Para quê |
|---|---|
| `cmd/nome/` | um `main` por executável |
| `internal/` | código privado do módulo ⭐ |
| `pkg/` | código que **você quer** que outros importem (opcional e controverso: muita gente não usa) |
| `api/` | especificações OpenAPI, arquivos `.proto` |
| `web/` ou `static/` | arquivos de front-end, templates |
| `configs/` | arquivos de configuração de exemplo |
| `scripts/` | scripts de build, deploy |
| `testdata/` | arquivos usados pelos testes (o Go **ignora** essa pasta na compilação) |

> ⚠️ O repositório "golang-standards/project-layout" **não é oficial**, apesar do nome. Use como inspiração, não como regra.

---

## 🧩 Organize por **domínio**, não por **camada técnica**

```
❌ Por camada                  ✅ Por domínio
├── models/                    ├── produto/
│   ├── produto.go             │   ├── produto.go
│   └── pedido.go              │   ├── repositorio.go
├── controllers/               │   └── handler.go
│   ├── produto.go             ├── pedido/
│   └── pedido.go              │   ├── pedido.go
└── services/                  │   ├── repositorio.go
    ├── produto.go             │   └── handler.go
    └── pedido.go              └── cliente/
```

Por quê?
- Nomes de pacote melhores: `produto.Repositorio` em vez de `models.Produto`, `services.ProdutoService`...
- Tudo sobre "produto" fica **junto**: mudar uma funcionalidade mexe em **uma** pasta
- Menos risco de **import cíclico** (`models` ↔ `services`)

---

## 🧱 Exemplo completo: a livraria organizada

```
livraria/
├── go.mod                          module github.com/Diego-Axel/livraria
├── cmd/
│   └── livraria/
│       └── main.go
└── internal/
    ├── produto/
    │   ├── produto.go
    │   └── produto_test.go
    └── loja/
        └── loja.go
```

**`internal/produto/produto.go`**

```go
package produto

import (
	"errors"
	"fmt"
)

var ErrPrecoInvalido = errors.New("preço deve ser positivo")

type Produto struct {
	Codigo  string
	Titulo  string
	Preco   float64
	Estoque int
}

func Novo(codigo, titulo string, preco float64, estoque int) (*Produto, error) {
	if preco <= 0 {
		return nil, ErrPrecoInvalido
	}
	return &Produto{Codigo: codigo, Titulo: titulo, Preco: preco, Estoque: estoque}, nil
}

func (p Produto) String() string {
	return fmt.Sprintf("[%s] %s: R$ %.2f (%d un.)", p.Codigo, p.Titulo, p.Preco, p.Estoque)
}
```

**`internal/loja/loja.go`**

```go
package loja

import (
	"errors"

	"github.com/Diego-Axel/livraria/internal/produto"
)

var ErrNaoEncontrado = errors.New("produto não encontrado")

type Loja struct {
	produtos map[string]*produto.Produto
}

func Nova() *Loja {
	return &Loja{produtos: make(map[string]*produto.Produto)}
}

func (l *Loja) Cadastrar(p *produto.Produto) {
	l.produtos[p.Codigo] = p
}

func (l *Loja) Buscar(codigo string) (*produto.Produto, error) {
	p, ok := l.produtos[codigo]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return p, nil
}
```

**`cmd/livraria/main.go`**

```go
package main

import (
	"fmt"
	"log"

	"github.com/Diego-Axel/livraria/internal/loja"
	"github.com/Diego-Axel/livraria/internal/produto"
)

func main() {
	l := loja.Nova()

	p, err := produto.Novo("L001", "O Hobbit", 49.9, 3)
	if err != nil {
		log.Fatal(err)
	}
	l.Cadastrar(p)

	encontrado, err := l.Buscar("L001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encontrado) // [L001] O Hobbit: R$ 49.90 (3 un.)
}
```

Repare como fica a leitura: `produto.Novo`, `loja.Nova`, `loja.ErrNaoEncontrado`. O nome do pacote **faz parte** do nome.

---

## ✅ Dicas finais

1. **Não crie pacotes cedo demais.** Um pacote com um arquivo e uma função costuma ser exagero.
2. **Pacote é uma API**: pense no que ele **exporta**. Quanto menos, melhor.
3. **`main` deve ser fino**: ler configurações, montar as dependências e iniciar. A lógica fica nos pacotes.
4. **Nada de `utils`**: se a função é sobre strings de CPF, crie um pacote `cpf`.
5. **Testes ficam ao lado** do código: `produto.go` + `produto_test.go` (módulo 15).

---

## ✍️ Exercícios

1. Transforme o seu `main.go` da livraria na estrutura do **Nível 2** (pacote `produto` + `main`).
2. Depois evolua para a estrutura com `cmd/` e `internal/` mostrada acima.
3. Crie um segundo módulo que tente importar um pacote `internal/` do primeiro e leia o erro.
4. Reorganize esta estrutura "por camada" em uma estrutura "por domínio":
   ```
   models/usuario.go, models/pedido.go
   services/usuario.go, services/pedido.go
   handlers/usuario.go, handlers/pedido.go
   ```
5. Crie um projeto com **dois executáveis** (`cmd/servidor` e `cmd/cliente`) que usam o mesmo pacote `internal/mensagem`.

---

⬅️ Anterior: [Módulos e go.mod](02-modulos-e-go-mod.md) · ➡️ Próximo: [Workspaces e ferramentas](04-workspaces-e-ferramentas.md)

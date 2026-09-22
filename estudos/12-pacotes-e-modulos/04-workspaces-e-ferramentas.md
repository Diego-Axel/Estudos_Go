# 04 — Workspaces, Build Tags, `embed` e Ferramentas

## 🗂️ Workspaces: vários módulos ao mesmo tempo (Go 1.18+)

Imagine que você está desenvolvendo **dois módulos** juntos: uma biblioteca e um app que a usa.

```
projetos/
├── minhalib/     (go.mod: github.com/Diego-Axel/minhalib)
└── meuapp/       (go.mod: github.com/Diego-Axel/meuapp, usa minhalib)
```

Sem workspace, o `meuapp` baixaria a versão **publicada** da `minhalib`. Para usar a versão **local** (com suas alterações), você teria que usar `replace` no `go.mod`, e lembrar de tirar depois.

Com **workspace**:

```bash
cd projetos
go work init ./minhalib ./meuapp
```

Isso cria o arquivo **`go.work`**:

```go
go 1.25

use (
	./meuapp
	./minhalib
)
```

Pronto: ao rodar `go run`/`go build` dentro de `projetos/`, o `meuapp` usa a `minhalib` **local**, sem tocar no `go.mod`. 🎉

```bash
go work use ./outromodulo   # adiciona mais um módulo ao workspace
```

> ⚠️ Normalmente o `go.work` **não** vai pro Git (é uma configuração do **seu** ambiente). Adicione no `.gitignore`.

---

## 🏷️ Build constraints (compilação condicional)

Às vezes um arquivo só faz sentido em **um sistema operacional** ou **uma situação**.

### Por nome de arquivo

O Go reconhece sufixos especiais:

```
arquivo_windows.go   → só compila no Windows
arquivo_linux.go     → só no Linux
arquivo_darwin.go    → só no macOS
arquivo_amd64.go     → só em processadores amd64
arquivo_test.go      → só nos testes (go test)
```

```go
// caminho_windows.go
package config

const PastaPadrao = `C:\ProgramData\livraria`
```

```go
// caminho_linux.go
package config

const PastaPadrao = "/etc/livraria"
```

### Com a diretiva `//go:build`

Na **primeira linha** do arquivo (antes do `package`):

```go
//go:build linux || darwin

package config
```

```go
//go:build !windows

package config
```

### Tags personalizadas

```go
//go:build debug

package main

func init() {
	fmt.Println("🐞 modo debug ativado")
}
```

```bash
go build              # não inclui o arquivo acima
go build -tags debug  # inclui
```

---

## 📎 `embed`: arquivos dentro do executável (Go 1.16+)

Lembra que o Go gera um **binário único**? Com `embed`, até arquivos (HTML, JSON, imagens, SQL) vão **dentro** dele:

```go
package main

import (
	"embed"
	"fmt"
)

//go:embed versao.txt
var versao string

//go:embed config.json
var configPadrao []byte

//go:embed templates/*.html
var templates embed.FS

func main() {
	fmt.Println("Versão:", versao)
	fmt.Println(len(configPadrao), "bytes de config")

	dados, _ := templates.ReadFile("templates/index.html")
	fmt.Println(string(dados))
}
```

- `//go:embed` fica **logo acima** de uma variável de pacote
- Tipos aceitos: `string`, `[]byte` ou `embed.FS` (para vários arquivos/pastas)
- Os caminhos são **relativos ao arquivo `.go`** e não podem subir com `..`

> Perfeito para servidores web: o HTML/CSS vai junto no executável, sem precisar copiar pastas no deploy.

---

## 🧰 Ferramentas do dia a dia

| Comando | O que faz |
|---|---|
| `go fmt ./...` | formata o código (padrão oficial) |
| `goimports -w .` | formata **e** organiza os imports |
| `go vet ./...` | encontra erros comuns (Printf com argumentos errados, cópia de mutex...) |
| `staticcheck ./...` | análise estática mais completa (ferramenta externa) |
| `golangci-lint run` | roda dezenas de linters de uma vez (ferramenta externa) |
| `go doc pacote.Função` | documentação no terminal |
| `go build -o nome` | compila com um nome de saída |
| `go build -ldflags="-s -w"` | binário menor (remove informações de debug) |
| `go version -m ./binario` | mostra com qual versão do Go e dependências um binário foi feito |
| `go generate ./...` | roda comandos `//go:generate` escritos no código |
| `go clean -modcache` | limpa o cache de módulos baixados |

### Injetando a versão na compilação

```go
package main

var Versao = "dev" // valor padrão

func main() {
	fmt.Println("versão", Versao)
}
```

```bash
go build -ldflags="-X main.Versao=1.2.0"
./app # versão 1.2.0
```

### `go doc`: documentação sem sair do terminal

```bash
go doc strings.Cut
go doc -all errors
go doc github.com/google/uuid
```

E lembre-se: **comentários acima de itens exportados viram documentação** automaticamente (inclusive no site `pkg.go.dev`).

```go
// Package produto define os produtos vendidos pela livraria.
package produto

// Novo cria um produto validando o preço.
// Retorna ErrPrecoInvalido se o preço não for positivo.
func Novo(...) (*Produto, error) { ... }
```

---

## 🧾 Resumão do módulo

| Conceito | Resumo |
|---|---|
| Pacote | pasta com `.go` do mesmo `package` |
| `package main` | gera executável |
| Maiúscula/minúscula | exportado/privado entre pacotes |
| `init()` | roda antes do `main` (use pouco) |
| Módulo | conjunto de pacotes com um `go.mod` |
| `go.sum` | hashes das dependências (commitar!) |
| `go mod tidy` | sincroniza dependências |
| SemVer + `/v2` | versão maior = caminho novo |
| `internal/` | pacotes privados do módulo |
| `cmd/` | um `main` por executável |
| `go.work` | vários módulos locais juntos |
| `//go:build` | compilação condicional |
| `//go:embed` | arquivos dentro do binário |

---

## ✍️ Exercícios

1. Crie dois módulos locais (`calc` e `app`) e use um **workspace** para o `app` usar o `calc` local.
2. Crie os arquivos `so_windows.go` e `so_linux.go` com uma constante `NomeSO` e imprima no `main`.
3. Crie um arquivo com `//go:build debug` que imprime uma mensagem. Compile com e sem `-tags debug`.
4. Use `//go:embed` para incluir um arquivo `ajuda.txt` no executável e mostrá-lo quando o usuário digitar `ajuda`.
5. Injete a versão do programa com `-ldflags "-X main.Versao=..."`.
6. Escreva comentários de documentação no seu pacote `produto` e veja o resultado com `go doc`.
7. Rode `go vet ./...` no seu projeto e corrija o que aparecer.

---

⬅️ Anterior: [Organizando projetos](03-organizando-projetos.md) · ➡️ Próximo módulo: [Generics](../13-generics/README.md)

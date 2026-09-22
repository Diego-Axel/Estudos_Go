# 02 — Módulos e `go.mod`

## 📚 Pacote vs. módulo

| | Pacote | Módulo |
|---|---|---|
| O que é | uma **pasta** com arquivos `.go` | uma **coleção de pacotes** versionada junta |
| Definido por | `package nome` | o arquivo **`go.mod`** na raiz |
| Analogia | um capítulo | o livro inteiro |

Normalmente, **1 repositório = 1 módulo**.

```
meuprojeto/          ← módulo (tem go.mod)
├── go.mod
├── main.go          ← pacote main
├── estoque/         ← pacote estoque
│   └── estoque.go
└── relatorio/       ← pacote relatorio
    └── relatorio.go
```

---

## 🆕 Criando um módulo

```bash
mkdir livraria && cd livraria
go mod init github.com/Diego-Axel/livraria
```

O argumento é o **caminho do módulo** (*module path*). Por convenção, é o endereço do repositório, pois é assim que os outros vão importar seu código:

```go
import "github.com/Diego-Axel/livraria/estoque"
```

> 💡 Para projetos locais que nunca serão publicados, pode ser qualquer nome: `go mod init livraria`.

---

## 📄 Anatomia do `go.mod`

```go
module github.com/Diego-Axel/livraria

go 1.25

require (
	github.com/google/uuid v1.6.0
	golang.org/x/text v0.21.0
)

require (
	golang.org/x/sys v0.28.0 // indirect
)
```

| Linha | Significado |
|---|---|
| `module` | o caminho do módulo |
| `go 1.25` | versão **mínima** do Go exigida (e quais recursos da linguagem valem) |
| `toolchain` | (opcional) versão do Go preferida para compilar |
| `require` | dependências e suas versões |
| `// indirect` | dependência **da dependência** (você não importa direto) |
| `replace` | troca uma dependência por outra (ex: uma pasta local) |
| `exclude` | proíbe uma versão específica |
| `tool` | ferramentas usadas no projeto (Go 1.24+) |

---

## 🔐 O `go.sum`

Junto do `go.mod` aparece o **`go.sum`**, com os **hashes criptográficos** de cada dependência (os valores abaixo são só ilustrativos):

```
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

Ele garante que **todo mundo** baixe exatamente o **mesmo código**, e que ninguém adulterou a dependência.

> ✅ **Sempre** faça commit do `go.mod` **e** do `go.sum`. Nunca edite o `go.sum` na mão.

---

## ➕ Adicionando dependências

### Jeito 1: importe e rode `go mod tidy` ⭐

```go
package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(uuid.NewString()) // ex: 3f8a2b1c-...
}
```

```bash
go mod tidy
```

O `tidy` **adiciona** o que está faltando e **remove** o que não é mais usado.

### Jeito 2: `go get`

```bash
go get github.com/google/uuid           # última versão
go get github.com/google/uuid@v1.5.0    # versão específica
go get github.com/google/uuid@latest    # atualiza para a mais recente
go get -u ./...                         # atualiza TODAS as dependências (minor/patch)
go get github.com/google/uuid@none      # remove a dependência
```

---

## 🔢 Versionamento semântico (SemVer)

Toda versão de módulo segue o formato **`vMAIOR.MENOR.CORREÇÃO`**:

```
v1.6.2
│ │ └── CORREÇÃO (patch): correções de bugs, compatível
│ └──── MENOR (minor): novos recursos, compatível
└────── MAIOR (major): mudanças que QUEBRAM compatibilidade
```

- `v0.x.x`: em desenvolvimento, **qualquer coisa** pode mudar
- `v1.0.0` em diante: promessa de **estabilidade**

### A regra das versões maiores (v2+)

Em Go, uma versão maior **nova** é tratada como **outro módulo**, e o caminho ganha o sufixo `/vN`:

```go
import "github.com/exemplo/lib"    // v0 ou v1
import "github.com/exemplo/lib/v2" // v2 (caminho diferente!)
```

Isso permite usar `v1` e `v2` **ao mesmo tempo** no mesmo projeto. É por isso que existe `math/rand/v2` na biblioteca padrão.

### Seleção de versão mínima (MVS)

Se duas dependências pedem versões diferentes de uma terceira (`v1.2.0` e `v1.4.0`), o Go escolhe a **menor versão que satisfaz todas**, ou seja, `v1.4.0`. Nada de "pegar a mais nova que existir": os builds são **reprodutíveis**.

---

## 🔧 Comandos úteis de módulo

| Comando | O que faz |
|---|---|
| `go mod init caminho` | cria o `go.mod` |
| `go mod tidy` | sincroniza dependências com o código ⭐ |
| `go get pkg@versão` | adiciona/atualiza/remove dependência |
| `go list -m all` | lista todos os módulos usados |
| `go list -m -u all` | mostra quais têm atualização disponível |
| `go mod why pkg` | explica **por que** um pacote é necessário |
| `go mod graph` | mostra o grafo de dependências |
| `go mod download` | baixa as dependências (útil em Docker/CI) |
| `go mod vendor` | copia as dependências para a pasta `vendor/` |
| `go mod verify` | confere se os downloads batem com o `go.sum` |

---

## 🔁 `replace`: usando uma cópia local

Quando você está alterando uma biblioteca e quer testá-la no seu projeto **sem publicar**:

```go
// go.mod
require github.com/Diego-Axel/minhalib v1.0.0

replace github.com/Diego-Axel/minhalib => ../minhalib
```

> ⚠️ Lembre de **remover** o `replace` antes de publicar. Para esse caso, **workspaces** (arquivo 04) costumam ser uma opção melhor.

---

## 🌐 De onde vêm os módulos?

Por padrão, o Go baixa pelo **proxy** oficial (`proxy.golang.org`) e confere os hashes num banco público (`sum.golang.org`). Tudo vai para um cache local (`go env GOMODCACHE`).

Para repositórios **privados** da sua empresa:

```bash
go env -w GOPRIVATE=github.com/minhaempresa/*
```

---

## 🛠️ Instalando ferramentas: `go install`

Compila e instala um **executável** em `$GOPATH/bin` (lembre de colocar essa pasta no `PATH`):

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```

### Ferramentas do projeto (Go 1.24+)

Para registrar uma ferramenta **no `go.mod`**, com versão fixa para todo o time:

```bash
go get -tool golang.org/x/tools/cmd/stringer
go tool stringer -type=Status   # executa a ferramenta
```

---

## ✍️ Exercícios

1. Crie um módulo `github.com/SEU-USUARIO/ola`, adicione `github.com/google/uuid` e imprima 3 UUIDs.
2. Abra o `go.mod` e o `go.sum` gerados e identifique cada parte.
3. Remova o uso do `uuid` do código, rode `go mod tidy` e veja o que mudou nos arquivos.
4. Rode `go list -m all` e `go mod graph` num projeto com dependências.
5. Explique por que `github.com/exemplo/lib` e `github.com/exemplo/lib/v2` podem ser usados juntos.
6. **Projeto:** este repositório de estudos tem um `main.go` na raiz, mas **não tem `go.mod`**. Rode `go mod init github.com/Diego-Axel/Estudos_Go` e confira que `go run .` funciona.

---

⬅️ Anterior: [Pacotes](01-pacotes.md) · ➡️ Próximo: [Organizando projetos](03-organizando-projetos.md)

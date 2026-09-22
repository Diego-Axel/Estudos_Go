# 02 — Ferramentas de Qualidade

Um dos maiores trunfos de Go é o **ecossistema de ferramentas**. Várias vêm junto com a linguagem; outras são instaladas com um `go install`.

---

## 🎨 Formatação: `gofmt` e `goimports`

```bash
gofmt -l .      # lista arquivos fora do padrão
gofmt -w .      # formata e salva
go fmt ./...    # atalho (roda gofmt -w em todos os pacotes)
```

O `goimports` faz o mesmo **e** organiza os imports (adiciona os que faltam, remove os não usados):

```bash
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

> ✅ Configure o editor para **formatar ao salvar**. No VS Code, com a extensão Go, isso já vem ligado (usando o `gopls`).

---

## 🔍 Análise estática

### `go vet` (embutido)

Encontra erros que **compilam** mas quase certamente são bugs:

```bash
go vet ./...
```

Exemplos do que ele pega:

```go
fmt.Printf("%d\n", "texto")      // verbo errado para o tipo
fmt.Println("valor: %d", x)      // Println com verbo de formatação
var mu sync.Mutex; m2 := mu      // cópia de mutex
if x == x { }                    // comparação sempre verdadeira
defer resp.Body.Close()          // antes de checar o erro do http.Get
```

> 🧠 O `go test` já roda um subconjunto do `go vet` automaticamente.

### `staticcheck`

Análise mais profunda: código morto, uso de funções obsoletas, simplificações, bugs sutis:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

```
main.go:15:2: this value of err is never used (SA4006)
main.go:22:5: should omit comparison to bool constant, can be simplified to !ativo (S1002)
main.go:30:9: strings.Title has been deprecated (SA1019)
```

### `golangci-lint`: dezenas de linters de uma vez ⭐

O padrão em projetos profissionais. Roda `govet`, `staticcheck`, `errcheck` e muitos outros, em paralelo:

```bash
# instalação: veja https://golangci-lint.run (há instaladores para cada sistema)
golangci-lint run
```

Configuração no arquivo `.golangci.yml` na raiz do projeto:

```yaml
version: "2"

linters:
  enable:
    - errcheck      # erros não verificados
    - govet
    - staticcheck
    - unused        # código não usado
    - ineffassign   # atribuições inúteis
    - misspell      # erros de digitação em comentários
    - gocritic      # várias sugestões de estilo e bugs
    - revive        # estilo (substituto do antigo golint)
```

---

## 🛡️ Segurança: `govulncheck`

Verifica se seu código usa dependências (ou a própria versão do Go) com **vulnerabilidades conhecidas**, e só avisa se a função vulnerável for **realmente chamada** pelo seu código:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

```
Vulnerability #1: GO-2024-XXXX
    Uma falha em golang.org/x/net/html permite ...
  Found in: golang.org/x/net@v0.20.0
  Fixed in: golang.org/x/net@v0.23.0
  Example traces found:
    #1: main.go:42:15: main.processar calls html.Parse
```

> ✅ Rode no CI e mantenha as dependências atualizadas (`go get -u ./...` + `go mod tidy` + testes).

---

## 🧠 `gopls`: o "cérebro" do editor

O `gopls` é o servidor de linguagem oficial. Ele dá ao VS Code (e outros editores):

- autocompletar e documentação ao passar o mouse
- ir para definição / encontrar referências
- renomear com segurança em todo o projeto
- avisos em tempo real (inclui análises do `go vet`)
- *code actions*: preencher struct, extrair função, organizar imports...

> A extensão **Go** do VS Code instala e configura o `gopls` automaticamente.

---

## 🔧 Automatizando com `Makefile` (ou Taskfile)

Evite decorar comandos: centralize-os.

```makefile
.PHONY: fmt lint test cover run build

fmt:
	gofmt -w .

lint:
	go vet ./...
	golangci-lint run

test:
	go test -race -count=1 ./...

cover:
	go test -coverprofile=cobertura.out ./...
	go tool cover -html=cobertura.out

run:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api ./cmd/api
```

```bash
make test
make lint
```

> No Windows, o `make` não vem instalado. Alternativas: [Task](https://taskfile.dev) (`Taskfile.yml`), scripts `.ps1`, ou instalar o `make` via Chocolatey/Scoop.

---

## 🤖 Integração contínua (CI) completa

Juntando tudo num workflow do GitHub Actions (`.github/workflows/ci.yml`):

```yaml
name: CI
on: [push, pull_request]

jobs:
  qualidade:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Formatação
        run: test -z "$(gofmt -l .)"

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v8

      - name: Testes
        run: go test -race -coverprofile=cobertura.out ./...

      - name: Vulnerabilidades
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

> A cada push, o GitHub verifica formatação, lint, testes e segurança automaticamente. ✅ ou ❌ no commit.

---

## 📋 Resumo das ferramentas

| Ferramenta | Para quê | Vem com o Go? |
|---|---|---|
| `gofmt` / `go fmt` | formatação | ✅ |
| `goimports` | formatação + imports | ❌ (`go install`) |
| `go vet` | bugs comuns | ✅ |
| `staticcheck` | análise profunda | ❌ |
| `golangci-lint` | vários linters juntos | ❌ |
| `govulncheck` | vulnerabilidades | ❌ |
| `gopls` | inteligência do editor | ❌ (a extensão instala) |
| `go test -race` | condições de corrida | ✅ |
| `go test -cover` | cobertura | ✅ |
| `go doc` | documentação | ✅ |

---

## ✍️ Exercícios

1. Escreva um código com `fmt.Printf("%d", "texto")` e rode `go vet`. Depois corrija.
2. Instale o `staticcheck` e rode num projeto seu. Corrija os avisos.
3. Instale o `golangci-lint`, crie o `.golangci.yml` acima e rode no projeto da livraria.
4. Rode o `govulncheck` num projeto com dependências.
5. Crie um `Makefile` (ou `Taskfile.yml`) com as tarefas `fmt`, `lint`, `test` e `run`.
6. **Projeto:** adicione o workflow de CI completo a este repositório (depois de criar o `go.mod`, módulo 12).

---

⬅️ Anterior: [Go idiomático](01-go-idiomatico.md) · ➡️ Próximo: [Debug e profiling](03-debug-e-profiling.md)

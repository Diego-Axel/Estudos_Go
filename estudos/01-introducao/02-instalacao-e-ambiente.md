# 02 — Instalação e Ambiente

## 📥 Instalando o Go

1. Acesse **https://go.dev/dl/**
2. Baixe o instalador do seu sistema:
   - **Windows:** arquivo `.msi` → só ir clicando em "Next"
   - **Linux:** arquivo `.tar.gz` (ou pelo gerenciador de pacotes)
   - **macOS:** arquivo `.pkg` (ou `brew install go`)
3. Abra um **novo** terminal e confira:

```bash
go version
```

Saída esperada (a versão pode ser diferente):

```
go version go1.25.0 windows/amd64
```

---

## 🧰 Editor recomendado: VS Code

1. Instale a extensão **Go** (publicada pela *Go Team at Google*).
2. Abra um arquivo `.go` → o VS Code vai sugerir instalar as ferramentas (`gopls`, `dlv`, etc.). Clique em **Install All**.

Com isso você ganha: autocompletar, formatação automática ao salvar, ir para definição, debug...

> Outras opções: **GoLand** (JetBrains, pago), **Neovim** com `gopls`.

---

## 🌍 Variáveis de ambiente importantes

Veja todas com:

```bash
go env
```

| Variável | Para que serve |
|---|---|
| `GOROOT` | Onde o Go está instalado. **Não mexa**, o instalador cuida disso. |
| `GOPATH` | Pasta onde ficam pacotes baixados e binários instalados (`go install`). Padrão: `~/go` (no Windows `%USERPROFILE%\go`). |
| `GOOS` / `GOARCH` | Sistema e arquitetura alvo da compilação. |
| `GOPROXY` | De onde os módulos são baixados. |

> 📌 Hoje em dia, com **Go Modules**, você pode criar seu projeto **em qualquer pasta**. Não precisa ficar dentro do `GOPATH`.

---

## 🛠️ Comandos essenciais do `go`

| Comando | O que faz |
|---|---|
| `go run arquivo.go` | Compila e executa na hora (não deixa binário) |
| `go build` | Compila e gera o executável |
| `go mod init nome` | Cria um novo módulo (projeto) → gera `go.mod` |
| `go mod tidy` | Adiciona dependências que faltam e remove as não usadas |
| `go get pacote` | Adiciona/atualiza uma dependência |
| `go install pacote@latest` | Compila e instala um binário no `GOPATH/bin` |
| `go fmt ./...` | Formata todo o código no padrão oficial |
| `go vet ./...` | Procura erros suspeitos no código |
| `go test ./...` | Roda os testes |
| `go doc fmt.Println` | Mostra a documentação de algo no terminal |

> O `./...` significa "este diretório e todos os subdiretórios".

---

## 🔀 Compilação cruzada (bônus)

Você pode gerar executável para **outro sistema** sem sair do seu:

```bash
# No PowerShell (Windows) -> gerando binário para Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o app

# No Bash (Linux/macOS/Git Bash) -> gerando .exe para Windows
GOOS=windows GOARCH=amd64 go build -o app.exe
```

Ver todas as combinações possíveis:

```bash
go tool dist list
```

---

## ✍️ Exercícios

1. Instale o Go e rode `go version`. Anote a versão.
2. Rode `go env GOPATH` e descubra onde fica seu `GOPATH`.
3. Qual a diferença entre `go run` e `go build`?
4. Para que serve o `go mod tidy`?

---

⬅️ Anterior: [O que é Go?](01-o-que-e-go.md) · ➡️ Próximo: [Primeiro programa](03-primeiro-programa.md)

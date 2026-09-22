# 04 — Armadilhas Comuns, Deploy e Próximos Passos

## 🪤 As armadilhas mais comuns de Go (revisão)

Uma coletânea das pegadinhas vistas ao longo dos módulos. Se você entende **todas** estas, já está acima da média! 💪

| # | Armadilha | Onde vimos | Resumo |
|---|---|---|---|
| 1 | `len("ação")` é 6, não 4 | [Módulo 07](../07-strings-runes-bytes/01-utf8-bytes-e-runes.md) | `len` conta **bytes**; use `range` ou `[]rune` |
| 2 | `string(65)` é `"A"` | [Módulo 02](../02-fundamentos/04-conversao-de-tipos.md) | use `strconv.Itoa` |
| 3 | Divisão inteira `7/2 == 3` | [Módulo 03](../03-operadores/01-aritmeticos.md) | converta para `float64` |
| 4 | Imprecisão de float | [Módulo 02](../02-fundamentos/02-tipos-de-dados.md) | dinheiro em centavos (`int64`) |
| 5 | Shadowing com `:=` | [Módulo 02](../02-fundamentos/01-variaveis.md) | `:=` num bloco interno cria **outra** variável |
| 6 | `break` dentro de `switch` num `for` | [Módulo 04](../04-controle-de-fluxo/04-break-continue-labels-goto.md) | use label ou `return` |
| 7 | Valor do `range` é cópia | [Módulo 04](../04-controle-de-fluxo/02-for.md) | altere com `s[i]` |
| 8 | `defer` em laço | [Módulo 05](../05-funcoes/05-defer.md) | só roda no fim da **função** |
| 9 | Argumentos do `defer` avaliados na hora | [Módulo 05](../05-funcoes/05-defer.md) | use closure para o valor final |
| 10 | `append` que sobrescreve o original | [Módulo 06](../06-arrays-slices-maps/03-slices-avancado.md) | slices compartilham o array; use `s[a:b:c]` ou copie |
| 11 | Escrever em map `nil` | [Módulo 06](../06-arrays-slices-maps/04-maps.md) | inicialize com `make` |
| 12 | Ordem de map é aleatória | [Módulo 06](../06-arrays-slices-maps/04-maps.md) | ordene as chaves |
| 13 | Alterar campo de struct dentro de map | [Módulo 06](../06-arrays-slices-maps/04-maps.md) | pegue, altere, guarde; ou use ponteiros |
| 14 | Desreferenciar ponteiro `nil` | [Módulo 08](../08-ponteiros/01-o-que-sao-ponteiros.md) | cheque `!= nil` |
| 15 | Método com receptor valor não altera | [Módulo 09](../09-structs-e-metodos/02-metodos.md) | use `*T` |
| 16 | `String()` com recursão infinita | [Módulo 09](../09-structs-e-metodos/02-metodos.md) | não use `%v` do próprio receptor |
| 17 | Receptor `*T` e interfaces | [Módulo 10](../10-interfaces/01-o-que-sao-interfaces.md) | só `*T` satisfaz a interface |
| 18 | Interface "nil" que não é nil | [Módulo 10](../10-interfaces/04-nil-e-boas-praticas.md) | retorne `nil` literal |
| 19 | `==` em vez de `errors.Is` | [Módulo 11](../11-tratamento-de-erros/03-wrapping-is-as-join.md) | erros embrulhados não batem com `==` |
| 20 | `os.Exit`/`log.Fatal` pulam os `defer` | [Módulo 11](../11-tratamento-de-erros/04-panic-recover-e-boas-praticas.md) | padrão `main → run()` |
| 21 | `main` termina antes das goroutines | [Módulo 14](../14-concorrencia/01-goroutines.md) | `WaitGroup` |
| 22 | Race condition | [Módulo 14](../14-concorrencia/04-sync-e-race-conditions.md) | mutex/atomic/channels + `-race` |
| 23 | Goroutine leak | [Módulo 14](../14-concorrencia/05-padroes-de-concorrencia.md) | toda goroutine precisa de um fim |
| 24 | Esquecer `Flush()` | [Módulo 16](../16-biblioteca-padrao/01-arquivos-e-io.md) | `bufio.Writer` e `csv.Writer` |
| 25 | Números JSON viram `float64` em `any` | [Módulo 16](../16-biblioteca-padrao/02-json-e-csv.md) | use structs |
| 26 | `http.Get` sem timeout | [Módulo 17](../17-web-net-http/04-cliente-http-e-templates.md) | `http.Client{Timeout: ...}` |
| 27 | Esquecer `resp.Body.Close()` | [Módulo 17](../17-web-net-http/04-cliente-http-e-templates.md) | `defer resp.Body.Close()` |
| 28 | `err == nil` num 404 | [Módulo 17](../17-web-net-http/04-cliente-http-e-templates.md) | verifique `resp.StatusCode` |

---

## 🚀 Deploy: colocando seu programa no mundo

### Compilando para produção

```bash
# Binário menor, sem informações de debug
go build -ldflags="-s -w" -o livraria ./cmd/api

# Binário 100% estático (sem depender de bibliotecas C do sistema)
CGO_ENABLED=0 go build -o livraria ./cmd/api
```

### Compilação cruzada (módulo 01)

```bash
# Do Windows (PowerShell) para um servidor Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o livraria ./cmd/api
```

Depois é só copiar **um único arquivo** para o servidor e rodar. Sem instalar Go, sem dependências. 🎉

### 🐳 Docker com build em múltiplos estágios

`Dockerfile`:

```dockerfile
# Estágio 1: compilar
FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/api

# Estágio 2: imagem final mínima (só o binário!)
FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
```

```bash
docker build -t livraria .
docker run -p 8080:8080 livraria
```

A imagem final fica com **poucos MB** (em vez de ~800 MB com o Go inteiro), e sem shell nem ferramentas que um invasor poderia usar.

> 💡 Se o projeto não tiver dependências externas, ainda não haverá `go.sum`: ajuste o `COPY` para copiar só o `go.mod`.

### Configuração por variáveis de ambiente

Seguindo o [12-Factor App](https://12factor.net/pt_br/), leia configurações do ambiente, não de valores fixos no código:

```go
func obterEnv(chave, padrao string) string {
	if v, ok := os.LookupEnv(chave); ok {
		return v
	}
	return padrao
}

porta := obterEnv("PORTA", "8080")
```

### Onde hospedar

- **Fly.io**, **Render**, **Railway**: deploy simples a partir do Dockerfile ou do repositório
- **Google Cloud Run**, **AWS App Runner / ECS**, **Azure Container Apps**: contêineres gerenciados
- **Uma VPS** (Hetzner, DigitalOcean...): copie o binário e rode como serviço do `systemd`

---

## 🧭 Próximos passos

Você terminou o roteiro! 🎉 O que vem agora?

### 📚 Para aprofundar

| Recurso | Tipo |
|---|---|
| [A Tour of Go](https://go.dev/tour/) | tutorial interativo oficial |
| [Go by Example](https://gobyexample.com/) | exemplos curtos |
| [Effective Go](https://go.dev/doc/effective_go) | estilo idiomático |
| [Go Blog](https://go.dev/blog/) | artigos oficiais (leia os de erros, slices, concorrência) |
| *The Go Programming Language* (Donovan & Kernighan) | o livro clássico |
| *Learning Go* (Jon Bodner) | livro moderno, com generics |
| *100 Go Mistakes and How to Avoid Them* (Teiva Harsanyi) | armadilhas em profundidade |
| *Concurrency in Go* (Katherine Cox-Buday) | concorrência a fundo |
| [Exercism: trilha de Go](https://exercism.org/tracks/go) | exercícios com mentoria gratuita |

### 🛠️ Bibliotecas para conhecer (quando precisar)

| Área | Bibliotecas |
|---|---|
| Banco de dados | `database/sql` + `pgx` (PostgreSQL), `sqlc` (gera código a partir de SQL), `GORM` (ORM) |
| Web | `chi`, `echo`, `gin` (routers/frameworks) |
| CLI | `cobra`, `urfave/cli`, `bubbletea` (interfaces no terminal) |
| Configuração | `envconfig`, `viper` |
| Testes | `go-cmp`, `testify` |
| gRPC / Protobuf | `google.golang.org/grpc`, `connect-go` |
| Observabilidade | `OpenTelemetry`, `prometheus/client_golang` |

> 🧠 Lembre do provérbio: *"A little copying is better than a little dependency."* Use a biblioteca padrão até ela não bastar.

### 💡 Ideias de projetos para praticar

1. **Livraria completa**: API REST + PostgreSQL + Docker + testes + CI (evolução do seu `main.go`)
2. **Encurtador de URLs**: API, banco, redirecionamento, estatísticas de acesso
3. **CLI de tarefas (to-do)**: com `flag` ou `cobra`, salvando em JSON
4. **Web scraper concorrente**: worker pool baixando e analisando páginas
5. **Chat em tempo real**: WebSockets + goroutines + channels
6. **Monitor de sites**: verifica URLs periodicamente e avisa quando caem
7. **Contribua com open source**: procure issues marcadas como `good first issue` em projetos Go no GitHub

---

## 🏁 Checklist final: você sabe...

- [ ] explicar a diferença entre array, slice e map, e as pegadinhas de cada um?
- [ ] quando usar receptor por valor ou por ponteiro?
- [ ] como interfaces são implementadas implicitamente, e a pegadinha do nil?
- [ ] tratar erros com `%w`, `errors.Is` e `errors.As`?
- [ ] organizar um projeto com módulos, `cmd/` e `internal/`?
- [ ] escrever funções e tipos genéricos?
- [ ] usar goroutines, channels, `select`, `context` e mutex sem race conditions?
- [ ] escrever table tests, benchmarks e usar `-race` e `-cover`?
- [ ] ler/escrever arquivos, JSON e trabalhar com datas?
- [ ] construir uma API HTTP com middlewares e graceful shutdown?
- [ ] usar `go vet`, `golangci-lint`, o debugger e o `pprof`?

Se marcou tudo: **parabéns, você é um(a) Gopher!** 🐹🎉

---

## ✍️ Exercícios finais

1. Revise a tabela de armadilhas e, para cada uma, escreva **do zero** um exemplo curto que a reproduza e a correção.
2. Crie o `Dockerfile` para a API da livraria e rode com Docker.
3. Compile a API para Linux a partir do Windows e confira o tamanho do binário com e sem `-ldflags="-s -w"`.
4. Escolha um dos projetos sugeridos e construa aplicando **tudo** do roteiro: módulos, testes, erros bem tratados, concorrência (se fizer sentido), lint e CI.
5. Escreva no `README.md` deste repositório o que você aprendeu e quais projetos construiu. 😉

---

⬅️ Anterior: [Debug e profiling](03-debug-e-profiling.md) · 🏠 [Voltar ao roteiro](../README.md)

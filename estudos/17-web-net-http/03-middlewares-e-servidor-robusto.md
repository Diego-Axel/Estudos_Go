# 03 — Middlewares e Servidor Robusto

## 🧅 O que é um middleware?

Um **middleware** é uma função que **envolve** um handler, executando código **antes** e/ou **depois** dele. É como as camadas de uma cebola:

```
requisição ──► [ log ─► [ recover ─► [ auth ─► HANDLER ] ] ] ──► resposta
```

Em Go, um middleware tem esta cara:

```go
func Middleware(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ... antes ...
		proximo.ServeHTTP(w, r) // chama o próximo da corrente
		// ... depois ...
	})
}
```

> É o padrão de **funções que recebem e retornam funções** (módulo 05) + a interface **`http.Handler`** (módulo 10).

---

## 📋 Middleware de log

```go
func Logger(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		proximo.ServeHTTP(w, r)
		slog.Info("requisição",
			"metodo", r.Method,
			"caminho", r.URL.Path,
			"duracao", time.Since(inicio),
		)
	})
}
```

### Capturando o status da resposta

O `ResponseWriter` não informa qual status foi enviado. A solução: **embutir** (módulo 09) e interceptar o `WriteHeader`:

```go
type gravadorDeStatus struct {
	http.ResponseWriter // embutido: todos os métodos são promovidos
	status int
}

func (g *gravadorDeStatus) WriteHeader(codigo int) {
	g.status = codigo
	g.ResponseWriter.WriteHeader(codigo)
}

func Logger(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &gravadorDeStatus{ResponseWriter: w, status: http.StatusOK}
		proximo.ServeHTTP(g, r)
		slog.Info("requisição",
			"metodo", r.Method,
			"caminho", r.URL.Path,
			"status", g.status,
			"duracao", time.Since(inicio),
		)
	})
}
```

---

## 🛟 Middleware de recuperação (recover)

Um panic num handler não deve derrubar o servidor inteiro (módulo 11):

```go
func Recuperar(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recuperado", "erro", rec, "caminho", r.URL.Path)
				http.Error(w, "erro interno do servidor", http.StatusInternalServerError)
			}
		}()
		proximo.ServeHTTP(w, r)
	})
}
```

---

## 🔐 Middleware de autenticação

```go
func ExigirToken(tokenValido string) func(http.Handler) http.Handler {
	return func(proximo http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token != tokenValido {
				http.Error(w, "não autorizado", http.StatusUnauthorized)
				return // ⚠️ NÃO chama o próximo: a requisição para aqui
			}
			proximo.ServeHTTP(w, r)
		})
	}
}
```

> Repare: é uma função que **retorna um middleware**, para poder receber configuração (o token).

### Passando dados adiante com `context`

Um middleware pode colocar informações no **contexto** da requisição (módulo 14):

```go
type chaveContexto string

const chaveUsuario chaveContexto = "usuario"

func Autenticar(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usuario := "diego" // (na vida real: validar um token JWT, sessão...)
		ctx := context.WithValue(r.Context(), chaveUsuario, usuario)
		proximo.ServeHTTP(w, r.WithContext(ctx))
	})
}

func perfil(w http.ResponseWriter, r *http.Request) {
	usuario, _ := r.Context().Value(chaveUsuario).(string)
	fmt.Fprintf(w, "Olá, %s!\n", usuario)
}
```

> ✅ Use um **tipo próprio** para a chave (`chaveContexto`), nunca uma `string` solta. Isso evita colisões com outros pacotes.

---

## 🌍 Middleware de CORS

Para uma API ser chamada por um front-end em **outro domínio** (ex: React em `localhost:3000`):

```go
func CORS(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions { // requisição "preflight" do navegador
			w.WriteHeader(http.StatusNoContent)
			return
		}
		proximo.ServeHTTP(w, r)
	})
}
```

---

## 🔗 Encadeando middlewares

```go
handler := Recuperar(Logger(CORS(api.Rotas())))
```

Fica ilegível com muitos. Uma função auxiliar ajuda:

```go
func Encadear(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- { // de trás pra frente
		h = middlewares[i](h)
	}
	return h
}

handler := Encadear(api.Rotas(), Recuperar, Logger, CORS)
// executa na ordem: Recuperar → Logger → CORS → rotas
```

Middleware só em **algumas** rotas:

```go
mux.Handle("GET /admin/relatorio", ExigirToken("segredo")(http.HandlerFunc(relatorio)))
```

---

## 🏗️ `http.Server` com timeouts ⭐

O `http.ListenAndServe` usa configurações **sem timeout**: um cliente lento (ou mal-intencionado) pode prender conexões para sempre. Em produção, configure um `http.Server`:

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           handler,
	ReadHeaderTimeout: 5 * time.Second,  // tempo para ler os cabeçalhos
	ReadTimeout:       10 * time.Second, // tempo para ler a requisição toda
	WriteTimeout:      10 * time.Second, // tempo para escrever a resposta
	IdleTimeout:       60 * time.Second, // conexões keep-alive ociosas
}
log.Fatal(srv.ListenAndServe())
```

---

## 🛑 Graceful shutdown (desligamento elegante)

Quando o servidor recebe `Ctrl+C` (ou um `SIGTERM` do Docker/Kubernetes), ele deve:
1. **parar de aceitar** novas conexões
2. **terminar** as requisições em andamento
3. só então encerrar

```go
func main() {
	api := &API{loja: NovaLoja()}
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           Encadear(api.Rotas(), Recuperar, Logger),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Contexto cancelado quando chegar Ctrl+C (SIGINT) ou SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Servidor numa goroutine
	go func() {
		slog.Info("servidor iniciado", "endereco", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("erro no servidor", "erro", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done() // espera o sinal
	slog.Info("desligando... aguardando requisições em andamento")

	ctxDesligar, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxDesligar); err != nil {
		slog.Error("desligamento forçado", "erro", err)
	}
	slog.Info("servidor encerrado ✅")
}
```

> Junta **goroutines**, **context**, **sinais do sistema** e **tratamento de erros**: vários módulos trabalhando juntos. 🎉

---

## ⏳ Respeitando o cancelamento da requisição

Se o cliente **desistir** (fechou o navegador, timeout), `r.Context()` é cancelado. Operações demoradas devem respeitar isso:

```go
func relatorioLento(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(5 * time.Second): // simula trabalho pesado
		fmt.Fprintln(w, "relatório pronto")
	case <-r.Context().Done():
		slog.Warn("cliente desistiu", "erro", r.Context().Err())
		return // não gasta mais recursos
	}
}
```

> Ao chamar bancos de dados ou outras APIs, **passe `r.Context()` adiante**: `db.QueryContext(r.Context(), ...)`.

---

## ✍️ Exercícios

1. Adicione os middlewares `Logger` (com status) e `Recuperar` à API da livraria.
2. Crie uma rota `/panico` que dá `panic` de propósito e confirme que o servidor continua no ar.
3. Proteja as rotas `POST`, `PUT` e `DELETE` com `ExigirToken`, deixando os `GET` públicos.
4. Crie um middleware que adiciona um cabeçalho `X-Request-ID` único a cada resposta e coloca esse ID no contexto e nos logs.
5. Crie um middleware de **limite de requisições** (rate limit): no máximo 10 requisições por segundo por IP (`map` + `Mutex` + `time`).
6. Implemente o **graceful shutdown** e teste: faça uma requisição a uma rota lenta e aperte `Ctrl+C` no meio.

---

⬅️ Anterior: [API REST com JSON](02-api-rest-com-json.md) · ➡️ Próximo: [Cliente HTTP e templates](04-cliente-http-e-templates.md)

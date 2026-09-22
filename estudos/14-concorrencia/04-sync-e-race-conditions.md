# 04 — O Pacote `sync` e Race Conditions

## 🏁 O que é uma race condition?

Uma **condição de corrida** acontece quando **duas ou mais goroutines acessam a mesma variável** ao mesmo tempo, e **pelo menos uma escreve**.

```go
contador := 0
var wg sync.WaitGroup

for range 1000 {
	wg.Go(func() {
		contador++ // ⚠️ leitura + soma + escrita: NÃO é uma operação única!
	})
}
wg.Wait()
fmt.Println(contador) // 1000? 972? 🎲
```

Por que dá errado? `contador++` são **3 passos**:

```
goroutine A: lê 5          goroutine B: lê 5
goroutine A: soma → 6      goroutine B: soma → 6
goroutine A: escreve 6     goroutine B: escreve 6   ← um incremento se perdeu!
```

> ⚠️ Race conditions são **traiçoeiras**: o programa pode funcionar 99 vezes e falhar na 100ª, em produção, às 3 da manhã.

---

## 🔎 O detector de corridas: `-race` ⭐

Go vem com um detector embutido:

```bash
go run -race main.go
go test -race ./...
go build -race
```

Saída (resumida):

```
==================
WARNING: DATA RACE
Read at 0x00c000014108 by goroutine 8:
  main.main.func1()
      /caminho/main.go:12 +0x...

Previous write at 0x00c000014108 by goroutine 7:
  main.main.func1()
      /caminho/main.go:12 +0x...
==================
```

> ✅ **Rode seus testes com `-race` sempre**, especialmente no CI. (No Windows, o `-race` precisa do CGO habilitado e de um compilador C instalado.)

---

## 🔒 `sync.Mutex`: exclusão mútua

Um **mutex** é uma "tranca": só **uma** goroutine por vez pode estar entre `Lock()` e `Unlock()`.

```go
var (
	mu       sync.Mutex
	contador int
)

var wg sync.WaitGroup
for range 1000 {
	wg.Go(func() {
		mu.Lock()
		contador++ // ✅ só uma goroutine por vez aqui
		mu.Unlock()
	})
}
wg.Wait()
fmt.Println(contador) // 1000 (sempre!)
```

### O padrão: mutex dentro da struct

```go
type Contador struct {
	mu    sync.Mutex
	valor map[string]int
}

func NovoContador() *Contador {
	return &Contador{valor: make(map[string]int)}
}

func (c *Contador) Incrementar(chave string) {
	c.mu.Lock()
	defer c.mu.Unlock() // ✅ garante o Unlock mesmo com return/panic
	c.valor[chave]++
}

func (c *Contador) Valor(chave string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.valor[chave] // até LEITURAS precisam da trava!
}
```

Regras de ouro:
- Use `defer mu.Unlock()` logo depois do `Lock()`
- **Nunca copie** uma struct com mutex: use **ponteiro** (`*Contador`) nos métodos. O `go vet` avisa.
- Mantenha a região travada **pequena**: nada de I/O ou `time.Sleep` com o lock pego
- O valor zero de `sync.Mutex` já está pronto: não precisa inicializar

---

## 📖 `sync.RWMutex`: muitos leitores, um escritor

Quando há **muito mais leituras** que escritas:

```go
type Config struct {
	mu    sync.RWMutex
	dados map[string]string
}

func (c *Config) Get(chave string) string {
	c.mu.RLock()         // vários leitores ao mesmo tempo ✅
	defer c.mu.RUnlock()
	return c.dados[chave]
}

func (c *Config) Set(chave, valor string) {
	c.mu.Lock()          // escritor: exclusivo
	defer c.mu.Unlock()
	c.dados[chave] = valor
}
```

---

## ⚛️ `sync/atomic`: operações atômicas

Para contadores e flags simples, operações **atômicas** são mais leves que um mutex:

```go
import "sync/atomic"

var contador atomic.Int64

var wg sync.WaitGroup
for range 1000 {
	wg.Go(func() {
		contador.Add(1) // atômico: nunca se perde
	})
}
wg.Wait()
fmt.Println(contador.Load()) // 1000
```

Tipos disponíveis: `atomic.Int32`, `atomic.Int64`, `atomic.Uint64`, `atomic.Bool`, `atomic.Pointer[T]`, `atomic.Value`.

> Use `atomic` para **um único valor** simples. Para proteger **várias** variáveis juntas (ou um map), use mutex.

---

## 1️⃣ `sync.Once`: executar uma única vez

Garante que algo rode **uma vez só**, mesmo chamado por várias goroutines:

```go
var (
	once   sync.Once
	config map[string]string
)

func carregarConfig() map[string]string {
	once.Do(func() {
		fmt.Println("carregando config... (só uma vez)")
		config = map[string]string{"porta": "8080"}
	})
	return config
}
```

Também existem `sync.OnceFunc`, `sync.OnceValue` e `sync.OnceValues` (Go 1.21+):

```go
var obterConfig = sync.OnceValue(func() map[string]string {
	fmt.Println("carregando...")
	return map[string]string{"porta": "8080"}
})

obterConfig() // carregando...
obterConfig() // (usa o valor já carregado)
```

---

## 🗺️ `sync.Map`

Um map seguro para concorrência, **otimizado para casos específicos** (chaves escritas uma vez e lidas muitas, ou goroutines trabalhando em chaves diferentes):

```go
var m sync.Map
m.Store("go", 2009)
v, ok := m.Load("go")
fmt.Println(v, ok) // 2009 true
```

> Na maioria dos casos, `map` + `sync.Mutex` é **mais simples e igualmente rápido**. E o `sync.Map` não é tipado (usa `any`).

---

## 🆚 Channel ou mutex?

| Use **channels** para... | Use **mutex** para... |
|---|---|
| **passar a posse** de dados entre goroutines | proteger um **estado compartilhado** (cache, contador) |
| coordenar etapas (pipelines, workers) | acessos curtos e simples |
| sinalizar eventos (done, cancelamento) | quando channels deixariam o código complicado |

> 🧠 Não existe resposta única. Use o que deixar o código **mais simples e claro**.

---

## 📦 Bônus: `errgroup` (goroutines que podem falhar)

O pacote `golang.org/x/sync/errgroup` junta **WaitGroup + erro + cancelamento**:

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())

for _, url := range urls {
	g.Go(func() error {
		return baixar(ctx, url) // se uma falhar, o ctx é cancelado para as outras
	})
}

if err := g.Wait(); err != nil { // retorna o PRIMEIRO erro
	fmt.Println("falhou:", err)
}
```

---

## ✍️ Exercícios

1. Rode o contador sem proteção com `go run -race` e leia o relatório.
2. Corrija o contador de 3 formas: `sync.Mutex`, `atomic.Int64` e um channel. Compare o código.
3. Crie uma struct `ContaBancaria` segura para concorrência com `Depositar`, `Sacar` e `Saldo`. Faça 1.000 depósitos e 500 saques concorrentes e confira o saldo final.
4. Crie um cache `map[string]string` com `RWMutex` e teste com 10 leitores e 2 escritores.
5. Use `sync.Once` para inicializar uma "conexão" que é chamada por 10 goroutines. Confira que inicializou uma vez só.
6. Explique por que copiar uma struct que contém um `sync.Mutex` é um bug. Rode `go vet` num código que faz isso.

---

⬅️ Anterior: [select, timeouts e context](03-select-timeouts-e-context.md) · ➡️ Próximo: [Padrões de concorrência](05-padroes-de-concorrencia.md)

# 03 — Debug e Profiling

## 🐞 Debug com o Delve

O **Delve** (`dlv`) é o depurador oficial da comunidade Go. Com ele, você pausa o programa, inspeciona variáveis e executa linha por linha.

### No VS Code (o jeito mais fácil) ⭐

1. Clique na margem à esquerda do número da linha para criar um **breakpoint** (🔴)
2. Aperte **F5** (ou "Run and Debug")
3. O programa para no breakpoint. Agora você pode:

| Tecla | Ação |
|---|---|
| **F5** | continuar até o próximo breakpoint |
| **F10** | próxima linha (*step over*) |
| **F11** | entrar na função (*step into*) |
| **Shift+F11** | sair da função (*step out*) |
| **Shift+F5** | parar |

No painel lateral você vê **variáveis**, a **pilha de chamadas** e as **goroutines**, e pode adicionar expressões para **observar**.

Configuração em `.vscode/launch.json` (para programas com argumentos ou entrada do teclado):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Rodar livraria",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "args": ["-porta", "8080"],
      "console": "integratedTerminal"
    }
  ]
}
```

> 💡 `"console": "integratedTerminal"` é necessário para programas que **leem do teclado**, como o menu do seu `main.go`.

**Breakpoints condicionais:** clique com o botão direito no breakpoint → "Edit Condition" → ex: `codigo == "L003"`. O programa só para quando a condição for verdadeira.

### No terminal

```bash
go install github.com/go-delve/delve/cmd/dlv@latest

dlv debug .              # compila e inicia o debug
dlv test ./produto       # debuga os testes de um pacote
```

Comandos dentro do `dlv`:

```
(dlv) break main.go:25        # breakpoint na linha 25
(dlv) break produto.Novo      # breakpoint numa função
(dlv) continue                # roda até o breakpoint
(dlv) next                    # próxima linha
(dlv) step                    # entra na função
(dlv) print p                 # mostra a variável
(dlv) locals                  # todas as variáveis locais
(dlv) stack                   # pilha de chamadas
(dlv) goroutines              # lista as goroutines
(dlv) quit
```

---

## 🖨️ Debug "raiz": prints bem feitos

Às vezes um print resolve mais rápido:

```go
fmt.Printf("%+v\n", produto)  // struct com nomes dos campos
fmt.Printf("%#v\n", produto)  // sintaxe Go completa (tipos incluídos)
fmt.Printf("%T\n", valor)     // o tipo

slog.Debug("estado", "produto", produto, "estoque", estoque)
```

Imprimir a pilha de chamadas de onde você está:

```go
import "runtime/debug"

debug.PrintStack()
```

> ✅ Remova os prints de debug antes do commit (ou use `slog.Debug`, que fica desligado por padrão).

---

## 🔥 Profiling com `pprof`

Seu programa está **lento** ou usando **memória demais**? Não chute: **meça**. O `pprof` mostra **onde** o tempo e a memória estão sendo gastos.

### 1. A partir de benchmarks (o jeito mais simples)

```bash
go test -bench=. -cpuprofile=cpu.out
go test -bench=. -memprofile=mem.out

go tool pprof -http=:8081 cpu.out  # abre uma interface web interativa
```

A interface mostra:
- **Top**: as funções que mais consomem
- **Graph**: o grafo de chamadas (caixas maiores = mais tempo)
- **Flame Graph**: o "gráfico de chamas", ótimo para ver o caminho quente
- **Source**: o código fonte com o custo de cada linha

### 2. Em um servidor rodando: `net/http/pprof`

```go
import _ "net/http/pprof" // registra rotas /debug/pprof/ no DefaultServeMux

func main() {
	// Servidor separado só para diagnóstico (NUNCA exponha publicamente!)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// ... seu servidor principal ...
}
```

Com o programa rodando:

```bash
# CPU por 30 segundos
go tool pprof -http=:8081 "http://localhost:6060/debug/pprof/profile?seconds=30"

# Memória
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap

# Goroutines (ótimo para achar vazamentos!)
curl "http://localhost:6060/debug/pprof/goroutine?debug=1"
```

> ⚠️ As rotas do `pprof` expõem detalhes internos. Sirva-as só em `localhost` ou atrás de autenticação.

### 3. No próprio código

```go
import "runtime/pprof"

f, _ := os.Create("cpu.out")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()

// ... código a ser medido ...
```

---

## 📈 Trace de execução

O `pprof` mostra **onde** o tempo vai. O **trace** mostra **quando** as coisas acontecem: goroutines, bloqueios, GC, uso de cada núcleo ao longo do tempo. Ótimo para problemas de **concorrência**.

```bash
go test -trace=trace.out
go tool trace trace.out
```

---

## 🧮 Entendendo memória e GC

```go
import "runtime"

var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("heap em uso: %d KB\n", m.HeapAlloc/1024)
fmt.Printf("coletas de lixo: %d\n", m.NumGC)
fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
```

Ver as decisões do compilador (módulo 08):

```bash
go build -gcflags=-m ./...   # o que escapa para o heap
```

Acompanhar o GC em tempo real:

```bash
GODEBUG=gctrace=1 go run .
```

---

## 🧭 Roteiro para otimizar

1. **Faça funcionar** e **escreva testes** primeiro
2. **Meça** com benchmarks (`-benchmem`) e `pprof`: descubra o gargalo **real**
3. Otimize **só** o gargalo (normalmente 1 ou 2 funções respondem pela maior parte do tempo)
4. **Meça de novo** (`benchstat`) para provar que melhorou
5. Garanta que os **testes continuam passando**

Otimizações comuns que o `pprof` costuma apontar:
- pré-alocar slices e maps (`make([]T, 0, n)`)
- `strings.Builder` em vez de `+=` em laços
- evitar conversões `[]byte` ↔ `string` repetidas
- reutilizar buffers (`sync.Pool`)
- trocar busca em slice por busca em map

> 🧠 *"Premature optimization is the root of all evil."* (Donald Knuth). Meça antes, sempre.

---

## ✍️ Exercícios

1. Coloque um breakpoint na função de cadastro da livraria e inspecione o `Produto` antes de ser salvo.
2. Crie um breakpoint **condicional** que só para quando o preço for maior que 100.
3. Depure um teste que está falhando com `dlv test` ou pelo VS Code.
4. Escreva um benchmark para uma função lenta (ex: concatenar com `+=`), gere o `cpuprofile` e abra no `pprof -http`. Qual função aparece no topo?
5. Adicione o `net/http/pprof` à API da livraria, gere carga com um laço de requisições e analise o perfil de CPU.
6. Provoque um vazamento de goroutines de propósito e encontre-o com `/debug/pprof/goroutine`.

---

⬅️ Anterior: [Ferramentas de qualidade](02-ferramentas-de-qualidade.md) · ➡️ Próximo: [Armadilhas, deploy e próximos passos](04-armadilhas-deploy-e-proximos-passos.md)

# 05 — Padrões de Concorrência

Agora que você conhece as peças (goroutines, channels, `select`, `sync`, `context`), vamos montar os **padrões** mais usados em projetos reais.

---

## 👷 1. Worker pool (grupo de trabalhadores) ⭐

**Problema:** você tem 1.000 tarefas, mas não quer criar 1.000 goroutines ao mesmo tempo (ex: não pode abrir 1.000 conexões com o banco). Quer processar com, digamos, **5 trabalhadores**.

```
            ┌──► worker 1 ──┐
tarefas ────┼──► worker 2 ──┼────► resultados
 (channel)  ├──► worker 3 ──┤      (channel)
            └──► ...      ──┘
```

```go
type Tarefa struct {
	ID     int
	Numero int
}

type Resultado struct {
	TarefaID int
	Valor    int
	Worker   int
}

func worker(id int, tarefas <-chan Tarefa, resultados chan<- Resultado) {
	for t := range tarefas { // pega tarefas até o channel ser fechado
		time.Sleep(100 * time.Millisecond) // simula trabalho pesado
		resultados <- Resultado{TarefaID: t.ID, Valor: t.Numero * t.Numero, Worker: id}
	}
}

func main() {
	const numWorkers = 3
	const numTarefas = 9

	tarefas := make(chan Tarefa, numTarefas)
	resultados := make(chan Resultado, numTarefas)

	// 1. Inicia os workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Go(func() { worker(w, tarefas, resultados) })
	}

	// 2. Envia as tarefas
	for i := 1; i <= numTarefas; i++ {
		tarefas <- Tarefa{ID: i, Numero: i}
	}
	close(tarefas) // avisa: não tem mais tarefa

	// 3. Fecha resultados quando todos os workers terminarem
	go func() {
		wg.Wait()
		close(resultados)
	}()

	// 4. Coleta os resultados
	for r := range resultados {
		fmt.Printf("tarefa %d → %d (worker %d)\n", r.TarefaID, r.Valor, r.Worker)
	}
}
```

9 tarefas de 100 ms com 3 workers levam **~300 ms**, em vez de 900 ms.

---

## 🚰 2. Pipeline (linha de montagem)

Cada **etapa** é uma goroutine que recebe de um channel, processa e envia para o próximo:

```
gerar ──► quadrado ──► filtrarPares ──► imprimir
```

```go
func gerar(nums ...int) <-chan int {
	saida := make(chan int)
	go func() {
		defer close(saida)
		for _, n := range nums {
			saida <- n
		}
	}()
	return saida
}

func quadrado(entrada <-chan int) <-chan int {
	saida := make(chan int)
	go func() {
		defer close(saida)
		for n := range entrada {
			saida <- n * n
		}
	}()
	return saida
}

func filtrarPares(entrada <-chan int) <-chan int {
	saida := make(chan int)
	go func() {
		defer close(saida)
		for n := range entrada {
			if n%2 == 0 {
				saida <- n
			}
		}
	}()
	return saida
}

func main() {
	for v := range filtrarPares(quadrado(gerar(1, 2, 3, 4, 5, 6))) {
		fmt.Print(v, " ")
	}
	// 4 16 36
}
```

O padrão de cada etapa: **cria o channel de saída → inicia a goroutine → `defer close` → retorna o channel**. As etapas trabalham **ao mesmo tempo**, como numa linha de produção.

---

## 🌬️ 3. Fan-out / Fan-in

- **Fan-out:** várias goroutines lendo do **mesmo** channel (dividindo o trabalho)
- **Fan-in:** juntar **vários** channels em **um só**

```go
// Fan-in: junta vários channels em um
func juntar[T any](canais ...<-chan T) <-chan T {
	saida := make(chan T)
	var wg sync.WaitGroup

	for _, c := range canais {
		wg.Go(func() {
			for v := range c {
				saida <- v
			}
		})
	}

	go func() {
		wg.Wait()
		close(saida)
	}()
	return saida
}

func main() {
	entrada := gerar(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out: 3 goroutines de "quadrado" lendo da mesma entrada
	q1 := quadrado(entrada)
	q2 := quadrado(entrada)
	q3 := quadrado(entrada)

	// Fan-in: junta os resultados
	for v := range juntar(q1, q2, q3) {
		fmt.Print(v, " ") // todos os quadrados, em ordem imprevisível
	}
}
```

---

## 🚦 4. Semáforo: limitando a concorrência

Um channel **com buffer** funciona como um "limite de vagas":

```go
func main() {
	urls := make([]string, 20)
	limite := make(chan struct{}, 3) // no máximo 3 ao mesmo tempo

	var wg sync.WaitGroup
	for i := range urls {
		wg.Go(func() {
			limite <- struct{}{}        // pega uma vaga (bloqueia se não houver)
			defer func() { <-limite }() // devolve a vaga

			fmt.Println("baixando", i)
			time.Sleep(200 * time.Millisecond)
		})
	}
	wg.Wait()
}
```

> Mais simples que um worker pool quando cada tarefa é independente.

---

## 📬 5. Gerador com cancelamento

Uma goroutine que produz valores **até ser cancelada**:

```go
func contador(ctx context.Context) <-chan int {
	saida := make(chan int)
	go func() {
		defer close(saida)
		for i := 0; ; i++ {
			select {
			case saida <- i:
			case <-ctx.Done():
				return // ✅ sai quando o contexto é cancelado
			}
		}
	}()
	return saida
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	for n := range contador(ctx) {
		fmt.Print(n, " ")
		if n == 5 {
			cancel() // chega! a goroutine vai terminar
			break
		}
	}
}
```

---

## 🕳️ Goroutine leak (vazamento)

Uma goroutine que **fica bloqueada para sempre** nunca é liberada: é um **vazamento de memória**.

```go
func buscar() string {
	ch := make(chan string) // ❌ sem buffer
	go func() {
		ch <- consultaLenta() // se ninguém receber, fica presa AQUI para sempre
	}()

	select {
	case r := <-ch:
		return r
	case <-time.After(time.Second):
		return "timeout" // saímos... mas a goroutine continua esperando!
	}
}
```

Correções:
- dar **buffer** ao channel (`make(chan string, 1)`) para o envio não bloquear
- usar **`context`** e fazer a goroutine desistir quando o contexto for cancelado

> ✅ **Regra:** ao iniciar uma goroutine, saiba **como e quando ela vai terminar**.

---

## ✅ Resumo: boas práticas de concorrência

1. **Não use concorrência sem motivo.** Código sequencial é mais simples. Use quando houver espera (I/O, rede) ou trabalho paralelizável.
2. Toda goroutine precisa de um **fim claro**: channel fechado, `context` cancelado ou trabalho concluído.
3. **Quem envia fecha** o channel.
4. Use **channels direcionais** nos parâmetros.
5. Passe **`context.Context`** para operações demoradas e respeite o cancelamento.
6. Proteja estado compartilhado com **mutex** ou **atomic**, ou evite compartilhar usando channels.
7. **Sempre** teste com **`-race`**.
8. **Limite** a concorrência (worker pool, semáforo) quando acessar recursos externos.

---

## 🧾 Resumão do módulo

| Ferramenta | Para quê |
|---|---|
| `go f()` | iniciar goroutine |
| `sync.WaitGroup` / `wg.Go` | esperar goroutines terminarem |
| `make(chan T)` | channel sem buffer (sincroniza) |
| `make(chan T, n)` | channel com buffer (fila) |
| `close(ch)` + `range ch` | sinalizar fim / consumir até o fim |
| `chan<- T` / `<-chan T` | channels direcionais |
| `select` | esperar vários channels |
| `time.After` / `Ticker` | timeouts / repetição |
| `context` | cancelamento e prazos |
| `sync.Mutex` / `RWMutex` | proteger estado compartilhado |
| `sync/atomic` | contadores/flags sem lock |
| `sync.Once` | inicializar uma vez |
| `-race` | detectar condições de corrida |

---

## ✍️ Exercícios

1. Implemente um **worker pool** com 4 workers que calculam se os números de 1 a 100 são primos. Imprima só os primos.
2. Monte um **pipeline** de 3 etapas: gerar palavras → converter para maiúsculas → filtrar as com mais de 4 letras.
3. Use **fan-out/fan-in** para calcular o quadrado de 1 a 20 com 4 goroutines e some os resultados.
4. Simule 30 "downloads" de 300 ms cada com um **semáforo** de 5. Meça o tempo total (deve ficar perto de 1,8 s).
5. Crie um **gerador** infinito de números aleatórios que para quando o contexto expira (`WithTimeout` de 1 s). Conte quantos números foram gerados.
6. Encontre e corrija o **goroutine leak** do exemplo `buscar()`. Use `runtime.NumGoroutine()` para provar que corrigiu.
7. **Projeto:** na livraria, simule 50 clientes comprando o mesmo livro (estoque 10) ao mesmo tempo. Garanta com mutex que o estoque nunca fique negativo e rode com `-race`.

---

⬅️ Anterior: [sync e race conditions](04-sync-e-race-conditions.md) · ➡️ Próximo módulo: [Testes](../15-testes/README.md)

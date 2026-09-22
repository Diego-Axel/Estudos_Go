# 03 — `select`, Timeouts e `context`

## 🔀 `select`: esperando vários channels

O `select` é como um `switch`, mas para **operações em channels**. Ele espera até **uma** delas estar pronta e executa aquele caso:

```go
func main() {
	rapido := make(chan string)
	lento := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		rapido <- "🐇 rápido"
	}()
	go func() {
		time.Sleep(time.Second)
		lento <- "🐢 lento"
	}()

	for range 2 {
		select {
		case msg := <-rapido:
			fmt.Println(msg)
		case msg := <-lento:
			fmt.Println(msg)
		}
	}
}
// 🐇 rápido
// 🐢 lento
```

Regras:
- Bloqueia até **algum** caso estar pronto
- Se **vários** estiverem prontos ao mesmo tempo, escolhe **um aleatoriamente** (justo)
- Os casos podem ser **envios** (`ch <- v`) ou **recebimentos** (`<-ch`)

---

## 🚫 `default`: operação sem bloquear

Com `default`, o `select` **não espera**: se nada estiver pronto, cai no `default`:

```go
msgs := make(chan string)

select {
case m := <-msgs:
	fmt.Println("recebi", m)
default:
	fmt.Println("nenhuma mensagem agora") // ← cai aqui
}
```

Enviar sem bloquear (descartando se estiver cheio):

```go
select {
case fila <- tarefa:
	// enviado
default:
	fmt.Println("fila cheia, descartando tarefa")
}
```

---

## ⏰ Timeouts com `time.After`

`time.After(d)` retorna um channel que recebe um valor **depois** de `d`. Combinado com `select`, cria um **timeout**:

```go
func buscarDados() <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(3 * time.Second) // simula algo lento
		ch <- "dados"
	}()
	return ch
}

func main() {
	select {
	case d := <-buscarDados():
		fmt.Println("recebi:", d)
	case <-time.After(1 * time.Second):
		fmt.Println("⏰ timeout! desisti de esperar")
	}
}
// ⏰ timeout! desisti de esperar
```

> 💡 Repare no `make(chan string, 1)`: com buffer de 1, a goroutine consegue enviar e **terminar** mesmo que ninguém receba mais. Sem o buffer, ela ficaria presa para sempre (um *goroutine leak*).

---

## ⏱️ Tarefas periódicas: `time.Ticker`

```go
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop()

fim := time.After(2 * time.Second)

for {
	select {
	case t := <-ticker.C:
		fmt.Println("tick", t.Format("15:04:05.000"))
	case <-fim:
		fmt.Println("fim!")
		return
	}
}
```

---

## 🛑 Canal de cancelamento (`done`)

Um padrão clássico: um channel que, ao ser **fechado**, avisa **todas** as goroutines para pararem (lembra: receber de channel fechado nunca bloqueia):

```go
func trabalhador(id int, done <-chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("trabalhador", id, "parando")
			return
		default:
			fmt.Println("trabalhador", id, "trabalhando...")
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan struct{})
	for i := 1; i <= 3; i++ {
		go trabalhador(i, done)
	}

	time.Sleep(time.Second)
	close(done) // 📢 avisa TODOS de uma vez
	time.Sleep(100 * time.Millisecond)
}
```

> `chan struct{}` é usado porque não carrega dados (ocupa 0 bytes): é só um **sinal**.

Esse padrão é tão comum que virou um pacote da biblioteca padrão: o **`context`**.

---

## 🎛️ O pacote `context` ⭐

O `context.Context` carrega **prazos**, **sinais de cancelamento** e **valores** entre funções e goroutines. É usado em **todo lugar** em Go moderno: HTTP, bancos de dados, chamadas de API...

### Cancelamento manual: `WithCancel`

```go
func trabalhador(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(id, "parando:", ctx.Err()) // context canceled
			return
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go trabalhador(ctx, i)
	}

	time.Sleep(time.Second)
	cancel() // cancela todos
	time.Sleep(100 * time.Millisecond)
}
```

### Tempo limite: `WithTimeout`

```go
func consultarAPI(ctx context.Context) (string, error) {
	select {
	case <-time.After(2 * time.Second): // simula uma API lenta
		return "resposta", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // ✅ SEMPRE chame cancel para liberar recursos

	resp, err := consultarAPI(ctx)
	if err != nil {
		fmt.Println("erro:", err) // erro: context deadline exceeded
		return
	}
	fmt.Println(resp)
}
```

### Resumo das funções

| Função | Cancela quando... |
|---|---|
| `context.Background()` | nunca (é a raiz, usada no `main` e nos testes) |
| `context.TODO()` | nunca (marcador para "ainda não sei qual contexto usar") |
| `context.WithCancel(pai)` | você chama `cancel()` |
| `context.WithTimeout(pai, d)` | passa o tempo `d` (ou `cancel()`) |
| `context.WithDeadline(pai, t)` | chega o horário `t` (ou `cancel()`) |
| `context.WithValue(pai, k, v)` | não cancela: carrega um valor |

Os contextos formam uma **árvore**: cancelar um pai cancela **todos os filhos**.

```
Background
   └── WithTimeout(5s)          ← requisição HTTP
          ├── WithTimeout(1s)   ← consulta ao banco
          └── WithCancel        ← chamada a outra API
```

### ✅ Boas práticas com `context`

1. O `ctx` é **sempre o primeiro parâmetro**, chamado `ctx`: `func Buscar(ctx context.Context, id int)`
2. **Não guarde** contexto dentro de structs; passe como parâmetro
3. Sempre `defer cancel()` logo depois de criar
4. `WithValue` só para dados **da requisição** (ID de rastreio, usuário autenticado), **nunca** para parâmetros opcionais de funções
5. Nunca passe `nil` como contexto. Na dúvida, `context.TODO()`

> Em servidores HTTP, cada requisição já vem com um contexto: `r.Context()`. Se o cliente desconectar, ele é cancelado automaticamente.

---

## ✍️ Exercícios

1. Crie duas goroutines que enviam mensagens em intervalos diferentes (200 ms e 500 ms) e use `select` num laço para imprimir as 10 primeiras mensagens.
2. Faça uma função que simula uma consulta com tempo aleatório (0 a 2 s) e use `time.After` para desistir após 1 s.
3. Use um `Ticker` para imprimir um relógio (`HH:MM:SS`) a cada segundo, parando após 5 segundos.
4. Implemente o padrão `done` com 5 trabalhadores e pare todos depois de 2 segundos.
5. Refaça o exercício 4 usando `context.WithTimeout`.
6. Crie uma função `baixar(ctx, url)` que respeite o cancelamento, e cancele o contexto quando **o primeiro** de 3 downloads terminar.
7. Tente enviar sem bloquear para um channel com buffer 2, 5 vezes seguidas, usando `select` + `default`. Quantos foram descartados?

---

⬅️ Anterior: [Channels](02-channels.md) · ➡️ Próximo: [sync e race conditions](04-sync-e-race-conditions.md)

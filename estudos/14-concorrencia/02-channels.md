# 02 — Channels

> *"Não comunique compartilhando memória; compartilhe memória comunicando."*
> (*Don't communicate by sharing memory; share memory by communicating.*)

Um **channel** é um **cano** por onde goroutines **enviam e recebem** valores de forma segura.

```
goroutine A  ──── valor ────►  [ channel ]  ──── valor ────►  goroutine B
```

---

## 🔧 Criando e usando

```go
ch := make(chan int) // channel de int

go func() {
	ch <- 42 // ENVIA 42 no channel
}()

v := <-ch      // RECEBE do channel
fmt.Println(v) // 42
```

| Operação | Sintaxe |
|---|---|
| criar | `ch := make(chan T)` |
| enviar | `ch <- valor` |
| receber | `v := <-ch` |
| receber e descartar | `<-ch` |
| fechar | `close(ch)` |

> 🧠 Dica para lembrar: a **seta aponta** para onde o dado vai. `ch <- v`: o valor entra no channel. `v := <-ch`: o valor sai do channel.

---

## 🚦 Channels sem buffer: sincronização

Um `make(chan T)` **sem tamanho** não guarda nada. Por isso:

- **Enviar bloqueia** até alguém receber
- **Receber bloqueia** até alguém enviar

É um **encontro marcado**: as duas goroutines se "encontram" no momento da troca.

```go
func trabalhar(pronto chan bool) {
	fmt.Println("trabalhando...")
	time.Sleep(time.Second)
	fmt.Println("terminei!")
	pronto <- true // avisa
}

func main() {
	pronto := make(chan bool)
	go trabalhar(pronto)
	<-pronto // espera o aviso (sem WaitGroup, sem Sleep!)
	fmt.Println("main pode continuar")
}
```

---

## 📦 Channels com buffer

Com tamanho, o channel vira uma **fila** com capacidade limitada:

```go
ch := make(chan string, 3) // cabe 3 valores

ch <- "a" // não bloqueia
ch <- "b" // não bloqueia
ch <- "c" // não bloqueia
// ch <- "d" // bloquearia: buffer cheio

fmt.Println(<-ch, <-ch) // a b (ordem FIFO: primeiro a entrar, primeiro a sair)
fmt.Println(len(ch), cap(ch)) // 1 3
```

- **Enviar** só bloqueia com o buffer **cheio**
- **Receber** só bloqueia com o buffer **vazio**

> 🧠 **Sem buffer** = sincronização garantida. **Com buffer** = desacopla quem produz de quem consome. Na dúvida, comece **sem buffer**.

---

## 🔒 Fechando channels: `close`

`close(ch)` avisa: **"não vou mandar mais nada"**.

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
close(ch)

fmt.Println(<-ch) // 1
fmt.Println(<-ch) // 2
fmt.Println(<-ch) // 0 ← channel fechado e vazio: retorna o valor zero

v, ok := <-ch
fmt.Println(v, ok) // 0 false ← "ok" diz se o valor veio de um envio real
```

### `range` em channel ⭐

Recebe valores **até o channel ser fechado**:

```go
func gerar(n int, ch chan int) {
	for i := 1; i <= n; i++ {
		ch <- i
	}
	close(ch) // ⚠️ sem isso, o range do main esperaria para sempre
}

func main() {
	ch := make(chan int)
	go gerar(5, ch)

	for v := range ch {
		fmt.Print(v, " ")
	}
	// 1 2 3 4 5
}
```

### Regras do `close`

| Situação | Resultado |
|---|---|
| receber de channel fechado | valor zero, `ok = false` (não bloqueia) |
| **enviar** para channel fechado | 💥 `panic: send on closed channel` |
| fechar channel já fechado | 💥 `panic: close of closed channel` |
| fechar channel `nil` | 💥 `panic: close of nil channel` |

> ✅ **Regra:** quem **envia** é quem **fecha**. Nunca feche do lado de quem recebe.
>
> 💡 Não é obrigatório fechar todo channel. Feche quando o receptor **precisa saber** que acabou (ex: para o `range` terminar).

---

## ➡️ Channels direcionais

Em parâmetros de funções, você pode restringir o channel a **só enviar** ou **só receber**:

```go
func produtor(saida chan<- int) { // só ENVIA
	for i := range 3 {
		saida <- i
	}
	close(saida)
	// v := <-saida // ❌ erro: não pode receber de um channel send-only
}

func consumidor(entrada <-chan int) { // só RECEBE
	for v := range entrada {
		fmt.Println("recebi", v)
	}
	// entrada <- 1 // ❌ erro: não pode enviar para um channel receive-only
}

func main() {
	ch := make(chan int) // bidirecional
	go produtor(ch)      // convertido automaticamente para chan<- int
	consumidor(ch)       // convertido automaticamente para <-chan int
}
```

| Tipo | Pode |
|---|---|
| `chan T` | enviar e receber |
| `chan<- T` | só enviar |
| `<-chan T` | só receber |

> ✅ Use direcionais sempre que possível: o compilador impede usos errados e o código documenta a intenção.

---

## 💀 Deadlock

Se **todas** as goroutines estão bloqueadas esperando umas pelas outras, o runtime detecta e encerra o programa:

```go
func main() {
	ch := make(chan int)
	ch <- 1 // bloqueia para sempre: ninguém vai receber
}
// fatal error: all goroutines are asleep - deadlock!
```

Causas comuns:
- enviar num channel sem buffer sem ter outra goroutine recebendo
- `range` num channel que nunca é fechado
- esquecer o `wg.Done()`

---

## 🕳️ Channel `nil`

O valor zero de um channel é `nil`. Enviar ou receber de um channel `nil` **bloqueia para sempre**:

```go
var ch chan int // nil
// ch <- 1   // bloqueia para sempre
// <-ch      // bloqueia para sempre
```

Parece inútil, mas é um truque útil com `select` (próximo arquivo), para "desligar" um caso.

---

## 🧪 Exemplo: somando em paralelo

```go
func somar(nums []int, resultado chan<- int) {
	total := 0
	for _, n := range nums {
		total += n
	}
	resultado <- total
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resultado := make(chan int)

	meio := len(nums) / 2
	go somar(nums[:meio], resultado) // primeira metade
	go somar(nums[meio:], resultado) // segunda metade

	a, b := <-resultado, <-resultado
	fmt.Println(a + b) // 55
}
```

---

## ✍️ Exercícios

1. Crie uma goroutine que envia seu nome por um channel e imprima no `main`.
2. Crie um produtor que envia os números de 1 a 10 e fecha o channel, e um consumidor com `range` que imprime o quadrado de cada um.
3. Mostre a diferença de comportamento entre `make(chan int)` e `make(chan int, 3)` enviando 3 valores sem ninguém recebendo.
4. Provoque um deadlock de propósito e leia a mensagem.
5. Divida um slice de 1.000 números em 4 partes, some cada parte numa goroutine e junte os resultados.
6. Reescreva o exercício 2 usando channels **direcionais** nos parâmetros.
7. O que acontece ao enviar para um channel fechado? E ao receber? Teste os dois.

---

⬅️ Anterior: [Goroutines](01-goroutines.md) · ➡️ Próximo: [select, timeouts e context](03-select-timeouts-e-context.md)

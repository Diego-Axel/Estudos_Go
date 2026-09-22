# 01 — Goroutines

## 🤹 Concorrência vs. paralelismo

| **Concorrência** | **Paralelismo** |
|---|---|
| **lidar** com várias coisas ao mesmo tempo | **fazer** várias coisas ao mesmo tempo |
| é sobre **estrutura** do programa | é sobre **execução** |
| funciona até com 1 núcleo de CPU (alternando) | precisa de vários núcleos |
| 👨‍🍳 um cozinheiro alternando entre 3 panelas | 👨‍🍳👩‍🍳👨‍🍳 três cozinheiros, cada um numa panela |

> *"Concurrency is not parallelism."* (Rob Pike)

Go foi feito para **concorrência**, e o runtime distribui o trabalho entre os núcleos disponíveis, dando **paralelismo** de brinde.

---

## 🚀 O que é uma goroutine?

Uma **goroutine** é uma função rodando **de forma concorrente** com o resto do programa. Para criar uma, basta colocar **`go`** antes da chamada:

```go
func dizer(msg string) {
	for i := range 3 {
		fmt.Println(msg, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go dizer("goroutine") // roda "em paralelo"
	dizer("main")         // roda na goroutine principal
}
```

Saída (a ordem pode variar):

```
main 0
goroutine 0
goroutine 1
main 1
main 2
goroutine 2
```

### Por que goroutines são especiais?

| | Thread do sistema operacional | Goroutine |
|---|---|---|
| Memória inicial | ~1 MB a 8 MB | **~2 KB** (cresce se precisar) |
| Criar/destruir | caro | **baratíssimo** |
| Quantas cabem | milhares | **centenas de milhares / milhões** |
| Quem gerencia | o sistema operacional | o **runtime do Go** |

O runtime do Go distribui milhares de goroutines entre poucas threads reais (o escalonador **M:N**).

```go
for i := range 100_000 {
	go func() { _ = i * 2 }() // 100 mil goroutines, sem problema
}
```

---

## ⚠️ Quando o `main` termina, tudo termina

```go
func main() {
	go fmt.Println("será que imprime?")
	// main termina aqui → o programa acaba → a goroutine morre junto
}
// (provavelmente não imprime nada)
```

O `main` **não espera** as goroutines. Colocar `time.Sleep` "resolve", mas é **gambiarra**: você não sabe quanto tempo esperar.

---

## ⏳ Esperando goroutines: `sync.WaitGroup`

O `WaitGroup` é um **contador** de goroutines em andamento:

```go
import "sync"

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // +1 goroutine para esperar
		go func() {
			defer wg.Done() // -1 quando terminar
			fmt.Println("tarefa", i)
		}()
	}

	wg.Wait() // bloqueia até o contador chegar a zero
	fmt.Println("todas terminaram ✅")
}
```

Regras:
- `Add` **antes** de iniciar a goroutine (não dentro dela)
- `Done` com **`defer`**, para rodar mesmo se a função retornar cedo ou der panic
- Passe o `WaitGroup` por **ponteiro** se for para outra função (`*sync.WaitGroup`)

### `wg.Go` (Go 1.25+) ⭐

Uma forma mais curta que faz o `Add(1)` e o `Done()` por você:

```go
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
	wg.Go(func() {
		fmt.Println("tarefa", i)
	})
}
wg.Wait()
```

---

## 🔁 Goroutines em laços (Go 1.22+)

Desde o Go 1.22, cada volta do laço tem **sua própria** variável, então isso funciona como esperado:

```go
for i := range 3 {
	go func() {
		fmt.Println(i) // 0, 1, 2 (em alguma ordem)
	}()
}
```

> Em versões antigas, todas as goroutines viam o **mesmo** `i` (e imprimiam `3 3 3`). Por isso, em código antigo, você vai ver `go func(i int) { ... }(i)`.

---

## 🧮 Exemplo: baixando várias "páginas" ao mesmo tempo

```go
func baixar(url string) {
	inicio := time.Now()
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // simula a rede
	fmt.Printf("%-25s %v\n", url, time.Since(inicio).Round(time.Millisecond))
}

func main() {
	urls := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://gobyexample.com",
		"https://github.com",
	}

	inicio := time.Now()
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Go(func() { baixar(url) })
	}
	wg.Wait()

	fmt.Println("total:", time.Since(inicio).Round(time.Millisecond))
	// total ≈ tempo da MAIS LENTA, e não a soma de todas! 🚀
}
```

---

## 🖥️ Quantos núcleos o Go usa?

```go
import "runtime"

fmt.Println(runtime.NumCPU())       // núcleos disponíveis
fmt.Println(runtime.GOMAXPROCS(0))  // quantos o Go está usando (padrão: todos)
fmt.Println(runtime.NumGoroutine()) // goroutines vivas agora
```

> Normalmente você **não precisa** mexer no `GOMAXPROCS`.

---

## ⚠️ Goroutines precisam se comunicar com cuidado

Se várias goroutines **mexem na mesma variável** ao mesmo tempo, o resultado é imprevisível:

```go
contador := 0
var wg sync.WaitGroup
for range 1000 {
	wg.Go(func() {
		contador++ // ⚠️ condição de corrida (race condition)!
	})
}
wg.Wait()
fmt.Println(contador) // 1000? 987? 954? 😱
```

Para resolver isso existem **channels** (próximo arquivo) e **mutex** (arquivo 04).

---

## ✍️ Exercícios

1. Crie 5 goroutines que imprimem "Olá da goroutine N" e use um `WaitGroup` para esperar todas.
2. Faça um programa sem `WaitGroup` e veja as goroutines "sumirem". Depois corrija.
3. Simule 10 downloads com tempos aleatórios. Compare o tempo total **sequencial** vs. **com goroutines**.
4. Crie 100.000 goroutines que só dormem 1 segundo e meça o tempo e a memória (`runtime.NumGoroutine()`).
5. Rode o exemplo do contador várias vezes e anote os resultados. Por que mudam?
6. Reescreva o exercício 1 usando `wg.Go` (Go 1.25+).

---

🏠 [Módulo 14](README.md) · ➡️ Próximo: [Channels](02-channels.md)

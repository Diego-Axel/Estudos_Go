# 03 — Composição de Interfaces e Interfaces Famosas

## 🧩 Compondo interfaces

Assim como structs, interfaces podem **embutir** outras interfaces:

```go
type Leitor interface {
	Ler() string
}

type Escritor interface {
	Escrever(s string)
}

type LeitorEscritor interface {
	Leitor   // tem todos os métodos de Leitor
	Escritor // e todos os de Escritor
}
```

Um tipo satisfaz `LeitorEscritor` se tiver `Ler()` **e** `Escrever()`.

> 🧠 A filosofia de Go: **interfaces pequenas** (1 ou 2 métodos) que se **combinam**. *"The bigger the interface, the weaker the abstraction."* (Rob Pike)

Convenção de nomes: interfaces de **um método** terminam em **-er**: `Reader`, `Writer`, `Stringer`, `Closer`, `Formatter`...

---

## ⭐ As interfaces mais importantes da biblioteca padrão

### 1. `fmt.Stringer`

```go
type Stringer interface {
	String() string
}
```

Você já conhece (módulo 09): define como o valor aparece no `fmt.Println`.

---

### 2. `error`

A interface de erro é **nativa** da linguagem:

```go
type error interface {
	Error() string
}
```

Então, **qualquer tipo** com `Error() string` é um erro:

```go
type ErroSaldo struct {
	Disponivel, Solicitado float64
}

func (e ErroSaldo) Error() string {
	return fmt.Sprintf("saldo insuficiente: disponível %.2f, solicitado %.2f",
		e.Disponivel, e.Solicitado)
}

func sacar(saldo, valor float64) error {
	if valor > saldo {
		return ErroSaldo{saldo, valor}
	}
	return nil
}

if err := sacar(100, 150); err != nil {
	fmt.Println(err) // saldo insuficiente: disponível 100.00, solicitado 150.00
}
```

> O módulo 11 (Tratamento de erros) aprofunda isso com `errors.Is`, `errors.As` e wrapping.

---

### 3. `io.Reader` e `io.Writer` 🔥

As duas interfaces mais poderosas do Go:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
```

- **`Reader`**: algo de onde se **lê** bytes (arquivo, rede, string, teclado, resposta HTTP, arquivo zip...)
- **`Writer`**: algo onde se **escreve** bytes (arquivo, terminal, rede, buffer, hash...)

Quem implementa:

| `io.Reader` | `io.Writer` |
|---|---|
| `os.Stdin`, `*os.File` | `os.Stdout`, `os.Stderr`, `*os.File` |
| `strings.NewReader(s)` | `*strings.Builder` |
| `*bytes.Buffer` | `*bytes.Buffer` |
| `http.Response.Body` | `http.ResponseWriter` |
| `gzip.Reader` | `gzip.Writer`, `sha256.New()` |

A mágica: uma função que recebe `io.Writer` funciona com **todos** eles:

```go
func gerarRelatorio(w io.Writer, itens []string) {
	fmt.Fprintln(w, "=== RELATÓRIO ===")
	for i, item := range itens {
		fmt.Fprintf(w, "%d. %s\n", i+1, item)
	}
}

itens := []string{"Caneta", "Caderno"}

gerarRelatorio(os.Stdout, itens) // imprime no terminal

var sb strings.Builder
gerarRelatorio(&sb, itens)       // monta uma string

f, _ := os.Create("relatorio.txt")
defer f.Close()
gerarRelatorio(f, itens)         // grava num arquivo
```

**Uma função, três destinos, zero alteração.** 🎉

#### Criando seu próprio `Writer`

```go
// ContadorDeBytes conta quantos bytes passaram por ele.
type ContadorDeBytes struct {
	Total int
}

func (c *ContadorDeBytes) Write(p []byte) (int, error) {
	c.Total += len(p)
	return len(p), nil
}

func main() {
	c := &ContadorDeBytes{}
	fmt.Fprintf(c, "Olá, %s!", "Go") // o fmt escreve no nosso tipo
	fmt.Println(c.Total)             // 9 ("Olá, Go!" tem 9 bytes: o 'á' ocupa 2)
}
```

#### Combinando com `io.MultiWriter` e `io.Copy`

```go
// Escreve no terminal E no arquivo ao mesmo tempo
w := io.MultiWriter(os.Stdout, f)
fmt.Fprintln(w, "log duplicado")

// Copia tudo de um Reader para um Writer
io.Copy(os.Stdout, strings.NewReader("copiado!\n"))
```

---

### 4. `sort.Interface`

Para ordenar coleções próprias com `sort.Sort`:

```go
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}
```

```go
type Produto struct {
	Nome  string
	Preco float64
}

type PorPreco []Produto

func (p PorPreco) Len() int           { return len(p) }
func (p PorPreco) Less(i, j int) bool { return p[i].Preco < p[j].Preco }
func (p PorPreco) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	produtos := []Produto{{"Caderno", 25}, {"Caneta", 2.5}, {"Mochila", 120}}
	sort.Sort(PorPreco(produtos))
	fmt.Println(produtos) // [{Caneta 2.5} {Caderno 25} {Mochila 120}]
}
```

> 💡 Hoje em dia, o mais prático é `slices.SortFunc` (módulo 06). Mas `sort.Interface` é ótima para **entender** interfaces, e você vai encontrá-la em muito código.

---

### 5. `io.Closer` e amigos

```go
type Closer interface {
	Close() error
}

type ReadCloser interface {
	Reader
	Closer
}
```

É por isso que você faz `defer resp.Body.Close()`: o `Body` é um `io.ReadCloser`.

---

### 6. `http.Handler` (spoiler do módulo Web)

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

Qualquer tipo com esse método pode responder requisições HTTP.

---

## 📋 Resumo das interfaces famosas

| Interface | Método(s) | Para quê |
|---|---|---|
| `fmt.Stringer` | `String() string` | imprimir bonito |
| `error` | `Error() string` | representar erros |
| `io.Reader` | `Read([]byte) (int, error)` | ler bytes de qualquer fonte |
| `io.Writer` | `Write([]byte) (int, error)` | escrever bytes em qualquer destino |
| `io.Closer` | `Close() error` | liberar recursos |
| `sort.Interface` | `Len`, `Less`, `Swap` | ordenar |
| `http.Handler` | `ServeHTTP(w, r)` | responder HTTP |

---

## ✍️ Exercícios

1. Crie as interfaces `Ligavel` (`Ligar()`), `Desligavel` (`Desligar()`) e a composta `Interruptor`. Implemente em `Lampada` e `Ventilador`.
2. Escreva `salvarLista(w io.Writer, nomes []string)` e use com `os.Stdout`, `strings.Builder` e um arquivo.
3. Crie um `Writer` que transforme todo texto em **MAIÚSCULAS** antes de escrever em `os.Stdout`.
4. Crie o erro personalizado `ErroValidacao{Campo, Motivo string}` e use numa função `validarEmail`.
5. Ordene um slice de `Aluno{Nome string; Nota float64}` pela nota (decrescente) implementando `sort.Interface`. Depois refaça com `slices.SortFunc`.
6. Use `io.MultiWriter` para gravar um log no terminal e num arquivo `app.log` ao mesmo tempo.

---

⬅️ Anterior: [any, type assertion e type switch](02-any-type-assertion-type-switch.md) · ➡️ Próximo: [nil e boas práticas](04-nil-e-boas-praticas.md)

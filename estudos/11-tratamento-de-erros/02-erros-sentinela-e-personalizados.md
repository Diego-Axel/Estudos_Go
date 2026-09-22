# 02 — Erros Sentinela e Erros Personalizados

Às vezes **só a mensagem** não basta: quem chama precisa **saber qual erro** aconteceu para decidir o que fazer. Go tem duas formas principais para isso.

---

## 🚩 Erros sentinela

Um **erro sentinela** é uma **variável de pacote** com um erro pré-definido, que serve de "marcador" para comparação:

```go
package estoque

import "errors"

var (
	ErrProdutoNaoEncontrado = errors.New("produto não encontrado")
	ErrEstoqueInsuficiente  = errors.New("estoque insuficiente")
)

func Vender(codigo string, qtd int) error {
	p, ok := produtos[codigo]
	if !ok {
		return ErrProdutoNaoEncontrado
	}
	if p.Estoque < qtd {
		return ErrEstoqueInsuficiente
	}
	p.Estoque -= qtd
	return nil
}
```

Quem chama pode **comparar**:

```go
err := estoque.Vender("L001", 5)
switch {
case err == nil:
	fmt.Println("venda realizada ✅")
case errors.Is(err, estoque.ErrProdutoNaoEncontrado):
	fmt.Println("código inválido, confira e tente de novo")
case errors.Is(err, estoque.ErrEstoqueInsuficiente):
	fmt.Println("sem estoque, deseja encomendar?")
default:
	fmt.Println("erro inesperado:", err)
}
```

> 📌 Convenção de nome: começa com **`Err`** (`ErrNotFound`, `ErrTimeout`...).
>
> 💡 Prefira `errors.Is(err, ErrX)` a `err == ErrX`. O `errors.Is` também funciona quando o erro foi **embrulhado** (próximo arquivo).

### Sentinelas famosos da biblioteca padrão

| Erro | Pacote | Significado |
|---|---|---|
| `io.EOF` | `io` | fim dos dados (não é bem um "erro"!) |
| `fs.ErrNotExist` / `os.ErrNotExist` | `io/fs`, `os` | arquivo não existe |
| `fs.ErrPermission` | `io/fs` | sem permissão |
| `sql.ErrNoRows` | `database/sql` | consulta sem resultados |
| `http.ErrServerClosed` | `net/http` | servidor foi desligado |
| `context.DeadlineExceeded` | `context` | tempo esgotado |

```go
_, err := os.Open("nao-existe.txt")
if errors.Is(err, fs.ErrNotExist) {
	fmt.Println("arquivo não existe, criando um novo...")
}
```

### `io.EOF`: o erro "esperado"

```go
r := strings.NewReader("abc")
buf := make([]byte, 2)

for {
	n, err := r.Read(buf)
	if err == io.EOF {
		break // chegou ao fim: comportamento normal!
	}
	if err != nil {
		log.Fatal(err) // aí sim, erro de verdade
	}
	fmt.Printf("%q\n", buf[:n])
}
// "ab"
// "c"
```

> `io.EOF` é uma das poucas exceções onde se usa `==` direto, pois por convenção ele **nunca** é embrulhado.

---

## 🧱 Erros personalizados (tipos de erro)

Quando o erro precisa carregar **dados** (qual campo? qual valor? qual código?), crie um **tipo**:

```go
type ErroValidacao struct {
	Campo  string
	Valor  any
	Motivo string
}

func (e *ErroValidacao) Error() string {
	return fmt.Sprintf("campo %q (valor %v): %s", e.Campo, e.Valor, e.Motivo)
}

func validarProduto(nome string, preco float64) error {
	if nome == "" {
		return &ErroValidacao{Campo: "nome", Valor: nome, Motivo: "não pode ser vazio"}
	}
	if preco <= 0 {
		return &ErroValidacao{Campo: "preco", Valor: preco, Motivo: "deve ser positivo"}
	}
	return nil
}
```

Para acessar os **campos** do erro, use `errors.As`:

```go
err := validarProduto("Caneta", -5)

var ev *ErroValidacao
if errors.As(err, &ev) {
	fmt.Println("Corrija o campo:", ev.Campo) // Corrija o campo: preco
	fmt.Println("Motivo:", ev.Motivo)         // Motivo: deve ser positivo
}
```

> 📌 Convenção de nome para **tipos** de erro: terminam em **`Error`** em inglês (`PathError`, `SyntaxError`). Em português, algo como `ErroValidacao` fica claro.

### Exemplo real da biblioteca padrão: `*fs.PathError`

```go
_, err := os.Open("/segredo/dados.txt")

var pe *fs.PathError
if errors.As(err, &pe) {
	fmt.Println("operação:", pe.Op)   // open
	fmt.Println("caminho:", pe.Path)  // /segredo/dados.txt
	fmt.Println("causa:", pe.Err)     // no such file or directory
}
```

E `strconv.Atoi` retorna um `*strconv.NumError`:

```go
_, err := strconv.Atoi("99999999999999999999")

var ne *strconv.NumError
if errors.As(err, &ne) {
	fmt.Println(ne.Func, ne.Num)                 // Atoi 99999999999999999999
	fmt.Println(errors.Is(ne, strconv.ErrRange)) // true (número grande demais)
}
```

---

## 🆚 Sentinela ou tipo personalizado?

| | Sentinela (`var ErrX = errors.New(...)`) | Tipo (`type XError struct{...}`) |
|---|---|---|
| Carrega dados extras? | ❌ só a mensagem | ✅ quantos campos quiser |
| Como verificar | `errors.Is(err, ErrX)` | `errors.As(err, &alvo)` |
| Quando usar | condição simples e conhecida | precisa de detalhes (campo, código HTTP...) |
| Exemplo | `ErrNaoEncontrado` | `ErroValidacao{Campo, Motivo}` |

> ⚠️ Ambos viram parte da **API pública** do seu pacote: depois que alguém depende deles, mudar quebra código. Crie só os que fazem sentido o chamador **tratar de forma diferente**.

---

## ✍️ Exercícios

1. Crie os sentinelas `ErrSaldoInsuficiente` e `ErrValorInvalido` e uma função `Sacar` que os use. No `main`, trate cada um com uma mensagem diferente.
2. Leia um arquivo; se ele não existir (`fs.ErrNotExist`), crie-o com um conteúdo padrão.
3. Crie o tipo `ErroHTTP{Codigo int; Mensagem string}` e uma função que simule requisições retornando 404 ou 500. Use `errors.As` para mostrar o código.
4. Leia um `strings.NewReader` de 3 em 3 bytes até o `io.EOF`, montando o texto de volta.
5. Use `errors.As` com `*strconv.NumError` para diferenciar "não é número" (`ErrSyntax`) de "número grande demais" (`ErrRange`).
6. No seu `main.go` da livraria, crie `ErrCodigoDuplicado` e `ErrProdutoNaoEncontrado` e use-os nas funções de cadastro e busca.

---

⬅️ Anterior: [Erros são valores](01-erros-sao-valores.md) · ➡️ Próximo: [Wrapping, Is, As e Join](03-wrapping-is-as-join.md)

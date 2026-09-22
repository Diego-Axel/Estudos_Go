# 03 — Ponteiros com Structs e Coleções

> Structs serão vistas em detalhe no **módulo 09**. Aqui vamos só usar o básico: uma struct é um "agrupamento de campos", como o `Produto` do seu `main.go`.

```go
type Produto struct {
	Nome    string
	Preco   float64
	Estoque int
}
```

---

## 🧱 Ponteiro para struct

```go
p := Produto{Nome: "Caneta", Preco: 2.5, Estoque: 10}
ptr := &p

(*ptr).Estoque = 5 // forma "completa"
ptr.Estoque = 5    // ✅ forma usada: Go desreferencia sozinho!

fmt.Println(p.Estoque) // 5
```

Go faz a **desreferência automática** para acessar campos: `ptr.Campo` é o mesmo que `(*ptr).Campo`.

### Criando direto como ponteiro ⭐

```go
ptr := &Produto{Nome: "Lápis", Preco: 1.2}
fmt.Println(ptr.Nome) // Lápis
```

Esse é o jeito **mais comum** de criar structs que vão ser modificadas ou compartilhadas.

---

## 🛠️ Funções que alteram structs

```go
// ❌ Recebe cópia: a alteração se perde
func darDesconto(p Produto, pct float64) {
	p.Preco *= 1 - pct/100
}

// ✅ Recebe ponteiro: altera o original
func darDescontoPtr(p *Produto, pct float64) {
	p.Preco *= 1 - pct/100
}

func main() {
	caneta := Produto{Nome: "Caneta", Preco: 10}

	darDesconto(caneta, 10)
	fmt.Println(caneta.Preco) // 10 😕

	darDescontoPtr(&caneta, 10)
	fmt.Println(caneta.Preco) // 9 ✅
}
```

### Construtor: função `NovoX` que retorna ponteiro

Go não tem construtores, mas existe a **convenção** de uma função `NewX` (ou `NovoX`) que cria e valida:

```go
func NovoProduto(nome string, preco float64) (*Produto, error) {
	if preco < 0 {
		return nil, errors.New("preço não pode ser negativo")
	}
	return &Produto{Nome: nome, Preco: preco}, nil
}

p, err := NovoProduto("Caderno", 25)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(p.Nome) // Caderno
```

---

## 📚 Slice de valores vs. slice de ponteiros

### `[]Produto`: o `range` entrega **cópias**

```go
produtos := []Produto{
	{"Caneta", 2.5, 10},
	{"Lápis", 1.2, 0},
}

for _, p := range produtos {
	p.Estoque = 100 // ❌ altera a cópia
}
fmt.Println(produtos[0].Estoque) // 10 😕
```

Soluções:

```go
// ✅ 1. Usar o índice
for i := range produtos {
	produtos[i].Estoque = 100
}

// ✅ 2. Pegar o endereço do elemento
for i := range produtos {
	p := &produtos[i]
	p.Estoque = 100
}
```

### `[]*Produto`: cada item já é um ponteiro

```go
produtos := []*Produto{
	{Nome: "Caneta", Preco: 2.5},  // Go entende &Produto{...}
	{Nome: "Lápis", Preco: 1.2},
}

for _, p := range produtos {
	p.Estoque = 100 // ✅ p é uma cópia do PONTEIRO, aponta pro mesmo produto
}
fmt.Println(produtos[0].Estoque) // 100
```

> ⚠️ **Cuidado:** guardar `&produtos[i]` e depois dar `append` no slice pode deixar o ponteiro apontando para o **array antigo** (se o `append` realocar). Com `[]*Produto` isso não acontece.

---

## 🗺️ Map de valores vs. map de ponteiros

Lembra da pegadinha do módulo 06?

```go
estoque := map[string]Produto{"p1": {"Caneta", 2.5, 10}}
// estoque["p1"].Estoque = 5 // ❌ erro: cannot assign to struct field in map
```

Com ponteiros, funciona direto:

```go
estoque := map[string]*Produto{
	"p1": {Nome: "Caneta", Preco: 2.5, Estoque: 10},
}
estoque["p1"].Estoque = 5 // ✅
fmt.Println(estoque["p1"].Estoque) // 5
```

> ⚠️ Mas se a chave **não existir**, `estoque["xyz"]` retorna `nil` → acessar `.Estoque` dá **panic**. Use o *comma ok*:

```go
if p, ok := estoque["xyz"]; ok {
	p.Estoque--
}
```

> 💡 O seu `main.go` usa `map[string]Produto`. Um bom exercício é refatorar para `map[string]*Produto` e criar uma opção "Vender produto" que diminua o estoque.

---

## ❓ Ponteiros para representar "valor opcional"

Com `int`, não dá pra diferenciar **"não informado"** de **zero**. Com `*int`, dá:

```go
type Filtro struct {
	PrecoMax *float64 // nil = sem filtro de preço
}

func filtrar(produtos []Produto, f Filtro) []Produto {
	var res []Produto
	for _, p := range produtos {
		if f.PrecoMax != nil && p.Preco > *f.PrecoMax {
			continue
		}
		res = append(res, p)
	}
	return res
}

limite := 5.0
baratos := filtrar(produtos, Filtro{PrecoMax: &limite})
todos := filtrar(produtos, Filtro{})
```

Isso é muito usado com **JSON** e **bancos de dados** (campo ausente vs. campo com zero).

---

## 🔁 Ponteiro para variável de laço (Go 1.22+)

```go
nums := []int{1, 2, 3}
var ptrs []*int

for _, n := range nums {
	ptrs = append(ptrs, &n)
}

for _, p := range ptrs {
	fmt.Print(*p, " ")
}
// Go 1.22+: 1 2 3   (cada volta tem um n novo)
// Go 1.21-: 3 3 3   (todas apontavam para o MESMO n)
```

---

## ✍️ Exercícios

1. Crie a struct `Conta{Titular string; Saldo float64}` e as funções `depositar(c *Conta, v float64)` e `sacar(c *Conta, v float64) error` (erro se o saldo for insuficiente).
2. Crie `NovaConta(titular string) *Conta` e use com as funções acima.
3. Dado `[]Produto`, zere o estoque de todos usando `range` **corretamente**.
4. Crie um `map[string]*Conta` e faça uma transferência entre duas contas.
5. Refatore o seu `main.go` da livraria para usar `map[string]*Produto` e adicione a opção **"Vender produto"** que diminui o estoque (sem deixar ficar negativo).
6. Crie `type Atualizacao struct { Nome *string; Preco *float64 }` e uma função que aplique **só os campos não-nil** a um `*Produto`.

---

⬅️ Anterior: [Ponteiros e funções](02-ponteiros-e-funcoes.md) · ➡️ Próximo: [Stack, heap e boas práticas](04-stack-heap-e-boas-praticas.md)

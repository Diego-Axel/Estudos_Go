# 01 — Structs

Uma **struct** agrupa vários valores (os **campos**) em um único tipo. É a forma de Go modelar "coisas" do mundo real: um produto, um usuário, um pedido...

> Go **não tem classes**. Structs + métodos + interfaces fazem esse papel.

---

## 📝 Declarando

```go
type Produto struct {
	Codigo  string
	Titulo  string
	Autor   string
	Preco   float64
	Estoque int
}
```

(Essa é exatamente a struct do seu `main.go` da livraria!)

Campos do mesmo tipo podem ser agrupados:

```go
type Ponto struct {
	X, Y float64
}
```

---

## 🏗️ Criando valores (literais)

```go
// 1. Com nomes dos campos ⭐ (recomendado)
p1 := Produto{
	Codigo: "L001",
	Titulo: "O Hobbit",
	Preco:  49.90,
} // campos omitidos ficam com valor zero (Autor "", Estoque 0)

// 2. Posicional: TODOS os campos, na ordem
p2 := Produto{"L002", "Duna", "Frank Herbert", 79.90, 5}

// 3. Valor zero
var p3 Produto // {"" "" "" 0 0}

// 4. Ponteiro
p4 := &Produto{Codigo: "L003"}
```

> ⚠️ Evite a forma **posicional**: se alguém adicionar um campo na struct, todo o código quebra (ou pior, fica com valores trocados). O `go vet` avisa quando você usa literal posicional de struct de **outro pacote**.

---

## 🎯 Acessando e alterando campos

```go
p := Produto{Titulo: "O Hobbit", Preco: 49.90}

fmt.Println(p.Titulo) // O Hobbit
p.Estoque = 10
p.Preco *= 0.9        // 10% de desconto

fmt.Printf("%v\n", p)  // { O Hobbit  44.91 10}
fmt.Printf("%+v\n", p) // {Codigo: Titulo:O Hobbit Autor: Preco:44.91 Estoque:10}
```

> 💡 `%+v` mostra os nomes dos campos: ótimo pra debug.

---

## 📋 Structs são **valores** (são copiadas)

```go
a := Ponto{1, 2}
b := a      // cópia completa
b.X = 99

fmt.Println(a) // {1 2}
fmt.Println(b) // {99 2}
```

Para compartilhar/alterar, use **ponteiros** (visto no [módulo 08](../08-ponteiros/03-ponteiros-com-structs-e-colecoes.md)).

---

## ⚖️ Comparando structs

Se **todos os campos** forem comparáveis, dá pra usar `==`:

```go
a := Ponto{1, 2}
b := Ponto{1, 2}
fmt.Println(a == b) // true
```

Se tiver um slice, map ou função como campo, **não** compila:

```go
type Pedido struct {
	Itens []string
}
// fmt.Println(Pedido{} == Pedido{}) // ❌ erro: struct containing []string cannot be compared
```

Structs comparáveis podem ser **chave de map**:

```go
visitas := map[Ponto]int{}
visitas[Ponto{0, 0}]++
```

---

## 🪆 Structs aninhadas

```go
type Endereco struct {
	Rua    string
	Cidade string
	UF     string
}

type Cliente struct {
	Nome     string
	Endereco Endereco   // struct dentro de struct
	Contatos []string
}

c := Cliente{
	Nome: "Ana",
	Endereco: Endereco{
		Rua:    "Rua das Flores, 10",
		Cidade: "Natal",
		UF:     "RN",
	},
	Contatos: []string{"ana@email.com"},
}

fmt.Println(c.Endereco.Cidade) // Natal
```

---

## 👻 Structs anônimas

Structs **sem nome de tipo**, para uso rápido e local:

```go
config := struct {
	Porta int
	Debug bool
}{
	Porta: 8080,
	Debug: true,
}
fmt.Println(config.Porta) // 8080
```

Muito usadas em **testes** (table tests, módulo 15) e para montar JSON rapidinho:

```go
casos := []struct {
	entrada  int
	esperado int
}{
	{2, 4},
	{3, 9},
}
```

---

## 🔓 Campos exportados vs. não exportados

Mesma regra de sempre: **maiúscula = visível fora do pacote**.

```go
type Usuario struct {
	Nome  string // exportado
	Email string // exportado
	senha string // privado: só o pacote atual acessa
}
```

> ⚠️ Pacotes como `encoding/json` **só enxergam campos exportados**. Um campo `nome` (minúsculo) **não** aparece no JSON!

---

## 🕳️ Struct vazia: `struct{}`

Uma struct sem campos ocupa **0 bytes**. É usada como "sinal" ou em *sets*:

```go
conjunto := map[string]struct{}{}
conjunto["go"] = struct{}{}

pronto := make(chan struct{}) // canal só para sinalizar (Concorrência)
```

---

## 🧭 `new` vs. `&T{}`

```go
p1 := new(Produto)      // *Produto com todos os campos zerados
p2 := &Produto{}        // idem, mais comum
p3 := &Produto{Preco: 10} // já com valores ✅ o mais usado
```

---

## ✍️ Exercícios

1. Crie a struct `Livro` com `Titulo`, `Autor`, `Ano` e `Paginas`. Crie 3 livros e imprima com `%+v`.
2. Crie um `[]Livro` e mostre o livro com **mais páginas**.
3. Crie as structs `Aluno` (Nome, Notas `[]float64`) e `Turma` (Nome, Alunos `[]Aluno`). Calcule a média de cada aluno.
4. Mostre, com código, que structs são copiadas na atribuição.
5. Crie a struct `Coordenada{Lat, Lng float64}` e use como chave de um `map[Coordenada]string` com nomes de lugares.
6. Crie uma struct anônima com os dados de conexão de um banco (host, porta, usuário) e imprima.

---

🏠 [Módulo 09](README.md) · ➡️ Próximo: [Métodos](02-metodos.md)

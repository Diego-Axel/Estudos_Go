# 02 — Métodos

Um **método** é uma função **ligada a um tipo**. Ele tem um parâmetro especial, o **receptor** (*receiver*), que fica **antes** do nome:

```go
func (receptor Tipo) NomeDoMetodo(parametros) retorno {
	// ...
}
```

---

## 🧩 Primeiro método

```go
type Retangulo struct {
	Largura, Altura float64
}

func (r Retangulo) Area() float64 {
	return r.Largura * r.Altura
}

func (r Retangulo) Perimetro() float64 {
	return 2 * (r.Largura + r.Altura)
}

func main() {
	ret := Retangulo{Largura: 3, Altura: 4}
	fmt.Println(ret.Area())      // 12
	fmt.Println(ret.Perimetro()) // 14
}
```

> 📌 Convenção: o nome do receptor é **curto** (1 ou 2 letras, geralmente a inicial do tipo): `r`, `p`, `c`. Nada de `this` ou `self`.

É equivalente a uma função comum, só que mais organizado:

```go
func Area(r Retangulo) float64 { ... }  // função
func (r Retangulo) Area() float64 { ... } // método
```

---

## 🎯 Receptor por **valor** vs. por **ponteiro** ⭐

### Receptor por valor: recebe uma **cópia**

```go
type Contador struct {
	Valor int
}

func (c Contador) Incrementar() {
	c.Valor++ // altera a CÓPIA
}

func main() {
	c := Contador{}
	c.Incrementar()
	fmt.Println(c.Valor) // 0 😕
}
```

### Receptor por ponteiro: altera o **original**

```go
func (c *Contador) Incrementar() {
	c.Valor++ // altera o original
}

func main() {
	c := Contador{}
	c.Incrementar() // Go faz (&c).Incrementar() automaticamente
	c.Incrementar()
	fmt.Println(c.Valor) // 2 ✅
}
```

### Qual usar?

| Use **ponteiro** (`*T`) quando... | Use **valor** (`T`) quando... |
|---|---|
| o método **altera** o receptor | o método só **lê** |
| a struct é **grande** | a struct é **pequena** (ex: `Ponto`, `Cor`) |
| a struct contém `sync.Mutex` ou algo que não deve ser copiado | o tipo é imutável por natureza (ex: `time.Time`) |

> 🧠 **Regra de consistência:** se **algum** método do tipo precisa de ponteiro, use ponteiro em **todos** os métodos desse tipo.

---

## 🪄 Go ajusta `&` e `*` automaticamente

```go
c := Contador{}      // valor
c.Incrementar()      // método com *Contador → Go faz (&c).Incrementar()

p := &Contador{}     // ponteiro
fmt.Println(p.Valor) // campo → Go faz (*p).Valor
```

Mas só funciona se o valor for **endereçável** (uma variável). Isso **não compila**:

```go
Contador{}.Incrementar() // ❌ erro: cannot call pointer method Incrementar on Contador

m := map[string]Contador{"a": {}}
m["a"].Incrementar()     // ❌ erro: cannot call pointer method Incrementar on Contador
```

---

## 🏷️ Métodos em qualquer tipo **seu**

Não precisa ser struct! Qualquer tipo **declarado no seu pacote** pode ter métodos:

```go
type Celsius float64

func (c Celsius) ParaFahrenheit() float64 {
	return float64(c)*9/5 + 32
}

type ListaDeCompras []string

func (l ListaDeCompras) Total() int {
	return len(l)
}

func main() {
	temp := Celsius(100)
	fmt.Println(temp.ParaFahrenheit()) // 212

	lista := ListaDeCompras{"arroz", "feijão"}
	fmt.Println(lista.Total()) // 2
}
```

> ⚠️ Você **não** pode adicionar métodos a tipos de **outros pacotes** nem a tipos nativos (`int`, `string`...). A solução é criar um tipo seu: `type MeuInt int`.

---

## 🖨️ O método `String()`: a interface `Stringer` ⭐

Se o seu tipo tiver um método `String() string`, o `fmt` usa ele automaticamente para imprimir:

```go
type Produto struct {
	Titulo string
	Preco  float64
}

func (p Produto) String() string {
	return fmt.Sprintf("%s (R$ %.2f)", p.Titulo, p.Preco)
}

func main() {
	p := Produto{"O Hobbit", 49.9}
	fmt.Println(p) // O Hobbit (R$ 49.90)
}
```

### Aplicando nos "enums" do módulo 02

```go
type Status int

const (
	Pendente Status = iota
	Pago
	Enviado
	Entregue
)

func (s Status) String() string {
	switch s {
	case Pendente:
		return "Pendente"
	case Pago:
		return "Pago"
	case Enviado:
		return "Enviado"
	case Entregue:
		return "Entregue"
	}
	return fmt.Sprintf("Status(%d)", int(s))
}

func main() {
	s := Enviado
	fmt.Println(s)               // Enviado (em vez de 2!)
	fmt.Printf("%v %d\n", s, s)  // Enviado 2
}
```

### ⚠️ Pegadinhas do `String()`

**1. Recursão infinita:** não use `%v` ou `%s` do próprio receptor dentro de `String()`:

```go
func (p Produto) String() string {
	return fmt.Sprintf("%v", p) // 💥 chama String() de novo... para sempre
}
```

**2. Receptor por ponteiro:** se `String()` tiver receptor `*Produto`, só **ponteiros** usam ele:

```go
func (p *Produto) String() string { ... }

p := Produto{"Duna", 79.9}
fmt.Println(p)  // {Duna 79.9}      ← não usou String()!
fmt.Println(&p) // Duna (R$ 79.90)  ← usou
```

> ✅ Para `String()`, prefira **receptor por valor**.

---

## 🧮 Exemplo completo: conta bancária

```go
type Conta struct {
	Titular string
	saldo   float64 // privado: só muda pelos métodos
}

func (c *Conta) Depositar(valor float64) error {
	if valor <= 0 {
		return errors.New("valor de depósito inválido")
	}
	c.saldo += valor
	return nil
}

func (c *Conta) Sacar(valor float64) error {
	if valor > c.saldo {
		return fmt.Errorf("saldo insuficiente: disponível R$ %.2f", c.saldo)
	}
	c.saldo -= valor
	return nil
}

func (c *Conta) Saldo() float64 {
	return c.saldo
}

func main() {
	conta := &Conta{Titular: "Diego"}
	conta.Depositar(100)

	if err := conta.Sacar(150); err != nil {
		fmt.Println("Erro:", err) // Erro: saldo insuficiente: disponível R$ 100.00
	}
	fmt.Println(conta.Saldo()) // 100
}
```

> 📌 Getter em Go **não** leva "Get": o método é `Saldo()`, não `GetSaldo()`. O setter, quando existe, é `SetSaldo()`.

---

## 🔗 Métodos como valores

```go
ret := Retangulo{3, 4}
f := ret.Area        // "method value": já vinculado a ret
fmt.Println(f())     // 12

g := Retangulo.Area  // "method expression": recebe o receptor como 1º parâmetro
fmt.Println(g(ret))  // 12
```

---

## ✍️ Exercícios

1. Crie `Circulo{Raio float64}` com os métodos `Area()` e `Perimetro()`.
2. Crie `Conta` com `Depositar`, `Sacar` e `Extrato()` (que retorna as últimas operações, guardadas num slice).
3. Adicione um método `String()` na sua struct `Livro` do exercício anterior.
4. Crie `type Temperatura float64` com os métodos `Celsius()`, `Fahrenheit()` e `Kelvin()`.
5. Crie o "enum" `DiaDaSemana` com `iota` e um método `String()`, e um método `FimDeSemana() bool`.
6. Explique por que `Contador{}.Incrementar()` não compila quando o receptor é `*Contador`.
7. No seu `main.go`, transforme as funções `cadastrarProduto`, `listarProdutos` e `buscarProduto` em **métodos** de `*Livraria`.

---

⬅️ Anterior: [Structs](01-structs.md) · ➡️ Próximo: [Embedding e composição](03-embedding-e-composicao.md)

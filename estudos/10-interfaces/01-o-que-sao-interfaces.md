# 01 — O que são Interfaces?

Uma **interface** define um **comportamento**: uma lista de **métodos** que um tipo precisa ter. Ela diz **o que** algo faz, não **como** faz, nem **o que** ele é.

```go
type Forma interface {
	Area() float64
	Perimetro() float64
}
```

Leia como: *"qualquer tipo que tenha `Area() float64` e `Perimetro() float64` é uma `Forma`"*.

---

## ✨ Implementação **implícita** ⭐

Em Java você escreve `class Circulo implements Forma`. Em Go, **não existe `implements`**: se o tipo tem os métodos, ele **automaticamente** satisfaz a interface.

```go
type Retangulo struct {
	Largura, Altura float64
}

func (r Retangulo) Area() float64      { return r.Largura * r.Altura }
func (r Retangulo) Perimetro() float64 { return 2 * (r.Largura + r.Altura) }

type Circulo struct {
	Raio float64
}

func (c Circulo) Area() float64      { return math.Pi * c.Raio * c.Raio }
func (c Circulo) Perimetro() float64 { return 2 * math.Pi * c.Raio }
```

`Retangulo` e `Circulo` **nunca mencionam** `Forma`, mas **são** `Forma`s. 🦆

> 🦆 *"Se anda como um pato e grasna como um pato, então é um pato."* (*duck typing*, mas verificado **em tempo de compilação**)

Vantagem enorme: você pode criar uma interface para tipos que **já existem**, até de outros pacotes, sem mexer neles.

---

## 🎭 Polimorfismo

Agora, uma função pode trabalhar com **qualquer** `Forma`:

```go
func descrever(f Forma) {
	fmt.Printf("%T → área %.2f, perímetro %.2f\n", f, f.Area(), f.Perimetro())
}

func main() {
	descrever(Retangulo{3, 4})
	descrever(Circulo{1})
}
// main.Retangulo → área 12.00, perímetro 14.00
// main.Circulo → área 3.14, perímetro 6.28
```

### Slice de interfaces

```go
formas := []Forma{
	Retangulo{3, 4},
	Circulo{2},
	Retangulo{1, 1},
}

total := 0.0
for _, f := range formas {
	total += f.Area()
}
fmt.Printf("Área total: %.2f\n", total) // Área total: 25.57
```

Adicionou um `Triangulo` com os dois métodos? Ele entra no slice **sem mudar nada** no resto do código.

---

## 🧬 O que tem dentro de uma variável de interface?

Uma variável de interface guarda **dois** dados:

```
var f Forma = Circulo{2}

┌──────────────────────────┐
│ tipo:  main.Circulo      │  ← tipo dinâmico
│ valor: {Raio: 2}         │  ← valor dinâmico
└──────────────────────────┘
```

```go
var f Forma
fmt.Printf("%T %v\n", f, f) // <nil> <nil>

f = Circulo{2}
fmt.Printf("%T %v\n", f, f) // main.Circulo {2}

f = Retangulo{3, 4}
fmt.Printf("%T %v\n", f, f) // main.Retangulo {3 4}
```

Pela interface, você **só** pode chamar os métodos **da interface**:

```go
var f Forma = Circulo{2}
fmt.Println(f.Area()) // ✅
// fmt.Println(f.Raio) // ❌ erro: f.Raio undefined (type Forma has no field or method Raio)
```

Para acessar o `Raio`, é preciso "extrair" o tipo concreto (*type assertion*, no próximo arquivo).

---

## ⚠️ Receptor por ponteiro e interfaces

Se um método tem receptor **por ponteiro**, só o **ponteiro** satisfaz a interface:

```go
type Falante interface {
	Falar() string
}

type Cachorro struct{ Nome string }

func (c *Cachorro) Falar() string { return c.Nome + ": Au au!" }

func main() {
	d := Cachorro{"Rex"}

	// var f Falante = d  // ❌ erro: Cachorro does not implement Falante (method Falar has pointer receiver)
	var f Falante = &d    // ✅
	fmt.Println(f.Falar()) // Rex: Au au!
}
```

| Métodos definidos com... | `T` satisfaz? | `*T` satisfaz? |
|---|---|---|
| receptor valor `(t T)` | ✅ | ✅ |
| receptor ponteiro `(t *T)` | ❌ | ✅ |

> 🧠 Motivo: um valor guardado numa interface **não é endereçável**, então Go não consegue fazer o `&` automático.

---

## ✅ Verificando a implementação em tempo de compilação

Como a implementação é implícita, você pode "garantir" que um tipo satisfaz uma interface com este truque:

```go
var _ Forma = Retangulo{}       // falha na compilação se Retangulo não for Forma
var _ Falante = (*Cachorro)(nil) // versão para receptor ponteiro
```

É comum ver isso logo abaixo da declaração do tipo, em bibliotecas.

---

## 🧪 Exemplo prático: sistema de pagamentos

```go
type MeioDePagamento interface {
	Pagar(valor float64) error
}

type Pix struct{ Chave string }
type Cartao struct {
	Numero string
	Limite float64
}

func (p Pix) Pagar(valor float64) error {
	fmt.Printf("Pix de R$ %.2f para %s\n", valor, p.Chave)
	return nil
}

func (c *Cartao) Pagar(valor float64) error {
	if valor > c.Limite {
		return fmt.Errorf("limite insuficiente (R$ %.2f)", c.Limite)
	}
	c.Limite -= valor
	fmt.Printf("Cartão final %s: R$ %.2f\n", c.Numero[len(c.Numero)-4:], valor)
	return nil
}

func finalizarCompra(m MeioDePagamento, total float64) {
	if err := m.Pagar(total); err != nil {
		fmt.Println("Falha:", err)
		return
	}
	fmt.Println("Compra aprovada ✅")
}

func main() {
	finalizarCompra(Pix{"diego@email.com"}, 50)
	finalizarCompra(&Cartao{"1234567812345678", 100}, 150)
}
// Pix de R$ 50.00 para diego@email.com
// Compra aprovada ✅
// Falha: limite insuficiente (R$ 100.00)
```

`finalizarCompra` **não sabe** (nem precisa saber) qual meio de pagamento está usando. Amanhã você cria `Boleto` e nada nela muda.

---

## ✍️ Exercícios

1. Crie a interface `Forma` e os tipos `Quadrado`, `Circulo` e `Triangulo` (use a fórmula de Heron para a área). Some as áreas de um `[]Forma`.
2. Crie a interface `Animal` com `Som() string` e `Nome() string`. Implemente 3 animais e faça um "coral" percorrendo um slice.
3. Crie `Cachorro` com `Falar()` de receptor **ponteiro** e tente atribuir um `Cachorro` (valor) a uma variável `Falante`. Leia o erro e corrija.
4. Adicione `var _ Forma = Quadrado{}` ao seu código. Depois remova o método `Area` de `Quadrado` e veja o erro.
5. Crie a interface `Notificador` com `Enviar(msg string) error`, e as implementações `Email`, `SMS` e `Console`. Faça uma função que envie para uma lista de notificadores.
6. Adicione um `Boleto` ao exemplo de pagamentos sem alterar `finalizarCompra`.

---

🏠 [Módulo 10](README.md) · ➡️ Próximo: [any, type assertion e type switch](02-any-type-assertion-type-switch.md)

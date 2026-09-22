# 02 — `any`, Type Assertion e Type Switch

## 🌌 A interface vazia: `any`

Uma interface **sem métodos** é satisfeita por **todos** os tipos (todo tipo tem "pelo menos zero métodos"):

```go
interface{} // forma antiga
any         // apelido, desde o Go 1.18 ⭐
```

```go
var x any

x = 42
x = "texto"
x = []int{1, 2, 3}
x = Circulo{2}

fmt.Printf("%T\n", x) // main.Circulo
```

Por isso `fmt.Println` aceita qualquer coisa: a assinatura é `func Println(a ...any)`.

Exemplos de uso:

```go
dados := map[string]any{
	"nome":   "Ana",
	"idade":  30,
	"ativo":  true,
	"tags":   []string{"go", "dev"},
}
```

> ⚠️ **Use `any` com moderação.** Ao usar `any`, você **perde a verificação de tipos** do compilador e passa a precisar de *type assertions* em tempo de execução. Prefira interfaces com métodos ou **Generics** (módulo 13).

---

## 🔍 Type assertion: "extraindo" o tipo concreto

Sintaxe: `valor.(Tipo)`

```go
var x any = "Olá, Go"

s := x.(string)       // "eu garanto que x é uma string"
fmt.Println(len(s))   // 7
```

### 💥 Se estiver errado: panic

```go
var x any = "Olá"
n := x.(int)
// panic: interface conversion: interface {} is string, not int
```

### ✅ Forma segura: *comma ok*

```go
if n, ok := x.(int); ok {
	fmt.Println("é int:", n)
} else {
	fmt.Println("não é int") // ← cai aqui
}
```

> Mesmo idioma dos maps: `v, ok := m[chave]`. Se `ok` for `false`, `n` recebe o **valor zero** do tipo e **não** há panic.

### Com interfaces "normais"

```go
var f Forma = Circulo{2}

if c, ok := f.(Circulo); ok {
	fmt.Println("raio:", c.Raio) // raio: 2
}
```

### Verificando se um valor tem **outro comportamento**

Type assertion também funciona com **interfaces**, perguntando "esse valor também sabe fazer X?":

```go
type Volume interface {
	Volume() float64
}

func info(f Forma) {
	fmt.Printf("área: %.2f\n", f.Area())
	if v, ok := f.(Volume); ok { // "essa forma TAMBÉM tem volume?"
		fmt.Printf("volume: %.2f\n", v.Volume())
	}
}
```

A biblioteca padrão usa muito esse truque (ex: `fmt` verificando se o valor é um `Stringer`).

---

## 🔀 Type switch

Quando há **vários tipos** possíveis, use um `switch` especial com `.(type)`:

```go
func descrever(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int: %d (dobro %d)", x, x*2)
	case float64:
		return fmt.Sprintf("float64: %.2f", x)
	case string:
		return fmt.Sprintf("string de %d bytes: %q", len(x), x)
	case bool:
		return fmt.Sprintf("bool: %t", x)
	case []int:
		return fmt.Sprintf("slice de %d ints", len(x))
	case error:
		return "erro: " + x.Error()
	case fmt.Stringer:
		return "Stringer: " + x.String()
	default:
		return fmt.Sprintf("tipo desconhecido: %T", x)
	}
}

func main() {
	fmt.Println(descrever(42))
	fmt.Println(descrever(3.14))
	fmt.Println(descrever("Go"))
	fmt.Println(descrever([]int{1, 2}))
	fmt.Println(descrever(errors.New("falhou")))
	fmt.Println(descrever(nil))
	fmt.Println(descrever(Circulo{1}))
}
// int: 42 (dobro 84)
// float64: 3.14
// string de 2 bytes: "Go"
// slice de 2 ints
// erro: falhou
// nil
// tipo desconhecido: main.Circulo
```

Detalhes:
- Dentro de cada `case`, `x` **já tem o tipo** daquele case (`int`, `string`...)
- Os `case` podem ser tipos concretos **ou interfaces** (`error`, `fmt.Stringer`)
- A ordem importa: o **primeiro** que bater vence
- `case int, int64:` com vários tipos → `x` continua sendo `any` (não dá pra saber qual)
- **Não** existe `fallthrough` em type switch

---

## 🧪 Exemplo: somando valores de tipos variados

```go
func somar(valores ...any) (float64, error) {
	total := 0.0
	for _, v := range valores {
		switch n := v.(type) {
		case int:
			total += float64(n)
		case float64:
			total += n
		case string:
			f, err := strconv.ParseFloat(n, 64)
			if err != nil {
				return 0, fmt.Errorf("valor inválido %q", n)
			}
			total += f
		default:
			return 0, fmt.Errorf("tipo não suportado: %T", v)
		}
	}
	return total, nil
}

fmt.Println(somar(1, 2.5, "3.5")) // 7 <nil>
fmt.Println(somar(1, true))       // 0 tipo não suportado: bool
```

---

## ⚖️ Comparando valores de interface

Dois valores de interface são iguais se têm o **mesmo tipo dinâmico** e **valores iguais**:

```go
var a any = 10
var b any = 10
var c any = int64(10)

fmt.Println(a == b) // true
fmt.Println(a == c) // false (int ≠ int64)
```

> ⚠️ Se o tipo dinâmico **não for comparável** (slice, map, função), o `==` compila, mas dá **panic** ao executar:

```go
var x any = []int{1}
var y any = []int{1}
fmt.Println(x == y) // 💥 panic: runtime error: comparing uncomparable type []int
```

---

## ✍️ Exercícios

1. Crie um `[]any` com valores de 5 tipos diferentes e imprima o tipo de cada um com `%T`.
2. Escreva `paraString(v any) string` usando **type switch**, tratando `int`, `float64`, `bool`, `string` e `nil`.
3. Faça uma type assertion **errada** sem *comma ok* e leia o panic. Depois corrija com *comma ok*.
4. Dado um `[]Forma`, conte quantos são `Circulo` e quantos são `Retangulo`.
5. Crie a interface `Descontavel` com `AplicarDesconto(pct float64)`. Num carrinho de `[]Produto` (interface), aplique desconto **só** nos itens que também são `Descontavel`.
6. Por que usar `any` em todo lugar é uma má ideia? Cite o que se perde.

---

⬅️ Anterior: [O que são interfaces?](01-o-que-sao-interfaces.md) · ➡️ Próximo: [Composição e interfaces famosas](03-composicao-e-interfaces-famosas.md)

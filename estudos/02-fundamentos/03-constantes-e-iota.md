# 03 — Constantes e `iota`

## 🔒 O que é uma constante?

Um valor **fixo**, definido em **tempo de compilação**, que **nunca muda**.

```go
const Pi = 3.14159
const Saudacao string = "Olá"

// Pi = 3 // ❌ erro: cannot assign to Pi
```

Constantes podem ser declaradas **dentro ou fora** de funções.

### Em bloco

```go
const (
	NomeApp  = "Livraria"
	Versao   = "1.0.0"
	MaxItens = 100
)
```

---

## ⚠️ O que pode ser constante?

Só valores que o compilador consegue calcular **antes** do programa rodar:

- números, strings, booleanos (e runes)
- expressões com eles: `const Dia = 24 * 60 * 60`

**Não podem** ser constantes: slices, maps, structs, resultado de funções comuns...

```go
const hoje = time.Now()    // ❌ erro: não é conhecido em tempo de compilação
const lista = []int{1, 2}  // ❌ erro
const tamanho = len("Go")  // ✅ ok! len de string constante é resolvido na compilação
```

---

## 🧪 Constantes *untyped* (sem tipo)

Esse é um detalhe **muito legal** de Go. Uma constante declarada sem tipo é **"sem tipo"** (*untyped*) e se adapta ao contexto:

```go
const Dez = 10 // untyped

var a int = Dez      // ✅ vira int
var b float64 = Dez  // ✅ vira float64
var c int64 = Dez    // ✅ vira int64
```

Compare com uma **variável** ou constante **tipada**:

```go
const DezInt int = 10
var d float64 = DezInt // ❌ erro: cannot use DezInt (constant of type int) as float64
```

Constantes *untyped* também têm **precisão altíssima** durante os cálculos do compilador:

```go
const Enorme = 1 << 100          // ✅ (seria overflow em qualquer int)
const Pequeno = Enorme >> 98     // 4
fmt.Println(Pequeno)             // 4
// fmt.Println(Enorme)           // ❌ erro: overflows int (ao virar valor de verdade)
```

---

## 🔢 `iota`: o gerador de sequências

Dentro de um bloco `const`, o **`iota`** vale `0` na primeira linha e **aumenta 1 a cada linha**.

```go
const (
	A = iota // 0
	B = iota // 1
	C = iota // 2
)
```

Se você omitir a expressão, Go **repete a anterior**, então basta escrever uma vez:

```go
const (
	Domingo = iota // 0
	Segunda        // 1
	Terca          // 2
	Quarta         // 3
	Quinta         // 4
	Sexta          // 5
	Sabado         // 6
)
```

> 🔁 O `iota` **reinicia em 0** a cada novo bloco `const`.

### Começando em 1 (pulando o zero)

```go
const (
	_       = iota // 0 descartado
	Janeiro        // 1
	Fevereiro      // 2
	Marco          // 3
)
```

Ou:

```go
const (
	Janeiro = iota + 1 // 1
	Fevereiro          // 2
	Marco              // 3
)
```

### Criando um "enum" tipado ⭐

Go não tem `enum`, mas esse padrão faz o mesmo papel:

```go
type Status int

const (
	Pendente Status = iota // 0
	Pago                   // 1
	Enviado                // 2
	Entregue               // 3
)

func main() {
	var s Status = Enviado
	fmt.Println(s) // 2
}
```

Mais pra frente (no módulo de Interfaces) você vai ver como fazer `fmt.Println(s)` imprimir `"Enviado"` com o método `String()`.

### Potências (tamanhos de arquivo)

```go
const (
	_  = iota             // ignora 0
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20
	GB                    // 1 << 30
	TB                    // 1 << 40
)

fmt.Println(KB, MB, GB) // 1024 1048576 1073741824
```

### Flags de bits (permissões)

```go
const (
	Ler      = 1 << iota // 1  (001)
	Escrever             // 2  (010)
	Executar             // 4  (100)
)

permissao := Ler | Escrever          // 3 (011)
podeEscrever := permissao&Escrever != 0
fmt.Println(podeEscrever)            // true
```

> Os operadores `<<`, `|` e `&` serão explicados no módulo de **Operadores**.

---

## ✍️ Exercícios

1. Crie as constantes `Gravidade = 9.8` e `NomePlaneta = "Terra"` e imprima.
2. Tente atribuir um novo valor a uma constante. Qual o erro?
3. Por que `const agora = time.Now()` não compila?
4. Crie um "enum" `Nivel` com os valores `Iniciante`, `Intermediario`, `Avancado`, **começando em 1**.
5. Usando `iota`, crie constantes para os tamanhos de camiseta `P`, `M`, `G`, `GG` e imprima todos.
6. Explique por que isto funciona:
   ```go
   const x = 5
   var f float64 = x
   ```
   mas isto não:
   ```go
   const x int = 5
   var f float64 = x
   ```

---

⬅️ Anterior: [Tipos de dados](02-tipos-de-dados.md) · ➡️ Próximo: [Conversão de tipos](04-conversao-de-tipos.md)

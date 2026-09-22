# 02 — Operadores de Comparação e Lógicos

## ⚖️ Operadores de comparação (relacionais)

Sempre retornam um **`bool`** (`true` ou `false`).

| Operador | Significado | Exemplo | Resultado |
|---|---|---|---|
| `==` | igual | `5 == 5` | `true` |
| `!=` | diferente | `5 != 3` | `true` |
| `<` | menor | `3 < 5` | `true` |
| `<=` | menor ou igual | `5 <= 5` | `true` |
| `>` | maior | `3 > 5` | `false` |
| `>=` | maior ou igual | `3 >= 5` | `false` |

```go
idade := 20
maior := idade >= 18
fmt.Println(maior) // true
```

### Só compara tipos iguais

```go
var a int = 5
var b int64 = 5
// fmt.Println(a == b) // ❌ erro: mismatched types int and int64
fmt.Println(int64(a) == b) // ✅ true
```

Constantes *untyped* se adaptam, então `a == 5` funciona normalmente.

### Comparando strings

Strings são comparadas **byte a byte** (ordem lexicográfica):

```go
fmt.Println("go" == "go")      // true
fmt.Println("Go" == "go")      // false (maiúscula ≠ minúscula)
fmt.Println("abc" < "abd")     // true
fmt.Println("Z" < "a")         // true  ('Z' = 90, 'a' = 97)
fmt.Println("10" < "9")        // true  😱 ('1' vem antes de '9')
```

Para comparar **ignorando maiúsculas/minúsculas**:

```go
strings.EqualFold("GoLang", "golang") // true
```

### O que pode e o que não pode ser comparado com `==`

| Tipo | `==` funciona? |
|---|---|
| números, strings, bool | ✅ |
| ponteiros, channels | ✅ (compara se apontam pro mesmo lugar) |
| arrays | ✅ se os elementos forem comparáveis |
| structs | ✅ se **todos** os campos forem comparáveis |
| **slices, maps, funções** | ❌ só podem ser comparados com `nil` |

```go
s1 := []int{1, 2}
s2 := []int{1, 2}
// fmt.Println(s1 == s2) // ❌ erro: slice can only be compared to nil
fmt.Println(slices.Equal(s1, s2)) // ✅ true (pacote "slices", Go 1.21+)
```

### ⚠️ Cuidado com floats

```go
x, y := 0.1, 0.2
fmt.Println(x+y == 0.3) // false 😱

// Compare com uma margem de tolerância:
const eps = 1e-9
fmt.Println(math.Abs((x+y)-0.3) < eps) // true

// Curiosidade: com CONSTANTES o compilador calcula com precisão exata,
// por isso a linha abaixo dá true. O problema aparece com VARIÁVEIS.
fmt.Println(0.1+0.2 == 0.3) // true

nan := math.NaN()
fmt.Println(nan == nan) // false! NaN nunca é igual a nada
fmt.Println(math.IsNaN(nan)) // true
```

---

## 🧠 Operadores lógicos

Trabalham **apenas com `bool`**.

| Operador | Nome | Verdadeiro quando... |
|---|---|---|
| `&&` | E (AND) | **os dois** lados são `true` |
| `\|\|` | OU (OR) | **pelo menos um** lado é `true` |
| `!` | NÃO (NOT) | inverte o valor |

### Tabela-verdade

| `a` | `b` | `a && b` | `a \|\| b` | `!a` |
|---|---|---|---|---|
| true | true | true | true | false |
| true | false | false | true | false |
| false | true | false | true | true |
| false | false | false | false | true |

```go
idade := 25
temCNH := true

podeDirigir := idade >= 18 && temCNH
fimDeSemana := dia == "sábado" || dia == "domingo"
menor := !(idade >= 18)
```

> ⚠️ Nada de `if 1` ou `if nome` como em Python/JS. A condição **tem** que ser `bool`: use `if nome != ""`.

---

## ⚡ Avaliação em curto-circuito

Go **para de avaliar** assim que já sabe o resultado:

- `a && b`: se `a` é `false`, **nem olha** `b` (o resultado já é `false`)
- `a || b`: se `a` é `true`, **nem olha** `b` (o resultado já é `true`)

Isso é muito usado para **proteger** a segunda condição:

```go
// Se o divisor for zero, a divisão NUNCA é executada (evita o panic)
if divisor != 0 && total/divisor > 10 {
	fmt.Println("média alta")
}

// Se o slice estiver vazio, não acessa lista[0] (evita o panic de índice)
if len(lista) > 0 && lista[0] == "admin" {
	fmt.Println("primeiro é admin")
}

// Se o ponteiro for nil, não acessa o campo
if usuario != nil && usuario.Ativo {
	// ...
}
```

Demonstração:

```go
func verdade() bool {
	fmt.Println("chamou verdade()")
	return true
}

func main() {
	_ = false && verdade() // não imprime nada
	_ = true || verdade()  // não imprime nada
	_ = true && verdade()  // imprime "chamou verdade()"
}
```

---

## 🚫 Não existe operador ternário

Em Go não tem `condicao ? a : b`. Use `if/else`:

```go
// Em outras linguagens: status = idade >= 18 ? "adulto" : "menor"
status := "menor"
if idade >= 18 {
	status = "adulto"
}
```

---

## ✍️ Exercícios

1. Leia a idade de uma pessoa e imprima `true`/`false` para: é maior de idade? É idoso (60+)? Está entre 13 e 17 anos?
2. Um ano é **bissexto** se for divisível por 4 **e** não por 100, **ou** se for divisível por 400. Escreva essa expressão lógica e teste com 1900, 2000, 2024 e 2025.
3. Por que `"10" < "9"` é `true`? Como comparar esses valores **numericamente**?
4. Explique por que este código **não** dá panic mesmo com `lista` vazia:
   ```go
   var lista []int
   if len(lista) > 0 && lista[0] == 1 { }
   ```
   E se trocar a ordem das condições?
5. Com `x, y := 0.1, 0.2`, por que `x+y == 0.3` é `false`? Como comparar floats corretamente?
6. Reescreva sem ternário: `desconto = vip ? 0.2 : 0.05`.

---

⬅️ Anterior: [Aritméticos](01-aritmeticos.md) · ➡️ Próximo: [Bit a bit](03-bit-a-bit.md)

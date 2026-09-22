# 01 — O que são Ponteiros?

## 🏠 A analogia da casa

Toda variável mora em algum lugar da **memória**, e esse lugar tem um **endereço** (tipo `0xc000012080`).

- A **variável** é a **casa**, com coisas dentro (o valor).
- O **ponteiro** é um **papel com o endereço da casa** anotado.

Com o endereço em mãos, você pode ir até a casa e **ver ou mudar** o que tem lá dentro, sem precisar de uma cópia da casa.

```
   variável x                ponteiro p
┌───────────────┐        ┌───────────────┐
│      42       │ ◄───── │ 0xc000012080  │
└───────────────┘        └───────────────┘
 endereço: 0xc000012080
```

---

## 🔣 Os dois operadores

| Operador | Nome | O que faz | Leia como |
|---|---|---|---|
| `&x` | endereço | pega o **endereço** de `x` | "endereço de x" |
| `*p` | desreferência | acessa o **valor** que está no endereço `p` | "valor apontado por p" |

```go
x := 42
p := &x            // p guarda o endereço de x

fmt.Println(x)     // 42          (valor)
fmt.Println(&x)    // 0xc000012080 (endereço, vai variar)
fmt.Println(p)     // 0xc000012080 (mesmo endereço)
fmt.Println(*p)    // 42          (valor no endereço)

*p = 100           // muda o valor NA CASA
fmt.Println(x)     // 100 ← x mudou!
```

---

## 🏷️ O tipo ponteiro: `*T`

Um ponteiro para `int` tem tipo **`*int`**. Para `string`, **`*string`**, e assim por diante:

```go
var p *int         // ponteiro para int
var s *string      // ponteiro para string

x := 10
p = &x
fmt.Printf("%T\n", p) // *int
```

> ⚠️ O `*` tem **dois papéis** diferentes:
> - no **tipo** (`*int`) → "ponteiro para int"
> - na **expressão** (`*p`) → "valor apontado por p"

Um ponteiro só aponta para o **seu tipo**:

```go
x := 10
var p *float64
// p = &x // ❌ erro: cannot use &x (value of type *int) as *float64 value
```

---

## 🕳️ Valor zero: `nil`

Um ponteiro que não aponta para nada vale **`nil`**:

```go
var p *int
fmt.Println(p)        // <nil>
fmt.Println(p == nil) // true
```

### 💥 Desreferenciar `nil` = panic

```go
var p *int
fmt.Println(*p)
// panic: runtime error: invalid memory address or nil pointer dereference
```

É **o erro mais comum** envolvendo ponteiros. Proteja-se:

```go
if p != nil {
	fmt.Println(*p)
}
```

---

## 🆕 `new`: criando um valor e pegando o ponteiro

`new(T)` aloca um valor **zero** do tipo `T` e retorna um **ponteiro** para ele:

```go
p := new(int)     // *int apontando para um int com valor 0
fmt.Println(*p)   // 0
*p = 7
fmt.Println(*p)   // 7
```

É o mesmo que:

```go
var tmp int
p := &tmp
```

> 💡 Na prática, `new` é pouco usado. Com structs, o comum é `&Pessoa{Nome: "Ana"}` (próximos arquivos).

### Ponteiro para um literal?

```go
// p := &42        // ❌ erro: não dá pra pegar o endereço de uma constante
p := new(int)
*p = 42            // ✅ jeito clássico

func ptr[T any](v T) *T { return &v } // ✅ função auxiliar genérica (Generics, módulo 13)
p2 := ptr(42)
```

> 🆕 A partir do **Go 1.26**, `new` também aceita uma expressão: `p := new(42)`.

---

## 🔗 Ponteiro para ponteiro

Possível, mas raro:

```go
x := 5
p := &x    // *int
pp := &p   // **int

**pp = 10
fmt.Println(x) // 10
```

---

## ⚖️ Comparando ponteiros

Dois ponteiros são iguais se apontam para **o mesmo endereço** (não se os valores são iguais):

```go
a, b := 5, 5
pa, pb := &a, &b
pa2 := &a

fmt.Println(pa == pb)   // false (endereços diferentes)
fmt.Println(*pa == *pb) // true  (valores iguais)
fmt.Println(pa == pa2)  // true  (mesmo endereço)
```

---

## 🚫 Sem aritmética de ponteiros

Diferente de C, em Go **não dá** pra fazer `p++` ou `p + 1` para "andar" na memória:

```go
arr := [3]int{1, 2, 3}
p := &arr[0]
// p++ // ❌ erro: invalid operation: p++ (non-numeric type *int)
```

Isso elimina toda uma categoria de bugs e falhas de segurança. (Existe o pacote `unsafe` para casos extremos, mas ele é, como o nome diz, **inseguro**.)

---

## 🖨️ Imprimindo endereços

```go
x := 42
fmt.Printf("%p\n", &x) // 0xc000012080
fmt.Printf("%v\n", &x) // 0xc000012080
```

---

## ✍️ Exercícios

1. Declare `idade := 25`, crie um ponteiro para ela, e use **só o ponteiro** para mudar a idade para 26. Imprima a variável original.
2. Imprima o endereço de três variáveis diferentes com `%p`. Os endereços são próximos?
3. Crie um ponteiro com `new(string)`, atribua "Go" e imprima o valor.
4. O que acontece aqui? Por quê?
   ```go
   var p *int
   *p = 10
   ```
5. Explique a diferença entre `*int` (na declaração) e `*p` (numa expressão).
6. Crie `a := 1` e `b := 1`. Mostre que `&a == &b` é `false` mas `a == b` é `true`.

---

🏠 [Módulo 08](README.md) · ➡️ Próximo: [Ponteiros e funções](02-ponteiros-e-funcoes.md)

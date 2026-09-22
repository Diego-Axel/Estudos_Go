# 04 — Atribuição, Precedência e Outros Operadores

## 📥 Operadores de atribuição

| Operador | Equivale a |
|---|---|
| `=` | atribuição simples |
| `:=` | declaração + atribuição |
| `+=` | `x = x + y` |
| `-=` | `x = x - y` |
| `*=` | `x = x * y` |
| `/=` | `x = x / y` |
| `%=` | `x = x % y` |
| `&=` | `x = x & y` |
| `\|=` | `x = x \| y` |
| `^=` | `x = x ^ y` |
| `&^=` | `x = x &^ y` |
| `<<=` | `x = x << y` |
| `>>=` | `x = x >> y` |

```go
saldo := 100
saldo += 50  // 150
saldo -= 30  // 120
saldo *= 2   // 240
saldo /= 4   // 60
saldo %= 7   // 4

nome := "Go"
nome += "lang" // "Golang"
```

> Assim como `++`/`--`, atribuições são **comandos**, não expressões. `if x = 5 { }` não compila (e ainda bem!).

### Atribuição múltipla

O lado direito é **todo avaliado primeiro**, depois atribuído:

```go
a, b := 1, 2
a, b = b, a+b // a = 2, b = 3  (usa os valores ANTIGOS de a e b)
```

Isso deixa sequências como Fibonacci bem limpas:

```go
a, b := 0, 1
for i := 0; i < 10; i++ {
	fmt.Print(a, " ")
	a, b = b, a+b
}
// 0 1 1 2 3 5 8 13 21 34
```

---

## 🏔️ Precedência de operadores

Go tem só **5 níveis** de precedência para operadores binários (bem mais simples que C):

| Nível | Operadores |
|---|---|
| **5** (maior) | `*` `/` `%` `<<` `>>` `&` `&^` |
| **4** | `+` `-` `\|` `^` |
| **3** | `==` `!=` `<` `<=` `>` `>=` |
| **2** | `&&` |
| **1** (menor) | `\|\|` |

- Operadores **unários** (`!`, `-x`, `^x`, `&x`, `*p`, `<-ch`) têm a **maior prioridade** de todas.
- No mesmo nível, avalia da **esquerda para a direita**.

```go
fmt.Println(2 + 3*4)     // 14  (* antes de +)
fmt.Println((2 + 3) * 4) // 20
fmt.Println(10 - 4 - 3)  // 3   (esquerda → direita)
```

### ⚠️ Diferenças em relação a C/Java/JS

Em Go, `<<`, `>>` e `&` estão no **mesmo nível da multiplicação**, e `|`/`^` no nível da **soma**:

```go
fmt.Println(1 + 2<<3) // 17 → 1 + (2<<3) = 1 + 16
                      // (em C seria (1+2)<<3 = 24!)

x := 6
fmt.Println(x&1 == 0) // true → (x&1) == 0
                      // (em C, & tem precedência MENOR que ==!)
```

Lógicos: `&&` vem antes de `||`:

```go
a, b, c := true, false, false
fmt.Println(a || b && c) // true → a || (b && c)
```

> ✅ **Dica de ouro:** na dúvida, **use parênteses**. O `gofmt` até ajusta os espaços para evidenciar a precedência: `2 + 3*4`.

---

## 🧩 Outros operadores (spoiler dos próximos módulos)

| Operador | Uso | Módulo |
|---|---|---|
| `&x` | pega o **endereço** de `x` | Ponteiros |
| `*p` | acessa o valor apontado por `p` | Ponteiros |
| `<-` | envia/recebe em um **channel** | Concorrência |
| `...` | parâmetro variádico / "espalhar" slice | Funções |
| `.` | acesso a campo/método/pacote | Structs |
| `[]` | índice / fatia | Arrays e Slices |

Uma prévia rápida:

```go
x := 10
p := &x        // p guarda o endereço de x
*p = 20        // altera x através do ponteiro
fmt.Println(x) // 20

ch := make(chan int, 1)
ch <- 42         // envia
v := <-ch        // recebe
fmt.Println(v)   // 42
```

---

## 🧾 Resumão do módulo

```go
// Aritméticos
+  -  *  /  %  ++  --

// Comparação
==  !=  <  <=  >  >=

// Lógicos
&&  ||  !

// Bit a bit
&  |  ^  &^  <<  >>   ^x (NOT)

// Atribuição
=  :=  +=  -=  *=  /=  %=  &=  |=  ^=  &^=  <<=  >>=

// Outros
&x  *p  <-  ...
```

---

## ✍️ Exercícios

1. Sem rodar, diga o resultado e depois confira:
   - `10 + 2*3 - 4/2`
   - `(10 + 2) * (3 - 4) / 2`
   - `1 + 2<<2`
   - `true || false && false`
   - `!true || !false && true`
2. Comece com `x := 10` e aplique, em sequência: `+= 5`, `*= 3`, `-= 5`, `/= 4`, `%= 3`. Qual o valor final?
3. Usando atribuição múltipla, imprima os 15 primeiros números de Fibonacci.
4. Por que `if x = 5 { }` não compila em Go? Que bug comum de C isso evita?
5. Monte uma "calculadora": leia dois números e um operador (`+`, `-`, `*`, `/`, `%`) e mostre o resultado. Trate a divisão por zero.

---

⬅️ Anterior: [Bit a bit](03-bit-a-bit.md) · 🏠 [Voltar ao roteiro](../README.md)

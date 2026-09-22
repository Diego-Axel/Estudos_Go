# 03 — Operadores Bit a Bit (Bitwise)

Esses operadores trabalham **direto nos bits** dos números **inteiros**. Parecem coisa de baixo nível, mas aparecem bastante em flags, permissões, criptografia, hashing e otimizações.

> 💡 Use `%b` no `Printf` para ver os bits: `fmt.Printf("%04b\n", 5)` → `0101`

---

## 📋 Os operadores

| Operador | Nome | Regra (bit a bit) |
|---|---|---|
| `&` | AND | 1 se **os dois** bits forem 1 |
| `\|` | OR | 1 se **pelo menos um** for 1 |
| `^` | XOR | 1 se os bits forem **diferentes** |
| `&^` | AND NOT (bit clear) | zera em `a` os bits que estão ligados em `b` ⭐ exclusivo de Go |
| `<<` | deslocamento à esquerda | empurra os bits para a esquerda |
| `>>` | deslocamento à direita | empurra os bits para a direita |
| `^x` (unário) | NOT | inverte todos os bits |

---

## 🔍 Exemplos com `a = 6` e `b = 3`

```
a = 6  →  0110
b = 3  →  0011
```

### `&` (AND)

```
  0110
& 0011
------
  0010  → 2
```

### `|` (OR)

```
  0110
| 0011
------
  0111  → 7
```

### `^` (XOR)

```
  0110
^ 0011
------
  0101  → 5
```

### `&^` (AND NOT / bit clear)

"Pegue `a` e **desligue** tudo que está ligado em `b`":

```
   0110
&^ 0011
-------
   0100  → 4
```

Em código:

```go
a, b := 6, 3
fmt.Println(a & b)  // 2
fmt.Println(a | b)  // 7
fmt.Println(a ^ b)  // 5
fmt.Println(a &^ b) // 4
```

---

## 🙃 NOT unário: `^x`

Go não tem `~` como C. O NOT bit a bit é o **`^` na frente** do número:

```go
var u uint8 = 5     // 00000101
fmt.Println(^u)     // 250 → 11111010

x := 5
fmt.Println(^x)     // -6 (em inteiros com sinal: ^x == -x - 1)
```

---

## ↔️ Deslocamentos: `<<` e `>>`

```go
fmt.Println(1 << 3)  // 8    (0001 → 1000)
fmt.Println(5 << 1)  // 10   (0101 → 1010)
fmt.Println(16 >> 2) // 4    (10000 → 00100)
fmt.Println(-8 >> 1) // -4   (com sinal: mantém o sinal)
```

**Atalho mental:**
- `x << n` = `x × 2ⁿ`
- `x >> n` = `x ÷ 2ⁿ` (arredondando para baixo)

```go
fmt.Println(1 << 10) // 1024
```

---

## 🛠️ Usos práticos

### 1. Flags / permissões (com `iota`)

```go
const (
	Ler      = 1 << iota // 001 = 1
	Escrever             // 010 = 2
	Executar             // 100 = 4
)

perm := Ler | Escrever // liga Ler e Escrever → 011 = 3

// Verificar se um bit está ligado
fmt.Println(perm&Escrever != 0) // true
fmt.Println(perm&Executar != 0) // false

// Ligar um bit
perm |= Executar  // 111 = 7

// Desligar um bit
perm &^= Escrever // 101 = 5

// Alternar (liga se estiver desligado, desliga se estiver ligado)
perm ^= Ler       // 100 = 4
```

### 2. Par ou ímpar (o último bit diz tudo)

```go
n := 7
fmt.Println(n&1 == 1) // true → ímpar
```

### 3. Trocar dois valores com XOR (curiosidade)

```go
a, b := 5, 9
a ^= b
b ^= a
a ^= b
fmt.Println(a, b) // 9 5
```

> Em Go, prefira `a, b = b, a`. O truque do XOR é só curiosidade de entrevista. 😉

### 4. Potência de 2?

```go
func ehPotenciaDe2(n int) bool {
	return n > 0 && n&(n-1) == 0
}
// 8 → 1000, 7 → 0111, 8 & 7 = 0 → true
```

---

## ⚠️ Pegadinha: `^` não é potência!

```go
fmt.Println(2 ^ 3) // 1  (0010 XOR 0011 = 0001), não 8!
fmt.Println(math.Pow(2, 3)) // 8 ✅
```

---

## ✍️ Exercícios

1. Calcule **no papel** e depois confira em código: `12 & 10`, `12 | 10`, `12 ^ 10`, `12 &^ 10`.
2. Imprima os números de 0 a 15 em binário com 4 dígitos (`%04b`).
3. Quanto vale `3 << 4`? E `100 >> 3`? Explique usando potências de 2.
4. Crie flags `Admin`, `Editor`, `Leitor` com `iota`. Dê a um usuário `Editor | Leitor`, depois:
   - verifique se ele é `Admin`;
   - adicione `Admin`;
   - remova `Leitor`.
5. Escreva uma função que diga se um número é ímpar usando **apenas** operador bit a bit.
6. Por que `2 ^ 10` não dá `1024` em Go? Como calcular corretamente?

---

⬅️ Anterior: [Comparação e lógicos](02-comparacao-e-logicos.md) · ➡️ Próximo: [Atribuição e precedência](04-atribuicao-e-precedencia.md)

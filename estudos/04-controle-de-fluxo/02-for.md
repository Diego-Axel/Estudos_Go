# 02 — `for`: o único laço de Go

Go tem **um único** comando de repetição: o **`for`**. Ele faz o papel de `for`, `while`, `do-while` e `foreach` das outras linguagens.

---

## 1️⃣ `for` clássico (3 partes)

```go
for i := 0; i < 5; i++ {
	fmt.Println(i)
}
// 0 1 2 3 4
```

```
for  inicialização ; condição ; pós {
```

- **inicialização:** roda uma vez, antes de tudo
- **condição:** checada antes de cada volta; se `false`, para
- **pós:** roda no fim de cada volta

Contando para trás, de 2 em 2:

```go
for i := 10; i > 0; i -= 2 {
	fmt.Print(i, " ")
}
// 10 8 6 4 2
```

Duas variáveis:

```go
for i, j := 0, 10; i < j; i, j = i+1, j-1 {
	fmt.Println(i, j)
}
```

---

## 2️⃣ `for` como `while`

Só a condição:

```go
n := 1
for n < 100 {
	n *= 2
}
fmt.Println(n) // 128
```

---

## 3️⃣ Laço infinito

Sem nada:

```go
for {
	fmt.Println("rodando pra sempre...")
}
```

Normalmente combinado com `break` ou `return`. É o que faz o **menu** do seu `main.go` da livraria:

```go
for {
	fmt.Print("Opção (0 para sair): ")
	var op int
	fmt.Scan(&op)

	if op == 0 {
		break
	}
	fmt.Println("você escolheu", op)
}
```

### Simulando `do-while`

Go não tem `do-while`, mas dá pra imitar (executa **pelo menos uma vez**):

```go
for {
	// corpo
	if !condicao {
		break
	}
}
```

---

## 4️⃣ `for range`: percorrendo coleções ⭐

O `range` itera sobre vários tipos e devolve **dois valores** por volta:

| Tipo | 1º valor | 2º valor |
|---|---|---|
| slice / array | índice | elemento |
| string | índice do **byte** | **rune** (caractere) |
| map | chave | valor |
| channel | elemento | *(só um valor)* |
| inteiro `n` (Go 1.22+) | `0` até `n-1` | *(só um valor)* |

### Slice

```go
frutas := []string{"maçã", "banana", "uva"}

for i, fruta := range frutas {
	fmt.Println(i, fruta)
}
// 0 maçã
// 1 banana
// 2 uva
```

Só o valor (ignora o índice com `_`):

```go
for _, fruta := range frutas {
	fmt.Println(fruta)
}
```

Só o índice:

```go
for i := range frutas {
	fmt.Println(i)
}
```

### String (percorre caracteres, não bytes!)

```go
for i, r := range "Olá!" {
	fmt.Printf("%d:%c ", i, r)
}
// 0:O 1:l 2:á 4:!   ← o índice pula o 3 porque 'á' ocupa 2 bytes
```

### Map

```go
idades := map[string]int{"Ana": 30, "Bia": 25, "Caio": 40}

for nome, idade := range idades {
	fmt.Println(nome, idade)
}
```

> ⚠️ A ordem de iteração de um **map é aleatória** e muda a cada execução. Se precisar de ordem, ordene as chaves antes (veremos no módulo de Maps).

### Inteiro (Go 1.22+)

```go
for i := range 3 {
	fmt.Println(i)
}
// 0 1 2
```

Repetir algo N vezes sem precisar do índice:

```go
for range 3 {
	fmt.Println("Go!")
}
```

---

## ⚠️ O valor do `range` é uma **cópia**

```go
nums := []int{1, 2, 3}

for _, n := range nums {
	n *= 10 // altera só a cópia!
}
fmt.Println(nums) // [1 2 3]

for i := range nums {
	nums[i] *= 10 // altera o slice de verdade
}
fmt.Println(nums) // [10 20 30]
```

---

## 🆕 Variável do laço é nova a cada volta (Go 1.22+)

Até o Go 1.21, a variável do laço era **a mesma** em todas as voltas, o que causava bugs com goroutines e closures. Desde o **Go 1.22**, **cada volta tem sua própria cópia**:

```go
var funcs []func()
for i := range 3 {
	funcs = append(funcs, func() { fmt.Print(i, " ") })
}
for _, f := range funcs {
	f()
}
// Go 1.22+: 0 1 2
// Go 1.21-: 2 2 2  (bug clássico)
```

> Isso depende da versão declarada no `go.mod` (`go 1.22` ou maior).

---

## 🔁 Laços aninhados

```go
for i := 1; i <= 3; i++ {
	for j := 1; j <= 3; j++ {
		fmt.Printf("%d×%d=%d\t", i, j, i*j)
	}
	fmt.Println()
}
```

---

## ✍️ Exercícios

1. Imprima os números de 1 a 100 que são múltiplos de 7.
2. Some todos os números de 1 a 1000 com `for` clássico. (Resposta: 500500)
3. Mostre a **tabuada** de um número lido do teclado (1 a 10).
4. Usando `for` estilo `while`, descubra quantas vezes dá pra dividir 1.000.000 por 2 até o resultado ficar menor que 1.
5. Percorra `"Programação"` com `range` e conte quantas **vogais** ela tem.
6. Dado `notas := []float64{7.5, 8, 6.5, 9, 10}`, calcule a média usando `range`.
7. Faça um menu em laço infinito com as opções: 1 - Dizer olá, 2 - Mostrar a hora (`time.Now()`), 0 - Sair.
8. Imprima um triângulo de asteriscos com altura 5:
   ```
   *
   **
   ***
   ****
   *****
   ```

---

⬅️ Anterior: [if / else](01-if-else.md) · ➡️ Próximo: [switch](03-switch.md)

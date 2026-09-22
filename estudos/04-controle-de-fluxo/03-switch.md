# 03 — `switch`

O `switch` de Go é **muito mais amigável** que o de C/Java:

- **Não precisa de `break`**: cada `case` para sozinho
- Um `case` pode ter **vários valores**
- Os valores **não precisam ser constantes** nem inteiros
- Pode ser usado **sem expressão**, como um `if/else if` mais limpo

---

## 1️⃣ `switch` com expressão

```go
dia := 3

switch dia {
case 1:
	fmt.Println("Domingo")
case 2:
	fmt.Println("Segunda")
case 3:
	fmt.Println("Terça")
default:
	fmt.Println("Outro dia")
}
// Terça
```

- Os `case` são testados **de cima para baixo**; o primeiro que bater é executado.
- `default` roda se **nenhum** bater (é opcional e pode ficar em qualquer posição).

---

## 2️⃣ Vários valores num `case`

```go
switch dia {
case 1, 7:
	fmt.Println("Fim de semana")
case 2, 3, 4, 5, 6:
	fmt.Println("Dia útil")
}
```

Com strings:

```go
switch opcao {
case "s", "S", "sim":
	fmt.Println("Confirmado")
case "n", "N", "nao", "não":
	fmt.Println("Cancelado")
default:
	fmt.Println("Opção inválida")
}
```

---

## 3️⃣ `switch` sem expressão (`switch true`) ⭐

Cada `case` é uma **condição booleana**. É a forma mais limpa de escrever vários `else if`:

```go
nota := 8.5

switch {
case nota >= 9:
	fmt.Println("A")
case nota >= 7:
	fmt.Println("B")
case nota >= 5:
	fmt.Println("C")
default:
	fmt.Println("Reprovado")
}
// B
```

---

## 4️⃣ `switch` com declaração curta

Igual ao `if`, dá pra declarar uma variável antes:

```go
switch hora := time.Now().Hour(); {
case hora < 12:
	fmt.Println("Bom dia")
case hora < 18:
	fmt.Println("Boa tarde")
default:
	fmt.Println("Boa noite")
}
```

Ou com expressão:

```go
switch so := runtime.GOOS; so {
case "windows":
	fmt.Println("Windows 🪟")
case "linux":
	fmt.Println("Linux 🐧")
case "darwin":
	fmt.Println("macOS 🍎")
default:
	fmt.Println(so)
}
```

---

## 5️⃣ `fallthrough`

Por padrão, cada `case` para sozinho. O `fallthrough` força a execução a **cair no próximo `case`**, **sem testar a condição dele**:

```go
switch n := 1; n {
case 1:
	fmt.Println("um")
	fallthrough
case 2:
	fmt.Println("dois")
	fallthrough
case 3:
	fmt.Println("três")
case 4:
	fmt.Println("quatro")
}
// um
// dois
// três
```

Regras:
- `fallthrough` tem que ser o **último comando** do `case`
- Não pode ser usado no **último** `case`
- Não funciona em *type switch*

> 💡 É pouco usado na prática. Quase sempre `case a, b:` resolve de forma mais clara.

---

## 6️⃣ *Type switch* (spoiler)

Um `switch` especial que testa o **tipo** de uma interface. Vai ser aprofundado no módulo de **Interfaces**:

```go
func descrever(v any) {
	switch x := v.(type) {
	case int:
		fmt.Println("inteiro:", x*2)
	case string:
		fmt.Println("texto com", len(x), "bytes")
	case bool:
		fmt.Println("booleano:", x)
	case nil:
		fmt.Println("nil")
	default:
		fmt.Printf("tipo desconhecido: %T\n", x)
	}
}

descrever(21)      // inteiro: 42
descrever("Go")    // texto com 2 bytes
descrever(3.14)    // tipo desconhecido: float64
```

---

## ⚠️ `break` dentro de `switch`

`break` dentro de um `switch` **sai do `switch`, não do `for`**! Pegadinha clássica em menus:

```go
for {
	switch op {
	case 0:
		break // ❌ sai só do switch, o for continua para sempre!
	}
}
```

A solução (com *label* ou `return`) está no próximo arquivo: [break, continue, labels e goto](04-break-continue-labels-goto.md).

---

## 🆚 `if/else if` ou `switch`?

| Use `switch` quando... | Use `if` quando... |
|---|---|
| compara **uma variável** com vários valores | tem 1 ou 2 condições |
| tem **muitas** faixas/condições (`switch` sem expressão) | as condições envolvem variáveis diferentes e complexas |
| quer testar o **tipo** de uma interface | precisa de *init statement* com `err` |

---

## ✍️ Exercícios

1. Leia um número de 1 a 12 e mostre o **nome do mês**.
2. Leia um mês (1 a 12) e mostre a **estação do ano** usando `case` com vários valores (considere o hemisfério Sul).
3. Refaça o exercício de classificação de IMC usando `switch` sem expressão.
4. Faça uma calculadora com `switch` no operador (`+`, `-`, `*`, `/`).
5. Leia uma letra e diga se é **vogal** ou **consoante** (trate maiúsculas também).
6. Qual a saída deste código? Explique.
   ```go
   switch x := 5; {
   case x > 3:
       fmt.Println("maior que 3")
       fallthrough
   case x > 10:
       fmt.Println("maior que 10")
   default:
       fmt.Println("default")
   }
   ```

---

⬅️ Anterior: [for](02-for.md) · ➡️ Próximo: [break, continue, labels e goto](04-break-continue-labels-goto.md)

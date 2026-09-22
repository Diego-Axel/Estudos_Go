# 01 — `if` / `else`

## 🔀 Sintaxe básica

```go
idade := 20

if idade >= 18 {
	fmt.Println("Maior de idade")
}
```

Regras de Go:

- **Sem parênteses** em volta da condição (o `gofmt` até remove se você colocar)
- **Chaves `{ }` obrigatórias**, mesmo com uma linha só
- A `{` fica **na mesma linha** do `if`
- A condição **tem que ser `bool`** (nada de `if 1` ou `if nome`)

```go
if (idade >= 18) { }       // funciona, mas não é o estilo Go
if idade >= 18 fmt.Println() // ❌ erro: chaves obrigatórias
```

---

## ↔️ `else` e `else if`

```go
nota := 7.5

if nota >= 9 {
	fmt.Println("A")
} else if nota >= 7 {
	fmt.Println("B")
} else if nota >= 5 {
	fmt.Println("C")
} else {
	fmt.Println("Reprovado")
}
```

> ⚠️ O `else` precisa ficar **na mesma linha** do `}` que fecha o bloco anterior:

```go
if x > 0 {
	// ...
}
else {        // ❌ erro de sintaxe
	// ...
}
```

---

## ⭐ `if` com declaração curta (*init statement*)

Você pode declarar uma variável **antes** da condição, separada por `;`. Ela só existe **dentro do `if`/`else`**:

```go
if n := len(nome); n > 10 {
	fmt.Println("nome longo:", n)
} else {
	fmt.Println("nome curto:", n) // n também existe aqui
}
// fmt.Println(n) // ❌ undefined: n
```

Esse é **o padrão mais usado de Go**, principalmente com erros:

```go
if err := salvar(dados); err != nil {
	fmt.Println("erro ao salvar:", err)
	return
}
```

E com conversões:

```go
if idade, err := strconv.Atoi(texto); err != nil {
	fmt.Println("idade inválida")
} else {
	fmt.Println("daqui a 10 anos você terá", idade+10)
}
```

E com maps (verificar se a chave existe):

```go
estoque := map[string]int{"caneta": 10}

if qtd, ok := estoque["lapis"]; ok {
	fmt.Println("tem", qtd)
} else {
	fmt.Println("produto não encontrado")
}
```

> Seu `main.go` da livraria usa exatamente isso: `if produto, existe := livraria.Produtos[codigo]; existe { ... }` 😉

---

## 🚪 Retorno antecipado (*early return*)

Em Go, o estilo é **tratar os erros/casos especiais primeiro e sair cedo**, deixando o "caminho feliz" sem indentação:

```go
// ❌ Estilo "pirâmide", difícil de ler
func processar(idade int, nome string) {
	if idade >= 0 {
		if nome != "" {
			fmt.Println("processando", nome)
		} else {
			fmt.Println("nome vazio")
		}
	} else {
		fmt.Println("idade inválida")
	}
}

// ✅ Estilo Go: guard clauses
func processar(idade int, nome string) {
	if idade < 0 {
		fmt.Println("idade inválida")
		return
	}
	if nome == "" {
		fmt.Println("nome vazio")
		return
	}
	fmt.Println("processando", nome)
}
```

Também evite `else` depois de um `return`:

```go
// ❌
if err != nil {
	return err
} else {
	fazAlgo()
}

// ✅
if err != nil {
	return err
}
fazAlgo()
```

---

## 🧩 Condições compostas

```go
if idade >= 18 && temCNH {
	fmt.Println("pode dirigir")
}

if dia == "sábado" || dia == "domingo" {
	fmt.Println("fim de semana")
}

if !(usuario == "admin") { // prefira: usuario != "admin"
	fmt.Println("acesso restrito")
}
```

---

## ✍️ Exercícios

1. Leia um número e diga se é **positivo, negativo ou zero**.
2. Leia uma nota (0 a 10) e mostre o conceito: A (9+), B (7+), C (5+), D (abaixo de 5). Se a nota estiver fora de 0 a 10, mostre "nota inválida".
3. Leia um texto e use `if` com *init statement* + `strconv.Atoi` para dizer se é um número válido.
4. Leia três números e mostre o **maior** deles.
5. Calcule o IMC (`peso / altura²`) e classifique: abaixo de 18,5 (abaixo do peso), até 24,9 (normal), até 29,9 (sobrepeso), acima (obesidade).
6. Reescreva sem aninhamento, usando retorno antecipado:
   ```go
   func sacar(saldo, valor float64) {
       if valor > 0 {
           if valor <= saldo {
               fmt.Println("saque realizado")
           } else {
               fmt.Println("saldo insuficiente")
           }
       } else {
           fmt.Println("valor inválido")
       }
   }
   ```

---

🏠 [Módulo 04](README.md) · ➡️ Próximo: [for](02-for.md)

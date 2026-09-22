# 04 — `break`, `continue`, labels e `goto`

## 🛑 `break`: sai do laço

```go
for i := 1; i <= 10; i++ {
	if i == 5 {
		break
	}
	fmt.Print(i, " ")
}
// 1 2 3 4
```

Exemplo real: procurar um item e parar quando achar.

```go
nomes := []string{"Ana", "Bia", "Caio", "Duda"}
alvo := "Caio"
encontrado := false

for i, nome := range nomes {
	if nome == alvo {
		fmt.Println("achei na posição", i)
		encontrado = true
		break
	}
}
if !encontrado {
	fmt.Println("não encontrado")
}
```

---

## ⏭️ `continue`: pula para a próxima volta

```go
for i := 1; i <= 10; i++ {
	if i%2 == 0 {
		continue // pula os pares
	}
	fmt.Print(i, " ")
}
// 1 3 5 7 9
```

Útil para **filtrar** e evitar aninhamento:

```go
for _, produto := range produtos {
	if produto.Estoque == 0 {
		continue
	}
	if produto.Preco > 100 {
		continue
	}
	fmt.Println(produto.Titulo) // só os em estoque e até R$ 100
}
```

---

## 🏷️ Labels (rótulos)

Com laços **aninhados**, `break` e `continue` afetam **só o laço mais interno**. Para agir no laço de fora, use um **label**:

```go
externo:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				continue externo // vai para a próxima volta do laço de FORA
			}
			if i == 2 {
				break externo // sai dos DOIS laços
			}
			fmt.Println(i, j)
		}
	}
// 0 0
// 0 1
// 1 0
// 1 1
```

- O label é um nome seguido de `:` logo **antes** do `for` (ou `switch`/`select`)
- Convenção: nomes curtos e descritivos (`externo`, `loop`, `busca`...)
- Label **declarado e não usado** é erro de compilação

### Exemplo: procurar numa matriz

```go
matriz := [][]int{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 9},
}
alvo := 5

busca:
	for i, linha := range matriz {
		for j, v := range linha {
			if v == alvo {
				fmt.Printf("achei %d em [%d][%d]\n", alvo, i, j)
				break busca
			}
		}
	}
// achei 5 em [1][1]
```

### ⭐ Resolvendo a pegadinha do `break` dentro do `switch`

```go
menu:
	for {
		fmt.Print("Opção: ")
		var op int
		fmt.Scan(&op)

		switch op {
		case 1:
			fmt.Println("Cadastrar")
		case 2:
			fmt.Println("Listar")
		case 0:
			fmt.Println("Saindo...")
			break menu // ✅ sai do FOR, não só do switch
		default:
			fmt.Println("Opção inválida")
		}
	}
```

> Alternativa: se o laço estiver dentro de uma função, um simples `return` também resolve. É o que o `main.go` da livraria faz no `case "0"`.

---

## 🦘 `goto`

Pula para um label **dentro da mesma função**. Existe, mas é **raramente usado**, porque deixa o código difícil de acompanhar.

```go
func main() {
	i := 0
inicio:
	if i < 3 {
		fmt.Println(i)
		i++
		goto inicio
	}
}
// 0 1 2
```

Restrições:
- Não pode pular **para dentro** de um bloco (`{ }`) de fora dele
- Não pode pular **por cima de uma declaração de variável** que ainda estaria em uso

```go
	goto fim
	x := 10 // ❌ erro: goto fim jumps over variable declaration
fim:
	fmt.Println(x)
```

> 🧠 **Regra prática:** use `for`, `break`, `continue`, labels e `return`. O `goto` fica para casos muito específicos (código gerado automaticamente, algumas rotinas de baixo nível na biblioteca padrão).

---

## 🧾 Resumão do módulo

| Comando | Faz o quê |
|---|---|
| `if cond { }` | executa se a condição for verdadeira |
| `if x := f(); cond { }` | declara e testa (escopo só do if/else) |
| `for i := 0; i < n; i++ { }` | laço clássico |
| `for cond { }` | laço tipo `while` |
| `for { }` | laço infinito |
| `for i, v := range x { }` | percorre slice, string, map, channel, int |
| `switch x { case ...: }` | compara um valor com vários casos |
| `switch { case cond: }` | vários `if/else if` de forma limpa |
| `fallthrough` | força cair no próximo `case` |
| `break` | sai do laço/switch mais interno |
| `continue` | pula para a próxima volta |
| `break label` / `continue label` | age no laço marcado pelo label |
| `goto label` | salto incondicional (evite) |

---

## ✍️ Exercícios

1. Imprima os números de 1 a 50, **pulando** os múltiplos de 3 (use `continue`).
2. Leia números até o usuário digitar `0` e mostre a **soma** de todos (use `for` infinito + `break`).
3. Verifique se um número é **primo**: teste os divisores de 2 até √n e use `break` ao achar o primeiro.
4. Em uma matriz 3×3, encontre o **primeiro número negativo** e pare as buscas nos dois laços com `break` + label.
5. Refaça o menu do exercício de `for` usando `switch` + `break menu`.
6. Imprima todos os pares `(i, j)` com `i` e `j` de 1 a 5, mas quando `i + j == 6` pule para o próximo `i` (use `continue` + label).
7. Jogo de adivinhação: gere um número aleatório de 1 a 100 (`rand.IntN(100) + 1`, pacote `math/rand/v2`) e deixe o usuário chutar até acertar, dizendo "maior" ou "menor" a cada chute. No final, mostre quantas tentativas foram.

---

⬅️ Anterior: [switch](03-switch.md) · ➡️ Próximo módulo: [Funções](../05-funcoes/README.md)

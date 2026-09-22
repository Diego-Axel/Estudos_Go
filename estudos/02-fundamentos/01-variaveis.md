# 01 — Variáveis

Uma **variável** é um nome que guarda um valor na memória. Em Go, **toda variável tem um tipo fixo**: depois de declarada como `int`, ela será `int` para sempre.

---

## 📝 Formas de declarar

### 1. `var` com tipo explícito

```go
var nome string
var idade int
nome = "Diego"
idade = 25
```

### 2. `var` com tipo e valor

```go
var nome string = "Diego"
var idade int = 25
```

### 3. `var` com inferência de tipo

O compilador **descobre o tipo** pelo valor:

```go
var nome = "Diego" // string
var idade = 25     // int
var altura = 1.75  // float64
```

### 4. Declaração curta `:=` ⭐ (a mais usada)

```go
nome := "Diego"
idade := 25
ativo := true
```

> ⚠️ O `:=` **só funciona dentro de funções**. Fora delas (no nível do pacote), use `var`.

```go
package main

var global = "ok"   // ✅
// outra := "erro"  // ❌ syntax error: non-declaration statement outside function body

func main() {
	local := "ok"   // ✅
	_ = local
}
```

---

## 🧮 Declarando várias de uma vez

```go
// Na mesma linha
var a, b, c int = 1, 2, 3
x, y := 10, "dez"

// Em bloco (muito usado para variáveis de pacote)
var (
	nome   string = "Ana"
	idade  int    = 30
	ativo  bool   = true
)
```

---

## `=` vs `:=`

| Operador | O que faz |
|---|---|
| `:=` | **Declara** uma nova variável **e** atribui |
| `=` | Só **atribui** a uma variável que já existe |

```go
x := 10  // declara
x = 20   // atribui (ok)
x := 30  // ❌ erro: no new variables on left side of :=
```

**Exceção:** o `:=` é permitido se **pelo menos uma** variável da esquerda for nova. Isso é muito comum com erros:

```go
arquivo, err := os.Open("a.txt")
dados, err := io.ReadAll(arquivo) // ✅ 'dados' é nova; 'err' é só reatribuída
```

---

## 0️⃣ Zero values (valores zero)

Em Go **não existe variável "sem valor"** (nem lixo de memória como em C). Toda variável declarada sem valor recebe o **valor zero** do seu tipo:

| Tipo | Valor zero |
|---|---|
| `int`, `float64`, etc. | `0` |
| `string` | `""` (string vazia) |
| `bool` | `false` |
| ponteiros, slices, maps, channels, funções, interfaces | `nil` |
| structs | cada campo com seu valor zero |

```go
var i int
var s string
var b bool
var p *int

fmt.Println(i, s, b, p) // 0  false <nil>
fmt.Printf("%q\n", s)   // ""
```

---

## 🚫 Variável não usada = erro

```go
func main() {
	x := 10 // ❌ declared and not used: x
}
```

Se você **precisa** receber um valor mas não vai usá-lo, use o **identificador em branco** `_`:

```go
_, err := fmt.Println("oi") // ignora o primeiro retorno
```

> 📌 Essa regra vale só para variáveis **locais**. Variáveis de pacote podem ficar sem uso.

---

## 🔄 Trocando valores (swap)

Go permite **atribuição múltipla**, então trocar valores é trivial:

```go
a, b := 1, 2
a, b = b, a
fmt.Println(a, b) // 2 1
```

---

## 🔭 Escopo

Uma variável só existe **dentro do bloco `{ }`** onde foi declarada.

```go
var pacote = "visível em todo o pacote"

func main() {
	funcao := "visível em toda a main"

	if true {
		bloco := "visível só dentro do if"
		fmt.Println(pacote, funcao, bloco)
	}

	// fmt.Println(bloco) // ❌ undefined: bloco
}
```

### ⚠️ Cuidado com *shadowing* (sombreamento)

Declarar com `:=` num bloco interno cria uma **nova** variável que "esconde" a de fora:

```go
x := 1
if true {
	x := 2         // NOVA variável x, só dentro do if
	fmt.Println(x) // 2
}
fmt.Println(x)     // 1  <- a de fora não mudou!
```

Se a intenção era alterar a de fora, use `=`:

```go
x := 1
if true {
	x = 2
}
fmt.Println(x) // 2
```

---

## 🏷️ Regras e convenções de nomes

**Regras (obrigatórias):**
- Começa com **letra** ou `_`; depois letras, dígitos ou `_`
- Não pode ser palavra-chave (`func`, `var`, `if`...)
- Diferencia maiúsculas de minúsculas (`nome` ≠ `Nome`)

**Convenções (estilo Go):**
- Use **camelCase**: `nomeCompleto`, `totalDeVendas` (nada de `nome_completo`)
- Siglas ficam inteiras: `userID`, `urlHTTP`, `parseJSON`
- **Letra maiúscula = exportada** (pública fora do pacote); minúscula = privada
- Nomes **curtos** para escopos pequenos (`i`, `n`, `err`), nomes **descritivos** para escopos grandes

```go
var Total int   // exportada: outros pacotes enxergam
var total int   // privada: só este pacote enxerga
```

---

## ✍️ Exercícios

1. Declare as variáveis `nome`, `idade` e `altura` usando as **4 formas** de declaração que você aprendeu.
2. Declare `var x int`, `var s string`, `var b bool` sem valor e imprima com `%v` e `%q`. O que aparece?
3. Troque os valores de duas variáveis `a := "primeiro"` e `b := "segundo"` em **uma linha**.
4. Explique por que este código imprime `10` e não `20`:
   ```go
   n := 10
   if n > 5 {
       n := 20
       _ = n
   }
   fmt.Println(n)
   ```
5. O que está errado aqui?
   ```go
   package main
   idade := 30
   func main() {}
   ```

---

🏠 [Módulo 02](README.md) · ➡️ Próximo: [Tipos de dados](02-tipos-de-dados.md)

# 04 — Maps

Um **map** guarda pares **chave → valor**, com busca muito rápida pela chave. É o equivalente ao *dicionário* do Python, ao *HashMap* do Java e ao *objeto* do JavaScript.

```go
map[TipoDaChave]TipoDoValor
```

---

## 📝 Criando

```go
// 1. Literal
idades := map[string]int{
	"Ana":  30,
	"Bia":  25,
	"Caio": 40, // ← a vírgula no último item é obrigatória
}

// 2. make
estoque := make(map[string]int)

// 3. make com tamanho inicial estimado (otimização)
cache := make(map[int]string, 1000)

// 4. Map vazio (literal)
vazio := map[string]bool{}
```

### ⚠️ Map `nil`: lê, mas não escreve

```go
var m map[string]int   // nil
fmt.Println(m["x"])    // 0 ✅ ler funciona
fmt.Println(len(m))    // 0 ✅

m["x"] = 1             // 💥 panic: assignment to entry in nil map
```

> ✅ Sempre inicialize com `make` ou `{}` antes de escrever.

---

## 🔑 Operações básicas

```go
estoque := make(map[string]int)

// Inserir / atualizar
estoque["caneta"] = 10
estoque["lapis"] = 5
estoque["caneta"] = 12   // atualiza

// Ler
fmt.Println(estoque["caneta"]) // 12

// Tamanho
fmt.Println(len(estoque))      // 2

// Remover
delete(estoque, "lapis")
delete(estoque, "borracha")    // não existe: não faz nada (sem erro)

// Limpar tudo (Go 1.21+)
clear(estoque)
```

---

## ❓ A chave existe? O idioma *"comma ok"* ⭐

Buscar uma chave que **não existe** retorna o **valor zero** do tipo, sem erro:

```go
notas := map[string]float64{"Ana": 9.5, "Bia": 0}

fmt.Println(notas["Bia"])  // 0
fmt.Println(notas["Caio"]) // 0  ← mesmo resultado! Caio tirou 0 ou não existe?
```

Para diferenciar, use o **segundo retorno**:

```go
nota, ok := notas["Caio"]
if !ok {
	fmt.Println("Caio não está no mapa")
} else {
	fmt.Println("nota:", nota)
}

// Forma mais comum, com if + init statement
if nota, ok := notas["Bia"]; ok {
	fmt.Println("Bia tirou", nota) // Bia tirou 0
}
```

> É exatamente o que o seu `main.go` faz: `if produto, existe := livraria.Produtos[codigo]; existe`.

---

## 🔁 Percorrendo

```go
precos := map[string]float64{"café": 5.5, "pão": 0.8, "leite": 4.9}

for produto, preco := range precos {
	fmt.Printf("%s: R$ %.2f\n", produto, preco)
}

for produto := range precos { // só as chaves
	fmt.Println(produto)
}
```

### ⚠️ A ordem é aleatória!

Go **embaralha de propósito** a ordem de iteração dos maps. Rode duas vezes e pode sair diferente. **Nunca** conte com a ordem.

### Percorrendo em ordem

Pegue as chaves, ordene, e depois acesse:

```go
import (
	"maps"
	"slices"
)

// Go 1.23+
for _, produto := range slices.Sorted(maps.Keys(precos)) {
	fmt.Println(produto, precos[produto])
}
// café 5.5
// leite 4.9
// pão 0.8
```

Forma clássica (funciona em qualquer versão):

```go
chaves := make([]string, 0, len(precos))
for k := range precos {
	chaves = append(chaves, k)
}
sort.Strings(chaves)
for _, k := range chaves {
	fmt.Println(k, precos[k])
}
```

---

## 🗝️ Quais tipos podem ser chave?

A chave precisa ser **comparável com `==`**:

| ✅ Pode ser chave | ❌ Não pode |
|---|---|
| `string`, `int`, `float64`, `bool` | slices |
| arrays (`[2]int`) | maps |
| structs (com campos comparáveis) | funções |
| ponteiros | |

```go
// Coordenadas como chave
tabuleiro := map[[2]int]string{
	{0, 0}: "X",
	{1, 1}: "O",
}
fmt.Println(tabuleiro[[2]int{1, 1}]) // O
```

---

## 🧷 Maps são "referências"

Diferente de arrays, passar um map para uma função **não copia os dados**. A função altera o map original:

```go
func adicionarBonus(salarios map[string]float64) {
	for nome := range salarios {
		salarios[nome] *= 1.1
	}
}

s := map[string]float64{"Ana": 1000}
adicionarBonus(s)
fmt.Println(s) // map[Ana:1100]
```

E atribuir um map a outra variável **não cria cópia**:

```go
a := map[string]int{"x": 1}
b := a
b["x"] = 99
fmt.Println(a["x"]) // 99 😱

c := maps.Clone(a) // ✅ cópia de verdade (Go 1.21+)
```

---

## 🛠️ Padrões úteis

### Contador de ocorrências

```go
texto := "go é legal e go é rápido"
contagem := make(map[string]int)

for _, palavra := range strings.Fields(texto) {
	contagem[palavra]++ // chave inexistente começa em 0!
}
fmt.Println(contagem["go"], contagem["é"]) // 2 2
```

### Agrupar (map de slices)

```go
alunos := []string{"Ana", "Bruno", "Alice", "Beatriz", "Carlos"}
porLetra := make(map[string][]string)

for _, nome := range alunos {
	inicial := string(nome[0])
	porLetra[inicial] = append(porLetra[inicial], nome) // append em nil funciona
}
fmt.Println(porLetra["A"]) // [Ana Alice]
```

### Conjunto (*set*)

Go não tem tipo `set`. Usa-se um map com valor **vazio**:

```go
visto := make(map[string]struct{}) // struct{} ocupa 0 bytes

visto["go"] = struct{}{}
visto["rust"] = struct{}{}

if _, ok := visto["go"]; ok {
	fmt.Println("já vi go")
}
```

Ou, mais legível, `map[string]bool`:

```go
visto := map[string]bool{}
visto["go"] = true
if visto["go"] { // chave ausente retorna false
	fmt.Println("já vi go")
}
```

### Map de maps

```go
notas := map[string]map[string]float64{
	"Ana": {"mat": 9, "port": 8},
}

// ⚠️ o map interno também precisa ser criado antes de escrever
if notas["Bia"] == nil {
	notas["Bia"] = make(map[string]float64)
}
notas["Bia"]["mat"] = 7
```

---

## 🧱 Map de structs: uma pegadinha

```go
type Produto struct {
	Nome    string
	Estoque int
}

produtos := map[string]Produto{"p1": {"Caneta", 10}}

// produtos["p1"].Estoque = 5 // ❌ erro: cannot assign to struct field produtos["p1"].Estoque in map

// ✅ Solução: pegar, alterar e guardar de volta
p := produtos["p1"]
p.Estoque = 5
produtos["p1"] = p
```

> Outra opção é usar `map[string]*Produto` (ponteiros). Veremos nos módulos **Ponteiros** e **Structs**.

---

## ⚡ Maps e concorrência

Maps **não são seguros** para várias goroutines escreverem ao mesmo tempo. O programa pode quebrar com `fatal error: concurrent map writes`. A solução (`sync.Mutex` ou `sync.Map`) está no módulo de **Concorrência**.

---

## 🧾 Resumão do módulo

| | Array | Slice | Map |
|---|---|---|---|
| Tamanho | fixo | dinâmico | dinâmico |
| Sintaxe | `[3]int{}` | `[]int{}` | `map[string]int{}` |
| Valor zero | array de zeros | `nil` | `nil` |
| Ao atribuir/passar | **copia** tudo | copia o cabeçalho (mesmo array) | mesma estrutura |
| Comparável com `==` | ✅ | ❌ (só com `nil`) | ❌ (só com `nil`) |
| Acesso | índice | índice | chave |
| Ordem | mantida | mantida | **aleatória** |

---

## ✍️ Exercícios

1. Crie um map `capitais` (estado → capital) com 5 estados. Leia uma sigla do teclado e mostre a capital ou "estado não cadastrado" (use *comma ok*).
2. Conte quantas vezes cada **letra** aparece em uma frase (ignore espaços). Mostre em ordem alfabética.
3. Dado um slice de números com repetições, crie um slice **sem duplicados mantendo a ordem original** (use um map como *set*).
4. Agrupe um slice de palavras pelo **tamanho**: `map[int][]string`.
5. Crie um "carrinho de compras": `map[string]int` (produto → quantidade) e funções `adicionar`, `remover` e `total` (com preços em outro map).
6. Por que `var m map[string]int; m["a"] = 1` dá panic? Corrija.
7. Inverta um map: de `map[string]int{"a": 1, "b": 2}` gere `map[int]string{1: "a", 2: "b"}`.

---

⬅️ Anterior: [Slices: avançado](03-slices-avancado.md) · 🏠 [Voltar ao roteiro](../README.md)

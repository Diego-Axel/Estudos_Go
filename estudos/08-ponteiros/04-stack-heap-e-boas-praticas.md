# 04 — Stack, Heap e Boas Práticas

## 🗄️ Onde as variáveis moram?

Um programa usa duas grandes áreas de memória:

| | **Stack (pilha)** | **Heap** |
|---|---|---|
| O que guarda | variáveis locais de cada chamada de função | valores que precisam viver além da função |
| Alocação | automática e **muito rápida** | mais lenta |
| Liberação | automática ao **sair da função** | pelo **Garbage Collector (GC)** |
| Tamanho | pequena (cresce sob demanda em Go) | grande |

```
Stack                          Heap
┌──────────────────┐          ┌───────────────────────┐
│ main()           │          │                       │
│   x = 10         │          │   Produto{...}  ◄──┐  │
│   p = 0xc000... ─┼──────────┼────────────────────┘  │
├──────────────────┤          │                       │
│ calcular()       │          │                       │
│   a = 3, b = 4   │          │                       │
└──────────────────┘          └───────────────────────┘
```

---

## 🔍 Análise de escape (*escape analysis*)

Em Go, **você não escolhe** stack ou heap: o **compilador decide**. Se ele prova que a variável **não é usada depois** que a função termina, ela fica na **stack**. Se ela **"escapa"** (ex: você retorna o endereço dela), vai para o **heap**:

```go
func naStack() int {
	x := 10     // fica na stack
	return x    // retorna uma CÓPIA do valor
}

func noHeap() *int {
	x := 10     // escapa para o heap
	return &x   // o endereço é usado fora da função
}
```

### Vendo as decisões do compilador

```bash
go build -gcflags=-m main.go
```

Saída (exemplo):

```
./main.go:8:2: moved to heap: x
./main.go:15:13: ... argument does not escape
```

Outros motivos comuns para escapar:
- guardar o ponteiro em uma variável global, slice ou map que sobrevive à função
- passar para uma interface (ex: `fmt.Println(x)` pode fazer o valor escapar)
- valores muito grandes, ou de tamanho desconhecido em tempo de compilação

---

## ♻️ Garbage Collector

Tudo que vai para o heap é liberado pelo **GC** quando **ninguém mais aponta** para o valor. Você **nunca** chama `free()`.

```go
func main() {
	p := &Produto{Nome: "Caneta"}
	p = nil // o Produto não tem mais referências → o GC vai liberar
}
```

O GC de Go é **concorrente** e tem pausas bem curtas (geralmente frações de milissegundo). Ainda assim, **menos alocações no heap = menos trabalho pro GC = programa mais rápido**.

---

## ⚡ Ponteiro nem sempre é mais rápido

Uma crença comum: "passar ponteiro é sempre mais eficiente que copiar". **Nem sempre!**

- Copiar um valor **pequeno** (um `int`, uma struct com poucos campos) é **baratíssimo**
- Um ponteiro pode fazer o valor **escapar para o heap**, o que custa alocação + trabalho de GC
- Acessar dados via ponteiro pode ser pior para o **cache** do processador

> 🧠 Otimize com **medições** (benchmarks, módulo 15), não por intuição.

---

## ✅ Boas práticas com ponteiros

### 1. Use ponteiro quando precisar **alterar**

```go
func (c *Conta) Depositar(v float64) { c.Saldo += v } // spoiler: métodos (módulo 09)
```

### 2. Use ponteiro para structs **grandes** ou que **não devem ser copiadas**

Ex: structs que contêm um `sync.Mutex`, conexões, buffers.

### 3. Use valor para tipos **pequenos e imutáveis**

```go
type Ponto struct{ X, Y float64 }

func distancia(a, b Ponto) float64 { ... } // ✅ valor: pequeno e só leitura
```

### 4. Não use ponteiro para slices, maps, channels, funções e interfaces

Eles já são leves e carregam referências internas. `*[]int` e `*map[string]int` são raros e quase sempre desnecessários.

### 5. Sempre verifique `nil` quando o ponteiro puder ser `nil`

```go
func imprimir(p *Produto) {
	if p == nil {
		fmt.Println("produto inexistente")
		return
	}
	fmt.Println(p.Nome)
}
```

### 6. Seja consistente

Se um tipo é normalmente manipulado por ponteiro (`*Conta`), use ponteiro em **todo lugar** que lida com ele. Misturar confunde.

---

## 🧾 Resumão do módulo

| Conceito | Código |
|---|---|
| Pegar endereço | `p := &x` |
| Ler/alterar pelo ponteiro | `*p`, `*p = 10` |
| Tipo ponteiro | `*int`, `*Produto` |
| Valor zero | `nil` (desreferenciar dá panic!) |
| Criar valor zero + ponteiro | `new(int)` |
| Struct como ponteiro | `&Produto{...}` |
| Acesso a campo | `p.Nome` (desreferência automática) |
| Função altera o original | `func f(p *T)` + `f(&x)` |
| Retornar ponteiro de local | seguro (escape analysis) |
| Aritmética de ponteiros | ❌ não existe |

---

## ✍️ Exercícios

1. Escreva as duas funções abaixo, rode `go build -gcflags=-m` e veja qual variável vai para o heap:
   ```go
   func a() int  { x := 1; return x }
   func b() *int { x := 1; return &x }
   ```
2. Explique com suas palavras a diferença entre stack e heap.
3. Por que "sempre usar ponteiro para ter performance" é um mito?
4. Para cada caso, diga se usaria **valor** ou **ponteiro** e por quê:
   - `func area(r Retangulo) float64`
   - `func (c Conta) Sacar(v float64)`
   - `func processar(lista []Pedido)`
   - `func atualizarConfig(cfg Config)` (Config tem 40 campos)
5. Escreva uma função `buscar(produtos []*Produto, nome string) *Produto` que retorne `nil` se não encontrar, e trate o `nil` no `main`.

---

⬅️ Anterior: [Ponteiros com structs e coleções](03-ponteiros-com-structs-e-colecoes.md) · 🏠 [Voltar ao roteiro](../README.md)

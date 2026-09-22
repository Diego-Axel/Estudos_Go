# 03 — Tipos Genéricos

Além de funções, **tipos** (structs, slices, maps...) também podem ter parâmetros de tipo. Perfeito para **estruturas de dados**.

---

## 📚 Uma pilha genérica

```go
type Pilha[T any] struct {
	itens []T
}

func (p *Pilha[T]) Empilhar(v T) {
	p.itens = append(p.itens, v)
}

func (p *Pilha[T]) Desempilhar() (T, bool) {
	if len(p.itens) == 0 {
		var zero T
		return zero, false
	}
	ultimo := p.itens[len(p.itens)-1]
	p.itens = p.itens[:len(p.itens)-1]
	return ultimo, true
}

func (p *Pilha[T]) Tamanho() int {
	return len(p.itens)
}
```

Usando:

```go
func main() {
	var numeros Pilha[int]   // precisa dizer o tipo: Pilha[int]
	numeros.Empilhar(1)
	numeros.Empilhar(2)
	v, _ := numeros.Desempilhar()
	fmt.Println(v) // 2

	nomes := &Pilha[string]{}
	nomes.Empilhar("Ana")
	// nomes.Empilhar(42) // ❌ erro: cannot use 42 as string value
	fmt.Println(nomes.Tamanho()) // 1
}
```

Detalhes da sintaxe:
- Na **declaração** do tipo: `type Pilha[T any] struct`
- No **receptor** dos métodos: `(p *Pilha[T])`, repetindo o `T` (sem a restrição)
- Ao **usar**: `Pilha[int]`, `Pilha[string]`: o tipo precisa ser informado (não há inferência para tipos)

---

## 🚶 Uma fila genérica

```go
type Fila[T any] struct {
	itens []T
}

func (f *Fila[T]) Entrar(v T) {
	f.itens = append(f.itens, v)
}

func (f *Fila[T]) Sair() (T, bool) {
	var zero T
	if len(f.itens) == 0 {
		return zero, false
	}
	primeiro := f.itens[0]
	f.itens[0] = zero // ajuda o GC a liberar a memória
	f.itens = f.itens[1:]
	return primeiro, true
}

func (f *Fila[T]) Vazia() bool {
	return len(f.itens) == 0
}
```

---

## 🧺 Um conjunto (*set*) genérico

```go
type Conjunto[T comparable] struct {
	itens map[T]struct{}
}

func NovoConjunto[T comparable](valores ...T) *Conjunto[T] {
	c := &Conjunto[T]{itens: make(map[T]struct{})}
	for _, v := range valores {
		c.Adicionar(v)
	}
	return c
}

func (c *Conjunto[T]) Adicionar(v T)      { c.itens[v] = struct{}{} }
func (c *Conjunto[T]) Remover(v T)        { delete(c.itens, v) }
func (c *Conjunto[T]) Contem(v T) bool    { _, ok := c.itens[v]; return ok }
func (c *Conjunto[T]) Tamanho() int       { return len(c.itens) }

func main() {
	linguagens := NovoConjunto("go", "rust", "go", "python") // T inferido: string
	fmt.Println(linguagens.Tamanho())    // 3 (o "go" repetido não conta)
	fmt.Println(linguagens.Contem("go")) // true
}
```

> Repare: na **função construtora** o tipo é inferido (`NovoConjunto("go", ...)`), por isso construtores genéricos são tão comuns.

---

## 🤝 Tipos com vários parâmetros

```go
type Par[K comparable, V any] struct {
	Chave K
	Valor V
}

func (p Par[K, V]) String() string {
	return fmt.Sprintf("%v=%v", p.Chave, p.Valor)
}

func ParesDe[K comparable, V any](m map[K]V) []Par[K, V] {
	pares := make([]Par[K, V], 0, len(m))
	for k, v := range m {
		pares = append(pares, Par[K, V]{k, v})
	}
	return pares
}

p := Par[string, int]{"idade", 30}
fmt.Println(p) // idade=30
```

---

## 🗃️ Um cache genérico

```go
type Cache[K comparable, V any] struct {
	dados map[K]V
}

func NovoCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{dados: make(map[K]V)}
}

func (c *Cache[K, V]) Guardar(k K, v V) { c.dados[k] = v }

func (c *Cache[K, V]) Obter(k K) (V, bool) {
	v, ok := c.dados[k]
	return v, ok
}

func main() {
	usuarios := NovoCache[int, string]() // aqui é obrigatório informar: não há argumentos
	usuarios.Guardar(1, "Ana")

	if nome, ok := usuarios.Obter(1); ok {
		fmt.Println(nome) // Ana
	}
}
```

---

## 🚫 Limitação: métodos não têm parâmetros de tipo próprios

Um método pode **usar** os parâmetros do tipo, mas **não pode declarar novos**:

```go
func (p *Pilha[T]) Mapear[U any](f func(T) U) *Pilha[U] { // ❌ erro de sintaxe: método não pode ter parâmetros de tipo
	...
}
```

Solução: use uma **função** comum:

```go
func MapearPilha[T, U any](p *Pilha[T], f func(T) U) *Pilha[U] {
	nova := &Pilha[U]{}
	for _, v := range p.itens {
		nova.Empilhar(f(v))
	}
	return nova
}
```

---

## 🏷️ Outros tipos genéricos

Não precisa ser struct:

```go
type Lista[T any] []T

func (l Lista[T]) Primeiro() T { return l[0] }

type Dicionario[V any] map[string]V

type Transformador[T any] func(T) T

nums := Lista[int]{1, 2, 3}
fmt.Println(nums.Primeiro()) // 1
```

### Alias de tipo genérico (Go 1.24+)

```go
type Conjunto[T comparable] = map[T]struct{} // alias: é o mesmo tipo que o map
```

---

## 🌳 Exemplo avançado: árvore binária de busca

```go
type Arvore[T cmp.Ordered] struct {
	raiz *no[T]
}

type no[T cmp.Ordered] struct {
	valor    T
	esq, dir *no[T]
}

func (a *Arvore[T]) Inserir(v T) {
	a.raiz = inserir(a.raiz, v)
}

func inserir[T cmp.Ordered](n *no[T], v T) *no[T] {
	if n == nil {
		return &no[T]{valor: v}
	}
	if v < n.valor {
		n.esq = inserir(n.esq, v)
	} else {
		n.dir = inserir(n.dir, v)
	}
	return n
}

func (a *Arvore[T]) EmOrdem(visitar func(T)) {
	var percorrer func(*no[T])
	percorrer = func(n *no[T]) {
		if n == nil {
			return
		}
		percorrer(n.esq)
		visitar(n.valor)
		percorrer(n.dir)
	}
	percorrer(a.raiz)
}

func main() {
	var arv Arvore[string]
	for _, s := range []string{"manga", "banana", "uva", "abacaxi"} {
		arv.Inserir(s)
	}
	arv.EmOrdem(func(s string) { fmt.Print(s, " ") })
	// abacaxi banana manga uva
}
```

---

## ✍️ Exercícios

1. Implemente a `Pilha[T]` e adicione os métodos `Topo() (T, bool)` (sem remover) e `Vazia() bool`.
2. Implemente a `Fila[T]` e use-a para simular uma fila de atendimento com nomes.
3. Adicione ao `Conjunto[T]` os métodos `Uniao`, `Intersecao` e `Valores() []T`.
4. Crie `Resultado[T any]` com os campos `Valor T` e `Err error` e um método `Ok() bool`.
5. Crie uma `ListaEncadeada[T any]` com `Adicionar`, `Remover` e `ParaSlice() []T`.
6. Tente criar um método com parâmetro de tipo próprio, leia o erro e resolva com uma função.
7. Adicione à `Arvore[T]` um método `Contem(v T) bool`.

---

⬅️ Anterior: [Restrições](02-restricoes.md) · ➡️ Próximo: [Funções utilitárias, iteradores e boas práticas](04-utilitarios-iteradores-e-boas-praticas.md)

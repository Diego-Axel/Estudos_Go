# 04 — Tags de Struct e Boas Práticas

## 🏷️ Tags de struct

Uma **tag** é um texto entre crases depois do tipo do campo. Ela não muda nada no Go em si, mas **bibliotecas** leem essas tags para saber como tratar cada campo.

```go
type Produto struct {
	Codigo  string  `json:"codigo"`
	Titulo  string  `json:"titulo"`
	Preco   float64 `json:"preco"`
	Estoque int     `json:"estoque,omitempty"`
	Custo   float64 `json:"-"`
}
```

O formato é `chave:"valor"`, e pode haver várias separadas por espaço:

```go
Email string `json:"email" db:"email_usuario" validate:"required,email"`
```

---

## 📦 Tags na prática: JSON

```go
import "encoding/json"

p := Produto{Codigo: "L001", Titulo: "O Hobbit", Preco: 49.9, Custo: 30}

dados, err := json.Marshal(p) // struct → JSON
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(dados))
// {"codigo":"L001","titulo":"O Hobbit","preco":49.9}
```

Repare:
- os nomes vieram das **tags** (minúsculos)
- `Estoque` sumiu por causa do **`omitempty`** (estava com valor zero)
- `Custo` sumiu por causa do **`-`** (nunca vai pro JSON)

E o caminho inverso:

```go
entrada := `{"codigo":"L002","titulo":"Duna","preco":79.9,"estoque":3}`

var p2 Produto
if err := json.Unmarshal([]byte(entrada), &p2); err != nil { // JSON → struct
	log.Fatal(err)
}
fmt.Printf("%+v\n", p2)
// {Codigo:L002 Titulo:Duna Preco:79.9 Estoque:3 Custo:0}
```

### Opções comuns da tag `json`

| Tag | Efeito |
|---|---|
| `json:"nome"` | usa "nome" como chave |
| `json:"nome,omitempty"` | omite se for valor "vazio" (0, "", nil, slice/map vazio) |
| `json:"nome,omitzero"` | omite se for valor **zero** (Go 1.24+; funciona melhor com structs e `time.Time`) |
| `json:"-"` | nunca inclui |
| `json:",string"` | codifica número/bool como string (`"42"`) |

> ⚠️ Lembre: **só campos exportados** (maiúsculos) entram no JSON, com ou sem tag.

Veremos JSON em detalhe no módulo **Biblioteca padrão**.

### Lendo tags você mesmo (reflexão)

```go
import "reflect"

t := reflect.TypeOf(Produto{})
campo, _ := t.FieldByName("Preco")
fmt.Println(campo.Tag.Get("json")) // preco
```

> É assim que as bibliotecas fazem por dentro. No dia a dia, você raramente vai precisar.

---

## ✅ Boas práticas com structs

### 1. Use literais com **nomes de campos**

```go
p := Produto{Titulo: "Duna", Preco: 79.9} // ✅
p := Produto{"L002", "Duna", 79.9, 0, 0}  // ❌ frágil
```

### 2. Faça o **valor zero ser útil**

Um bom tipo Go funciona **sem inicialização especial**. Exemplos da biblioteca padrão: `sync.Mutex`, `bytes.Buffer`, `strings.Builder` funcionam direto com `var x T`.

```go
type Carrinho struct {
	itens []string // slice nil: append funciona sem make!
}

func (c *Carrinho) Adicionar(item string) {
	c.itens = append(c.itens, item)
}

var c Carrinho       // ✅ já pronto pra usar
c.Adicionar("café")
```

⚠️ Maps **não** têm valor zero útil para escrita. Se a struct tem map, inicialize no construtor ou "preguiçosamente":

```go
func (e *Estoque) Adicionar(nome string, qtd int) {
	if e.itens == nil {
		e.itens = make(map[string]int)
	}
	e.itens[nome] += qtd
}
```

### 3. Use um construtor quando houver **validação** ou **inicialização obrigatória**

```go
func NovaLivraria() *Livraria {
	return &Livraria{Produtos: make(map[string]Produto)}
}
```

> Seu `main.go` faz `Livraria{Produtos: make(map[string]Produto)}` direto no `main`. Um construtor `NovaLivraria()` deixaria isso reutilizável.

### 4. **Encapsule** com campos privados + métodos

```go
type Conta struct {
	titular string
	saldo   float64
}

func (c *Conta) Saldo() float64 { return c.saldo } // getter sem "Get"
```

Assim, **ninguém** consegue fazer `conta.saldo = 1000000` de fora do pacote. 😄

### 5. Ordem dos campos afeta o tamanho (curiosidade)

O compilador alinha os campos na memória, e pode sobrar "espaço vazio" (*padding*):

```go
import "unsafe"

type Ruim struct {
	a bool  // 1 byte + 7 de padding
	b int64 // 8 bytes
	c bool  // 1 byte + 7 de padding
}

type Bom struct {
	b int64 // 8 bytes
	a bool  // 1 byte
	c bool  // 1 byte + 6 de padding
}

fmt.Println(unsafe.Sizeof(Ruim{})) // 24
fmt.Println(unsafe.Sizeof(Bom{}))  // 16
```

> Só vale se preocupar com isso quando tiver **milhões** de instâncias. No dia a dia, **priorize a legibilidade**.

### 6. Opções funcionais (padrão avançado)

Para construtores com **muitas configurações opcionais**:

```go
type Servidor struct {
	porta   int
	timeout time.Duration
}

type Opcao func(*Servidor)

func ComPorta(p int) Opcao             { return func(s *Servidor) { s.porta = p } }
func ComTimeout(t time.Duration) Opcao { return func(s *Servidor) { s.timeout = t } }

func NovoServidor(opcoes ...Opcao) *Servidor {
	s := &Servidor{porta: 8080, timeout: 30 * time.Second} // padrões
	for _, op := range opcoes {
		op(s)
	}
	return s
}

s1 := NovoServidor()                             // tudo padrão
s2 := NovoServidor(ComPorta(3000))               // só muda a porta
s3 := NovoServidor(ComPorta(3000), ComTimeout(5*time.Second))
```

> Junta **variádicas** + **closures** + **ponteiros**: tudo que você já viu! 🎉

---

## 🧾 Resumão do módulo

| Conceito | Exemplo |
|---|---|
| Declarar struct | `type Produto struct { Nome string }` |
| Literal | `Produto{Nome: "Duna"}` |
| Ponteiro | `&Produto{...}` |
| Struct anônima | `struct{ X int }{X: 1}` |
| Método (leitura) | `func (p Produto) Total() float64` |
| Método (alteração) | `func (p *Produto) Vender(n int)` |
| Imprimir bonito | `func (p Produto) String() string` |
| Embedding | `type Cliente struct { Pessoa; Pontos int }` |
| Tag | `` Nome string `json:"nome"` `` |
| Construtor | `func NovoProduto(...) *Produto` |

---

## ✍️ Exercícios

1. Adicione tags `json` à struct `Produto` do seu `main.go` e imprima a lista de produtos em JSON com `json.MarshalIndent(produtos, "", "  ")`.
2. Converta este JSON para struct e imprima com `%+v`: `{"nome":"Ana","idade":30,"email":"ana@email.com"}`.
3. Crie uma struct `Usuario` com um campo `Senha` que **nunca** vá para o JSON.
4. Crie o tipo `Pilha` (com slice interno privado) cujo **valor zero** já funcione, com os métodos `Empilhar`, `Desempilhar` e `Tamanho`.
5. Use `unsafe.Sizeof` para comparar duas structs com os mesmos campos em ordens diferentes.
6. Implemente o padrão de **opções funcionais** para um `Cliente HTTP` com `timeout`, `tentativas` e `baseURL`.
7. **Projeto:** refatore o `main.go` da livraria aplicando o que aprendeu: construtor `NovaLivraria()`, métodos em `*Livraria`, `String()` em `Produto` e tags JSON.

---

⬅️ Anterior: [Embedding e composição](03-embedding-e-composicao.md) · 🏠 [Voltar ao roteiro](../README.md)

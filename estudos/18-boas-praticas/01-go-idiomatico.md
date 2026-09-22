# 01 — Go Idiomático

Escrever Go **que compila** é fácil. Escrever Go **idiomático** (do jeito que a comunidade escreve e espera ler) é o que diferencia um código "traduzido de Java" de um código **Go de verdade**.

> 📖 Leitura obrigatória: [Effective Go](https://go.dev/doc/effective_go) e [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).

---

## 🧘 Os Provérbios de Go

Frases de Rob Pike que resumem a filosofia da linguagem:

| Provérbio | Significado |
|---|---|
| *Don't communicate by sharing memory, share memory by communicating.* | Prefira channels a variáveis compartilhadas (módulo 14) |
| *Concurrency is not parallelism.* | Concorrência é estrutura; paralelismo é execução |
| *The bigger the interface, the weaker the abstraction.* | Interfaces pequenas (módulo 10) |
| *Make the zero value useful.* | Tipos que funcionam sem inicialização (módulo 09) |
| *`interface{}` says nothing.* | `any` não diz nada sobre o valor: evite |
| *Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.* | Um estilo único acaba com discussões |
| *A little copying is better than a little dependency.* | Não importe uma biblioteca por causa de 5 linhas |
| *Clear is better than clever.* | Código claro vence código "esperto" |
| *Errors are values.* | Erros são valores comuns, trate-os com código comum (módulo 11) |
| *Don't just check errors, handle them gracefully.* | Adicione contexto, decida o que fazer |
| *Don't panic.* | `panic` só para o impossível |

---

## 🏷️ Nomes

### Curtos em escopos pequenos, descritivos em escopos grandes

```go
// ✅ escopo pequeno: nome curto
for i, p := range produtos {
	total += p.Preco
}

// ✅ escopo grande (pacote): nome descritivo
var tempoMaximoDeConexao = 30 * time.Second
```

### Convenções

| Situação | ✅ Idiomático | ❌ Evite |
|---|---|---|
| Composição de palavras | `nomeCompleto`, `NomeCompleto` | `nome_completo` |
| Siglas | `userID`, `HTTPServer`, `parseURL` | `userId`, `HttpServer`, `parseUrl` |
| Receptor | `func (p *Produto)` | `func (this *Produto)`, `func (self *Produto)` |
| Getter | `p.Nome()` | `p.GetNome()` |
| Interface de 1 método | `Reader`, `Validador` | `IReader`, `ReaderInterface` |
| Pacote | `produto`, `http` | `produtoUtils`, `common` |
| Repetição | `produto.Novo()` | `produto.NovoProduto()` |
| Erros | `ErrNaoEncontrado`, `ErroValidacao` | `NaoEncontradoError` misturado |
| Booleanos | `ativo`, `temEstoque`, `podeEditar` | `flag`, `status` |

---

## ↩️ Retorno antecipado e o "caminho feliz" à esquerda

```go
// ❌ Pirâmide
func processar(p *Produto) error {
	if p != nil {
		if p.Preco > 0 {
			if p.Estoque > 0 {
				// ... lógica ...
				return nil
			} else {
				return errors.New("sem estoque")
			}
		} else {
			return errors.New("preço inválido")
		}
	} else {
		return errors.New("produto nulo")
	}
}

// ✅ Guard clauses: trata os problemas e sai; a lógica principal fica sem indentação
func processar(p *Produto) error {
	if p == nil {
		return errors.New("produto nulo")
	}
	if p.Preco <= 0 {
		return errors.New("preço inválido")
	}
	if p.Estoque <= 0 {
		return errors.New("sem estoque")
	}
	// ... lógica ...
	return nil
}
```

---

## 🧱 Estruturas e tipos

```go
// ✅ Valor zero útil
var sb strings.Builder // pronto para uso
var mu sync.Mutex      // pronto para uso

// ✅ Construtor só quando necessário
func NovaLoja() *Loja { return &Loja{produtos: make(map[string]Produto)} }

// ✅ Aceite interfaces, retorne structs
func Salvar(w io.Writer, p Produto) error
func NovoServidor(cfg Config) *Servidor

// ✅ Interfaces definidas por quem usa, pequenas
type buscador interface {
	Buscar(codigo string) (Produto, error)
}
```

---

## ⚠️ Erros

```go
// ✅ Contexto com %w, minúsculas, sem pontuação
return fmt.Errorf("carregar produtos de %s: %w", caminho, err)

// ✅ errors.Is / errors.As
if errors.Is(err, fs.ErrNotExist) { ... }

// ✅ Trate uma vez: ou loga, ou retorna
// ✅ nil literal quando não há erro
// ❌ nunca ignore com _ sem motivo claro
```

---

## 💬 Comentários

```go
// Package produto define os produtos vendidos pela livraria.
package produto

// Produto representa um livro à venda.
type Produto struct { ... }

// Novo cria um produto validando o preço.
// Retorna ErrPrecoInvalido se o preço não for positivo.
func Novo(...) (*Produto, error)
```

Regras:
- Todo item **exportado** deve ter comentário
- O comentário **começa com o nome** do item
- Explique o **porquê**, não o óbvio:

```go
// ❌ incrementa i
i++

// ✅ Pulamos o cabeçalho do CSV exportado pelo sistema legado.
i++
```

---

## 🧹 Simplicidade acima de tudo

```go
// ❌ "Esperto"
func ehPar(n int) bool { return n&1 == 0 }

// ✅ Claro (e o compilador otimiza igual)
func ehPar(n int) bool { return n%2 == 0 }
```

```go
// ❌ Abstração prematura
type ProdutoFactory interface { Criar() Produto }
type ProdutoFactoryImpl struct{}
type ProdutoServiceInterface interface { ... }
type AbstractProdutoRepository struct { ... }

// ✅ Comece concreto. Abstraia quando houver uma segunda implementação real.
func NovoProduto(...) *Produto
```

> 🧠 Em Go, **menos é mais**: poucas camadas, poucas abstrações, nenhuma "mágica". Quem ler seu código daqui a 6 meses (inclusive você) agradece.

---

## 📋 Checklist de código idiomático

- [ ] Passou pelo `gofmt` / `goimports`
- [ ] `go vet` sem avisos
- [ ] Nomes curtos, em camelCase, siglas inteiras (`ID`, `URL`)
- [ ] Sem `Get` nos getters, sem `this`/`self`
- [ ] Erros tratados com contexto (`%w`), sem `_` escondendo falhas
- [ ] Guard clauses em vez de `if` aninhado
- [ ] Interfaces pequenas, no consumidor
- [ ] `context.Context` como primeiro parâmetro em operações demoradas
- [ ] Goroutines com fim claro, sem vazamentos
- [ ] Comentários de documentação nos itens exportados
- [ ] Testes (table tests) para a lógica importante

---

## ✍️ Exercícios

1. Pegue o `main.go` original da livraria e liste tudo que **não** está idiomático segundo este arquivo.
2. Renomeie para o estilo Go: `get_user_id`, `HttpClient`, `IRepository`, `UserServiceImpl`, `this.name`.
3. Reescreva uma função com 3 níveis de `if` aninhados usando guard clauses.
4. Escolha 3 provérbios de Go e escreva, com suas palavras, um exemplo de código que mostre cada um.
5. Leia a seção "Names" do *Effective Go* e anote 3 coisas que você não sabia.

---

🏠 [Módulo 18](README.md) · ➡️ Próximo: [Ferramentas de qualidade](02-ferramentas-de-qualidade.md)

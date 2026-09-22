# 04 — Interfaces e `nil` + Boas Práticas

## 🕳️ A pegadinha mais famosa de Go: interface "nil" que não é `nil`

Lembra que uma interface guarda **(tipo, valor)**? Ela só é `nil` quando **os dois** são `nil`:

```go
var e error                  // (tipo: nil, valor: nil)
fmt.Println(e == nil)        // true ✅

var p *MeuErro = nil
e = p                        // (tipo: *MeuErro, valor: nil)
fmt.Println(e == nil)        // false 😱
```

```
var e error          e = (*MeuErro)(nil)
┌──────────────┐     ┌────────────────────┐
│ tipo:  nil   │     │ tipo:  *MeuErro    │  ← não é nil!
│ valor: nil   │     │ valor: nil         │
└──────────────┘     └────────────────────┘
   == nil ✅             == nil ❌
```

### Como isso vira bug na vida real

```go
type MeuErro struct{ Msg string }

func (e *MeuErro) Error() string { return e.Msg }

func validar(idade int) error {
	var err *MeuErro // nil do tipo *MeuErro
	if idade < 0 {
		err = &MeuErro{"idade negativa"}
	}
	return err // ⚠️ SEMPRE retorna uma interface com tipo preenchido!
}

func main() {
	if err := validar(30); err != nil {
		fmt.Println("deu erro?!", err) // 😱 entra aqui mesmo com idade válida
	}
}
```

### ✅ A correção: retorne `nil` **literal**

```go
func validar(idade int) error {
	if idade < 0 {
		return &MeuErro{"idade negativa"}
	}
	return nil // ✅ interface realmente nil
}
```

> 🧠 **Regra:** funções que retornam `error` (ou qualquer interface) devem retornar **`nil` explícito**, nunca uma variável de ponteiro que "por acaso" é nil.

---

## 🧰 Chamando métodos em receptor `nil`

Curiosidade: métodos com receptor ponteiro **podem** ser chamados com ponteiro `nil`, desde que não acessem os campos:

```go
type Lista struct {
	valor int
	prox  *Lista
}

func (l *Lista) Soma() int {
	if l == nil { // trata o caso nil
		return 0
	}
	return l.valor + l.prox.Soma()
}

var l *Lista
fmt.Println(l.Soma()) // 0 (sem panic!)
```

Mas chamar um método numa **interface** `nil` (tipo e valor nil) **sempre** dá panic:

```go
var f Forma
f.Area() // 💥 panic: runtime error: invalid memory address or nil pointer dereference
```

---

## ✅ Boas práticas com interfaces

### 1. Interfaces pequenas

```go
// ❌ Interface "gorda": difícil de implementar, difícil de testar
type Repositorio interface {
	Criar(p Produto) error
	Buscar(id string) (Produto, error)
	Listar() ([]Produto, error)
	Atualizar(p Produto) error
	Remover(id string) error
	Contar() int
	Exportar(w io.Writer) error
}

// ✅ Interfaces pequenas, combináveis
type BuscadorDeProduto interface {
	Buscar(id string) (Produto, error)
}
```

### 2. Defina a interface **onde ela é usada** (no consumidor)

Em Java, quem **implementa** declara a interface. Em Go, o idioma é o contrário: quem **usa** declara só o que precisa.

```go
// pacote relatorio (o consumidor)
type fonteDeVendas interface {
	VendasDoMes(mes int) []Venda // só o que o relatório precisa
}

func Gerar(f fonteDeVendas, mes int) string { ... }
```

Qualquer tipo que tenha `VendasDoMes` (um banco de dados real, um CSV, um mock de teste) serve.

### 3. "Aceite interfaces, retorne structs"

```go
// ✅ Recebe o comportamento mínimo necessário
func Processar(r io.Reader) error { ... }

// ✅ Retorna o tipo concreto (quem chama tem acesso a tudo)
func NovoServidor(porta int) *Servidor { ... }
```

Exceção comum: retornar `error` (que é interface).

### 4. Não crie interface antes de precisar

```go
// ❌ Interface com uma única implementação e "por via das dúvidas"
type ServicoDeUsuario interface { ... }
type servicoDeUsuarioImpl struct { ... }
```

Comece com a **struct concreta**. Crie a interface quando surgir uma **segunda implementação** ou a necessidade de **testar com mock**. Como a implementação é implícita, adicionar a interface depois **não quebra nada**.

### 5. Interfaces facilitam testes

```go
type Relogio interface {
	Agora() time.Time
}

type relogioReal struct{}
func (relogioReal) Agora() time.Time { return time.Now() }

type relogioFalso struct{ t time.Time }
func (r relogioFalso) Agora() time.Time { return r.t }

func Saudacao(r Relogio) string {
	if r.Agora().Hour() < 12 {
		return "Bom dia"
	}
	return "Boa tarde"
}

// Em produção:
Saudacao(relogioReal{})

// Em teste: horário controlado!
manha := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
Saudacao(relogioFalso{manha}) // "Bom dia" (sempre)
```

---

## 🧾 Resumão do módulo

| Conceito | Resumo |
|---|---|
| Interface | conjunto de métodos (comportamento) |
| Implementação | **implícita**: basta ter os métodos |
| Valor de interface | par (tipo dinâmico, valor dinâmico) |
| `any` | aceita qualquer tipo (use com moderação) |
| Type assertion | `v, ok := x.(T)` |
| Type switch | `switch v := x.(type) { case T: ... }` |
| Composição | interfaces embutem interfaces |
| Receptor ponteiro | só `*T` satisfaz a interface |
| Pegadinha do nil | ponteiro nil dentro de interface ≠ nil |
| Filosofia | interfaces pequenas, no consumidor, "aceite interfaces, retorne structs" |

---

## ✍️ Exercícios

1. Reproduza a pegadinha do "nil que não é nil" com uma função que retorna `error`. Depois corrija.
2. Imprima `%T` e `%v` de uma interface `error` nos dois casos (nil de verdade e ponteiro nil dentro).
3. Crie uma lista encadeada com um método `Tamanho()` que funcione mesmo com receptor `nil`.
4. Refatore esta interface "gorda" em interfaces menores e mostre uma função que só precisa de uma delas:
   ```go
   type Arquivo interface {
       Abrir() error
       Ler() ([]byte, error)
       Escrever([]byte) error
       Fechar() error
   }
   ```
5. Crie a interface `Armazenamento` com `Salvar(chave, valor string)` e `Carregar(chave string) (string, bool)`, com duas implementações: `MemoriaStorage` (map) e `ArquivoStorage` (arquivo). Faça uma função que funcione com as duas.
6. **Projeto:** na livraria do seu `main.go`, crie a interface `Repositorio` (`Salvar`, `Buscar`, `Listar`) e uma implementação em memória. Faça a `Livraria` depender da **interface**, não do map.

---

⬅️ Anterior: [Composição e interfaces famosas](03-composicao-e-interfaces-famosas.md) · ➡️ Próximo módulo: [Tratamento de erros](../11-tratamento-de-erros/README.md)

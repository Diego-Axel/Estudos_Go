# 04 — Fuzzing, Mocks e Boas Práticas

## 🎲 Fuzzing (Go 1.18+)

Nos testes normais, **você** escolhe as entradas. No **fuzzing**, o Go **gera milhares de entradas aleatórias** automaticamente, tentando quebrar seu código. Ótimo para achar casos que você nunca imaginaria.

### Exemplo: uma função com bug escondido

```go
// Inverter inverte uma string.
func Inverter(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
```

Um teste comum passa tranquilo:

```go
func TestInverter(t *testing.T) {
	if got := Inverter("Go"); got != "oG" {
		t.Errorf("Inverter(\"Go\") = %q", got)
	}
}
```

### O fuzz test

Um fuzz test verifica **propriedades** que devem valer para **qualquer** entrada:

```go
func FuzzInverter(f *testing.F) {
	// Corpus inicial: exemplos para o fuzzer começar
	f.Add("Go")
	f.Add("olá")
	f.Add("")

	f.Fuzz(func(t *testing.T, original string) {
		invertida := Inverter(original)
		duasVezes := Inverter(invertida)

		// Propriedade 1: inverter duas vezes volta ao original
		if original != duasVezes {
			t.Errorf("antes: %q, depois: %q", original, duasVezes)
		}

		// Propriedade 2: se a entrada é UTF-8 válido, a saída também deve ser
		if utf8.ValidString(original) && !utf8.ValidString(invertida) {
			t.Errorf("Inverter gerou UTF-8 inválido a partir de %q: %q", original, invertida)
		}
	})
}
```

```bash
go test -fuzz=FuzzInverter
```

```
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, execs: 12 (...), new interesting: 1
--- FAIL: FuzzInverter (0.02s)
    --- FAIL: FuzzInverter (0.00s)
        Inverter gerou UTF-8 inválido a partir de "á": "\xa1\xc3"

    Failing input written to testdata/fuzz/FuzzInverter/a1b2c3...
```

Em milissegundos o fuzzer descobriu: inverter **bytes** quebra os acentos (lembra do módulo 07?). A correção é inverter **runes**:

```go
func Inverter(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```

> A entrada que falhou é salva em `testdata/fuzz/...` e passa a ser testada **sempre** pelo `go test` normal. Faça commit dela!

### Regras do fuzzing

- A função começa com **`Fuzz`** e recebe `f *testing.F`
- Tipos aceitos nos argumentos: `string`, `[]byte`, inteiros, floats, `bool`, `rune`, `byte`
- Sem `-fuzz`, o `go test` roda só o corpus inicial (como um teste comum)
- Use `-fuzztime=30s` para limitar o tempo (senão roda até achar erro ou você parar com Ctrl+C)

> 💡 Ótimo para testar **parsers**, **validadores**, **conversores** e qualquer coisa que receba entrada do usuário.

---

## 🎭 Mocks com interfaces

Como testar código que depende de **banco de dados**, **API externa**, **e-mail** ou **relógio**? Com **interfaces** (módulo 10)! O código depende da **interface**, e no teste você passa uma implementação **falsa**.

### O código de produção

```go
package pedido

type Notificador interface {
	Enviar(para, msg string) error
}

type Servico struct {
	notificador Notificador
}

func NovoServico(n Notificador) *Servico {
	return &Servico{notificador: n}
}

func (s *Servico) Finalizar(cliente string, total float64) error {
	if total <= 0 {
		return errors.New("total inválido")
	}
	msg := fmt.Sprintf("Pedido confirmado: R$ %.2f", total)
	if err := s.notificador.Enviar(cliente, msg); err != nil {
		return fmt.Errorf("notificar cliente: %w", err)
	}
	return nil
}
```

Em produção, o `Notificador` é um envio de e-mail de verdade. No teste:

### O mock

```go
package pedido

import (
	"errors"
	"testing"
)

// notificadorFalso registra as mensagens em vez de enviar de verdade.
type notificadorFalso struct {
	enviadas []string
	err      error // para simular falhas
}

func (n *notificadorFalso) Enviar(para, msg string) error {
	if n.err != nil {
		return n.err
	}
	n.enviadas = append(n.enviadas, para+": "+msg)
	return nil
}

func TestFinalizar_EnviaNotificacao(t *testing.T) {
	falso := &notificadorFalso{}
	s := NovoServico(falso)

	if err := s.Finalizar("ana@email.com", 99.9); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(falso.enviadas) != 1 {
		t.Fatalf("esperado 1 notificação, obtido %d", len(falso.enviadas))
	}
	esperado := "ana@email.com: Pedido confirmado: R$ 99.90"
	if falso.enviadas[0] != esperado {
		t.Errorf("mensagem = %q; esperado %q", falso.enviadas[0], esperado)
	}
}

func TestFinalizar_FalhaNaNotificacao(t *testing.T) {
	errSMTP := errors.New("servidor fora do ar")
	s := NovoServico(&notificadorFalso{err: errSMTP})

	err := s.Finalizar("ana@email.com", 50)
	if !errors.Is(err, errSMTP) {
		t.Errorf("esperado erro de SMTP embrulhado, obtido %v", err)
	}
}
```

Sem banco, sem internet, sem e-mail de verdade: **rápido e determinístico**. ✅

> 💡 Existem geradores de mocks (`gomock`, `mockery`), mas em Go é muito comum escrever os falsos **à mão**, como acima. São simples e claros.

---

## 🌐 Testando HTTP: `net/http/httptest`

A biblioteca padrão já traz ferramentas para testar servidores e clientes HTTP:

```go
func OlaHandler(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	fmt.Fprintf(w, "Olá, %s!", nome)
}

func TestOlaHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/ola?nome=Diego", nil)
	rec := httptest.NewRecorder() // um ResponseWriter "gravador"

	OlaHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; esperado 200", rec.Code)
	}
	if got := rec.Body.String(); got != "Olá, Diego!" {
		t.Errorf("corpo = %q", got)
	}
}
```

> Isso será útil no módulo 17 (Web com `net/http`). Também existe `httptest.NewServer` para subir um servidor de teste de verdade.

---

## ⏳ Testando código concorrente: `testing/synctest` (Go 1.25+)

Testar código com `time.Sleep` e timeouts costuma deixar os testes **lentos** e **instáveis**. O pacote `testing/synctest` cria uma "bolha" com um **relógio falso**: o tempo só avança quando todas as goroutines da bolha estão bloqueadas.

```go
func TestTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		<-ctx.Done() // "espera" 5 segundos... instantaneamente!
		if ctx.Err() != context.DeadlineExceeded {
			t.Error("esperado deadline exceeded")
		}
	})
}
```

---

## ✅ Boas práticas de testes

1. **Teste comportamento, não implementação.** Teste o que a função **faz** (entradas → saídas), não **como** faz. Assim você pode refatorar sem reescrever testes.
2. **Use table tests** para vários cenários.
3. **Teste os caminhos de erro**, não só o caminho feliz.
4. **Mensagens de falha claras**: `Funcao(entrada) = obtido; esperado X`.
5. **Testes independentes**: um teste não pode depender de outro ter rodado antes.
6. **Testes rápidos**: use mocks para dependências externas; separe testes lentos com `testing.Short()`.
7. **Sem `time.Sleep`** para "esperar" goroutines: use channels, `WaitGroup` ou `synctest`.
8. **Rode com `-race`** sempre que houver concorrência.
9. **Dependa de interfaces** no código de produção: é o que torna o código **testável**.
10. **Rode os testes no CI** (GitHub Actions, por exemplo) a cada push.

### Bônus: GitHub Actions para rodar os testes

Crie `.github/workflows/testes.yml` no seu repositório:

```yaml
name: Testes
on: [push, pull_request]

jobs:
  testes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: go vet ./...
      - run: go test -race -cover ./...
```

A cada push, o GitHub roda os testes e mostra ✅ ou ❌ no commit.

---

## 🧾 Resumão do módulo

| Recurso | Assinatura / Comando |
|---|---|
| Teste | `func TestX(t *testing.T)` |
| Rodar | `go test ./...` / `-v` / `-run` |
| Falha | `t.Errorf` (continua) / `t.Fatalf` (para) |
| Subteste | `t.Run(nome, func(t *testing.T) {...})` |
| Helper | `t.Helper()` |
| Paralelo | `t.Parallel()` |
| Limpeza | `t.Cleanup`, `t.TempDir()` |
| Cobertura | `go test -cover` / `-coverprofile` |
| Benchmark | `func BenchmarkX(b *testing.B)` + `for b.Loop()` |
| Exemplo | `func ExampleX()` + `// Output:` |
| Fuzz | `func FuzzX(f *testing.F)` + `go test -fuzz=FuzzX` |
| Mock | interface + implementação falsa |
| HTTP | `httptest.NewRecorder`, `httptest.NewRequest` |

---

## ✍️ Exercícios

1. Escreva a versão com bug do `Inverter`, o fuzz test, e deixe o fuzzer encontrar o problema. Depois corrija.
2. Escreva um fuzz test para `strconv.Itoa` + `strconv.Atoi`: converter um int para string e de volta deve dar o mesmo número.
3. Crie a interface `Repositorio` para a livraria e um repositório **falso** em memória. Teste o cadastro e a busca usando o falso.
4. Crie um mock de `Relogio` (módulo 10) e teste uma função de saudação às 9h e às 15h.
5. Escreva um handler HTTP que retorna JSON e teste com `httptest`.
6. Configure o GitHub Actions deste repositório para rodar `go test ./...` (depois de criar o `go.mod`, módulo 12).
7. **Projeto:** escreva testes para toda a livraria refatorada. Tente chegar a mais de 80% de cobertura **verificando os resultados** de verdade.

---

⬅️ Anterior: [Cobertura, benchmarks e exemplos](03-cobertura-benchmarks-e-exemplos.md) · 🏠 [Voltar ao roteiro](../README.md)

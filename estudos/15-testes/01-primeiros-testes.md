# 01 — Primeiros Testes

Go tem testes **embutidos na linguagem**: o pacote `testing` e o comando `go test`. Não precisa instalar framework nenhum.

---

## 📁 Convenções

| Regra | Exemplo |
|---|---|
| Arquivo de teste termina em **`_test.go`** | `calc.go` → `calc_test.go` |
| Fica **na mesma pasta** do código testado | `calc/calc.go` e `calc/calc_test.go` |
| Função começa com **`Test`** + nome com maiúscula | `TestSomar`, `TestDividirPorZero` |
| Recebe **`t *testing.T`** | `func TestSomar(t *testing.T)` |

> Os arquivos `_test.go` **não** entram no executável final: só são compilados pelo `go test`.

---

## ✍️ Primeiro teste

**`calc/calc.go`**

```go
package calc

import "errors"

var ErrDivisaoPorZero = errors.New("divisão por zero")

func Somar(a, b int) int {
	return a + b
}

func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisaoPorZero
	}
	return a / b, nil
}
```

**`calc/calc_test.go`**

```go
package calc

import "testing"

func TestSomar(t *testing.T) {
	resultado := Somar(2, 3)
	esperado := 5

	if resultado != esperado {
		t.Errorf("Somar(2, 3) = %d; esperado %d", resultado, esperado)
	}
}
```

Rodando:

```bash
go test ./calc
```

```
ok      meuprojeto/calc    0.002s
```

Se o teste falhar:

```
--- FAIL: TestSomar (0.00s)
    calc_test.go:9: Somar(2, 3) = 6; esperado 5
FAIL
FAIL    meuprojeto/calc    0.002s
```

> 🧠 Não existe `assertEquals` em Go: você usa um **`if` comum** e reporta a falha. Simples e explícito.
>
> 📌 Convenção da mensagem: **`Funcao(entrada) = obtido; esperado X`** (ou em inglês, `got X, want Y`).

---

## 📣 Reportando falhas

| Método | Marca como falha? | Para o teste? |
|---|---|---|
| `t.Error(...)` / `t.Errorf(...)` | ✅ | ❌ continua |
| `t.Fatal(...)` / `t.Fatalf(...)` | ✅ | ✅ para na hora |
| `t.Log(...)` / `t.Logf(...)` | ❌ | ❌ (só aparece com `-v` ou se falhar) |
| `t.Skip(...)` | ❌ | ✅ pula o teste |

Use **`Fatal`** quando não faz sentido continuar (ex: a preparação falhou):

```go
func TestDividir(t *testing.T) {
	resultado, err := Dividir(10, 2)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err) // sem resultado, não adianta seguir
	}
	if resultado != 5 {
		t.Errorf("Dividir(10, 2) = %v; esperado 5", resultado)
	}
}
```

### Testando o caso de erro

```go
func TestDividirPorZero(t *testing.T) {
	_, err := Dividir(10, 0)
	if !errors.Is(err, ErrDivisaoPorZero) {
		t.Errorf("esperado ErrDivisaoPorZero, obtido %v", err)
	}
}
```

> ✅ **Teste os caminhos de erro também**, não só o "caminho feliz".

---

## 🏃 O comando `go test`

| Comando | O que faz |
|---|---|
| `go test` | testa o pacote da pasta atual |
| `go test ./...` | testa **todos** os pacotes do módulo ⭐ |
| `go test -v` | mostra cada teste (verbose) |
| `go test -run TestSomar` | roda só os testes cujo nome bate com a regex |
| `go test -run 'Dividir'` | todos que contêm "Dividir" |
| `go test -count=1` | ignora o cache (força rodar de novo) |
| `go test -race` | ativa o detector de corrida (módulo 14) |
| `go test -short` | modo rápido (veja `testing.Short()` abaixo) |
| `go test -failfast` | para no primeiro teste que falhar |

Saída com `-v`:

```
=== RUN   TestSomar
--- PASS: TestSomar (0.00s)
=== RUN   TestDividir
--- PASS: TestDividir (0.00s)
=== RUN   TestDividirPorZero
--- PASS: TestDividirPorZero (0.00s)
PASS
ok      meuprojeto/calc    0.003s
```

> 💡 O Go **guarda em cache** resultados de testes que passaram. Se nada mudou, aparece `(cached)`. Use `-count=1` para forçar.

### Pulando testes lentos

```go
func TestIntegracaoBanco(t *testing.T) {
	if testing.Short() {
		t.Skip("pulando teste lento em modo -short")
	}
	// ... teste que demora ...
}
```

---

## 📦 `package calc` vs. `package calc_test`

Há duas formas de declarar o pacote do arquivo de teste:

| | `package calc` | `package calc_test` |
|---|---|---|
| Tipo de teste | **caixa-branca** | **caixa-preta** |
| Acessa itens **privados**? | ✅ sim | ❌ não, só os exportados |
| Como chama | `Somar(2, 3)` | `calc.Somar(2, 3)` (precisa importar) |
| Bom para | testar detalhes internos | testar a **API pública**, como um usuário faria |

```go
package calc_test

import (
	"testing"

	"meuprojeto/calc"
)

func TestSomarPublico(t *testing.T) {
	if calc.Somar(1, 1) != 2 {
		t.Error("1 + 1 deveria ser 2")
	}
}
```

> 🧠 Prefira `_test` (caixa-preta) quando possível: o teste fica menos acoplado à implementação.

---

## 🧪 Testando a livraria (exemplo)

```go
package produto

import "testing"

func TestNovo_PrecoInvalido(t *testing.T) {
	_, err := Novo("L001", "Livro", -10, 1)
	if err == nil {
		t.Fatal("esperado erro para preço negativo, mas veio nil")
	}
}

func TestNovo_Valido(t *testing.T) {
	p, err := Novo("L001", "O Hobbit", 49.9, 3)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if p.Titulo != "O Hobbit" {
		t.Errorf("Titulo = %q; esperado %q", p.Titulo, "O Hobbit")
	}
}
```

> O `_` no nome (`TestNovo_PrecoInvalido`) é uma convenção comum para indicar o **cenário** testado.

---

## ✍️ Exercícios

1. Crie um pacote `calc` com `Somar`, `Subtrair`, `Multiplicar` e `Dividir`, e um teste para cada.
2. Faça um teste falhar de propósito e leia a mensagem. Depois conserte.
3. Teste o caso de erro do `Dividir` usando `errors.Is`.
4. Rode com `go test -v`, depois `go test -run Dividir -v`. Compare.
5. Escreva `EhPalindromo(s string) bool` e teste com `"arara"`, `"go"` e `""`.
6. Crie um teste em `package calc_test` e tente acessar uma função privada. O que acontece?
7. Escreva testes para as funções de cadastro e busca da sua livraria.

---

🏠 [Módulo 15](README.md) · ➡️ Próximo: [Table tests e subtestes](02-table-tests-e-subtestes.md)

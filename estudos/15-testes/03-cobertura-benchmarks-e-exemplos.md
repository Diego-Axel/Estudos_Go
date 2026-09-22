# 03 — Cobertura, Benchmarks e Testes de Exemplo

## 📊 Cobertura de testes

Quanto do seu código os testes **executam**?

```bash
go test -cover ./...
```

```
ok      meuprojeto/calc      0.002s  coverage: 85.7% of statements
ok      meuprojeto/produto   0.003s  coverage: 100.0% of statements
```

### Relatório visual em HTML ⭐

```bash
go test -coverprofile=cobertura.out ./...
go tool cover -html=cobertura.out
```

Abre no navegador o código **colorido**: 🟩 verde = testado, 🟥 vermelho = nunca executado pelos testes.

Resumo por função no terminal:

```bash
go tool cover -func=cobertura.out
```

```
meuprojeto/calc/calc.go:7:    Somar      100.0%
meuprojeto/calc/calc.go:11:   Dividir     66.7%
total:                        (statements) 85.7%
```

> ⚠️ **Cobertura alta ≠ testes bons.** 100% de cobertura só diz que as linhas **rodaram**, não que os resultados foram **verificados**. Use a cobertura para achar **o que ficou sem teste** (principalmente caminhos de erro), não como meta cega.

---

## ⏱️ Benchmarks

Benchmarks medem **desempenho**. Ficam nos arquivos `_test.go` e começam com **`Benchmark`**:

### Forma moderna: `b.Loop()` (Go 1.24+) ⭐

```go
func BenchmarkSomar(b *testing.B) {
	for b.Loop() {
		Somar(2, 3)
	}
}
```

### Forma clássica: `b.N`

Você vai ver muito em código existente:

```go
func BenchmarkSomar(b *testing.B) {
	for i := 0; i < b.N; i++ { // b.N é ajustado automaticamente
		Somar(2, 3)
	}
}
```

> O `b.Loop()` é preferível: ele evita que o compilador "otimize fora" o código medido e exclui automaticamente a preparação feita **antes** do laço.

### Rodando

```bash
go test -bench=.                  # roda todos os benchmarks
go test -bench=Somar              # só os que batem com a regex
go test -bench=. -benchmem        # mostra também uso de memória
go test -bench=. -run=^$          # só benchmarks (nenhum teste comum)
```

```
goos: windows
goarch: amd64
BenchmarkSomar-8          1000000000     0.2531 ns/op     0 B/op    0 allocs/op
```

| Coluna | Significado |
|---|---|
| `-8` | GOMAXPROCS (núcleos usados) |
| `1000000000` | quantas vezes o laço rodou |
| `ns/op` | **tempo** médio por operação |
| `B/op` | **bytes** alocados por operação |
| `allocs/op` | **alocações** no heap por operação |

---

## 🏎️ Exemplo real: comparando formas de concatenar

Lembra do módulo 07? Vamos **provar** que `strings.Builder` é mais rápido que `+=`:

```go
package texto

import (
	"strings"
	"testing"
)

func concatenarMais(n int) string {
	s := ""
	for range n {
		s += "x"
	}
	return s
}

func concatenarBuilder(n int) string {
	var sb strings.Builder
	for range n {
		sb.WriteString("x")
	}
	return sb.String()
}

func BenchmarkConcatenarMais(b *testing.B) {
	for b.Loop() {
		concatenarMais(1000)
	}
}

func BenchmarkConcatenarBuilder(b *testing.B) {
	for b.Loop() {
		concatenarBuilder(1000)
	}
}
```

```bash
go test -bench=Concatenar -benchmem
```

Resultado típico (os números variam por máquina):

```
BenchmarkConcatenarMais-8         5000    250000 ns/op   530000 B/op   999 allocs/op
BenchmarkConcatenarBuilder-8    300000      3500 ns/op     2000 B/op    11 allocs/op
```

O `Builder` é **dezenas de vezes** mais rápido e faz **muito** menos alocações. 🚀

### Benchmarks com tamanhos diferentes

```go
func BenchmarkBuilder(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for b.Loop() {
				concatenarBuilder(n)
			}
		})
	}
}
```

### Comparando antes e depois

Para comparar com rigor estatístico, use a ferramenta `benchstat`:

```bash
go install golang.org/x/perf/cmd/benchstat@latest

go test -bench=. -count=10 > antes.txt
# ... faz a otimização ...
go test -bench=. -count=10 > depois.txt
benchstat antes.txt depois.txt
```

> 🧠 **Regra:** nunca otimize sem medir. Benchmarks transformam "acho que é mais rápido" em **dados**.

---

## 📖 Testes de exemplo (`Example`)

Funções que começam com **`Example`** servem **ao mesmo tempo** como:
1. **Documentação** (aparecem no `go doc` e no pkg.go.dev)
2. **Testes** (o `go test` compara a saída com o comentário `// Output:`)

```go
package calc_test

import (
	"fmt"

	"meuprojeto/calc"
)

func ExampleSomar() {
	fmt.Println(calc.Somar(2, 3))
	// Output: 5
}

func ExampleDividir() {
	r, err := calc.Dividir(10, 4)
	fmt.Println(r, err)

	_, err = calc.Dividir(1, 0)
	fmt.Println(err)
	// Output:
	// 2.5 <nil>
	// divisão por zero
}
```

- Se a saída for **diferente** do `// Output:`, o teste **falha**
- **Sem** o comentário `// Output:`, o exemplo só é compilado (não executado)
- Para saídas em ordem imprevisível (ex: maps), use `// Unordered output:`

Convenção de nomes:

| Nome | Documenta |
|---|---|
| `Example()` | o pacote |
| `ExampleSomar()` | a função `Somar` |
| `ExampleProduto()` | o tipo `Produto` |
| `ExampleProduto_String()` | o método `String` de `Produto` |
| `ExampleSomar_negativos()` | um segundo exemplo de `Somar` (sufixo minúsculo) |

> ✅ Exemplos são **documentação que nunca fica desatualizada**, porque se ficar, o teste falha!

---

## ✍️ Exercícios

1. Rode `go test -cover` nos seus testes da calculadora. Qual a cobertura?
2. Gere o relatório HTML e encontre as linhas **não** testadas. Escreva testes para elas.
3. Escreva os benchmarks de `concatenarMais` vs. `concatenarBuilder` e rode com `-benchmem`.
4. Compare com benchmarks: `fmt.Sprintf("%d", n)` vs. `strconv.Itoa(n)`. Qual é mais rápido?
5. Compare: busca num `[]int` com `slices.Contains` vs. busca num `map[int]bool`, para 10, 1.000 e 100.000 elementos (use `b.Run`).
6. Escreva `ExampleSomar` e `ExampleDividir` e veja-os no `go doc`.
7. Faça um `Example` falhar de propósito mudando o `// Output:`.

---

⬅️ Anterior: [Table tests e subtestes](02-table-tests-e-subtestes.md) · ➡️ Próximo: [Fuzzing, mocks e boas práticas](04-fuzzing-mocks-e-boas-praticas.md)

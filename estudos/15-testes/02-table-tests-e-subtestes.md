# 02 — Table Tests, Subtestes e Helpers

## 📋 Table-driven tests ⭐

Em vez de escrever um teste para cada caso, monte uma **tabela** de casos e percorra com um laço. É **o estilo mais usado** em Go.

```go
func TestSomar(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     int
		esperado int
	}{
		{"positivos", 2, 3, 5},
		{"com zero", 0, 7, 7},
		{"negativos", -2, -3, -5},
		{"misto", -4, 10, 6},
	}

	for _, c := range casos {
		resultado := Somar(c.a, c.b)
		if resultado != c.esperado {
			t.Errorf("%s: Somar(%d, %d) = %d; esperado %d",
				c.nome, c.a, c.b, resultado, c.esperado)
		}
	}
}
```

Vantagens:
- Adicionar um caso novo = adicionar **uma linha**
- Fácil de ver **quais cenários** estão cobertos
- A lógica de verificação é escrita **uma vez só**

> 💡 Isso usa uma **struct anônima** (módulo 09). Se preferir, `casos` também pode ser um `map[string]struct{...}` com o nome como chave.

---

## 🌿 Subtestes com `t.Run`

Com `t.Run`, cada caso vira um **subteste** com nome próprio:

```go
func TestSomar(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     int
		esperado int
	}{
		{"positivos", 2, 3, 5},
		{"com zero", 0, 7, 7},
		{"negativos", -2, -3, -5},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := Somar(c.a, c.b); got != c.esperado {
				t.Errorf("Somar(%d, %d) = %d; esperado %d", c.a, c.b, got, c.esperado)
			}
		})
	}
}
```

Saída com `-v`:

```
=== RUN   TestSomar
=== RUN   TestSomar/positivos
=== RUN   TestSomar/com_zero
=== RUN   TestSomar/negativos
--- PASS: TestSomar (0.00s)
    --- PASS: TestSomar/positivos (0.00s)
    --- PASS: TestSomar/com_zero (0.00s)
    --- PASS: TestSomar/negativos (0.00s)
```

Por que usar subtestes?
- Um `t.Fatal` dentro do subteste **para só aquele caso**, os outros continuam
- Dá pra rodar **um caso específico**:

```bash
go test -run 'TestSomar/negativos' -v
```

(Espaços no nome viram `_`.)

---

## 🧪 Table test com erros

```go
func TestDividir(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     float64
		esperado float64
		errEsp   error
	}{
		{"simples", 10, 2, 5, nil},
		{"decimal", 7, 2, 3.5, nil},
		{"por zero", 1, 0, 0, ErrDivisaoPorZero},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := Dividir(c.a, c.b)

			if !errors.Is(err, c.errEsp) {
				t.Fatalf("erro = %v; esperado %v", err, c.errEsp)
			}
			if got != c.esperado {
				t.Errorf("Dividir(%v, %v) = %v; esperado %v", c.a, c.b, got, c.esperado)
			}
		})
	}
}
```

> `errors.Is(nil, nil)` é `true`, então o mesmo código serve para casos com e sem erro.

---

## ⚖️ Comparando slices, maps e structs

`==` não funciona com slices e maps. Opções:

```go
// Slices (Go 1.21+)
if !slices.Equal(got, esperado) {
	t.Errorf("got %v; esperado %v", got, esperado)
}

// Maps (Go 1.21+)
if !maps.Equal(got, esperado) { ... }

// Qualquer coisa (inclusive structs com slices dentro)
if !reflect.DeepEqual(got, esperado) { ... }
```

> 📦 Muitos projetos usam `github.com/google/go-cmp/cmp`, que mostra **a diferença** exata: `if diff := cmp.Diff(esperado, got); diff != "" { t.Errorf("(-esperado +obtido):\n%s", diff) }`

---

## 🤝 `t.Helper()`: funções auxiliares

Quando você cria funções de verificação reutilizáveis, marque-as com `t.Helper()`. Assim, se falhar, o erro aponta a **linha do teste**, não a linha dentro do helper:

```go
func verificarIgual[T comparable](t *testing.T, got, esperado T) {
	t.Helper() // ✅ sem isso, o erro apontaria para esta linha
	if got != esperado {
		t.Errorf("obtido %v; esperado %v", got, esperado)
	}
}

func TestSomar(t *testing.T) {
	verificarIgual(t, Somar(2, 2), 4)
	verificarIgual(t, Somar(1, 1), 3) // ← o erro vai apontar AQUI
}
```

---

## ⚡ Testes em paralelo: `t.Parallel()`

```go
func TestLento(t *testing.T) {
	casos := []struct {
		nome string
		dur  time.Duration
	}{
		{"a", time.Second},
		{"b", time.Second},
		{"c", time.Second},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			t.Parallel() // roda junto com os outros subtestes paralelos
			time.Sleep(c.dur)
		})
	}
}
// ~1s no total, em vez de 3s
```

> ⚠️ Só use `t.Parallel()` se os testes **não compartilham estado** (variáveis globais, mesmo arquivo, mesmo banco...).

---

## 🧹 Preparação e limpeza

### `t.Cleanup`: executa ao fim do teste

```go
func TestComArquivo(t *testing.T) {
	f, err := os.CreateTemp("", "teste-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(f.Name()) // roda quando o teste (e seus subtestes) terminar
	})
	// ... usa o arquivo ...
}
```

### `t.TempDir()`: pasta temporária automática

```go
func TestSalvar(t *testing.T) {
	dir := t.TempDir() // criada agora e APAGADA automaticamente no fim
	caminho := filepath.Join(dir, "dados.json")

	if err := Salvar(caminho, dados); err != nil {
		t.Fatal(err)
	}
	// ... verifica o arquivo ...
}
```

### `t.Context()` (Go 1.24+)

Um `context` que é cancelado automaticamente quando o teste termina:

```go
func TestBuscar(t *testing.T) {
	resultado, err := Buscar(t.Context(), 42)
	// ...
}
```

### `TestMain`: preparação do pacote inteiro

Roda **uma vez** antes de todos os testes do pacote:

```go
func TestMain(m *testing.M) {
	fmt.Println("preparando o banco de testes...")
	codigo := m.Run() // roda todos os testes
	fmt.Println("limpando...")
	os.Exit(codigo)
}
```

> Use com moderação: na maioria dos casos, `t.Cleanup` e helpers resolvem.

---

## 📂 A pasta `testdata`

Arquivos usados pelos testes (JSONs de exemplo, imagens, entradas) ficam em **`testdata/`**. O Go **ignora** essa pasta na compilação, e o `go test` roda **dentro da pasta do pacote**, então os caminhos relativos funcionam:

```go
dados, err := os.ReadFile("testdata/produtos.json")
```

---

## ✍️ Exercícios

1. Reescreva os testes da calculadora como **table tests** com **subtestes**.
2. Escreva um table test para `EhPalindromo` com pelo menos 6 casos (inclua acentos e maiúsculas).
3. Escreva `Filtrar[T any]` (módulo 13) e teste com `slices.Equal`.
4. Crie um helper `verificarErro(t, err, esperado)` com `t.Helper()`. Faça falhar e veja qual linha é apontada, com e sem o `t.Helper()`.
5. Crie 5 subtestes que dormem 500 ms cada, com e sem `t.Parallel()`. Compare o tempo.
6. Escreva uma função que salva produtos em JSON num arquivo, e teste usando `t.TempDir()`.
7. Rode apenas um subteste específico com `-run 'TestX/nome'`.

---

⬅️ Anterior: [Primeiros testes](01-primeiros-testes.md) · ➡️ Próximo: [Cobertura, benchmarks e exemplos](03-cobertura-benchmarks-e-exemplos.md)

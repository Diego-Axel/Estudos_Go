# 03 — Wrapping, `errors.Is`, `errors.As` e `errors.Join`

## 🎁 O problema: erros sem contexto

```go
func carregarConfig() error {
	_, err := os.ReadFile("config.json")
	return err
}

func iniciar() error {
	return carregarConfig()
}

func main() {
	if err := iniciar(); err != nil {
		fmt.Println(err)
	}
}
// open config.json: no such file or directory
```

Numa aplicação grande, **quem** tentou abrir esse arquivo? **Por quê?** A mensagem não conta.

---

## 🧷 Embrulhando (*wrapping*) com `%w` ⭐

`fmt.Errorf` com o verbo **`%w`** cria um erro novo que **contém** o original, adicionando contexto:

```go
func carregarConfig() error {
	_, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("carregar configuração: %w", err)
	}
	return nil
}

func iniciar() error {
	if err := carregarConfig(); err != nil {
		return fmt.Errorf("iniciar servidor: %w", err)
	}
	return nil
}
// iniciar servidor: carregar configuração: open config.json: no such file or directory
```

Agora a mensagem conta a **história completa**, e o erro original continua **dentro**:

```
┌─ iniciar servidor ──────────────────────────────┐
│ ┌─ carregar configuração ──────────────────────┐ │
│ │ ┌─ *fs.PathError ──────────────────────────┐ │ │
│ │ │ open config.json: no such file ...       │ │ │
│ │ └──────────────────────────────────────────┘ │ │
│ └──────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────┘
```

### `%w` vs. `%v`

| Verbo | Mensagem | Erro original acessível? |
|---|---|---|
| `%w` | igual | ✅ sim (`errors.Is`/`As` enxergam) |
| `%v` | igual | ❌ não (vira só texto) |

Use `%v` quando **não quiser** expor o erro interno (ex: esconder detalhes do banco de dados na API pública).

---

## 🔍 `errors.Is`: "tem esse erro em algum lugar da corrente?"

Compara com **cada erro da corrente**, de fora para dentro:

```go
err := iniciar()

fmt.Println(err == fs.ErrNotExist)          // false ❌ (o de fora é outro erro)
fmt.Println(errors.Is(err, fs.ErrNotExist)) // true  ✅ (achou lá dentro)
```

> ✅ Por isso: **sempre** use `errors.Is` em vez de `==` para comparar erros.

---

## 🎯 `errors.As`: "tem um erro desse TIPO na corrente? Me dá ele!"

```go
err := iniciar()

var pe *fs.PathError
if errors.As(err, &pe) {
	fmt.Println("arquivo problemático:", pe.Path) // config.json
}
```

- O segundo argumento é um **ponteiro** para uma variável do tipo procurado
- Se achar, preenche a variável e retorna `true`

> ⚠️ Passar a variável **sem `&`** causa panic: `errors.As(err, pe)` ❌

---

## 📦 `errors.Unwrap`: tirando uma camada

```go
err := fmt.Errorf("camada 2: %w", fmt.Errorf("camada 1: %w", io.EOF))

fmt.Println(err)                               // camada 2: camada 1: EOF
fmt.Println(errors.Unwrap(err))                // camada 1: EOF
fmt.Println(errors.Unwrap(errors.Unwrap(err))) // EOF
```

> Na prática, raramente se usa `Unwrap` direto: `Is` e `As` já percorrem a corrente toda.

---

## ➕ Juntando vários erros: `errors.Join` (Go 1.20+)

Útil para **validações** que reportam todos os problemas de uma vez:

```go
var ErrIdadeNegativa = errors.New("idade não pode ser negativa")

func validar(nome, email string, idade int) error {
	var errs []error

	if nome == "" {
		errs = append(errs, errors.New("nome é obrigatório"))
	}
	if !strings.Contains(email, "@") {
		errs = append(errs, fmt.Errorf("email inválido: %q", email))
	}
	if idade < 0 {
		errs = append(errs, ErrIdadeNegativa)
	}

	return errors.Join(errs...) // retorna nil se a lista estiver vazia!
}

func main() {
	err := validar("", "ana.com", -1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(errors.Is(err, ErrIdadeNegativa)) // true
}
// nome é obrigatório
// email inválido: "ana.com"
// idade não pode ser negativa
// true
```

- As mensagens são unidas com **quebra de linha**
- `errors.Is` e `errors.As` procuram em **todos** os erros juntados
- Se todos forem `nil` (ou a lista estiver vazia), retorna `nil` ✅

Também dá pra usar **vários `%w`** num mesmo `fmt.Errorf`:

```go
err := fmt.Errorf("falha dupla: %w e %w", ErrA, ErrB)
errors.Is(err, ErrA) // true
errors.Is(err, ErrB) // true
```

---

## 🔧 Personalizando `Is` e `Unwrap` no seu tipo

Seu tipo de erro pode **embrulhar** outro erro:

```go
type ErroConsulta struct {
	Query string
	Err   error // causa
}

func (e *ErroConsulta) Error() string {
	return fmt.Sprintf("consulta %q: %v", e.Query, e.Err)
}

func (e *ErroConsulta) Unwrap() error { // permite Is/As enxergarem a causa
	return e.Err
}

err := &ErroConsulta{Query: "SELECT ...", Err: sql.ErrNoRows}
fmt.Println(errors.Is(err, sql.ErrNoRows)) // true
```

---

## 🧭 Onde adicionar contexto?

✅ Adicione contexto dizendo **o que a função atual estava tentando fazer**:

```go
return fmt.Errorf("buscar pedido %d: %w", id, err)
```

❌ Evite repetir palavras como "erro" ou "falha ao" em cada camada, senão vira:

```
erro ao iniciar: erro ao carregar: erro ao abrir: ...
```

❌ Não **registre (log) e retorne** o mesmo erro, senão ele aparece várias vezes no log. **Ou** trata (loga), **ou** propaga:

```go
// ❌
if err != nil {
	log.Println("erro ao salvar:", err)
	return err
}

// ✅
if err != nil {
	return fmt.Errorf("salvar pedido: %w", err)
}
```

---

## ✍️ Exercícios

1. Crie três funções encadeadas (`main → servico → repositorio`) em que o repositório retorna `ErrNaoEncontrado`. Embrulhe o erro em cada camada com `%w` e imprima a mensagem final.
2. No `main` do exercício 1, verifique com `errors.Is` se o erro é `ErrNaoEncontrado`. Depois troque `%w` por `%v` e veja o resultado mudar.
3. Use `errors.As` para extrair um `*fs.PathError` de um erro embrulhado duas vezes.
4. Escreva `validarUsuario(u Usuario) error` que use `errors.Join` para retornar **todos** os problemas encontrados.
5. Crie o tipo `ErroRequisicao{URL string; Err error}` com `Unwrap()`, e mostre que `errors.Is` enxerga o erro interno.
6. Explique por que "logar e retornar" o mesmo erro é considerado má prática.

---

⬅️ Anterior: [Erros sentinela e personalizados](02-erros-sentinela-e-personalizados.md) · ➡️ Próximo: [panic, recover e boas práticas](04-panic-recover-e-boas-praticas.md)

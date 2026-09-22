# 01 — Funções: o Básico

## 🧱 Anatomia de uma função

```go
func nome(parametro tipo) tipoDeRetorno {
	// corpo
	return valor
}
```

```go
func soma(a int, b int) int {
	return a + b
}

func main() {
	resultado := soma(3, 4)
	fmt.Println(resultado) // 7
}
```

### Parâmetros do mesmo tipo podem ser agrupados

```go
func soma(a, b int) int {            // a e b são int
	return a + b
}

func cadastrar(nome, email string, idade int) {
	// nome e email são string; idade é int
}
```

### Sem parâmetros e/ou sem retorno

```go
func saudacao() {
	fmt.Println("Olá!")
}
```

---

## 🎯 Múltiplos retornos ⭐

Uma das marcas registradas de Go: funções podem retornar **vários valores**.

```go
func dividir(a, b int) (int, int) {
	return a / b, a % b
}

func main() {
	quociente, resto := dividir(17, 5)
	fmt.Println(quociente, resto) // 3 2
}
```

Ignorando um dos retornos com `_`:

```go
q, _ := dividir(17, 5)
```

### O padrão `(resultado, error)`

É assim que Go trata erros: a função retorna o resultado **e** um `error` (que é `nil` quando deu tudo certo).

```go
import "errors"

func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divisão por zero")
	}
	return a / b, nil
}

func main() {
	r, err := dividir(10, 0)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	fmt.Println("Resultado:", r)
}
// Erro: divisão por zero
```

> 📌 Por convenção, o `error` é sempre o **último** retorno. O módulo 11 aprofunda tratamento de erros.

---

## 🏷️ Retornos nomeados

Você pode dar **nome** aos valores de retorno. Eles viram variáveis já declaradas (com valor zero) dentro da função:

```go
func dividir(a, b int) (quociente, resto int) {
	quociente = a / b
	resto = a % b
	return // "naked return": retorna quociente e resto
}
```

Vantagens:
- **Documenta** o que cada retorno significa (aparece no `go doc` e no autocompletar)
- Útil junto com `defer` (veremos em [defer](05-defer.md))

> ⚠️ Evite o `return` "pelado" (*naked return*) em funções longas: fica difícil saber o que está sendo retornado. Você pode usar nomes e **ainda assim** retornar explicitamente: `return quociente, resto`.

---

## 📋 Passagem por valor

Em Go, **todo argumento é passado por valor**: a função recebe uma **cópia**.

```go
func dobrar(n int) {
	n *= 2 // altera só a cópia
}

func main() {
	x := 10
	dobrar(x)
	fmt.Println(x) // 10 (não mudou!)
}
```

Para a função alterar a variável original, passe um **ponteiro**:

```go
func dobrar(n *int) {
	*n *= 2
}

func main() {
	x := 10
	dobrar(&x)
	fmt.Println(x) // 20
}
```

> 🔍 **Slices e maps** parecem ser passados "por referência" porque a cópia continua apontando para os mesmos dados por baixo. Alterar `s[0]` dentro da função altera o original. Isso será detalhado nos módulos **Arrays, Slices e Maps** e **Ponteiros**.

---

## 🚫 O que Go **não** tem em funções

| Recurso | Em Go |
|---|---|
| Parâmetros com valor padrão (`func f(x = 10)`) | ❌ não existe |
| Sobrecarga (duas funções com o mesmo nome) | ❌ não existe |
| Argumentos nomeados na chamada (`f(x: 1)`) | ❌ não existe |

Alternativas comuns:

```go
// Nomes diferentes em vez de sobrecarga
func somaInt(a, b int) int         { return a + b }
func somaFloat(a, b float64) float64 { return a + b }

// Struct de opções em vez de parâmetros padrão
type Config struct {
	Porta   int
	Debug   bool
}

func iniciar(cfg Config) {
	if cfg.Porta == 0 {
		cfg.Porta = 8080 // "valor padrão"
	}
	fmt.Println("porta", cfg.Porta, "debug", cfg.Debug)
}

iniciar(Config{Debug: true}) // porta 8080 debug true
```

> Para funções que funcionam com vários tipos, Go tem **Generics** (módulo 13).

---

## 🔓 Visibilidade: maiúscula ou minúscula

```go
func Calcular() {} // Exportada: outros pacotes podem chamar
func calcular() {} // Privada: só dentro deste pacote
```

---

## 📝 Comentário de documentação

Por convenção, toda função exportada tem um comentário **começando com o nome dela**:

```go
// Soma retorna a soma de a e b.
func Soma(a, b int) int {
	return a + b
}
```

Isso aparece no `go doc`, no VS Code (ao passar o mouse) e no site **pkg.go.dev**.

---

## ✍️ Exercícios

1. Crie uma função `maior(a, b int) int` que retorne o maior dos dois.
2. Crie `minMax(a, b, c int) (int, int)` que retorne o menor e o maior de três números.
3. Crie `raizQuadrada(n float64) (float64, error)` que retorne erro se `n` for negativo (use `math.Sqrt`).
4. Reescreva o exercício 2 usando **retornos nomeados** (`min, max int`).
5. Crie `trocar(a, b *int)` que troque os valores de duas variáveis do `main`.
6. Crie `ehPrimo(n int) bool` e use-a para imprimir todos os primos de 1 a 100.
7. Refatore seu programa de IMC em funções: `calcularIMC(peso, altura float64) float64` e `classificarIMC(imc float64) string`.

---

🏠 [Módulo 05](README.md) · ➡️ Próximo: [Funções variádicas](02-variadicas.md)

# 02 — Funções Variádicas

## 🔢 O que são?

Funções que aceitam **qualquer quantidade** de argumentos de um tipo. Usam `...` antes do tipo:

```go
func somar(numeros ...int) int {
	total := 0
	for _, n := range numeros {
		total += n
	}
	return total
}

func main() {
	fmt.Println(somar())              // 0
	fmt.Println(somar(1))             // 1
	fmt.Println(somar(1, 2, 3))       // 6
	fmt.Println(somar(1, 2, 3, 4, 5)) // 15
}
```

Dentro da função, `numeros` é um **slice** (`[]int`):

```go
func mostrar(numeros ...int) {
	fmt.Printf("%T %v %d\n", numeros, numeros, len(numeros))
}

mostrar(1, 2, 3) // []int [1 2 3] 3
mostrar()        // []int [] 0
```

---

## 📌 Regras

1. O parâmetro variádico tem que ser o **último**
2. Só pode existir **um** por função
3. Pode ter parâmetros "normais" **antes** dele

```go
func saudar(saudacao string, nomes ...string) {
	for _, nome := range nomes {
		fmt.Println(saudacao, nome)
	}
}

saudar("Olá", "Ana", "Bia", "Caio")
// Olá Ana
// Olá Bia
// Olá Caio

// func errada(nums ...int, nome string) {} // ❌ erro: só pode ser o último
```

---

## 🧨 "Espalhando" um slice com `...`

Se você já tem um slice, passe ele com `...` **depois**:

```go
valores := []int{10, 20, 30}

// fmt.Println(somar(valores))  // ❌ erro: cannot use valores (variable of type []int) as int value
fmt.Println(somar(valores...))  // ✅ 60
```

> ⚠️ Não dá pra misturar: `somar(1, 2, valores...)` **não compila**. Ou passa valores soltos, ou passa o slice espalhado.

### ⚠️ Cuidado: espalhar **não copia** o slice

```go
func zerarPrimeiro(nums ...int) {
	if len(nums) > 0 {
		nums[0] = 0
	}
}

s := []int{1, 2, 3}
zerarPrimeiro(s...)
fmt.Println(s) // [0 2 3] ← o original foi alterado!

zerarPrimeiro(1, 2, 3) // aqui Go cria um slice novo, sem efeito colateral
```

---

## 🧰 Você já usa variádicas!

Várias funções da biblioteca padrão são variádicas:

```go
fmt.Println(a ...any)                  // por isso aceita qualquer quantidade de coisas
fmt.Printf(format string, a ...any)
append(slice []T, elems ...T) []T
```

Juntar dois slices com `append` + `...`:

```go
a := []int{1, 2}
b := []int{3, 4, 5}
a = append(a, b...)
fmt.Println(a) // [1 2 3 4 5]
```

> 💡 `any` é um apelido para `interface{}` e significa "qualquer tipo". Veremos em **Interfaces**.

---

## 🛠️ Exemplo prático

```go
func media(notas ...float64) (float64, error) {
	if len(notas) == 0 {
		return 0, errors.New("nenhuma nota informada")
	}
	soma := 0.0
	for _, n := range notas {
		soma += n
	}
	return soma / float64(len(notas)), nil
}

func main() {
	m, err := media(7, 8.5, 9, 6.5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Média: %.2f\n", m) // Média: 7.75
}
```

---

## ✍️ Exercícios

1. Crie `maior(nums ...int) int` que retorne o maior valor. O que fazer se não receber nenhum número?
2. Crie `concatenar(sep string, partes ...string) string` que junte as partes com o separador. Ex: `concatenar("-", "a", "b", "c")` → `"a-b-c"`. (Depois compare com `strings.Join`.)
3. Crie `contem(alvo int, nums ...int) bool`.
4. Dado `notas := []float64{8, 9, 7}`, chame a função `media` passando o slice.
5. Crie `logar(nivel string, msgs ...any)` que imprima `[NIVEL]` seguido das mensagens (dica: repasse com `fmt.Println(msgs...)`).

---

⬅️ Anterior: [O básico](01-basico.md) · ➡️ Próximo: [Funções como valor e closures](03-funcoes-como-valor-e-closures.md)

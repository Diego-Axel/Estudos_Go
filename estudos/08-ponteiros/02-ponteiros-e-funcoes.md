# 02 — Ponteiros e Funções

## 📋 Relembrando: Go sempre passa **por valor**

Quando você chama uma função, ela recebe uma **cópia** de cada argumento:

```go
func zerar(n int) {
	n = 0 // muda só a cópia
}

func main() {
	x := 42
	zerar(x)
	fmt.Println(x) // 42 😕
}
```

```
main:   x = 42        zerar:  n = 42 → 0
        (intacto)             (cópia, some quando a função termina)
```

---

## 🎯 Passando um ponteiro

Se a função recebe o **endereço**, ela consegue alterar o original:

```go
func zerar(n *int) {
	*n = 0 // muda o valor NO ENDEREÇO recebido
}

func main() {
	x := 42
	zerar(&x)
	fmt.Println(x) // 0 ✅
}
```

> 🧠 Tecnicamente, o **ponteiro** também é passado por valor (a função recebe uma cópia do endereço), mas a cópia aponta para o **mesmo lugar**.

### Exemplo clássico: trocar valores

```go
func trocar(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	x, y := 1, 2
	trocar(&x, &y)
	fmt.Println(x, y) // 2 1
}
```

### Exemplo: incrementar um contador

```go
func incrementar(contador *int) {
	*contador++
}

func main() {
	visitas := 0
	for range 3 {
		incrementar(&visitas)
	}
	fmt.Println(visitas) // 3
}
```

> 💡 `*contador++` funciona porque `++` só se aplica depois de desreferenciar: é o mesmo que `(*contador)++`.

---

## ↩️ Retornando um ponteiro

Em C, retornar o endereço de uma variável local é um bug grave. **Em Go é perfeitamente seguro**:

```go
func novoContador() *int {
	c := 0
	return &c // ✅ Go percebe e mantém c viva enquanto houver referência
}

func main() {
	p := novoContador()
	*p++
	fmt.Println(*p) // 1
}
```

O compilador faz a **análise de escape** (*escape analysis*): se a variável "escapa" da função, ela é alocada no **heap** em vez da **pilha**. Mais detalhes em [Stack, heap e boas práticas](04-stack-heap-e-boas-praticas.md).

---

## 🧷 Slices e maps: já se comportam "como referência"

Lembra do módulo 06? Slices e maps **carregam um ponteiro interno**. Então, para **alterar elementos**, não precisa de ponteiro:

```go
func dobrar(nums []int) {
	for i := range nums {
		nums[i] *= 2 // ✅ altera o original
	}
}

func adicionarNota(notas map[string]float64) {
	notas["Bia"] = 9 // ✅ altera o original
}
```

Mas, para **mudar o slice em si** (tamanho, trocar por outro), a função precisa **retornar** o novo slice, ou receber um **ponteiro para o slice**:

```go
// ✅ Opção 1 (mais idiomática): retornar
func adicionar(s []int, v int) []int {
	return append(s, v)
}

// ✅ Opção 2: ponteiro para slice
func adicionarPtr(s *[]int, v int) {
	*s = append(*s, v)
}

func main() {
	nums := []int{1, 2}
	nums = adicionar(nums, 3)
	adicionarPtr(&nums, 4)
	fmt.Println(nums) // [1 2 3 4]
}
```

### Resumo: quem precisa de ponteiro para ser alterado?

| Tipo | Alterar o conteúdo dentro da função | Precisa de ponteiro? |
|---|---|---|
| `int`, `float64`, `bool`, `string` | mudar o valor | ✅ sim |
| arrays | mudar elementos | ✅ sim |
| structs | mudar campos | ✅ sim |
| slices | mudar **elementos** | ❌ não |
| slices | `append` / mudar tamanho | retorne o slice (ou use `*[]T`) |
| maps | adicionar/remover/alterar | ❌ não |
| channels | enviar/receber | ❌ não |

---

## 🤔 Quando usar ponteiro em parâmetros?

✅ **Use ponteiro quando:**
1. A função **precisa alterar** o valor original
2. O valor é **grande** (uma struct com muitos campos) e copiá-lo seria caro
3. Você quer representar **"ausência de valor"** com `nil`

❌ **Não use ponteiro quando:**
1. O valor é **pequeno** e a função só **lê** (`int`, `string`, structs pequenas)
2. O tipo **já** é "referência" (slices, maps, channels, funções)
3. Só por "performance" sem medir: copiar valores pequenos costuma ser **mais rápido** que lidar com ponteiros

> 🧠 **Regra de bolso:** comece passando **valores**. Troque para ponteiros quando houver um motivo claro.

---

## ✍️ Exercícios

1. Escreva `dobrar(n *int)` e teste com uma variável do `main`.
2. Escreva `dividir(a, b int, quociente, resto *int)` que preencha o quociente e o resto via ponteiros. (Depois compare com a versão de **múltiplos retornos**. Qual é mais idiomática em Go?)
3. Escreva `normalizar(s *string)` que aplique `strings.TrimSpace` e `strings.ToLower` na string original.
4. Escreva `resetar(nums []int)` que zere todos os elementos. Precisa de ponteiro? Por quê?
5. Escreva `limpar(s *[]int)` que deixe o slice original vazio (`len` 0).
6. Explique por que `func novo() *int { x := 5; return &x }` é seguro em Go.

---

⬅️ Anterior: [O que são ponteiros?](01-o-que-sao-ponteiros.md) · ➡️ Próximo: [Ponteiros com structs e coleções](03-ponteiros-com-structs-e-colecoes.md)

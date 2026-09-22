# 03 — Slices: Avançado

## ⚠️ Slices compartilham memória

Fatiar **não copia** os dados. O novo slice aponta para o **mesmo array** de baixo:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3] // [2 3]

b[0] = 99
fmt.Println(a) // [1 99 3 4 5] ← a também mudou!
fmt.Println(b) // [99 3]
```

```
a:  [ 1 | 99 | 3 | 4 | 5 ]
          ▲    ▲
b:        └────┘ (mesmo array!)
```

### A pegadinha do `append` que sobrescreve

`b` tem `len=2`, mas `cap=4` (vai até o fim de `a`). Então `append(b, ...)` **cabe no array** e escreve **por cima** de `a`:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]              // [2 3], len=2, cap=4

b = append(b, 100)
fmt.Println(b)           // [2 3 100]
fmt.Println(a)           // [1 2 3 100 5] 😱 o 4 virou 100!
```

### Soluções

**1. Limitar a capacidade** com `s[a:b:c]`: o `append` é obrigado a criar um array novo:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3:3]            // len=2, cap=2

b = append(b, 100)       // não cabe → array novo
fmt.Println(a)           // [1 2 3 4 5] ✅ intacto
fmt.Println(b)           // [2 3 100]
```

**2. Fazer uma cópia de verdade** (próxima seção).

---

## 📑 `copy`: copiando slices

```go
origem := []int{1, 2, 3}
destino := make([]int, len(origem)) // precisa ter espaço!

n := copy(destino, origem)          // copy(DESTINO, ORIGEM)
fmt.Println(n, destino)             // 3 [1 2 3]

destino[0] = 99
fmt.Println(origem)                 // [1 2 3] ✅ independente
```

> ⚠️ `copy` copia **o mínimo** entre `len(destino)` e `len(origem)`. Se o destino tiver `len` 0, **nada** é copiado!

```go
var errado []int
copy(errado, origem)
fmt.Println(errado) // [] 😬
```

Formas curtas de clonar:

```go
clone := append([]int(nil), origem...)
clone2 := slices.Clone(origem) // Go 1.21+
```

---

## ✂️ Operações comuns "na mão"

```go
s := []int{10, 20, 30, 40, 50}

// Remover o elemento do índice i (mantendo a ordem)
i := 2
s = append(s[:i], s[i+1:]...)
fmt.Println(s) // [10 20 40 50]

// Inserir v no índice i
i, v := 1, 15
s = append(s[:i], append([]int{v}, s[i:]...)...)
fmt.Println(s) // [10 15 20 40 50]

// Remover o último (pop)
ultimo := s[len(s)-1]
s = s[:len(s)-1]
fmt.Println(ultimo, s) // 50 [10 15 20 40]

// Remover o primeiro (shift)
primeiro := s[0]
s = s[1:]
fmt.Println(primeiro, s) // 10 [15 20 40]
```

> 😅 Confuso, né? Por isso o Go 1.21 trouxe o pacote **`slices`**.

---

## 📦 Pacote `slices` (Go 1.21+) ⭐

```go
import "slices"

nums := []int{5, 2, 8, 1, 9, 3}

slices.Sort(nums)                     // ordena (altera o próprio slice)
fmt.Println(nums)                     // [1 2 3 5 8 9]

fmt.Println(slices.Contains(nums, 8)) // true
fmt.Println(slices.Index(nums, 5))    // 3  (-1 se não achar)
fmt.Println(slices.Max(nums))         // 9
fmt.Println(slices.Min(nums))         // 1

idx, achou := slices.BinarySearch(nums, 8) // busca em slice ORDENADO
fmt.Println(idx, achou)               // 4 true

slices.Reverse(nums)
fmt.Println(nums)                     // [9 8 5 3 2 1]

nums = slices.Delete(nums, 1, 3)      // remove índices [1, 3)
fmt.Println(nums)                     // [9 3 2 1]

nums = slices.Insert(nums, 1, 7, 7)   // insere no índice 1
fmt.Println(nums)                     // [9 7 7 3 2 1]

nums = slices.Compact(nums)           // remove duplicados CONSECUTIVOS
fmt.Println(nums)                     // [9 7 3 2 1]

fmt.Println(slices.Equal([]int{1, 2}, []int{1, 2})) // true
```

Ordenando com critério próprio:

```go
nomes := []string{"Caio", "Ana", "Bernardo"}
slices.SortFunc(nomes, func(a, b string) int {
	return len(a) - len(b) // negativo: a vem antes; positivo: b vem antes
})
fmt.Println(nomes) // [Ana Caio Bernardo]
```

> Antes do Go 1.21 usava-se o pacote `sort` (`sort.Ints`, `sort.Strings`, `sort.Slice`). Você ainda vai ver muito em código existente.

---

## 🧹 `clear` (Go 1.21+)

Zera todos os elementos, mantendo `len`:

```go
s := []int{1, 2, 3}
clear(s)
fmt.Println(s, len(s)) // [0 0 0] 3
```

---

## 🧷 Slices em funções

A função recebe uma **cópia do cabeçalho** (ponteiro, len, cap), mas o **ponteiro aponta pro mesmo array**:

```go
func dobrar(s []int) {
	for i := range s {
		s[i] *= 2 // ✅ altera o original
	}
}

func adicionar(s []int) {
	s = append(s, 99) // ❌ o len do chamador NÃO muda
}

func main() {
	nums := []int{1, 2, 3}
	dobrar(nums)
	fmt.Println(nums) // [2 4 6]

	adicionar(nums)
	fmt.Println(nums) // [2 4 6] ← o 99 não aparece
}
```

✅ Se a função precisa **adicionar** elementos, **retorne** o novo slice (igual o `append` faz):

```go
func adicionar(s []int, v int) []int {
	return append(s, v)
}

nums = adicionar(nums, 99)
```

---

## ✍️ Exercícios

1. Explique a saída e depois corrija para que `a` não seja alterado:
   ```go
   a := []int{1, 2, 3, 4}
   b := a[:2]
   b = append(b, 50)
   fmt.Println(a)
   ```
2. Escreva `remover(s []int, i int) []int` "na mão" e depois com `slices.Delete`.
3. Dado `[]int{3, 1, 3, 2, 1, 2, 3}`, gere um slice **sem duplicados** (dica: `slices.Sort` + `slices.Compact`).
4. Escreva `inverter(s []int)` que inverta o slice **no lugar** (sem criar outro), trocando elementos das pontas. Compare com `slices.Reverse`.
5. Ordene um slice de strings pelo **tamanho decrescente** com `slices.SortFunc`.
6. Implemente uma **pilha** (push/pop) e uma **fila** (enqueue/dequeue) usando slices.
7. Por que esta função não funciona? Corrija.
   ```go
   func limpar(s []int) { s = nil }
   ```

---

⬅️ Anterior: [Slices](02-slices.md) · ➡️ Próximo: [Maps](04-maps.md)

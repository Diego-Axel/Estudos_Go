# 01 — Arrays

Um **array** é uma sequência de elementos do **mesmo tipo** com **tamanho fixo**, definido na declaração e que **nunca muda**.

> 💡 Na prática, em Go você vai usar **slices** 95% das vezes. Mas entender arrays é essencial, porque **todo slice é construído em cima de um array**.

---

## 📝 Declarando

```go
var notas [5]float64          // 5 posições, todas com valor zero
fmt.Println(notas)            // [0 0 0 0 0]

cores := [3]string{"vermelho", "verde", "azul"}
fmt.Println(cores)            // [vermelho verde azul]

// Deixa o compilador contar o tamanho
dias := [...]string{"seg", "ter", "qua", "qui", "sex"}
fmt.Println(len(dias))        // 5

// Inicializando posições específicas
parcial := [5]int{1: 10, 3: 30}
fmt.Println(parcial)          // [0 10 0 30 0]
```

---

## 🎯 Acessando e alterando

Índices começam em **0** e vão até **`len - 1`**:

```go
cores := [3]string{"vermelho", "verde", "azul"}

fmt.Println(cores[0])         // vermelho
fmt.Println(cores[len(cores)-1]) // azul (último)

cores[1] = "amarelo"
fmt.Println(cores)            // [vermelho amarelo azul]
```

### Fora dos limites

```go
// cores[3] = "rosa" // ❌ erro de compilação: invalid argument: index 3 out of bounds [0:3]

i := 3
cores[i] = "rosa"    // 💥 panic: runtime error: index out of range [3] with length 3
```

> Go **sempre verifica os limites**. Não existe acessar memória "fora do array" como em C.

---

## 🔁 Percorrendo

```go
notas := [4]float64{7.5, 8, 9.5, 6}

for i := 0; i < len(notas); i++ {
	fmt.Println(i, notas[i])
}

for i, nota := range notas {
	fmt.Println(i, nota)
}
```

---

## 🧬 O tamanho faz parte do tipo

`[3]int` e `[4]int` são **tipos diferentes**:

```go
var a [3]int
var b [4]int

// a = b // ❌ erro: cannot use b (variable of type [4]int) as [3]int value
fmt.Printf("%T %T\n", a, b) // [3]int [4]int
```

Por isso, uma função que recebe `[3]int` **não aceita** `[4]int`. Esse é um dos motivos de slices serem mais usados.

---

## 📋 Arrays são **valores** (são copiados!)

Atribuir ou passar um array para uma função **copia todos os elementos**:

```go
original := [3]int{1, 2, 3}
copia := original
copia[0] = 99

fmt.Println(original) // [1 2 3]  ← não mudou
fmt.Println(copia)    // [99 2 3]
```

```go
func zerar(arr [3]int) {
	arr[0] = 0 // altera só a cópia
}

func zerarDeVerdade(arr *[3]int) {
	arr[0] = 0 // altera o original (através do ponteiro)
}
```

> ⚠️ Copiar um array de 1 milhão de elementos a cada chamada é caro. Mais um motivo para usar slices.

---

## ⚖️ Arrays podem ser comparados

Se os elementos forem comparáveis, dá pra usar `==`:

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
fmt.Println(a == b) // true
```

Isso permite usar arrays como **chave de map** (slices não podem):

```go
visitados := map[[2]int]bool{}
visitados[[2]int{0, 0}] = true // coordenada (0, 0)
```

---

## 🧊 Arrays multidimensionais (matrizes)

```go
var tabuleiro [3][3]string

tabuleiro[0][0] = "X"
tabuleiro[1][1] = "O"
tabuleiro[2][2] = "X"

for _, linha := range tabuleiro {
	fmt.Println(linha)
}
// [X  ]
// [ O ]
// [  X]

matriz := [2][3]int{
	{1, 2, 3},
	{4, 5, 6},
}
fmt.Println(matriz[1][2]) // 6
```

---

## 🤔 Quando usar array?

- Quando o tamanho é **realmente fixo** e conhecido: dias da semana, coordenadas `[2]float64`, cores RGB `[3]uint8`, hashes (`[32]byte` do SHA-256)
- Quando precisa usar como **chave de map**
- Em todo o resto → **slice**

---

## ✍️ Exercícios

1. Crie um array com 5 notas, calcule a **média** e mostre quantas ficaram acima dela.
2. Crie um array `[7]string` com os dias da semana e imprima **de trás para frente**.
3. Encontre o **maior** e o **menor** valor de `[...]int{42, 7, 19, 88, 3, 56}`.
4. Mostre que arrays são copiados: crie `a`, faça `b := a`, altere `b` e imprima os dois.
5. Crie uma matriz `[3][3]int` e calcule a soma da **diagonal principal**.
6. Por que `[3]int` e `[5]int` são tipos diferentes? Que problema isso causa em funções?

---

🏠 [Módulo 06](README.md) · ➡️ Próximo: [Slices](02-slices.md)

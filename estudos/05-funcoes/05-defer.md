# 05 — `defer`

## ⏳ O que é?

`defer` **agenda** uma chamada de função para ser executada **quando a função atual terminar**, seja por `return`, por chegar ao fim ou por `panic`.

```go
func main() {
	defer fmt.Println("mundo")
	fmt.Println("olá")
}
// olá
// mundo
```

---

## 🧹 Para que serve? Limpeza de recursos

O uso principal é **garantir** que algo seja liberado/fechado, logo depois de ser aberto:

```go
func lerArquivo(caminho string) error {
	arquivo, err := os.Open(caminho)
	if err != nil {
		return err
	}
	defer arquivo.Close() // ✅ vai fechar não importa por onde a função saia

	// ... lê o arquivo, pode ter vários returns no meio ...
	return nil
}
```

Sem `defer`, você teria que lembrar de chamar `Close()` **antes de cada `return`**. Fácil de esquecer!

Outros usos comuns:

```go
mu.Lock()
defer mu.Unlock()        // libera o mutex (Concorrência)

resp, err := http.Get(url)
// trata err...
defer resp.Body.Close()  // fecha a resposta HTTP

inicio := time.Now()
defer func() {
	fmt.Println("levou", time.Since(inicio)) // mede tempo de execução
}()
```

---

## 📚 Ordem de execução: pilha (LIFO)

Vários `defer` são executados na **ordem inversa** em que foram declarados (o último a entrar é o primeiro a sair):

```go
func main() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
	fmt.Println("início")
}
// início
// 3
// 2
// 1
```

Faz sentido: se você abriu A, depois B, deve fechar B antes de A.

---

## ⚠️ Os argumentos são avaliados **na hora do `defer`**

```go
func main() {
	x := 10
	defer fmt.Println("valor no defer:", x) // x é avaliado AGORA (10)
	x = 20
	fmt.Println("valor atual:", x)
}
// valor atual: 20
// valor no defer: 10
```

Se quiser o valor **final**, use uma closure (ela lê a variável na hora que executa):

```go
func main() {
	x := 10
	defer func() {
		fmt.Println("valor no defer:", x)
	}()
	x = 20
}
// valor no defer: 20
```

---

## 🔧 `defer` pode alterar retornos nomeados

Como o `defer` roda **depois** do `return` definir o valor, mas **antes** da função devolver para quem chamou, ele pode modificar retornos **nomeados**:

```go
func dobrarNoFinal() (resultado int) {
	defer func() {
		resultado *= 2
	}()
	return 5 // resultado = 5, depois o defer faz resultado = 10
}

fmt.Println(dobrarNoFinal()) // 10
```

Uso prático: enriquecer um erro na saída:

```go
func processar(nome string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("processar %q: %w", nome, err)
		}
	}()
	// ...
	return errors.New("falhou")
}

fmt.Println(processar("dados.csv")) // processar "dados.csv": falhou
```

---

## 🔁 Cuidado: `defer` dentro de laço

O `defer` só executa quando a **função** termina, **não** no fim de cada volta do laço:

```go
// ❌ Todos os arquivos ficam abertos até o fim da função!
func processarTodos(caminhos []string) {
	for _, c := range caminhos {
		f, _ := os.Open(c)
		defer f.Close()
		// ...
	}
}
```

Solução: mover o corpo do laço para uma função própria:

```go
// ✅ Cada arquivo é fechado ao fim de processarUm
func processarTodos(caminhos []string) {
	for _, c := range caminhos {
		processarUm(c)
	}
}

func processarUm(caminho string) {
	f, err := os.Open(caminho)
	if err != nil {
		return
	}
	defer f.Close()
	// ...
}
```

---

## 💥 `defer` + `panic` + `recover` (spoiler)

`defer` **roda mesmo quando acontece um `panic`**. E, dentro de um `defer`, a função `recover()` pode **capturar** o panic e impedir que o programa quebre:

```go
func dividirSeguro(a, b int) (resultado int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recuperado: %v", r)
		}
	}()
	return a / b, nil
}

func main() {
	fmt.Println(dividirSeguro(10, 2)) // 5 <nil>
	fmt.Println(dividirSeguro(10, 0)) // 0 recuperado: runtime error: integer divide by zero
}
```

> Isso será aprofundado no módulo 11, **Tratamento de erros**. Em Go, `panic` é para situações **excepcionais**; erros esperados usam `error`.

---

## 🧾 Resumão do módulo

| Conceito | Exemplo |
|---|---|
| Função simples | `func soma(a, b int) int` |
| Múltiplos retornos | `func dividir(a, b int) (int, error)` |
| Retornos nomeados | `func f() (q, r int)` |
| Variádica | `func somar(nums ...int)` / `somar(s...)` |
| Função como valor | `f := soma` |
| Anônima | `func(x int) int { return x * 2 }` |
| Tipo de função | `type Op func(int, int) int` |
| Closure | função que captura variáveis externas |
| Recursão | função que chama a si mesma |
| `defer` | executa ao sair da função (LIFO) |

---

## ✍️ Exercícios

1. Qual a saída? Explique.
   ```go
   for i := range 3 {
       defer fmt.Print(i, " ")
   }
   ```
2. Crie uma função `medirTempo(nome string) func()` que imprima quanto tempo se passou quando a função retornada for chamada. Use assim: `defer medirTempo("processar")()`.
3. Qual a saída? Explique.
   ```go
   func f() (n int) {
       defer func() { n++ }()
       return 41
   }
   ```
4. Por que `defer f.Close()` dentro de um laço que abre 10.000 arquivos pode ser um problema?
5. Crie `acessarIndice(s []int, i int) (v int, err error)` que use `defer` + `recover` para retornar erro em vez de quebrar quando o índice for inválido.

---

⬅️ Anterior: [Recursão](04-recursao.md) · ➡️ Próximo módulo: [Arrays, Slices e Maps](../06-arrays-slices-maps/README.md)

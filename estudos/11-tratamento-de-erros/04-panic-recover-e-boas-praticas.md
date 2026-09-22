# 04 — `panic`, `recover` e Boas Práticas

## 💥 O que é `panic`?

Um `panic` **interrompe** a execução normal da função. Ele:

1. Para a função atual imediatamente
2. Executa os **`defer`** pendentes (de trás pra frente)
3. Sobe para a função que chamou e repete, até o topo da goroutine
4. Se ninguém recuperar: **o programa inteiro termina**, mostrando a mensagem e a pilha de chamadas

```go
func main() {
	defer fmt.Println("defer executado mesmo com panic")
	fmt.Println("antes")
	panic("algo muito errado aconteceu")
	fmt.Println("depois") // nunca executa
}
```

Saída:

```
antes
defer executado mesmo com panic
panic: algo muito errado aconteceu

goroutine 1 [running]:
main.main()
	/caminho/main.go:6 +0x...
exit status 2
```

---

## ⚙️ Panics gerados pelo runtime

Você já viu vários ao longo dos módulos:

| Situação | Mensagem |
|---|---|
| índice fora do limite | `index out of range [5] with length 3` |
| ponteiro nil | `invalid memory address or nil pointer dereference` |
| divisão inteira por zero | `integer divide by zero` |
| escrever em map nil | `assignment to entry in nil map` |
| type assertion errada | `interface conversion: interface {} is string, not int` |
| fechar channel já fechado | `close of closed channel` |

Todos indicam **bugs no programa**, não situações esperadas.

---

## 🤔 Quando usar `panic`?

> 🧠 **Regra de ouro:** erros **esperados** → `error`. Situações **impossíveis** (bugs) → `panic`.

✅ **Situações aceitáveis para `panic`:**

1. **Bug de programação**, estado que "nunca deveria acontecer":
   ```go
   switch status {
   case Ativo, Inativo:
   	// ...
   default:
   	panic(fmt.Sprintf("status desconhecido: %d", status))
   }
   ```

2. **Falha na inicialização** que impede o programa de funcionar:
   ```go
   var emailRegex = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
   ```
   Funções `MustXxx` dão panic se falharem. Use só com valores **fixos no código** (se a regex estiver errada, é bug e deve quebrar logo).

❌ **NÃO use `panic` para:**
- arquivo não encontrado, entrada inválida do usuário, falha de rede, registro não encontrado no banco...
- Tudo isso é **esperado** e deve ser um `error`.

### Criando seu próprio `Must`

```go
func Must[T any](v T, err error) T { // Generics (módulo 13)
	if err != nil {
		panic(err)
	}
	return v
}

porta := Must(strconv.Atoi("8080")) // ok para valores fixos/de teste
```

---

## 🛟 `recover`: capturando um panic

`recover()` **interrompe** o panic e devolve o valor passado ao `panic`. Mas ele **só funciona dentro de uma função `defer`**:

```go
func executarComSeguranca(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recuperado: %v", r)
		}
	}()
	f()
	return nil
}

func main() {
	err := executarComSeguranca(func() {
		var m map[string]int
		m["x"] = 1 // panic!
	})
	fmt.Println(err)                     // panic recuperado: assignment to entry in nil map
	fmt.Println("o programa continua 🎉")
}
```

Regras:
- Fora de um `defer`, `recover()` retorna `nil` e não faz nada
- Tem que ser chamado **diretamente** pela função adiada (não por uma função chamada por ela)
- Depois de recuperado, a função que tinha o `defer` **retorna normalmente** (com os retornos nomeados que o `defer` definiu)

### Onde `recover` é usado de verdade?

**Nas "bordas" do programa**, para que um bug em **uma** requisição não derrube **todo** o servidor:

```go
func middlewareRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic em %s: %v", r.URL.Path, rec)
				http.Error(w, "erro interno", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

> O próprio `net/http` já faz algo parecido em cada requisição.

### ⚠️ Panic em goroutine derruba **tudo**

Um `recover` só captura panics **da mesma goroutine**. Se uma goroutine sem `recover` entrar em panic, **o programa inteiro cai**, mesmo que o `main` tenha `recover`:

```go
func main() {
	defer func() { recover() }() // ❌ não ajuda!
	go func() {
		panic("boom") // derruba o programa
	}()
	time.Sleep(time.Second)
}
```

---

## 🚪 `log.Fatal` e `os.Exit`

```go
log.Fatal("não foi possível conectar:", err) // imprime e chama os.Exit(1)
os.Exit(1)                                    // encerra imediatamente
```

> ⚠️ `os.Exit` (e `log.Fatal`) **não executam os `defer`**! Arquivos podem não ser salvos, conexões não são fechadas. Use só no `main`, e com cuidado.

| | `panic` | `log.Fatal` / `os.Exit` |
|---|---|---|
| Executa `defer`? | ✅ sim | ❌ não |
| Pode ser recuperado? | ✅ com `recover` | ❌ não |
| Mostra a pilha de chamadas? | ✅ sim | ❌ não |
| Código de saída | 2 | o que você passar (Fatal = 1) |

### Padrão para o `main`

Uma forma limpa de garantir que os `defer` rodem:

```go
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := conectar()
	if err != nil {
		return fmt.Errorf("conectar ao banco: %w", err)
	}
	defer db.Close() // ✅ roda antes do os.Exit, pois run() retorna primeiro

	// ... toda a lógica do programa ...
	return nil
}
```

---

## ✅ Boas práticas: resumo

1. **Sempre verifique** o `error` retornado. Nunca ignore com `_` sem motivo claro.
2. **Adicione contexto** com `fmt.Errorf("fazendo X: %w", err)`.
3. Use **`errors.Is`** e **`errors.As`**, nunca `==` ou type assertion direta (exceto `io.EOF`).
4. **Trate o erro uma única vez**: ou loga, ou retorna. Não os dois.
5. Crie **sentinelas** (`ErrX`) para condições que o chamador precisa diferenciar.
6. Crie **tipos de erro** quando precisar carregar dados.
7. Mensagens em **minúsculas** e **sem pontuação final**.
8. Retorne **`nil` literal** quando não houver erro (lembra da pegadinha do módulo 10?).
9. **`panic` só para bugs** e falhas irrecuperáveis de inicialização.
10. **`recover` só nas bordas** (servidores, workers), nunca como "try/catch".

---

## 🧾 Resumão do módulo

| Ferramenta | Uso |
|---|---|
| `errors.New("msg")` | erro simples |
| `fmt.Errorf("ctx: %w", err)` | adicionar contexto e embrulhar |
| `var ErrX = errors.New(...)` | erro sentinela |
| `type XError struct{...}` | erro com dados |
| `errors.Is(err, ErrX)` | "a corrente contém ErrX?" |
| `errors.As(err, &alvo)` | "a corrente tem um erro desse tipo?" |
| `errors.Join(e1, e2)` | juntar vários erros |
| `panic(v)` | bug / situação impossível |
| `recover()` | capturar panic (só em `defer`) |

---

## ✍️ Exercícios

1. Escreva uma função que acesse um índice de um slice e use `defer` + `recover` para transformar o panic em `error`.
2. Crie um "executor de tarefas" que roda uma lista de `func()` e continua mesmo se alguma der panic, reportando quais falharam.
3. Mostre, com código, que `os.Exit` **não** executa os `defer`, mas `panic` executa.
4. Reescreva um programa que usa `log.Fatal` em vários lugares para o padrão `main → run() error`.
5. Para cada situação, diga se usaria `error` ou `panic`:
   - usuário digitou letras onde deveria ser um número
   - uma regex fixa no código está com a sintaxe errada
   - o banco de dados está fora do ar
   - um `switch` recebeu um valor de enum que não existe
6. **Projeto:** revise todo o seu `main.go` aplicando o que aprendeu: erros sentinela, contexto com `%w`, e o padrão `main → run()`.

---

⬅️ Anterior: [Wrapping, Is, As e Join](03-wrapping-is-as-join.md) · ➡️ Próximo módulo: [Pacotes e Módulos](../12-pacotes-e-modulos/README.md)

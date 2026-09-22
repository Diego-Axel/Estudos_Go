# 01 — Erros São Valores

## 🧠 A filosofia de Go

Em Java, Python ou JavaScript, erros são **exceções**: elas "explodem" e sobem pela pilha até alguém capturar com `try/catch`.

Em Go, **não existe `try/catch`**. Um erro é um **valor comum**, retornado pela função como qualquer outro:

```go
arquivo, err := os.Open("dados.txt")
if err != nil {
	// trata o erro aqui, AGORA
}
```

> *"Errors are values."* (Rob Pike)

Por quê?
- ✅ O erro fica **visível** na assinatura: `func Abrir() (*File, error)`. Você **sabe** que pode falhar.
- ✅ O fluxo é **explícito**: não tem salto escondido para um `catch` 5 funções acima.
- ✅ Você é **forçado a pensar** no que fazer quando algo dá errado.
- ❌ O custo: mais linhas de `if err != nil`. (Com o tempo você acostuma e até gosta. 😄)

---

## 🧩 A interface `error`

```go
type error interface {
	Error() string
}
```

Qualquer tipo com `Error() string` é um erro (visto no módulo 10). E **`nil` significa "sem erro"**.

---

## 🛠️ Criando erros

### `errors.New`: mensagem fixa

```go
import "errors"

func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divisão por zero")
	}
	return a / b, nil
}
```

### `fmt.Errorf`: mensagem formatada

```go
func buscarUsuario(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("id inválido: %d", id)
	}
	return "Ana", nil
}
```

---

## ✅ O padrão `if err != nil`

```go
func main() {
	resultado, err := dividir(10, 0)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	fmt.Println("Resultado:", resultado)
}
```

Regras de ouro:

1. O `error` é **sempre o último** valor de retorno
2. Se `err != nil`, os **outros** valores de retorno normalmente devem ser **ignorados** (geralmente são valores zero)
3. Trate o erro **logo depois** da chamada
4. Mantenha o "caminho feliz" **sem indentação** (retorno antecipado, módulo 04)

```go
// ✅ Estilo Go: erros tratados e "caminho feliz" alinhado à esquerda
func processar(caminho string) error {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}

	config, err := interpretar(dados)
	if err != nil {
		return err
	}

	return aplicar(config)
}
```

---

## 🤔 O que fazer com um erro?

| Estratégia | Quando | Exemplo |
|---|---|---|
| **Propagar** (retornar) | a função atual não sabe resolver | `return fmt.Errorf("carregar config: %w", err)` |
| **Tentar de novo** | falhas temporárias (rede) | laço com algumas tentativas |
| **Usar um valor padrão** | o erro não é grave | porta inválida → usa 8080 |
| **Registrar e continuar** | o programa pode seguir | `log.Println("aviso:", err)` |
| **Encerrar o programa** | só no `main`, erro fatal | `log.Fatal(err)` |

### Exemplo: valor padrão

```go
porta, err := strconv.Atoi(os.Getenv("PORTA"))
if err != nil {
	porta = 8080 // variável ausente ou inválida → usa o padrão
}
```

### Exemplo: tentar de novo

```go
var resp *http.Response
var err error

for tentativa := 1; tentativa <= 3; tentativa++ {
	resp, err = http.Get("https://api.exemplo.com")
	if err == nil {
		break
	}
	fmt.Printf("tentativa %d falhou: %v\n", tentativa, err)
	time.Sleep(time.Second * time.Duration(tentativa))
}
if err != nil {
	log.Fatal("desisti:", err)
}
defer resp.Body.Close()
```

---

## 🚫 Nunca ignore erros silenciosamente

```go
valor, _ := strconv.Atoi(entrada) // ❌ se falhar, valor = 0 e ninguém fica sabendo

os.Remove("temp.txt") // ❌ retorno de erro descartado
```

Se ignorar for **realmente** intencional, deixe isso claro:

```go
_ = os.Remove("temp.txt") // ignorado de propósito: tudo bem se o arquivo não existir
```

> 🛠️ Ferramentas como `errcheck` e `golangci-lint` apontam erros não verificados.

---

## ✍️ Estilo das mensagens de erro

Convenções da comunidade Go:

- **Minúsculas** no início: `"arquivo não encontrado"`, não `"Arquivo não encontrado"`
- **Sem pontuação** no final: nada de `.` ou `!`
- **Descreva o que falhou**, com contexto: `"abrir config.json: permissão negada"`

Motivo: erros costumam ser **encadeados**, e a mensagem final fica natural:

```
carregar configuração: abrir config.json: open config.json: permission denied
```

---

## 🧪 Erros da biblioteca padrão que você vai ver muito

```go
_, err := strconv.Atoi("abc")
fmt.Println(err) // strconv.Atoi: parsing "abc": invalid syntax

_, err = os.Open("nao-existe.txt")
fmt.Println(err) // open nao-existe.txt: The system cannot find the file specified. (Windows)
                 // open nao-existe.txt: no such file or directory             (Linux/macOS)
```

---

## ✍️ Exercícios

1. Escreva `raiz(n float64) (float64, error)` que retorne erro para números negativos.
2. Escreva `lerIdade(texto string) (int, error)` que retorne erro se o texto não for número **ou** se a idade estiver fora de 0 a 150.
3. Leia um número do teclado repetidamente **até** o usuário digitar um valor válido.
4. Escreva `sacar(saldo, valor float64) (float64, error)` com erros diferentes para valor negativo e saldo insuficiente.
5. Encontre no seu `main.go` os lugares onde erros de `strconv` são tratados. Estão sendo tratados de forma adequada? O que você melhoraria?
6. Reescreva estas mensagens seguindo as convenções de Go: `"Erro ao abrir o arquivo!"`, `"Usuário Não Encontrado."`, `"FALHA NA CONEXÃO"`.

---

🏠 [Módulo 11](README.md) · ➡️ Próximo: [Erros sentinela e personalizados](02-erros-sentinela-e-personalizados.md)

# 04 — Logs, Regex e Outros Pacotes Úteis

## 📝 `log`: o log clássico

```go
import "log"

log.Println("servidor iniciado")          // 2025/09/22 15:04:05 servidor iniciado
log.Printf("porta %d em uso", 8080)
log.Fatal("erro fatal")                   // imprime e chama os.Exit(1)
log.Panic("algo impossível")              // imprime e chama panic
```

Configurando:

```go
log.SetPrefix("[livraria] ")
log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
log.Println("teste")
// [livraria] 2025/09/22 15:04:05 main.go:12: teste

// Escrevendo em arquivo
f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
log.SetOutput(f)

// Um logger separado
erros := log.New(os.Stderr, "ERRO: ", log.LstdFlags)
erros.Println("falha ao conectar")
```

> O `log` escreve em `os.Stderr` por padrão e já é **seguro para concorrência**.

---

## 📊 `log/slog`: logs estruturados (Go 1.21+) ⭐

Em sistemas reais, logs são lidos por **ferramentas** (Grafana, Datadog, ELK). Para isso, eles precisam ser **estruturados**: pares de **chave=valor**, não só texto.

```go
import "log/slog"

slog.Info("pedido criado", "id", 42, "cliente", "Ana", "total", 99.9)
slog.Warn("estoque baixo", "produto", "L001", "restante", 2)
slog.Error("falha no pagamento", "err", err)
slog.Debug("detalhe interno") // não aparece por padrão (nível mínimo é Info)
```

Saída padrão:

```
2025/09/22 15:04:05 INFO pedido criado id=42 cliente=Ana total=99.9
2025/09/22 15:04:05 WARN estoque baixo produto=L001 restante=2
```

### Saída em JSON

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug, // mostra também Debug
}))
slog.SetDefault(logger)

slog.Info("pedido criado", "id", 42, "cliente", "Ana")
```

```json
{"time":"2025-09-22T15:04:05.123-03:00","level":"INFO","msg":"pedido criado","id":42,"cliente":"Ana"}
```

### Campos fixos com `With`

```go
reqLogger := slog.With("request_id", "abc-123", "usuario", "diego")
reqLogger.Info("buscando produtos")   // já inclui request_id e usuario
reqLogger.Info("produtos encontrados", "qtd", 12)
```

Também existem os handlers de texto (`slog.NewTextHandler`) e atributos tipados (`slog.Int("id", 42)`, `slog.String(...)`, `slog.Group(...)`).

> ✅ Para projetos novos, **prefira `slog`** ao `log`.

---

## 🔎 `regexp`: expressões regulares

```go
import "regexp"

// Compile uma vez (de preferência como variável de pacote)
var reEmail = regexp.MustCompile(`^[\w.+-]+@[\w-]+\.[\w.]+$`)

fmt.Println(reEmail.MatchString("ana@email.com")) // true
fmt.Println(reEmail.MatchString("ana@"))          // false
```

> 💡 Use **crases** (string crua) para regex: assim `\w` não precisa virar `\\w`.

### Operações comuns

```go
re := regexp.MustCompile(`\d+`)
texto := "Pedido 123 com 4 itens custou 250 reais"

re.FindString(texto)          // "123" (primeira ocorrência)
re.FindAllString(texto, -1)   // [123 4 250] (todas; -1 = sem limite)
re.ReplaceAllString(texto, "#") // "Pedido # com # itens custou # reais"
re.Split("a1b22c333d", -1)    // [a b c d]
```

### Grupos de captura

```go
reData := regexp.MustCompile(`(\d{2})/(\d{2})/(\d{4})`)

m := reData.FindStringSubmatch("Nascimento: 15/08/1998")
fmt.Println(m)             // [15/08/1998 15 08 1998]
fmt.Println(m[1], m[2], m[3]) // 15 08 1998

// Trocando o formato com os grupos
iso := reData.ReplaceAllString("15/08/1998", "$3-$2-$1")
fmt.Println(iso) // 1998-08-15
```

Grupos com nome:

```go
re := regexp.MustCompile(`(?P<ano>\d{4})-(?P<mes>\d{2})`)
m := re.FindStringSubmatch("2025-09")
fmt.Println(m[re.SubexpIndex("ano")]) // 2025
```

> ⚠️ O `regexp` de Go usa a sintaxe **RE2**: sem *lookahead*/*lookbehind* e sem *backreferences*, mas com **tempo linear garantido** (sem risco de uma regex travar o servidor). E, às vezes, `strings.Contains`/`HasPrefix` resolvem sem regex nenhuma, e são mais rápidos.

---

## 🧮 `math` e `math/rand/v2`

```go
import "math"

math.Sqrt(16)       // 4
math.Pow(2, 10)     // 1024
math.Abs(-3.5)      // 3.5
math.Round(2.5)     // 3
math.Floor(2.9)     // 2
math.Ceil(2.1)      // 3
math.Max(3, 7)      // 7 (ou o built-in max(3, 7))
math.Pi             // 3.141592653589793
math.Inf(1)         // +Inf
math.MaxInt64       // 9223372036854775807
```

```go
import "math/rand/v2"

rand.IntN(100)          // inteiro de 0 a 99
rand.Float64()          // float de 0.0 a 1.0
rand.N(10 * time.Second) // uma Duration aleatória (genérico!)

nomes := []string{"Ana", "Bia", "Caio"}
rand.Shuffle(len(nomes), func(i, j int) { nomes[i], nomes[j] = nomes[j], nomes[i] })
sorteado := nomes[rand.IntN(len(nomes))]
```

> ⚠️ `math/rand` **não** é seguro para senhas, tokens ou criptografia. Para isso, use **`crypto/rand`**:

```go
import "crypto/rand"

token := rand.Text() // Go 1.24+: string aleatória segura (base32)
fmt.Println(token)
```

---

## 🔐 Hashes e codificações

```go
import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

hash := sha256.Sum256([]byte("minha senha"))
fmt.Println(hex.EncodeToString(hash[:])) // 64 caracteres hexadecimais

codificado := base64.StdEncoding.EncodeToString([]byte("Olá, Go!"))
fmt.Println(codificado) // T2zDoSwgR28h
decodificado, _ := base64.StdEncoding.DecodeString(codificado)
fmt.Println(string(decodificado)) // Olá, Go!
```

> ⚠️ **Nunca** guarde senhas com SHA-256 puro. Use `golang.org/x/crypto/bcrypt` ou `argon2`, que são feitos para isso.

---

## ⚙️ `os/exec`: rodando outros programas

```go
import "os/exec"

saida, err := exec.Command("git", "log", "--oneline", "-5").Output()
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(saida))
```

Com contexto (timeout):

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := exec.CommandContext(ctx, "ping", "go.dev").Run()
```

---

## 🗺️ Mapa rápido da biblioteca padrão

| Pacote | Para quê |
|---|---|
| `fmt` | formatação e impressão |
| `strings`, `strconv`, `unicode`, `bytes` | texto (módulo 07) |
| `slices`, `maps`, `cmp`, `iter` | coleções genéricas (módulos 06 e 13) |
| `errors` | erros (módulo 11) |
| `os`, `io`, `bufio`, `path/filepath`, `io/fs` | arquivos e I/O |
| `flag` | opções de linha de comando |
| `encoding/json`, `encoding/csv`, `encoding/xml`, `encoding/base64` | formatos de dados |
| `time` | datas, durações, timers |
| `log`, `log/slog` | logs |
| `regexp` | expressões regulares |
| `math`, `math/rand/v2`, `math/big` | matemática |
| `crypto/*` | criptografia, hashes, `crypto/rand` |
| `sync`, `sync/atomic`, `context` | concorrência (módulo 14) |
| `net/http`, `net/url`, `html/template` | web (módulo 17) |
| `database/sql` | bancos de dados (com um driver externo) |
| `testing`, `net/http/httptest` | testes (módulo 15) |
| `embed` | arquivos dentro do binário (módulo 12) |
| `os/exec`, `os/signal` | processos e sinais do sistema |
| `archive/zip`, `compress/gzip` | compactação |
| `sort`, `container/heap`, `container/list` | ordenação e estruturas clássicas |

> 🧭 Documentação completa: **[pkg.go.dev/std](https://pkg.go.dev/std)**. Vale navegar de vez em quando: sempre tem algo útil que você não conhecia.

---

## ✍️ Exercícios

1. Troque todos os `fmt.Println` de erro do seu programa por `slog.Error`, com campos estruturados.
2. Configure o `slog` para escrever em JSON num arquivo `app.log`.
3. Escreva uma regex que valide **CPF** no formato `000.000.000-00` e outra que extraia todos os **telefones** `(00) 00000-0000` de um texto.
4. Converta todas as datas `dd/mm/aaaa` de um texto para `aaaa-mm-dd` usando grupos de captura.
5. Faça um **sorteador**: leia nomes de um arquivo e sorteie 3 sem repetição.
6. Gere um token seguro com `crypto/rand` e compare com um gerado por `math/rand/v2`. Por que um é seguro e o outro não?
7. Calcule o SHA-256 de um arquivo (use `sha256.New()` + `io.Copy`, já que ele é um `io.Writer`!).

---

⬅️ Anterior: [Datas e horários](03-datas-e-horarios.md) · 🏠 [Voltar ao roteiro](../README.md)

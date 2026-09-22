# 04 — Cliente HTTP e Templates HTML

## 📡 Fazendo requisições: o cliente HTTP

O mesmo pacote `net/http` também **consome** APIs.

### GET simples

```go
resp, err := http.Get("https://api.github.com/users/Diego-Axel")
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close() // ⚠️ SEMPRE feche o corpo da resposta

if resp.StatusCode != http.StatusOK {
	log.Fatalf("status inesperado: %s", resp.Status)
}

corpo, err := io.ReadAll(resp.Body)
fmt.Println(string(corpo))
```

> ⚠️ `err == nil` **não** significa sucesso! Um 404 ou 500 também retorna `err == nil`. **Sempre** verifique o `resp.StatusCode`.

### Decodificando JSON da resposta ⭐

```go
type UsuarioGitHub struct {
	Login       string `json:"login"`
	Nome        string `json:"name"`
	Repos       int    `json:"public_repos"`
	Seguidores  int    `json:"followers"`
}

resp, err := http.Get("https://api.github.com/users/Diego-Axel")
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()

var u UsuarioGitHub
if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
	log.Fatal(err)
}
fmt.Printf("%s tem %d repositórios públicos\n", u.Login, u.Repos)
```

> Só os campos que você declarar são lidos: o resto do JSON é ignorado.

---

## ⏱️ Sempre use um cliente com timeout

O `http.Get` usa o `http.DefaultClient`, que **não tem timeout**: se o servidor travar, seu programa espera **para sempre**.

```go
cliente := &http.Client{
	Timeout: 10 * time.Second, // ✅ tempo máximo da requisição inteira
}

resp, err := cliente.Get("https://go.dev")
```

> ✅ Crie **um** cliente e **reutilize** (ele mantém um pool de conexões). Não crie um novo a cada requisição.

---

## 🛠️ Requisições completas: `http.NewRequestWithContext`

Para escolher método, cabeçalhos, corpo e contexto:

### POST com JSON

```go
novo := map[string]any{
	"codigo": "L010",
	"titulo": "O Senhor dos Anéis",
	"preco":  129.9,
	"estoque": 2,
}
corpo, _ := json.Marshal(novo)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, http.MethodPost,
	"http://localhost:8080/produtos", bytes.NewReader(corpo))
if err != nil {
	log.Fatal(err)
}
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer meu-token")

resp, err := cliente.Do(req)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()

fmt.Println("status:", resp.Status) // status: 201 Created
```

### Query string com `net/url`

Monte parâmetros com segurança (acentos e espaços são codificados automaticamente):

```go
params := url.Values{}
params.Set("q", "linguagem go")
params.Set("pagina", "2")

endereco := "https://api.exemplo.com/busca?" + params.Encode()
fmt.Println(endereco) // https://api.exemplo.com/busca?pagina=2&q=linguagem+go
```

---

## 🧰 Um cliente reutilizável para uma API

Juntando tudo num tipo organizado:

```go
type ClienteLivraria struct {
	baseURL string
	http    *http.Client
}

func NovoClienteLivraria(baseURL string) *ClienteLivraria {
	return &ClienteLivraria{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *ClienteLivraria) Buscar(ctx context.Context, codigo string) (Produto, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/produtos/"+url.PathEscape(codigo), nil)
	if err != nil {
		return Produto{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Produto{}, fmt.Errorf("buscar produto %s: %w", codigo, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var p Produto
		err := json.NewDecoder(resp.Body).Decode(&p)
		return p, err
	case http.StatusNotFound:
		return Produto{}, ErrNaoEncontrado
	default:
		return Produto{}, fmt.Errorf("status inesperado: %s", resp.Status)
	}
}

// Uso:
cli := NovoClienteLivraria("http://localhost:8080")
p, err := cli.Buscar(context.Background(), "L001")
```

---

## 🖼️ Templates HTML com `html/template`

Para gerar **páginas HTML** no servidor (sem front-end separado):

```go
import "html/template"

const paginaProdutos = `<!DOCTYPE html>
<html lang="pt-BR">
<head><meta charset="utf-8"><title>{{.Titulo}}</title></head>
<body>
  <h1>{{.Titulo}}</h1>
  {{if .Produtos}}
    <ul>
    {{range .Produtos}}
      <li>{{.Titulo}}: R$ {{printf "%.2f" .Preco}}
        {{if eq .Estoque 0}}<strong>(esgotado)</strong>{{end}}
      </li>
    {{end}}
    </ul>
  {{else}}
    <p>Nenhum produto cadastrado.</p>
  {{end}}
</body>
</html>`

var tmpl = template.Must(template.New("produtos").Parse(paginaProdutos))

func (a *API) paginaInicial(w http.ResponseWriter, r *http.Request) {
	dados := struct {
		Titulo   string
		Produtos []Produto
	}{
		Titulo:   "📚 Livraria do Diego",
		Produtos: a.loja.Listar(),
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "erro ao gerar página", http.StatusInternalServerError)
	}
}
```

### Sintaxe dos templates

| Ação | Significado |
|---|---|
| `{{.Campo}}` | valor de um campo (o `.` é o dado atual) |
| `{{if .X}} ... {{else}} ... {{end}}` | condicional |
| `{{range .Lista}} ... {{end}}` | laço (dentro, o `.` vira cada item) |
| `{{printf "%.2f" .Preco}}` | chama funções |
| `{{eq .A .B}}`, `{{lt .A .B}}`, `{{len .Lista}}` | comparações e funções prontas |
| `{{template "nome" .}}` | inclui outro template |
| `{{/* comentário */}}` | comentário |

### Templates em arquivos (com `embed`)

```go
//go:embed templates/*.html
var arquivosTemplate embed.FS

var tmpls = template.Must(template.ParseFS(arquivosTemplate, "templates/*.html"))

tmpls.ExecuteTemplate(w, "produtos.html", dados)
```

### 🛡️ Segurança automática

O `html/template` **escapa** tudo automaticamente, protegendo contra **XSS** (injeção de scripts):

```go
dados := map[string]string{"Nome": "<script>alert('hack')</script>"}
// Template: <p>Olá, {{.Nome}}</p>
// Saída:    <p>Olá, &lt;script&gt;alert(&#39;hack&#39;)&lt;/script&gt;</p>  ✅ inofensivo
```

> ⚠️ Para HTML, use **sempre** `html/template`, **nunca** `text/template` (que não escapa nada). O `text/template` é para gerar texto comum: e-mails em texto puro, arquivos de configuração, código.

---

## 🧾 Resumão do módulo

| Recurso | Uso |
|---|---|
| `http.HandleFunc` / `ServeMux` | registrar rotas |
| `"GET /produtos/{id}"` + `r.PathValue` | rotas com método e parâmetros (Go 1.22+) |
| `w.Header().Set`, `w.WriteHeader`, `Write` | montar a resposta |
| `json.NewEncoder(w)` / `NewDecoder(r.Body)` | API JSON |
| middleware `func(http.Handler) http.Handler` | log, recover, auth, CORS |
| `http.Server{...Timeout}` | servidor pronto para produção |
| `srv.Shutdown(ctx)` + `signal.NotifyContext` | desligamento elegante |
| `http.Client{Timeout}` + `NewRequestWithContext` | consumir APIs |
| `html/template` | páginas HTML seguras |
| `httptest` | testar handlers |

---

## ✍️ Exercícios

1. Consuma a API do GitHub e mostre nome, repositórios e seguidores de um usuário lido do teclado.
2. Faça um programa que consulte o CEP na API pública ViaCEP (`https://viacep.com.br/ws/{cep}/json/`) e mostre o endereço.
3. Crie um cliente para a sua API da livraria com os métodos `Listar`, `Buscar` e `Criar`.
4. Faça 10 requisições **em paralelo** (módulo 14) a uma API e mostre o tempo total vs. sequencial.
5. Crie uma página HTML com template que lista os produtos da livraria, com um formulário para cadastrar novos.
6. Teste a proteção contra XSS: cadastre um produto com título `<b>negrito</b>` e veja como aparece na página.
7. **Projeto final:** junte tudo: API REST + página HTML + middlewares + graceful shutdown + persistência em JSON + testes.

---

⬅️ Anterior: [Middlewares e servidor robusto](03-middlewares-e-servidor-robusto.md) · 🏠 [Voltar ao roteiro](../README.md)

# 01 — Servidor HTTP

Go tem um servidor web **de produção** na biblioteca padrão: o pacote `net/http`. Muitas empresas usam **só ele**, sem framework nenhum.

---

## 👋 O menor servidor possível

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Olá, web! 🌐")
	})

	log.Println("servidor em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

```bash
go run .
# abra http://localhost:8080 no navegador
# ou, no terminal:
curl http://localhost:8080
```

Cada requisição roda na **sua própria goroutine** automaticamente. Milhares de acessos simultâneos? Sem problema. 🚀

---

## 🧩 As peças

### `http.ResponseWriter`: onde você escreve a resposta

É um `io.Writer` (módulo 10), então funciona com `fmt.Fprintf`, `json.NewEncoder`, `io.Copy`...

```go
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // 1º: cabeçalhos
	w.WriteHeader(http.StatusCreated)                           // 2º: status (opcional; padrão 200)
	fmt.Fprintln(w, "criado!")                                  // 3º: corpo
}
```

> ⚠️ A **ordem importa**: cabeçalhos → status → corpo. Depois que o corpo começa a ser escrito, mudar cabeçalhos ou status **não tem efeito**.

### `*http.Request`: tudo sobre a requisição

```go
func info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Método:", r.Method)                     // GET, POST...
	fmt.Fprintln(w, "Caminho:", r.URL.Path)                  // /info
	fmt.Fprintln(w, "Query:", r.URL.Query().Get("nome"))     // ?nome=Diego
	fmt.Fprintln(w, "User-Agent:", r.Header.Get("User-Agent"))
	fmt.Fprintln(w, "IP:", r.RemoteAddr)
}
```

```bash
curl "http://localhost:8080/info?nome=Diego"
```

### `http.Handler`: a interface por trás de tudo

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

Qualquer tipo com esse método responde requisições:

```go
type Contador struct {
	mu     sync.Mutex
	visitas int
}

func (c *Contador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	c.visitas++
	n := c.visitas
	c.mu.Unlock()
	fmt.Fprintf(w, "Você é o visitante nº %d\n", n)
}

http.Handle("/visitas", &Contador{})
```

> 💡 `http.HandlerFunc` é um tipo que transforma uma **função comum** em `Handler`. É o que o `HandleFunc` faz por baixo.

---

## 🗺️ Roteamento com `http.ServeMux` (Go 1.22+) ⭐

Desde o Go 1.22, o roteador padrão aceita **métodos** e **parâmetros no caminho**:

```go
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /produtos", listarProdutos)
	mux.HandleFunc("POST /produtos", criarProduto)
	mux.HandleFunc("GET /produtos/{codigo}", buscarProduto)
	mux.HandleFunc("DELETE /produtos/{codigo}", removerProduto)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func buscarProduto(w http.ResponseWriter, r *http.Request) {
	codigo := r.PathValue("codigo") // pega o {codigo} da URL
	fmt.Fprintf(w, "Você pediu o produto %s\n", codigo)
}
```

```bash
curl http://localhost:8080/produtos/L001    # Você pediu o produto L001
curl -X PUT http://localhost:8080/produtos  # 405 Method Not Allowed (automático!)
```

### Regras dos padrões

| Padrão | Casa com |
|---|---|
| `"/produtos"` | exatamente `/produtos` (qualquer método) |
| `"GET /produtos"` | `GET` (e `HEAD`) em `/produtos` |
| `"/produtos/{id}"` | `/produtos/42`, `/produtos/abc` |
| `"/arquivos/{caminho...}"` | `/arquivos/a/b/c.txt` (o resto todo) |
| `"/estaticos/"` | tudo que **começa** com `/estaticos/` |
| `"/{$}"` | **só** a raiz `/` |
| `"/"` | **qualquer** caminho (é o "pega-tudo") |

- Quando dois padrões casam, vence o **mais específico**
- Método errado → **405** automático; caminho inexistente → **404** automático

> Antes do Go 1.22, era preciso usar bibliotecas como `gorilla/mux` ou `chi` para isso. Hoje, para a maioria dos projetos, o `ServeMux` padrão basta.

---

## 🚦 Códigos de status

Use as **constantes** em vez de números soltos:

| Constante | Código | Quando |
|---|---|---|
| `http.StatusOK` | 200 | sucesso |
| `http.StatusCreated` | 201 | recurso criado (POST) |
| `http.StatusNoContent` | 204 | sucesso sem corpo (DELETE) |
| `http.StatusBadRequest` | 400 | dados inválidos |
| `http.StatusUnauthorized` | 401 | não autenticado |
| `http.StatusForbidden` | 403 | sem permissão |
| `http.StatusNotFound` | 404 | não encontrado |
| `http.StatusMethodNotAllowed` | 405 | método errado |
| `http.StatusConflict` | 409 | conflito (ex: código duplicado) |
| `http.StatusInternalServerError` | 500 | erro no servidor |

Atalhos para erros:

```go
http.Error(w, "produto não encontrado", http.StatusNotFound)
http.NotFound(w, r)
http.Redirect(w, r, "/novo-endereco", http.StatusMovedPermanently)
```

---

## 📝 Lendo formulários

```go
func contato(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulário inválido", http.StatusBadRequest)
		return
	}
	nome := r.FormValue("nome") // funciona para query string e corpo de formulário
	email := r.FormValue("email")
	fmt.Fprintf(w, "Obrigado, %s! Responderemos em %s\n", nome, email)
}
```

```bash
curl -X POST -d "nome=Diego&email=diego@email.com" http://localhost:8080/contato
```

---

## 📁 Servindo arquivos estáticos

```go
// Tudo da pasta ./public fica acessível em /static/...
fs := http.FileServer(http.Dir("./public"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

// Um arquivo específico
mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/favicon.ico")
})
```

Com `embed` (módulo 12), os arquivos vão **dentro** do executável:

```go
//go:embed public
var arquivos embed.FS

sub, _ := fs.Sub(arquivos, "public")
mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(sub)))
```

---

## ✍️ Exercícios

1. Crie um servidor com as rotas `/` (boas-vindas), `/hora` (hora atual) e `/ola/{nome}` (saudação personalizada).
2. Adicione `/soma?a=2&b=3` que retorna a soma. Retorne **400** se `a` ou `b` não forem números.
3. Crie um `Handler` (struct com `ServeHTTP`) que conta quantas vezes cada caminho foi acessado (use `map` + `Mutex`).
4. Crie uma página HTML simples em `public/index.html` e sirva com `http.FileServer`.
5. Teste o **405** automático: registre só `GET /produtos` e faça um `POST`.
6. Crie um formulário HTML de contato e um handler `POST /contato` que mostra os dados recebidos.

---

🏠 [Módulo 17](README.md) · ➡️ Próximo: [API REST com JSON](02-api-rest-com-json.md)

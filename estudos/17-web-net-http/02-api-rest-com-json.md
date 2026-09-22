# 02 — API REST com JSON

Vamos transformar a **livraria** do seu `main.go` numa **API REST** completa, só com a biblioteca padrão. Tudo que você aprendeu até aqui vai aparecer: structs, métodos, interfaces, erros, mutex, JSON, `net/http`.

---

## 🗺️ O que vamos construir

| Método | Rota | Ação | Sucesso |
|---|---|---|---|
| `GET` | `/produtos` | listar todos | 200 |
| `GET` | `/produtos/{codigo}` | buscar um | 200 / 404 |
| `POST` | `/produtos` | cadastrar | 201 / 400 / 409 |
| `PUT` | `/produtos/{codigo}` | atualizar | 200 / 400 / 404 |
| `DELETE` | `/produtos/{codigo}` | remover | 204 / 404 |

---

## 1️⃣ O modelo e o armazenamento

```go
package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
)

type Produto struct {
	Codigo  string  `json:"codigo"`
	Titulo  string  `json:"titulo"`
	Autor   string  `json:"autor"`
	Preco   float64 `json:"preco"`
	Estoque int     `json:"estoque"`
}

var (
	ErrNaoEncontrado = errors.New("produto não encontrado")
	ErrDuplicado     = errors.New("código já cadastrado")
)

// Loja guarda os produtos em memória, com segurança para concorrência.
type Loja struct {
	mu       sync.RWMutex
	produtos map[string]Produto
}

func NovaLoja() *Loja {
	return &Loja{produtos: make(map[string]Produto)}
}

func (l *Loja) Listar() []Produto {
	l.mu.RLock()
	defer l.mu.RUnlock()

	lista := make([]Produto, 0, len(l.produtos))
	for _, p := range l.produtos {
		lista = append(lista, p)
	}
	slices.SortFunc(lista, func(a, b Produto) int { return cmp.Compare(a.Codigo, b.Codigo) })
	return lista
}

func (l *Loja) Buscar(codigo string) (Produto, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	p, ok := l.produtos[codigo]
	if !ok {
		return Produto{}, ErrNaoEncontrado
	}
	return p, nil
}

func (l *Loja) Criar(p Produto) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, existe := l.produtos[p.Codigo]; existe {
		return ErrDuplicado
	}
	l.produtos[p.Codigo] = p
	return nil
}

func (l *Loja) Atualizar(codigo string, p Produto) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, existe := l.produtos[codigo]; !existe {
		return ErrNaoEncontrado
	}
	p.Codigo = codigo
	l.produtos[codigo] = p
	return nil
}

func (l *Loja) Remover(codigo string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, existe := l.produtos[codigo]; !existe {
		return ErrNaoEncontrado
	}
	delete(l.produtos, codigo)
	return nil
}
```

---

## 2️⃣ Validação

```go
func (p Produto) Validar() error {
	var errs []error
	if strings.TrimSpace(p.Codigo) == "" {
		errs = append(errs, errors.New("codigo é obrigatório"))
	}
	if strings.TrimSpace(p.Titulo) == "" {
		errs = append(errs, errors.New("titulo é obrigatório"))
	}
	if p.Preco <= 0 {
		errs = append(errs, errors.New("preco deve ser positivo"))
	}
	if p.Estoque < 0 {
		errs = append(errs, errors.New("estoque não pode ser negativo"))
	}
	return errors.Join(errs...) // nil se não houver problemas (módulo 11)
}
```

---

## 3️⃣ Funções auxiliares para responder JSON

Evitam repetir o mesmo código em todo handler:

```go
func responderJSON(w http.ResponseWriter, status int, dados any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(dados); err != nil {
		log.Printf("erro ao escrever JSON: %v", err)
	}
}

func responderErro(w http.ResponseWriter, status int, msg string) {
	responderJSON(w, status, map[string]string{"erro": msg})
}

func lerJSON(w http.ResponseWriter, r *http.Request, destino any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limita o corpo a 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(destino)
}
```

---

## 4️⃣ Os handlers

Os handlers ficam como **métodos** de uma struct que guarda as dependências (a `Loja`). Assim nada fica em variável global:

```go
type API struct {
	loja *Loja
}

func (a *API) listar(w http.ResponseWriter, r *http.Request) {
	responderJSON(w, http.StatusOK, a.loja.Listar())
}

func (a *API) buscar(w http.ResponseWriter, r *http.Request) {
	p, err := a.loja.Buscar(r.PathValue("codigo"))
	if errors.Is(err, ErrNaoEncontrado) {
		responderErro(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func (a *API) criar(w http.ResponseWriter, r *http.Request) {
	var p Produto
	if err := lerJSON(w, r, &p); err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if err := p.Validar(); err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	switch err := a.loja.Criar(p); {
	case errors.Is(err, ErrDuplicado):
		responderErro(w, http.StatusConflict, err.Error())
	case err != nil:
		responderErro(w, http.StatusInternalServerError, "erro interno")
	default:
		w.Header().Set("Location", "/produtos/"+p.Codigo)
		responderJSON(w, http.StatusCreated, p)
	}
}

func (a *API) atualizar(w http.ResponseWriter, r *http.Request) {
	codigo := r.PathValue("codigo")

	var p Produto
	if err := lerJSON(w, r, &p); err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	p.Codigo = codigo // o código vem da URL
	if err := p.Validar(); err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.loja.Atualizar(codigo, p); errors.Is(err, ErrNaoEncontrado) {
		responderErro(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func (a *API) remover(w http.ResponseWriter, r *http.Request) {
	if err := a.loja.Remover(r.PathValue("codigo")); errors.Is(err, ErrNaoEncontrado) {
		responderErro(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
```

---

## 5️⃣ Juntando tudo no `main`

```go
func (a *API) Rotas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /produtos", a.listar)
	mux.HandleFunc("POST /produtos", a.criar)
	mux.HandleFunc("GET /produtos/{codigo}", a.buscar)
	mux.HandleFunc("PUT /produtos/{codigo}", a.atualizar)
	mux.HandleFunc("DELETE /produtos/{codigo}", a.remover)
	return mux
}

func main() {
	api := &API{loja: NovaLoja()}

	log.Println("API da livraria em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", api.Rotas()))
}
```

---

## 🧪 Testando com `curl`

```bash
# Cadastrar
curl -i -X POST http://localhost:8080/produtos \
  -H "Content-Type: application/json" \
  -d '{"codigo":"L001","titulo":"O Hobbit","autor":"Tolkien","preco":49.9,"estoque":3}'
# HTTP/1.1 201 Created
# Location: /produtos/L001

# Listar
curl http://localhost:8080/produtos
# [{"codigo":"L001","titulo":"O Hobbit","autor":"Tolkien","preco":49.9,"estoque":3}]

# Buscar um que não existe
curl -i http://localhost:8080/produtos/X999
# HTTP/1.1 404 Not Found
# {"erro":"produto não encontrado"}

# Dados inválidos
curl -X POST http://localhost:8080/produtos -d '{"codigo":"L002","preco":-5}'
# {"erro":"titulo é obrigatório\npreco deve ser positivo"}

# Remover
curl -i -X DELETE http://localhost:8080/produtos/L001
# HTTP/1.1 204 No Content
```

> 💡 **No PowerShell do Windows**, `curl` pode ser um apelido para `Invoke-WebRequest`. Use `curl.exe` ou ferramentas como **Postman**, **Insomnia**, **Bruno** ou a extensão **REST Client** do VS Code.

---

## 🧪 Testando com `httptest` (módulo 15)

```go
func TestCriarEBuscar(t *testing.T) {
	api := &API{loja: NovaLoja()}
	srv := api.Rotas()

	// POST
	corpo := strings.NewReader(`{"codigo":"L001","titulo":"Duna","preco":79.9,"estoque":1}`)
	req := httptest.NewRequest("POST", "/produtos", corpo)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d; esperado 201. corpo: %s", rec.Code, rec.Body)
	}

	// GET
	req = httptest.NewRequest("GET", "/produtos/L001", nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	var p Produto
	json.NewDecoder(rec.Body).Decode(&p)
	if p.Titulo != "Duna" {
		t.Errorf("Titulo = %q; esperado Duna", p.Titulo)
	}
}
```

---

## ✍️ Exercícios

1. Monte a API completa acima e teste todas as rotas com `curl` (ou Postman).
2. Adicione `GET /produtos?autor=Tolkien` que filtra por autor (use `r.URL.Query()`).
3. Adicione `POST /produtos/{codigo}/vender` com corpo `{"quantidade": 2}`, que diminui o estoque e retorna **409** se não houver estoque suficiente.
4. Adicione **paginação**: `GET /produtos?pagina=2&tamanho=10`.
5. Faça a `Loja` **salvar em `produtos.json`** a cada alteração e carregar ao iniciar (módulo 16).
6. Troque a `Loja` concreta por uma **interface** `Repositorio` na `API` e escreva os testes com `httptest`.
7. **Projeto final:** transforme o `main.go` de terminal deste repositório nessa API, mantendo os testes passando.

---

⬅️ Anterior: [Servidor HTTP](01-servidor-http.md) · ➡️ Próximo: [Middlewares e servidor robusto](03-middlewares-e-servidor-robusto.md)

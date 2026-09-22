# 02 — JSON e CSV

## 📦 JSON com `encoding/json`

JSON é o formato mais usado para APIs, configurações e troca de dados. Go converte **structs ↔ JSON** usando as **tags** que você viu no módulo 09.

```go
type Produto struct {
	Codigo  string   `json:"codigo"`
	Titulo  string   `json:"titulo"`
	Preco   float64  `json:"preco"`
	Estoque int      `json:"estoque"`
	Tags    []string `json:"tags,omitempty"`
}
```

### Struct → JSON: `json.Marshal`

```go
p := Produto{Codigo: "L001", Titulo: "O Hobbit", Preco: 49.9, Estoque: 3}

dados, err := json.Marshal(p)
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(dados))
// {"codigo":"L001","titulo":"O Hobbit","preco":49.9,"estoque":3}
```

Formatado (bonito, com indentação):

```go
dados, _ := json.MarshalIndent(p, "", "  ")
fmt.Println(string(dados))
```

```json
{
  "codigo": "L001",
  "titulo": "O Hobbit",
  "preco": 49.9,
  "estoque": 3
}
```

### JSON → struct: `json.Unmarshal`

```go
entrada := []byte(`{"codigo":"L002","titulo":"Duna","preco":79.9,"estoque":5,"extra":"ignorado"}`)

var p Produto
if err := json.Unmarshal(entrada, &p); err != nil { // ⚠️ passe um PONTEIRO
	log.Fatal(err)
}
fmt.Printf("%+v\n", p)
// {Codigo:L002 Titulo:Duna Preco:79.9 Estoque:5 Tags:[]}
```

Comportamentos importantes:
- Campos do JSON que **não existem** na struct são **ignorados** (o `"extra"`)
- Campos da struct que **faltam** no JSON ficam com **valor zero**
- A correspondência de nomes **não diferencia** maiúsculas/minúsculas (`"TITULO"` também casaria)
- Só campos **exportados** participam

### Listas

```go
produtos := []Produto{
	{Codigo: "L001", Titulo: "O Hobbit", Preco: 49.9},
	{Codigo: "L002", Titulo: "Duna", Preco: 79.9},
}

dados, _ := json.Marshal(produtos)
// [{"codigo":"L001",...},{"codigo":"L002",...}]

var lidos []Produto
json.Unmarshal(dados, &lidos)
fmt.Println(len(lidos)) // 2
```

---

## 🌊 `Encoder` e `Decoder`: JSON direto de/para `io.Writer`/`io.Reader`

Em vez de passar por `[]byte`, você pode escrever direto num arquivo, resposta HTTP, etc.:

```go
// Salvar em arquivo
f, _ := os.Create("produtos.json")
defer f.Close()

enc := json.NewEncoder(f)
enc.SetIndent("", "  ")
if err := enc.Encode(produtos); err != nil {
	log.Fatal(err)
}

// Ler de arquivo
f2, _ := os.Open("produtos.json")
defer f2.Close()

var carregados []Produto
if err := json.NewDecoder(f2).Decode(&carregados); err != nil {
	log.Fatal(err)
}
```

> 🧠 Regra prática: tem um `[]byte` na mão → `Marshal`/`Unmarshal`. Tem um arquivo, rede ou HTTP → `Encoder`/`Decoder`.

Rejeitando campos desconhecidos (útil em APIs):

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields() // erro se vier campo que a struct não tem
```

---

## 🧩 JSON de formato desconhecido: `map[string]any`

```go
entrada := []byte(`{"nome":"Ana","idade":30,"ativo":true,"tags":["go","dev"]}`)

var dados map[string]any
json.Unmarshal(entrada, &dados)

fmt.Println(dados["nome"])            // Ana
idade := dados["idade"].(float64)     // ⚠️ números viram float64!
fmt.Println(int(idade))               // 30
tags := dados["tags"].([]any)         // arrays viram []any
fmt.Println(tags[0])                  // go
```

| JSON | Go (em `any`) |
|---|---|
| objeto `{}` | `map[string]any` |
| array `[]` | `[]any` |
| string | `string` |
| número | **`float64`** |
| `true`/`false` | `bool` |
| `null` | `nil` |

> ✅ Sempre que souber o formato, **use structs**: é mais seguro, mais rápido e dispensa type assertions.

### Structs aninhadas e anônimas

```go
var resposta struct {
	Status string `json:"status"`
	Dados  struct {
		Usuarios []struct {
			Nome string `json:"nome"`
		} `json:"usuarios"`
	} `json:"dados"`
}
json.Unmarshal(entrada, &resposta)
```

---

## 🎨 Personalizando: `MarshalJSON` e `UnmarshalJSON`

Se o tipo implementar a interface `json.Marshaler`, você controla a conversão:

```go
type Status int

const (
	Ativo Status = iota
	Inativo
)

func (s Status) MarshalJSON() ([]byte, error) {
	nomes := map[Status]string{Ativo: "ativo", Inativo: "inativo"}
	return json.Marshal(nomes[s])
}

type Usuario struct {
	Nome   string `json:"nome"`
	Status Status `json:"status"`
}

dados, _ := json.Marshal(Usuario{"Ana", Inativo})
fmt.Println(string(dados)) // {"nome":"Ana","status":"inativo"}  (em vez de 1)
```

---

## ⚠️ Pegadinhas do JSON

```go
type Config struct {
	porta int    // ❌ minúsculo: NUNCA aparece no JSON
	Host  string // ✅ sem tag: vira "Host" (com H maiúsculo)
}

var c Config
json.Unmarshal(dados, c)  // ❌ sem &: erro "non-pointer"
json.Unmarshal(dados, &c) // ✅
```

- `[]byte` vira **string base64** no JSON
- `time.Time` vira string no formato RFC 3339 (`"2025-01-15T10:30:00Z"`)
- Slice `nil` vira `null`; slice vazio `[]T{}` vira `[]`

---

## 📊 CSV com `encoding/csv`

### Lendo

```go
import "encoding/csv"

entrada := `codigo,titulo,preco
L001,O Hobbit,49.90
L002,"Duna, o livro",79.90`

r := csv.NewReader(strings.NewReader(entrada)) // ou um *os.File
registros, err := r.ReadAll()
if err != nil {
	log.Fatal(err)
}

for i, linha := range registros {
	if i == 0 {
		continue // pula o cabeçalho
	}
	preco, _ := strconv.ParseFloat(linha[2], 64)
	fmt.Printf("%s | %-15s | R$ %.2f\n", linha[0], linha[1], preco)
}
// L001 | O Hobbit        | R$ 49.90
// L002 | Duna, o livro   | R$ 79.90
```

> O pacote trata **aspas** e **vírgulas dentro de campos** corretamente, coisa que um simples `strings.Split(linha, ",")` **não** faz.

Para arquivos grandes, leia **registro por registro**:

```go
for {
	linha, err := r.Read()
	if err == io.EOF {
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(linha)
}
```

CSV com ponto e vírgula (comum no Excel em português):

```go
r.Comma = ';'
```

### Escrevendo

```go
f, _ := os.Create("produtos.csv")
defer f.Close()

w := csv.NewWriter(f)
w.Write([]string{"codigo", "titulo", "preco"})
for _, p := range produtos {
	w.Write([]string{p.Codigo, p.Titulo, strconv.FormatFloat(p.Preco, 'f', 2, 64)})
}
w.Flush() // ⚠️ igual ao bufio.Writer: não esqueça!
if err := w.Error(); err != nil {
	log.Fatal(err)
}
```

---

## ✍️ Exercícios

1. Converta um `[]Produto` para JSON formatado e imprima.
2. Leia este JSON para uma struct e mostre a média das notas:
   `{"aluno":"Ana","notas":[8.5,9,7.5],"aprovado":true}`
3. Leia um JSON **desconhecido** para `map[string]any` e imprima cada chave com o tipo do valor (`%T`).
4. Crie um `Status` com `MarshalJSON` **e** `UnmarshalJSON` que use texto (`"ativo"`) em vez de número.
5. Leia um CSV de vendas (produto, quantidade, preço) e calcule o faturamento total e por produto.
6. Converta um arquivo CSV em JSON (e vice-versa).
7. **Projeto:** faça a livraria **salvar em `produtos.json`** ao sair e **carregar** ao iniciar, usando `Encoder`/`Decoder`.

---

⬅️ Anterior: [Arquivos e I/O](01-arquivos-e-io.md) · ➡️ Próximo: [Datas e horários](03-datas-e-horarios.md)

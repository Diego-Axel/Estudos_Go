# 01 — Arquivos, I/O e Linha de Comando

A biblioteca padrão de Go é famosa por ser **completa**: dá pra fazer muita coisa sem instalar nenhuma dependência. Vamos começar pelo básico de qualquer programa: **arquivos**.

---

## 📄 Lendo e escrevendo arquivos inteiros (o jeito fácil)

```go
import "os"

// Escrever (cria ou sobrescreve)
conteudo := []byte("Olá, arquivo!\nSegunda linha\n")
if err := os.WriteFile("notas.txt", conteudo, 0644); err != nil {
	log.Fatal(err)
}

// Ler tudo de uma vez
dados, err := os.ReadFile("notas.txt")
if err != nil {
	log.Fatal(err)
}
fmt.Print(string(dados))
```

> `0644` são as **permissões** (padrão Unix): o dono lê e escreve; os outros só leem. No Windows, isso é praticamente ignorado, mas o parâmetro é obrigatório.

✅ Ótimo para arquivos **pequenos** (configurações, JSONs). Para arquivos **grandes**, leia aos poucos (abaixo).

---

## 📂 Abrindo arquivos: `os.Open`, `os.Create`, `os.OpenFile`

```go
// Só leitura
f, err := os.Open("notas.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close() // ✅ sempre feche

// Criar/sobrescrever para escrita
f2, err := os.Create("saida.txt")

// Controle total: adicionar ao final (append), criando se não existir
f3, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
defer f3.Close()
fmt.Fprintln(f3, "nova linha de log") // *os.File é um io.Writer!
```

| Flag | Significado |
|---|---|
| `os.O_RDONLY` | só leitura |
| `os.O_WRONLY` | só escrita |
| `os.O_RDWR` | leitura e escrita |
| `os.O_CREATE` | cria se não existir |
| `os.O_APPEND` | escreve no final |
| `os.O_TRUNC` | apaga o conteúdo ao abrir |

---

## 📖 Lendo linha por linha: `bufio.Scanner` ⭐

```go
import "bufio"

f, err := os.Open("notas.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

scanner := bufio.NewScanner(f)
numero := 1
for scanner.Scan() {
	fmt.Printf("%3d: %s\n", numero, scanner.Text())
	numero++
}
if err := scanner.Err(); err != nil { // ✅ verifique o erro no final
	log.Fatal(err)
}
```

É o mesmo `bufio.Scanner` do seu `main.go`, só que lendo de um **arquivo** em vez de `os.Stdin`. (Lembra do `io.Reader`? 😉)

Outros modos de divisão:

```go
scanner.Split(bufio.ScanWords) // palavra por palavra
scanner.Split(bufio.ScanRunes) // caractere por caractere
```

> ⚠️ Por padrão, o Scanner aceita linhas de até 64 KB. Para linhas maiores: `scanner.Buffer(make([]byte, 1024*1024), 1024*1024)`.

---

## ✍️ Escrevendo com buffer: `bufio.Writer`

Escrever muitas linhas pequenas direto no arquivo é lento (cada escrita vai ao disco). O `bufio.Writer` junta tudo num buffer:

```go
f, err := os.Create("numeros.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

w := bufio.NewWriter(f)
for i := range 10000 {
	fmt.Fprintln(w, i)
}
if err := w.Flush(); err != nil { // ⚠️ OBRIGATÓRIO: descarrega o que ficou no buffer
	log.Fatal(err)
}
```

> ⚠️ Esquecer o `Flush()` é um bug clássico: as últimas linhas **somem**.

---

## 🔁 Copiando dados: `io.Copy`

```go
origem, _ := os.Open("foto.jpg")
defer origem.Close()

destino, _ := os.Create("copia.jpg")
defer destino.Close()

n, err := io.Copy(destino, origem) // Reader → Writer, em pedaços, sem carregar tudo na memória
fmt.Println(n, "bytes copiados", err)
```

Outros utilitários de `io`: `io.ReadAll(r)`, `io.MultiWriter(w1, w2)`, `io.TeeReader(r, w)`, `io.LimitReader(r, n)`.

---

## 🗂️ Pastas e caminhos

```go
import (
	"os"
	"path/filepath"
)

// Montar caminhos de forma portável (\ no Windows, / no Linux)
caminho := filepath.Join("dados", "2025", "vendas.csv") // dados\2025\vendas.csv no Windows

filepath.Base(caminho) // vendas.csv
filepath.Dir(caminho)  // dados\2025
filepath.Ext(caminho)  // .csv
abs, _ := filepath.Abs(".") // caminho absoluto

// Criar pastas (inclusive as intermediárias)
os.MkdirAll(filepath.Join("dados", "2025"), 0755)

// Verificar se existe
if _, err := os.Stat("config.json"); errors.Is(err, fs.ErrNotExist) {
	fmt.Println("config.json não existe")
}

// Renomear / mover e apagar
os.Rename("velho.txt", "novo.txt")
os.Remove("novo.txt")      // arquivo ou pasta vazia
os.RemoveAll("temporario") // pasta com tudo dentro ⚠️ cuidado!
```

> ⚠️ Use `path/filepath` para caminhos do **sistema de arquivos**. O pacote `path` (sem "file") é para URLs e caminhos com `/`.

### Listando uma pasta

```go
entradas, err := os.ReadDir(".")
if err != nil {
	log.Fatal(err)
}
for _, e := range entradas {
	tipo := "📄"
	if e.IsDir() {
		tipo = "📁"
	}
	fmt.Println(tipo, e.Name())
}
```

### Percorrendo subpastas: `filepath.WalkDir`

```go
err := filepath.WalkDir(".", func(caminho string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir // não entra na pasta .git
	}
	if filepath.Ext(caminho) == ".md" {
		fmt.Println(caminho)
	}
	return nil
})
```

> Esse código lista todos os `.md` deste repositório de estudos. Experimente!

---

## 💻 Linha de comando

### `os.Args`: argumentos crus

```go
// go run . ola mundo
fmt.Println(os.Args)      // [C:\...\exe\main.exe ola mundo]
fmt.Println(os.Args[1:])  // [ola mundo] (o [0] é o próprio programa)
```

### Pacote `flag`: opções com nome ⭐

```go
import "flag"

func main() {
	porta := flag.Int("porta", 8080, "porta do servidor")
	nome := flag.String("nome", "mundo", "nome para saudar")
	verbose := flag.Bool("v", false, "modo detalhado")
	flag.Parse()

	fmt.Printf("Olá, %s! Porta %d, verbose=%t\n", *nome, *porta, *verbose)
	fmt.Println("argumentos extras:", flag.Args())
}
```

```bash
go run . -nome=Diego -porta 3000 -v arquivo1.txt
# Olá, Diego! Porta 3000, verbose=true
# argumentos extras: [arquivo1.txt]

go run . -h   # ajuda gerada automaticamente!
```

> As funções `flag.Int`, `flag.String`... retornam **ponteiros** (módulo 08), por isso o `*porta`.

### Variáveis de ambiente

```go
home := os.Getenv("HOME") // "" se não existir

if token, ok := os.LookupEnv("API_TOKEN"); ok { // diferencia "vazia" de "não existe"
	fmt.Println("token configurado:", len(token) > 0)
}

os.Setenv("MODO", "teste")
```

### Entrada, saída e erro padrão

```go
fmt.Fprintln(os.Stdout, "saída normal")
fmt.Fprintln(os.Stderr, "mensagem de erro") // vai para a saída de erro
os.Exit(1) // código de saída (0 = sucesso)
```

---

## ✍️ Exercícios

1. Escreva um programa que conte **linhas, palavras e caracteres** de um arquivo (um mini `wc`).
2. Faça um programa que leia um arquivo e escreva outro com as linhas **numeradas**.
3. Crie um "diário": cada execução adiciona uma linha com a data e um texto em `diario.txt` (use `O_APPEND`).
4. Liste todos os arquivos `.go` e `.md` de uma pasta e suas subpastas, com o tamanho de cada um (`d.Info()`).
5. Crie uma ferramenta `copiar -de origem.txt -para destino.txt` usando `flag` e `io.Copy`.
6. **Projeto:** faça a livraria **salvar os produtos em arquivo** ao sair e **carregar** ao iniciar (use o formato que quiser; no próximo arquivo veremos JSON).

---

🏠 [Módulo 16](README.md) · ➡️ Próximo: [JSON e CSV](02-json-e-csv.md)

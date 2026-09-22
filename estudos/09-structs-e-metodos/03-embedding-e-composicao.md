# 03 — Embedding e Composição

## 🚫 Go não tem herança

Não existe `class Cachorro extends Animal`. Em vez disso, Go usa **composição**: montar tipos maiores **juntando** tipos menores.

> 🧠 *"Prefira composição a herança"* é um princípio clássico de design, e Go o levou ao pé da letra.

---

## 🧩 Composição "normal" (campo nomeado)

```go
type Motor struct {
	Potencia int
}

func (m Motor) Ligar() {
	fmt.Println("Motor ligado:", m.Potencia, "cv")
}

type Carro struct {
	Modelo string
	Motor  Motor // campo com nome
}

func main() {
	c := Carro{Modelo: "Fusca", Motor: Motor{Potencia: 65}}
	c.Motor.Ligar()              // precisa passar por c.Motor
	fmt.Println(c.Motor.Potencia)
}
```

---

## 🪆 Embedding (campo embutido) ⭐

Se você colocar **só o tipo**, sem nome de campo, ele é **embutido**. Os campos e métodos dele são **promovidos**, ou seja, ficam acessíveis direto:

```go
type Carro struct {
	Modelo string
	Motor         // embutido: sem nome de campo
}

func main() {
	c := Carro{Modelo: "Fusca", Motor: Motor{Potencia: 65}}

	c.Ligar()               // ✅ método promovido
	fmt.Println(c.Potencia) // ✅ campo promovido

	c.Motor.Ligar()         // ✅ o caminho completo continua funcionando
}
```

> O nome do campo embutido é o **nome do tipo** (`Motor`), por isso `c.Motor` funciona.

### Exemplo mais realista

```go
type Pessoa struct {
	Nome  string
	Email string
}

func (p Pessoa) Contato() string {
	return fmt.Sprintf("%s <%s>", p.Nome, p.Email)
}

type Funcionario struct {
	Pessoa         // embutido
	Cargo   string
	Salario float64
}

type Cliente struct {
	Pessoa         // o mesmo tipo reaproveitado
	Pontos int
}

func main() {
	f := Funcionario{
		Pessoa: Pessoa{Nome: "Ana", Email: "ana@empresa.com"},
		Cargo:  "Dev",
	}
	fmt.Println(f.Nome)      // Ana
	fmt.Println(f.Contato()) // Ana <ana@empresa.com>
}
```

---

## 🎭 "Sobrescrevendo" (*shadowing*)

Se o tipo externo tiver um campo/método com o **mesmo nome**, o **externo vence**:

```go
type Animal struct{ Nome string }

func (a Animal) Falar() string { return "..." }

type Cachorro struct {
	Animal
	Raca string
}

func (c Cachorro) Falar() string { return "Au au!" } // "sobrescreve"

func main() {
	d := Cachorro{Animal: Animal{Nome: "Rex"}, Raca: "Vira-lata"}
	fmt.Println(d.Falar())        // Au au!
	fmt.Println(d.Animal.Falar()) // ... (o original continua acessível)
}
```

---

## ⚠️ Embedding **não** é herança

Duas diferenças importantes:

### 1. `Cachorro` **não é** um `Animal`

```go
func apresentar(a Animal) {
	fmt.Println(a.Nome)
}

d := Cachorro{Animal: Animal{Nome: "Rex"}}
// apresentar(d)      // ❌ erro: cannot use d (...Cachorro) as Animal value
apresentar(d.Animal)  // ✅ passa a parte Animal
```

> Para ter "polimorfismo" (uma função que aceita vários tipos), Go usa **interfaces**, no próximo módulo.

### 2. Não há "despacho virtual"

Um método de `Animal` **não enxerga** os métodos "sobrescritos" em `Cachorro`:

```go
func (a Animal) Apresentar() string {
	return a.Nome + " diz: " + a.Falar() // chama SEMPRE Animal.Falar
}

d := Cachorro{Animal: Animal{Nome: "Rex"}}
fmt.Println(d.Apresentar()) // Rex diz: ...   (e não "Au au!")
```

Em Java, o resultado seria "Au au!". Em Go, o `Animal` embutido **não sabe** que está dentro de um `Cachorro`.

---

## 💥 Conflitos de nomes

Se dois tipos embutidos **no mesmo nível** tiverem o mesmo campo/método, o acesso direto fica **ambíguo**:

```go
type A struct{ ID int }
type B struct{ ID int }

type C struct {
	A
	B
}

c := C{}
// fmt.Println(c.ID) // ❌ erro: ambiguous selector c.ID
fmt.Println(c.A.ID, c.B.ID) // ✅ seja explícito
```

> O erro só acontece se você **usar** o nome ambíguo.

---

## 🔗 Embutindo ponteiros

Também dá pra embutir um **ponteiro**:

```go
type Logger struct{ Prefixo string }

func (l *Logger) Log(msg string) {
	fmt.Println(l.Prefixo, msg)
}

type Servico struct {
	*Logger // compartilha o MESMO logger entre vários serviços
	Nome string
}

func main() {
	log := &Logger{Prefixo: "[APP]"}
	s1 := Servico{Logger: log, Nome: "pagamentos"}
	s2 := Servico{Logger: log, Nome: "estoque"}

	s1.Log("iniciado") // [APP] iniciado
	s2.Log("iniciado") // [APP] iniciado
}
```

> ⚠️ Se o ponteiro embutido for `nil`, chamar um método que acessa campos dele causa **panic**.

---

## 🧱 Exemplo: embutindo `sync.Mutex` (spoiler de concorrência)

Um padrão muito comum na biblioteca padrão e em projetos reais:

```go
type ContadorSeguro struct {
	sync.Mutex
	valor int
}

func (c *ContadorSeguro) Incrementar() {
	c.Lock()         // método promovido do Mutex
	defer c.Unlock()
	c.valor++
}
```

---

## ✍️ Exercícios

1. Crie `Endereco` e embuta em `Cliente` e `Fornecedor`. Acesse `Cidade` diretamente nos dois.
2. Crie `Veiculo{Marca string; Ano int}` com o método `Descricao()`, e embuta em `Carro` e `Moto`. Faça `Moto` "sobrescrever" `Descricao()`.
3. Mostre, com código, que um método do tipo embutido **não** chama a versão "sobrescrita" do tipo externo.
4. Crie dois tipos embutidos com um campo de mesmo nome e veja o erro de ambiguidade. Resolva.
5. Crie um `Logger` compartilhado (embutido como ponteiro) entre três structs diferentes.
6. Por que dizemos que embedding é **composição**, e não herança? Cite duas diferenças práticas.

---

⬅️ Anterior: [Métodos](02-metodos.md) · ➡️ Próximo: [Tags e boas práticas](04-tags-e-boas-praticas.md)

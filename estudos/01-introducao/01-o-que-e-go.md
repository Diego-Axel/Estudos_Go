# 01 — O que é Go?

## 📜 Um pouco de história

**Go** (ou **Golang**) foi criada no **Google** em **2007** por três pesos-pesados da computação:

- **Robert Griesemer** (trabalhou na V8, do JavaScript)
- **Rob Pike** (Unix, UTF-8)
- **Ken Thompson** (criador do Unix e da linguagem B, antecessora do C)

Foi lançada publicamente em **2009** e a versão **1.0** saiu em **2012**.

### Por que criaram Go?

O Google tinha sistemas gigantes em C++ e Java que:

- demoravam **muito** para compilar;
- eram complexos demais para ler e manter;
- tinham dificuldade para aproveitar **processadores com vários núcleos**.

Go nasceu para resolver isso: **simples, rápida de compilar e feita para concorrência**.

---

## ⭐ Características principais

| Característica | O que significa |
|---|---|
| **Compilada** | O código vira um executável nativo (binário). Não precisa de máquina virtual. |
| **Estaticamente tipada** | Os tipos são verificados na compilação — muitos erros aparecem antes de rodar. |
| **Simples** | Só **25 palavras-chave**. A linguagem cabe na cabeça. |
| **Garbage Collector** | A memória é liberada automaticamente, sem `free()` manual. |
| **Concorrência nativa** | *Goroutines* e *channels* fazem parte da linguagem. |
| **Compilação rápida** | Projetos grandes compilam em segundos. |
| **Binário único** | Gera um único arquivo executável, fácil de distribuir. |
| **Multiplataforma** | Compila para Windows, Linux, macOS, ARM... a partir de qualquer um deles. |
| **Ferramentas embutidas** | `go fmt`, `go test`, `go vet`, `go mod` já vêm junto. |

---

## 🚫 O que Go **não** tem (de propósito)

Go é minimalista por escolha. Não tem:

- **Classes e herança** → usa *structs* + *composição* + *interfaces*
- **Exceções (`try/catch`)** → erros são **valores** retornados pelas funções
- **Sobrecarga de funções/operadores**
- **Operador ternário** (`a ? b : c`) → usa `if/else`
- **`while`** → só existe `for` (que faz o papel de todos os laços)

> 🧠 A filosofia é: *"deve existir uma forma clara e óbvia de fazer as coisas"*.

---

## 🔑 As 25 palavras-chave

```
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

Você vai aprender todas elas ao longo dos módulos.

---

## 🏢 Onde Go é usado?

- **Docker** e **Kubernetes** — escritos em Go
- **Terraform**, **Prometheus**, **Grafana**, **Hugo**
- Back-ends de empresas como **Google, Uber, Twitch, Dropbox, Mercado Livre, Nubank (partes), PicPay**
- **APIs**, **microsserviços**, **CLIs**, **ferramentas de DevOps/infra**

---

## 🆚 Go comparado com outras linguagens

| | Go | Python | Java | C |
|---|---|---|---|---|
| Tipagem | Estática | Dinâmica | Estática | Estática |
| Execução | Compilada (nativa) | Interpretada | JVM | Compilada (nativa) |
| Memória | GC | GC | GC | Manual |
| Concorrência | Goroutines (leves) | Threads/async (GIL) | Threads | Threads (manual) |
| Verbosidade | Baixa/média | Baixa | Alta | Média |

---

## ✍️ Exercícios

1. Quem são os três criadores de Go e em que ano ela foi lançada publicamente?
2. Cite 3 problemas que Go foi criada para resolver.
3. Go tem `while`? E `try/catch`? Como ela resolve essas situações?
4. Cite 3 ferramentas famosas escritas em Go.

---

➡️ Próximo: [Instalação e ambiente](02-instalacao-e-ambiente.md)

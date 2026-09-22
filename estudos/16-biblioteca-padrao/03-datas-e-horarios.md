# 03 — Datas e Horários com `time`

## 🕐 O instante atual

```go
import "time"

agora := time.Now()
fmt.Println(agora) // 2025-09-22 15:04:05.123456 -0300 -03 m=+0.000012

fmt.Println(agora.Year(), agora.Month(), agora.Day()) // 2025 September 22
fmt.Println(agora.Hour(), agora.Minute(), agora.Second())
fmt.Println(agora.Weekday())                          // Monday
fmt.Println(agora.YearDay())                          // dia do ano (1 a 366)
```

### Criando uma data específica

```go
natal := time.Date(2025, time.December, 25, 20, 0, 0, 0, time.Local)
fmt.Println(natal) // 2025-12-25 20:00:00 -0300 -03
```

---

## 🎨 Formatando: o layout mágico de Go ⭐

Outras linguagens usam `%Y-%m-%d` ou `yyyy-MM-dd`. Go usa uma **data de referência**:

```
Mon Jan 2 15:04:05 MST 2006
```

Decore pela sequência: **1 2 3 4 5 6 7**

```
01/02 03:04:05PM '06 -0700
 │  │  │  │  │     │    └── 7: fuso horário (-0700)
 │  │  │  │  │     └─────── 6: ano (2006)
 │  │  │  │  └───────────── 5: segundo (05)
 │  │  │  └──────────────── 4: minuto (04)
 │  │  └─────────────────── 3: hora (03 ou 15 no formato 24h)
 │  └────────────────────── 2: dia (02)
 └───────────────────────── 1: mês (01)
```

Você **escreve a data de referência no formato que quer**:

```go
t := time.Date(2025, 3, 7, 14, 5, 9, 0, time.Local)

fmt.Println(t.Format("02/01/2006"))          // 07/03/2025  ← padrão brasileiro
fmt.Println(t.Format("02/01/2006 15:04"))    // 07/03/2025 14:05
fmt.Println(t.Format("2006-01-02"))          // 2025-03-07  ← ISO
fmt.Println(t.Format("15:04:05"))            // 14:05:09
fmt.Println(t.Format("03:04 PM"))            // 02:05 PM
fmt.Println(t.Format("Mon, 02 Jan 2006"))    // Fri, 07 Mar 2025
fmt.Println(t.Format("January 2, 2006"))     // March 7, 2025
```

| Elemento | Código |
|---|---|
| Ano (4 / 2 dígitos) | `2006` / `06` |
| Mês (número / com zero / nome / abreviado) | `1` / `01` / `January` / `Jan` |
| Dia (sem zero / com zero) | `2` / `02` |
| Dia da semana | `Monday` / `Mon` |
| Hora 24h / 12h | `15` / `03` ou `3` |
| Minuto / segundo | `04` / `05` |
| AM/PM | `PM` |
| Milissegundos | `.000` |
| Fuso | `-0700` / `MST` / `Z07:00` |

### Layouts prontos

```go
time.RFC3339    // "2006-01-02T15:04:05Z07:00" ← padrão de APIs/JSON
time.DateTime   // "2006-01-02 15:04:05"  (Go 1.20+)
time.DateOnly   // "2006-01-02"           (Go 1.20+)
time.TimeOnly   // "15:04:05"             (Go 1.20+)
time.Kitchen    // "3:04PM"

fmt.Println(t.Format(time.RFC3339)) // 2025-03-07T14:05:09-03:00
```

> ⚠️ Nomes de meses e dias saem **sempre em inglês**. Para português, use um `map`/slice com os nomes (exercício abaixo).

---

## 🔍 Convertendo texto em data: `time.Parse`

Mesmo layout, caminho inverso:

```go
t, err := time.Parse("02/01/2006", "25/12/2025")
if err != nil {
	log.Fatal(err)
}
fmt.Println(t) // 2025-12-25 00:00:00 +0000 UTC

_, err = time.Parse("02/01/2006", "31/02/2025")
fmt.Println(err) // parsing time "31/02/2025": day out of range
```

> ⚠️ `time.Parse` assume **UTC** se o texto não tiver fuso. Para interpretar no horário local: `time.ParseInLocation(layout, texto, time.Local)`.

---

## ⏳ Durações: `time.Duration`

```go
d := 90 * time.Minute
fmt.Println(d)           // 1h30m0s
fmt.Println(d.Hours())   // 1.5
fmt.Println(d.Minutes()) // 90

d2 := 2*time.Hour + 15*time.Second
fmt.Println(d2)          // 2h0m15s

// De texto
d3, _ := time.ParseDuration("1h45m")
fmt.Println(d3.Minutes()) // 105
```

Constantes: `time.Nanosecond`, `Microsecond`, `Millisecond`, `Second`, `Minute`, `Hour`. (Não existe `time.Day`: dias têm duração variável por causa do horário de verão.)

> ⚠️ Multiplicar por variável `int` precisa de conversão: `time.Duration(n) * time.Second`.

---

## ➕ Aritmética com datas

```go
agora := time.Now()

amanha := agora.Add(24 * time.Hour)
semanaPassada := agora.Add(-7 * 24 * time.Hour)
proximoMes := agora.AddDate(0, 1, 0)   // anos, meses, dias
daquiUmAno := agora.AddDate(1, 0, 0)

// Diferença entre datas
natal := time.Date(2025, 12, 25, 0, 0, 0, 0, time.Local)
faltam := natal.Sub(agora) // Duration
fmt.Printf("faltam %.0f dias para o Natal\n", faltam.Hours()/24)

// Tempo decorrido (muito usado para medir desempenho)
inicio := time.Now()
processar()
fmt.Println("levou", time.Since(inicio))

// Tempo até uma data
fmt.Println(time.Until(natal).Round(time.Hour))
```

### Comparando

```go
a := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
b := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

fmt.Println(a.Before(b)) // true
fmt.Println(a.After(b))  // false
fmt.Println(a.Equal(b))  // false
fmt.Println(a.Compare(b)) // -1 (Go 1.20+)
```

> ⚠️ Use **`Equal`**, não `==`, para comparar datas. O `==` também compara o fuso e detalhes internos, e pode dar `false` para o mesmo instante.

### Arredondando e truncando

```go
t := time.Date(2025, 3, 7, 14, 37, 50, 0, time.UTC)
fmt.Println(t.Truncate(time.Hour).Format("15:04")) // 14:00
fmt.Println(t.Round(time.Hour).Format("15:04"))    // 15:00

// Início do dia
inicioDoDia := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
```

---

## 🌎 Fusos horários

```go
sp, err := time.LoadLocation("America/Sao_Paulo")
if err != nil {
	log.Fatal(err)
}
ny, _ := time.LoadLocation("America/New_York")

agora := time.Now()
fmt.Println("São Paulo:", agora.In(sp).Format("15:04"))
fmt.Println("Nova York:", agora.In(ny).Format("15:04"))
fmt.Println("UTC:      ", agora.UTC().Format("15:04"))
```

> ⚠️ No **Windows**, o `LoadLocation` pode falhar por falta do banco de fusos. Solução: adicione `import _ "time/tzdata"` no programa (embute os fusos no executável, ~450 KB).

> ✅ **Boa prática:** guarde e transmita datas em **UTC** (bancos, APIs). Converta para o fuso local **só na hora de mostrar** ao usuário.

---

## 🔢 Timestamps Unix

```go
agora := time.Now()
fmt.Println(agora.Unix())      // segundos desde 01/01/1970 UTC
fmt.Println(agora.UnixMilli()) // milissegundos

t := time.Unix(1700000000, 0)
fmt.Println(t.UTC()) // 2023-11-14 22:13:20 +0000 UTC
```

---

## ⏲️ Esperar e agendar

```go
time.Sleep(500 * time.Millisecond) // pausa a goroutine atual

timer := time.AfterFunc(2*time.Second, func() {
	fmt.Println("executado depois de 2 segundos")
})
// timer.Stop() cancela

// Tickers, time.After e timeouts: veja o módulo 14
```

---

## ✍️ Exercícios

1. Mostre a data atual no formato `segunda-feira, 22 de setembro de 2025`, em português (crie slices com os nomes dos dias e meses).
2. Leia a data de nascimento do usuário (`dd/mm/aaaa`) e calcule a **idade exata** e quantos **dias** faltam para o próximo aniversário.
3. Escreva `diasUteis(inicio, fim time.Time) int` que conte os dias de segunda a sexta entre duas datas.
4. Mostre o horário atual em São Paulo, Lisboa, Tóquio e Nova York.
5. Crie um cronômetro: mostre o tempo decorrido a cada segundo, por 5 segundos, formatado como `00:00:05`.
6. Converta `"2025-03-07T14:05:09Z"` (RFC 3339) para o horário de Brasília e mostre em `dd/mm/aaaa HH:MM`.
7. **Projeto:** adicione ao `Produto` da livraria um campo `CadastradoEm time.Time`, preencha no cadastro e mostre formatado na listagem.

---

⬅️ Anterior: [JSON e CSV](02-json-e-csv.md) · ➡️ Próximo: [Logs, regex e outros pacotes](04-logs-regex-e-outros-pacotes.md)

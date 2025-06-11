package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Produto representa um livro da livraria
type Produto struct {
	Codigo  string
	Titulo  string
	Autor   string
	Preco   float64
	Estoque int
}

// Livraria contém os produtos cadastrados
type Livraria struct {
	Produtos map[string]Produto
}

func main() {
	livraria := Livraria{Produtos: make(map[string]Produto)}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n📚 Sistema Livraria")
		fmt.Println("1 - Cadastrar Produto")
		fmt.Println("2 - Listar Produtos")
		fmt.Println("3 - Buscar Produto por Código")
		fmt.Println("0 - Sair")
		fmt.Print("Escolha uma opção: ")

		scanner.Scan()
		opcao := scanner.Text()

		switch opcao {
		case "1":
			cadastrarProduto(scanner, &livraria)
		case "2":
			listarProdutos(livraria)
		case "3":
			buscarProduto(scanner, livraria)
		case "0":
			fmt.Println("👋 Encerrando o sistema.")
			return
		default:
			fmt.Println("❌ Opção inválida!")
		}
	}
}

func cadastrarProduto(scanner *bufio.Scanner, livraria *Livraria) {
	fmt.Println("\n--- Cadastro de Produto ---")

	fmt.Print("Código: ")
	scanner.Scan()
	codigo := strings.TrimSpace(scanner.Text())

	if _, existe := livraria.Produtos[codigo]; existe {
		fmt.Println("⚠️ Código já cadastrado.")
		return
	}

	fmt.Print("Título: ")
	scanner.Scan()
	titulo := strings.TrimSpace(scanner.Text())

	fmt.Print("Autor: ")
	scanner.Scan()
	autor := strings.TrimSpace(scanner.Text())

	fmt.Print("Preço (ex: 49.90): ")
	scanner.Scan()
	precoStr := strings.TrimSpace(scanner.Text())
	preco, err := strconv.ParseFloat(precoStr, 64)
	if err != nil {
		fmt.Println("Preço inválido.")
		return
	}

	fmt.Print("Quantidade em estoque: ")
	scanner.Scan()
	estoqueStr := strings.TrimSpace(scanner.Text())
	estoque, err := strconv.Atoi(estoqueStr)
	if err != nil {
		fmt.Println("Estoque inválido.")
		return
	}

	produto := Produto{
		Codigo:  codigo,
		Titulo:  titulo,
		Autor:   autor,
		Preco:   preco,
		Estoque: estoque,
	}

	livraria.Produtos[codigo] = produto
	fmt.Println("✅ Produto cadastrado com sucesso!")
}

func listarProdutos(livraria Livraria) {
	if len(livraria.Produtos) == 0 {
		fmt.Println("\n📭 Nenhum produto cadastrado.")
		return
	}

	fmt.Println("\n📚 Produtos Cadastrados:")
	for _, p := range livraria.Produtos {
		fmt.Printf("\nCódigo: %s\nTítulo: %s\nAutor: %s\nPreço: R$ %.2f\nEstoque: %d unidades\n",
			p.Codigo, p.Titulo, p.Autor, p.Preco, p.Estoque)
		fmt.Println("-------------------------------")
	}
}

func buscarProduto(scanner *bufio.Scanner, livraria Livraria) {
	fmt.Print("\n🔍 Digite o código do produto: ")
	scanner.Scan()
	codigo := strings.TrimSpace(scanner.Text())

	if produto, existe := livraria.Produtos[codigo]; existe {
		fmt.Printf("\n🔎 Produto Encontrado:\nCódigo: %s\nTítulo: %s\nAutor: %s\nPreço: R$ %.2f\nEstoque: %d unidades\n",
			produto.Codigo, produto.Titulo, produto.Autor, produto.Preco, produto.Estoque)
	} else {
		fmt.Println("❌ Produto não encontrado.")
	}
}




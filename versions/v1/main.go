// V1 processa exatamente um arquivo CSV de forma sequencial.
package main

import (
	"fmt"
	"os"

	"job-processing-system/internal/transactions"
)

// main é o ponto de entrada de um programa executável em Go.
func main() {
	// os.Args[0] é o programa e os.Args[1] deve ser o caminho do CSV.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run ./versions/v1 <file.csv>\n")
		os.Exit(2)
	}
	// Process devolve o resultado calculado e um possível erro.
	result, err := transactions.Process(os.Args[1])
	if err != nil {
		// Stderr é a saída reservada para erros; Exit(1) indica falha.
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	// Stdout representa a saída normal do terminal.
	transactions.PrintResult(os.Stdout, result)
}

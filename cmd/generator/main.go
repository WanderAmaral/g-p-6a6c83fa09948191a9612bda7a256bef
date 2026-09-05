package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// categories é um slice com as categorias possíveis.
var categories = []string{"books", "electronics", "food", "health", "transport"}

// main interpreta os argumentos do terminal e chama o gerador.
func main() {
	// As funções flag criam opções e retornam ponteiros, acessados com *.
	records := flag.Int("records", 10_000, "number of transactions")
	output := flag.String("output", "data/transactions-10k.csv", "output CSV path")
	seed := flag.Int64("seed", 42, "random seed (same seed produces the same data)")
	flag.Parse() // interpreta as opções escritas no comando

	// Não permitimos quantidade vazia ou negativa.
	if *records < 1 {
		fmt.Fprintln(os.Stderr, "records must be greater than zero")
		os.Exit(1)
	}
	if err := generate(*output, *records, *seed); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("generated %d records in %s\n", *records, *output)
}

// generate cria count transações no caminho informado.
func generate(path string, count int, seed int64) error {
	// Cria a pasta de destino caso ela ainda não exista.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Create cria o arquivo ou substitui outro de mesmo nome.
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close() // fecha o arquivo quando a função terminar

	// O buffer reúne pequenas escritas antes de enviá-las ao disco.
	buffer := bufio.NewWriterSize(file, 256*1024)
	writer := csv.NewWriter(buffer)
	// A primeira linha é o cabeçalho do CSV.
	if err := writer.Write([]string{"id", "timestamp", "category", "amount"}); err != nil {
		return err
	}

	// A mesma seed sempre produz a mesma sequência de dados.
	rng := rand.New(rand.NewSource(seed))
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	// Cada repetição cria uma transação.
	for i := 1; i <= count; i++ {
		amount := 100 + rng.Intn(99_900) // cents: 1.00 through 999.99
		// record contém as quatro colunas da nova linha.
		record := []string{
			strconv.Itoa(i),
			base.Add(time.Duration(i) * time.Second).Format(time.RFC3339),
			categories[rng.Intn(len(categories))],
			fmt.Sprintf("%d.%02d", amount/100, amount%100),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	// Flush envia os dados restantes dos buffers para o arquivo.
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return buffer.Flush()
}

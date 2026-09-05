package transactions

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// Result é uma struct que agrupa os dados produzidos por um job.
type Result struct {
	File       string
	Records    int64
	Total      float64
	ByCategory map[string]int64
	Duration   time.Duration
}

// Average calcula o valor médio das transações.
func (r Result) Average() float64 {
	// Evita divisão por zero quando o CSV não contém transações.
	if r.Records == 0 {
		return 0
	}
	return r.Total / float64(r.Records)
}

// RecordsPerSecond calcula quantos registros foram processados por segundo.
func (r Result) RecordsPerSecond() float64 {
	if r.Duration <= 0 {
		return 0
	}
	return float64(r.Records) / r.Duration.Seconds()
}

// Process lê o CSV linha por linha, sem carregar o arquivo inteiro na memória.
// Ela retorna dois valores: o Result e um possível error.
func Process(path string) (Result, error) {
	// Guardamos o instante inicial para calcular a duração depois.
	started := time.Now()
	// make inicializa o map que contará as categorias.
	result := Result{File: path, ByCategory: make(map[string]int64)}

	// Open abre o arquivo somente para leitura.
	file, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("open %s: %w", path, err)
	}
	// defer garante o fechamento do arquivo quando Process terminar.
	defer file.Close()

	// csv.Reader interpreta as colunas e regras do formato CSV.
	reader := csv.NewReader(file)
	// A primeira leitura corresponde ao cabeçalho.
	header, err := reader.Read()
	if err != nil {
		return result, fmt.Errorf("read header: %w", err)
	}
	if len(header) != 4 || header[0] != "id" || header[1] != "timestamp" || header[2] != "category" || header[3] != "amount" {
		return result, errors.New("invalid header: expected id,timestamp,category,amount")
	}

	// Começamos na linha 2 porque a linha 1 foi o cabeçalho.
	// O laço termina quando Read retorna io.EOF: fim do arquivo.
	for line := 2; ; line++ {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return result, fmt.Errorf("read line %d: %w", line, err)
		}
		if len(record) != 4 {
			return result, fmt.Errorf("line %d: expected 4 columns", line)
		}

		// O CSV contém texto; ParseFloat converte amount em número.
		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return result, fmt.Errorf("line %d: invalid amount %q: %w", line, record[3], err)
		}
		// Este é o núcleo da agregação.
		result.Records++               // conta o registro
		result.Total += amount         // acumula seu valor
		result.ByCategory[record[2]]++ // conta sua categoria
	}

	// Since calcula o tempo decorrido desde started.
	result.Duration = time.Since(started)
	// nil significa que nenhum erro ocorreu.
	return result, nil
}

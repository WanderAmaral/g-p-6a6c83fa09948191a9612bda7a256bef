// V4 usa um channel como fila de jobs em memória.
package main

import (
	"fmt"
	"os"

	"job-processing-system/internal/jobs"
	"job-processing-system/internal/transactions"
)

type outcome struct {
	result transactions.Result
	err    error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./versions/v4 <file.csv> [file.csv...]")
		os.Exit(2)
	}

	paths := os.Args[1:]
	jobQueue := make(chan string)
	results := make(chan outcome)te

	// O produtor coloca caminhos na fila e a fecha quando termina.
	go func() {
		for _, path := range paths {
			jobQueue <- path
		}
		close(jobQueue)
	}()

	// Existe um consumidor, então os jobs ainda rodam um por vez.
	go func() {
		for path := range jobQueue {
			result, err := transactions.Process(path)
			results <- outcome{result, err}
		}
		close(results)
	}()

	for item := range results {
		if !summary.Add(item.result, item.err) {
			fmt.Fprintln(os.Stderr, "job failed:", item.err)
			continue
		}
		transactions.PrintResult(os.Stdout, item.result)
	}

	summary.Print(os.Stdout, "channel queue with one consumer")
	os.Exit(summary.ExitCode())
}

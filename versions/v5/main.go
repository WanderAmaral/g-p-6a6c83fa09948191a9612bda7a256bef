// V5 processa a fila com uma quantidade configurável de workers.
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"job-processing-system/internal/jobs"
	"job-processing-system/internal/transactions"
)

type outcome struct {
	workerID int
	result   transactions.Result
	err      error
}

func main() {
	workers := flag.Int("workers", 4, "number of concurrent workers")
	flag.Parse()
	paths := flag.Args()

	if *workers < 1 || len(paths) < 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./versions/v5 -workers 4 <file.csv> [file.csv...]")
		os.Exit(2)
	}

	jobQueue := make(chan string)
	results := make(chan outcome)
	summary := jobs.NewSummary(len(paths))
	var wg sync.WaitGroup

	// Todos os workers retiram trabalho da mesma fila.
	for workerID := 1; workerID <= *workers; workerID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobQueue {
				result, err := transactions.Process(path)
				results <- outcome{workerID, result, err}
			}
		}()
	}

	go func() {
		for _, path := range paths {
			jobQueue <- path
		}
		close(jobQueue)
	}()

	// Fechamos results somente depois que todos os workers terminarem.
	go func() {
		wg.Wait()
		close(results)
	}()

	for item := range results {
		if !summary.Add(item.result, item.err) {
			fmt.Fprintf(os.Stderr, "worker %d: job failed: %v\n", item.workerID, item.err)
			continue
		}
		fmt.Printf("\nWorker %d completed a job\n", item.workerID)
		transactions.PrintResult(os.Stdout, item.result)
	}

	summary.Print(os.Stdout, fmt.Sprintf("worker pool (%d workers)", *workers))
	os.Exit(summary.ExitCode())
}

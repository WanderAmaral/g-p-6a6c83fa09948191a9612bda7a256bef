// V3 executa uma goroutine para cada arquivo recebido.
package main

import (
	"fmt"
	"os"
	"sync"

	"job-processing-system/internal/jobs"
	"job-processing-system/internal/transactions"
)

type outcome struct {
	result transactions.Result
	err    error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./versions/v3 <file.csv> [file.csv...]")
		os.Exit(2)
	}

	paths := os.Args[1:]
	summary := jobs.NewSummary(len(paths))
	outcomes := make([]outcome, len(paths))
	var wg sync.WaitGroup

	// A palavra go inicia cada Process sem esperar o anterior terminar.
	for index, path := range paths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			outcomes[index].result, outcomes[index].err = transactions.Process(path)
		}()
	}
	wg.Wait()

	// Imprimimos depois para manter a ordem original dos arquivos.
	for _, item := range outcomes {
		if !summary.Add(item.result, item.err) {
			fmt.Fprintln(os.Stderr, "job failed:", item.err)
			continue
		}
		transactions.PrintResult(os.Stdout, item.result)
	}

	summary.Print(os.Stdout, "one goroutine per job")
	os.Exit(summary.ExitCode())
}

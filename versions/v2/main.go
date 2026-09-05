// V2 processa vários CSVs sequencialmente, na ordem recebida.
package main

import (
	"fmt"
	"io"
	"os"

	"job-processing-system/internal/jobs"
	"job-processing-system/internal/transactions"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./versions/v2 <file.csv> [file.csv...]")
		os.Exit(2)
	}
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run usa um for comum: um job termina antes do próximo começar.
func run(paths []string, stdout, stderr io.Writer) int {
	summary := jobs.NewSummary(len(paths))

	for number, path := range paths {
		fmt.Fprintf(stdout, "\nStarting job %d/%d: %s\n", number+1, len(paths), path)
		result, err := transactions.Process(path)

		if !summary.Add(result, err) {
			fmt.Fprintf(stderr, "job failed: %v\n", err)
			continue
		}
		transactions.PrintResult(stdout, result)
	}

	summary.Print(stdout, "sequential")
	return summary.ExitCode()
}

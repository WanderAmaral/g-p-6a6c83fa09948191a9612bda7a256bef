package transactions

import (
	"fmt"
	"io"
	"sort"
)

// PrintResult escreve o resultado em w, que pode ser terminal, arquivo ou buffer.
func PrintResult(w io.Writer, result Result) {
	// Fprintf substitui os marcadores (%s, %d, %f) pelos valores.
	fmt.Fprintf(w, "\nJob: %s\n", result.File)
	fmt.Fprintf(w, "  records: %d\n", result.Records)
	fmt.Fprintf(w, "  total: %.2f\n", result.Total)
	fmt.Fprintf(w, "  average: %.2f\n", result.Average())
	fmt.Fprintf(w, "  duration: %s\n", result.Duration)
	fmt.Fprintf(w, "  throughput: %.0f records/s\n", result.RecordsPerSecond())

	// Maps não possuem ordem garantida; copiamos as chaves para um slice.
	categories := make([]string, 0, len(result.ByCategory))
	for category := range result.ByCategory {
		// append adiciona um item ao final do slice.
		categories = append(categories, category)
	}
	// A ordenação mantém a saída previsível.
	sort.Strings(categories)
	for _, category := range categories {
		fmt.Fprintf(w, "  category[%s]: %d\n", category, result.ByCategory[category])
	}
}

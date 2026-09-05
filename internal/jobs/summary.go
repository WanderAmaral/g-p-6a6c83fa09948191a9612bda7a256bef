// Package jobs contém estruturas compartilhadas pelas versões com vários jobs.
package jobs

import (
	"fmt"
	"io"
	"time"

	"job-processing-system/internal/transactions"
)

// Summary guarda as métricas de um lote inteiro.
type Summary struct {
	Jobs, Succeeded, Failed int
	Records                 int64
	Duration                time.Duration
	started                 time.Time
}

// NewSummary inicia o cronômetro e informa quantos jobs existem no lote.
func NewSummary(total int) *Summary {
	return &Summary{Jobs: total, started: time.Now()}
}

// Add contabiliza um job concluído. Retorna false quando ele falhou.
func (s *Summary) Add(result transactions.Result, err error) bool {
	if err != nil {
		s.Failed++
		return false
	}
	s.Succeeded++
	s.Records += result.Records
	return true
}

// Print encerra o cronômetro e mostra o resumo.
func (s *Summary) Print(w io.Writer, mode string) {
	s.Duration = time.Since(s.started)
	throughput := float64(0)
	if s.Duration > 0 {
		throughput = float64(s.Records) / s.Duration.Seconds()
	}

	fmt.Fprintln(w, "\nBatch summary")
	fmt.Fprintf(w, "  mode: %s\n", mode)
	fmt.Fprintf(w, "  jobs: %d\n", s.Jobs)
	fmt.Fprintf(w, "  succeeded: %d\n", s.Succeeded)
	fmt.Fprintf(w, "  failed: %d\n", s.Failed)
	fmt.Fprintf(w, "  records: %d\n", s.Records)
	fmt.Fprintf(w, "  duration: %s\n", s.Duration)
	fmt.Fprintf(w, "  throughput: %.0f records/s\n", throughput)
}

// ExitCode segue a convenção: zero é sucesso, diferente de zero é falha.
func (s Summary) ExitCode() int {
	if s.Failed > 0 {
		return 1
	}
	return 0
}

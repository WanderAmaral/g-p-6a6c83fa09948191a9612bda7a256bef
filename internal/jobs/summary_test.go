package jobs

import (
	"errors"
	"testing"

	"job-processing-system/internal/transactions"
)

func TestSummaryCountsSuccessAndFailure(t *testing.T) {
	summary := NewSummary(2)
	summary.Add(transactions.Result{Records: 10}, nil)
	summary.Add(transactions.Result{}, errors.New("test error"))

	if summary.Succeeded != 1 || summary.Failed != 1 || summary.Records != 10 || summary.ExitCode() != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

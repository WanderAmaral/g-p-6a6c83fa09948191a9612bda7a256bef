package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunContinuesAfterFailedJob confirma que um erro não cancela o lote.
func TestRunContinuesAfterFailedJob(t *testing.T) {
	dir := t.TempDir()
	validFile := filepath.Join(dir, "valid.csv")
	missingFile := filepath.Join(dir, "missing.csv")
	data := "id,timestamp,category,amount\n1,2026-01-01T00:00:00Z,food,10.50\n"
	if err := os.WriteFile(validFile, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	// Buffers capturam as saídas sem imprimi-las no terminal.
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{missingFile, validFile}, &stdout, &stderr)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	for _, expected := range []string{"jobs: 2", "succeeded: 1", "failed: 1", "records: 1"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Errorf("output does not contain %q:\n%s", expected, stdout.String())
		}
	}
	if !strings.Contains(stderr.String(), "job failed") {
		t.Errorf("error output does not describe failed job: %s", stderr.String())
	}
}

// Este teste protege o cálculo contra divisão por zero.

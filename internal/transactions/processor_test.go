package transactions

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProcess cria dados conhecidos e confirma os cálculos principais.
func TestProcess(t *testing.T) {
	// TempDir cria uma pasta temporária removida após o teste.
	path := filepath.Join(t.TempDir(), "transactions.csv")
	data := "id,timestamp,category,amount\n1,2026-01-01T00:00:00Z,food,10.50\n2,2026-01-01T00:00:01Z,books,20.00\n3,2026-01-01T00:00:02Z,food,4.50\n"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	// Chamamos a mesma função usada pelas versões reais.
	result, err := Process(path)
	if err != nil {
		t.Fatal(err)
	}
	if result.Records != 3 || result.Total != 35 || result.Average() != 35.0/3.0 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.ByCategory["food"] != 2 || result.ByCategory["books"] != 1 {
		t.Fatalf("unexpected categories: %+v", result.ByCategory)
	}
}

// BenchmarkProcess mede o desempenho do processador com 10 mil linhas.
func BenchmarkProcess(b *testing.B) {
	path := filepath.Join(b.TempDir(), "transactions.csv")
	file, err := os.Create(path)
	if err != nil {
		b.Fatal(err)
	}
	_, _ = file.WriteString("id,timestamp,category,amount\n")
	for i := 0; i < 10_000; i++ {
		_, _ = file.WriteString("1,2026-01-01T00:00:00Z,food,10.50\n")
	}
	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	// A criação do CSV não entra na medição.
	b.ResetTimer()
	// b.N é escolhido automaticamente para produzir uma média estável.
	for i := 0; i < b.N; i++ {
		if _, err := Process(path); err != nil {
			b.Fatal(err)
		}
	}
}

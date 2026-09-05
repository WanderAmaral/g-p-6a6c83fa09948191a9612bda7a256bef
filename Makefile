# .PHONY informa ao make que estes nomes são comandos, não arquivos.
.PHONY: test benchmark datasets

# Executa todos os testes do projeto.
test:
	go test ./...

# Mede tempo e alocações do processador compartilhado.
benchmark:
	go test -bench=. -benchmem ./internal/transactions

# Cria os três datasets usados nas demonstrações.
datasets:
	go run ./cmd/generator -records 10000 -output data/transactions-10k.csv
	go run ./cmd/generator -records 100000 -output data/transactions-100k.csv
	go run ./cmd/generator -records 1000000 -output data/transactions-1m.csv

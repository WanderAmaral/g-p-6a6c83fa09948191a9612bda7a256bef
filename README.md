# Job Processing System em Go — V1 a V5

Projeto pequeno e didático para mostrar a evolução de um processador de jobs. Cada job lê um CSV de transações em streaming e calcula quantidade, total, média, contagem por categoria, tempo e throughput.

O objetivo não é provar que concorrência sempre é mais rápida. É criar uma linha de base, aplicar carga e observar quando cada solução ajuda ou atrapalha.

## Requisitos

- Go 1.22 ou mais recente
- nenhum serviço externo ou dependência de terceiros

## Estrutura

```text
cmd/generator/          gerador determinístico de CSV
internal/transactions/ leitura, agregação, métricas e testes
internal/jobs/         resumo compartilhado dos lotes
versions/v1/            um job sequencial
versions/v2/            vários jobs sequenciais
versions/v3/            uma goroutine para cada job
versions/v4/            channel como fila, um consumidor
versions/v5/            channel com worker pool configurável
data/                   datasets gerados (ignorados pelo Git)
```

Cada versão resolve um novo problema. V1 cria a linha de base, V2 recebe vários jobs, V3 introduz concorrência, V4 separa produção e consumo com uma fila e V5 controla a concorrência com um worker pool.

## 1. Gerar os datasets

Gere os três tamanhos de uma vez:

```bash
make datasets
```

Ou gere um tamanho específico:

```bash
go run ./cmd/generator -records 10000 -output data/transactions-10k.csv
go run ./cmd/generator -records 100000 -output data/transactions-100k.csv
go run ./cmd/generator -records 1000000 -output data/transactions-1m.csv
```

O gerador usa uma seed fixa por padrão, portanto produz dados reproduzíveis. Use `-seed 123` para trocar a sequência.

## 2. Rodar a V1

```bash
go run ./versions/v1 data/transactions-10k.csv
```

Esta é a linha de base: lê uma linha por vez sem carregar o arquivo inteiro na memória, agrega os valores e imprime as métricas.

## 3. Rodar a V2

Passe dois ou mais arquivos no mesmo comando:

```bash
go run ./versions/v2 \
  data/transactions-10k.csv \
  data/transactions-100k.csv \
  data/transactions-1m.csv
```

A V2 processa os arquivos na ordem recebida. Cada job só começa depois que o anterior termina. Ao final, ela mostra quantidade de jobs, sucessos, falhas, total de registros, duração e throughput do lote.

Para verificar que uma falha não impede os outros jobs de rodar:

```bash
go run ./versions/v2 \
  data/transactions-10k.csv \
  data/nao-existe.csv \
  data/transactions-100k.csv
```

Nesse caso o programa continua, mostra `succeeded: 2` e `failed: 1`, mas termina com código de erro porque o lote não teve sucesso completo.

## 4. Rodar a V3 — uma goroutine por job

```bash
go run ./versions/v3 \
  data/transactions-10k.csv \
  data/transactions-100k.csv \
  data/transactions-1m.csv
```

Os três jobs começam concorrentemente. A V3 é fácil de escrever, mas não limita a quantidade de goroutines: mil arquivos criariam mil goroutines e poderiam causar muita disputa por CPU, memória e disco.

## 5. Rodar a V4 — channel como fila

```bash
go run ./versions/v4 \
  data/transactions-10k.csv \
  data/transactions-100k.csv \
  data/transactions-1m.csv
```

Um produtor coloca caminhos no channel `jobs`; um consumidor retira e processa. Como existe somente um consumidor, os jobs continuam sequenciais. O ganho desta versão é o desacoplamento entre criar e executar trabalho.

## 6. Rodar a V5 — worker pool

```bash
go run ./versions/v5 -workers 2 \
  data/transactions-10k.csv \
  data/transactions-100k.csv \
  data/transactions-1m.csv
```

Troque `-workers 2` por `1`, `4` e `8` e compare os tempos. Os workers compartilham a mesma fila, e cada job é entregue a apenas um deles. A ordem de conclusão pode mudar entre execuções.

## Testes e benchmark

Rode todos os testes e também a verificação de concorrência:

```bash
go test -race ./...
```

Rode o benchmark isolado do parser com um CSV temporário de 10 mil registros:

```bash
go test -bench=. -benchmem ./internal/transactions
```

Para comparar as versões de forma mais justa, execute cada comando várias vezes na mesma máquina, feche aplicações pesadas e anote a mediana. O cache do sistema operacional e a velocidade do disco influenciam o resultado.

## O que esperar — sem números inventados

- V1 fornece a referência mais simples para tempo, memória e registros por segundo.
- V2 deve produzir para cada arquivo os mesmos resultados da V1.
- Como a V2 é sequencial, a duração do lote tende a ficar próxima da soma dos jobs, além de um pequeno custo de controle e exibição.
- V3 pode terminar vários jobs em menos tempo, mas concorrência ilimitada pode gerar contenção.
- V4 não promete ser mais rápida: ela demonstra uma fila em memória e separação de responsabilidades.
- V5 limita o trabalho simultâneo. Mais workers não significa automaticamente mais velocidade.
- Memória deve crescer de forma moderada com arquivos maiores, pois o CSV é lido linha por linha.

Não há resultados numéricos registrados no repositório: eles dependem do hardware e devem ser medidos durante o vídeo.

## Roteiro sugerido para o primeiro vídeo/teste

1. Gere o CSV de 10k e apresente o formato.
2. Mostre a estrutura simples da V1 e explique a leitura linha por linha.
3. Rode V1 e apresente total, média, categorias, tempo e throughput.
4. Repita com 100k e 1M e registre os números reais.
5. Execute a V2 com vários arquivos e mostre que os jobs esperam uns pelos outros.
6. Use V3 para executar concorrentemente e explique o risco de criar goroutines sem limite.
7. Mostre o channel da V4 como fila entre produtor e consumidor.
8. Termine na V5 comparando quantidades diferentes de workers.

## Comparação conceitual

```text
V1: arquivo → Process

V2: arquivos → for → Process (um após o outro)

V3: arquivos → uma goroutine por arquivo → Process

V4: produtor → channel de jobs → um consumidor → Process

V5: produtor → channel de jobs → N workers → Process
```

As próximas evoluções possíveis são persistir a fila, criar retries, garantir idempotência e adicionar observabilidade. Elas ainda não fazem parte deste projeto.

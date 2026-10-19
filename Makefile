# Proofline — offline demo and tests.
# CGO is required because tree-sitter is linked into the binary.

export CGO_ENABLED := 1

BIN := bin/proofline

.PHONY: build test demo measure

build:
	go build -o $(BIN) ./cmd/proofline

test:
	go test ./...

demo: build
	$(BIN) demo --toy examples/toy-repo --out docs

measure: build
	$(BIN) measure --per-class 20 --reruns 2 --max-mutants 8 --out docs/measured-results.json

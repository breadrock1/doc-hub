BIN_FILE_PATH := "./bin/docs-hub"

build:
	go build -v -o $(BIN_FILE_PATH) ./cmd/docs-hub

run:
	$(BIN_FILE_PATH) -c ./configs/production.toml

test:
	go test -race ./...

.PHONY: build run test

.DEFAULT_GOAL := build

.PHONY: fmt lint run build clear
fmt:
	@go fmt ./...

lint: fmt
	@golangci-lint run ./...

build: lint
	@go build -race -o xkcd ./cmd/xkcd/

run: lint
	@go run -race ./cmd/xkcd/

clear:
	@rm ./xkcd
	@rm ./database.json

clearModel:
	@rm ./internal/resources/spellchecker/savedModel
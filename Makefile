.DEFAULT_GOAL := build

.PHONY: fmt lint run build clear
fmt:
	@go fmt ./...

lint: fmt
	@golangci-lint run ./...

build: lint
	@go build ./ -o xkcd -race

run: lint
	@go run ./

clear:
	@go clean
	@rm ./database.json

clearModel:
	@rm ./internal/resources/spellchecker/savedModel
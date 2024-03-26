.DEFAULT_GOAL := build

.PHONY: fmt lint run build clear
fmt:
	@go fmt ./...

lint: fmt
	@golangci-lint run ./...

build: lint
	@go build ./

run: lint
	@go run ./

clear:
	@go clean
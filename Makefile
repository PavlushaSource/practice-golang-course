.DEFAULT_GOAL := build

BINARYNAME := xkcd-server
SERVERPATH := ./cmd/http

.PHONY: style test/race build run
style:
	@go fmt ./...
	@golangci-lint run ./...

test/race:
	@go build -race -o ${BINARYNAME} ${SERVERPATH}
	@rm ./${BINARYNAME}

build:
	@go build -o ${BINARYNAME} ${SERVERPATH}

run:
	@go run ${SERVERPATH}
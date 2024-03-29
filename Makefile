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
	@rm spellcheck/savedModel

examples: build
	./yadro-practice-course -s "Филльм фильм Spiderman смотретть на кинопоиске"
	@echo
	./yadro-practice-course -s "follower brings bunch bunch bunch of questions"
	@echo
	# word 'long' in stopWords ISO639
	./yadro-practice-course -s "i'll follow you as long as you are following me"
	@echo
	./yadro-practice-course -s "Only four contestants remained: Louise, Jack, Michael and Ruby."
	@echo
	./yadro-practice-course -s "children's, doctor's, babies', co-operate, son-in-law"
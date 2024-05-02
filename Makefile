.PHONY:

build:
	@go mod download && go build -o bin/friendly-bot cmd/main.go

run: build
	@./bin/friendly-bot

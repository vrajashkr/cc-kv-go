.PHONY: tidy build run

tidy:
	@go mod tidy

build:
	@go build

run:
	@go run src/main.go

test:
	@go test ./...

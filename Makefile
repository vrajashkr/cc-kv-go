ENTRYPOINT := "src/main.go"

.PHONY: tidy build run run-race test test-race

tidy:
	@go mod tidy

build:
	@go build -o ./cc-kv-go ./src/main.go

run:
	@go run $(ENTRYPOINT)

run-race:
	@go run -race ${ENTRYPOINT}

test:
	@go test ./...

test-race:
	@go test -race ./...

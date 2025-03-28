ENTRYPOINT := "src/main.go"

.PHONY: tidy build run run-race test test-race test-race-coverage lint format

tidy:
	@go mod tidy

build:
	@go build -o ./cc-kv-go $(ENTRYPOINT)

run:
	@go run $(ENTRYPOINT)

run-race:
	@go run -race ${ENTRYPOINT}

test:
	@go test ./...

test-race:
	@go test -race ./...

test-race-coverage:
	@go test -race -coverprofile=coverage.txt ./...

lint:
	@golangci-lint run

format:
	@gofumpt -l -w .


ENTRYPOINT := "src/main.go"

.PHONY: tidy build run

tidy:
	@go mod tidy

build:
	@go build

run:
	@go run $(ENTRYPOINT)

run-race:
	@go run -race ${ENTRYPOINT}

test:
	@go test ./...

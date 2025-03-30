package main

import (
	"log/slog"
	"os"

	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/server"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func main() {
	slog.Info("starting cc-kv-go server")

	slog.Info("initializing storage engine")
	storageEngine := storage.NewMapStorageEngine()

	slog.Info("initializing command handler")
	commandHandler := handler.NewCommandHandler(&storageEngine)

	slog.Info("starting listener")
	listener, err := server.NewTcpServer("6379", commandHandler.ServeInput)
	if err != nil {
		slog.Error("failed to start listener", "error", err.Error())
		os.Exit(1)
	}
	defer listener.StopListen()

	listener.Serve()
}

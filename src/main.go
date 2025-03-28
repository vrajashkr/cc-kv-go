package main

import (
	"io"
	"log/slog"
	"net"

	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

const (
	LARGE_BUF_SIZE = 2048
	READ_BUF_SIZE  = 128
)

func main() {
	slog.Info("starting cc-kv-go server")

	slog.Info("initializing storage engine")
	storageEngine := storage.NewMapStorageEngine()

	slog.Info("starting listener")
	l, err := net.Listen("tcp4", ":6379")
	if err != nil {
		slog.Error("failed to start TCP server", "error", err.Error())
		return
	}
	defer func() {
		err := l.Close()
		if err != nil {
			slog.Error("failed to close listener", "error", err.Error())
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err.Error())
			return
		}
		go handleConnection(conn, &storageEngine)
	}
}

func handleConnection(c net.Conn, storage storage.StorageEngine) {
	defer func() {
		err := c.Close()
		if err != nil {
			slog.Error("failed to close connection", "error", err.Error())
		}
	}()

	reachedEnd := false
	for {
		buf := make([]byte, 0, LARGE_BUF_SIZE)
		tmp := make([]byte, READ_BUF_SIZE)
		for {
			numBytes, err := c.Read(tmp)
			if err != nil {
				if err != io.EOF {
					slog.Error("error while processing request", "error", err.Error())
				}
				reachedEnd = true
				break
			}
			buf = append(buf, tmp[:numBytes]...)
			if numBytes < READ_BUF_SIZE {
				result := handler.ServeInput(buf, storage)
				_, err := c.Write([]byte(result))
				if err != nil {
					slog.Error("failed to respond to client", "error", err.Error())
				}
				break
			}
		}
		if reachedEnd {
			break
		}
	}
}

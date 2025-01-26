package main

import (
	"io"
	"log/slog"
	"net"

	"github.com/vrajashkr/cc-kv-go/src/handler"
)

const (
	LARGE_BUF_SIZE = 2048
	READ_BUF_SIZE  = 256
)

func main() {
	slog.Info("starting cc-kv-go server")

	l, err := net.Listen("tcp4", ":6379")
	if err != nil {
		slog.Error("failed to start TCP server due to error: " + err.Error())
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("failed to accept connection due to error: " + err.Error())
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 0, LARGE_BUF_SIZE)
	tmp := make([]byte, READ_BUF_SIZE)
	for {
		numBytes, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				slog.Error("error while processing request: " + err.Error())
			}
			break
		}
		buf = append(buf, tmp[:numBytes]...)
		if numBytes < READ_BUF_SIZE {
			break
		}
	}

	result := handler.ServeInput(buf)
	_, err := c.Write([]byte(result))
	if err != nil {
		slog.Error("failed to respond to client due to error: " + err.Error())
	}
}

package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
)

const (
	LARGE_BUF_SIZE = 2048
	READ_BUF_SIZE  = 128
)

type TcpServer struct {
	listener       *net.Listener
	connHandleFunc func(msg []byte) string
}

func NewTcpServer(port string, handlerFunc func([]byte) string) (*TcpServer, error) {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}

	return &TcpServer{
		&listener,
		handlerFunc,
	}, nil
}

func (ts *TcpServer) Serve() {
	for {
		conn, err := (*ts.listener).Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err.Error())
			return
		}
		go ts.handleConnection(conn)
	}
}

func (ts *TcpServer) StopListen() {
	if ts.listener != nil {
		closeErr := (*ts.listener).Close()
		if closeErr != nil {
			slog.Error("failed to close listener", "error", closeErr.Error())
		}
	}
}

func (ts *TcpServer) handleConnection(c net.Conn) {
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
				result := ts.connHandleFunc(buf)
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

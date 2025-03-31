package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/server"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

type TestCase struct {
	input string
	want  string
}

func TestApplicationWithConcurrentClient(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	serverPort := "34565"

	strgEng := storage.NewMapStorageEngine()
	cmdHandler := handler.NewCommandHandler(&strgEng)
	listener, err := server.NewTcpServer(serverPort, cmdHandler.ServeInput)
	require.Nil(err)

	go listener.Serve()

	tcSet1 := []TestCase{
		{"*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n", "+OK\r\n"},
		{"*2\r\n$3\r\nGET\r\n$5\r\nworld\r\n", "$-1\r\n"},
		{"*2\r\n$3\r\nGET\r\n$5\r\nhello\r\n", "$5\r\nworld\r\n"},
		{"*5\r\n$5\r\nLPUSH\r\n$5\r\nlist1\r\n$2\r\nk1\r\n$2\r\nk2\r\n$2\r\nk3\r\n", ":3\r\n"},
		{"*4\r\n$5\r\nRPUSH\r\n$5\r\nlist1\r\n$2\r\nk4\r\n$2\r\nk5\r\n", ":2\r\n"},
		{"*2\r\n$4\r\nINCR\r\n$4\r\nctr1\r\n", ":1\r\n"},
		{"*2\r\n$4\r\nDECR\r\n$4\r\nctr1\r\n", ":0\r\n"},
		{"*2\r\n$3\r\nDEL\r\n$4\r\nctr1\r\n", ":1\r\n"},
		{"*4\r\n$3\r\nDEL\r\n$4\r\nctr1\r\n$5\r\nlist1\r\n$5\r\nhello\r\n", ":2\r\n"},
	}

	tcSet2 := []TestCase{
		{"*5\r\n$3\r\nDEL\r\n$2\r\nn1\r\n$2\r\nn1\r\n$2\r\nn2\r\n$2\r\nn2\r\n", ":0\r\n"},
		{"*3\r\n$3\r\nSET\r\n$2\r\nn1\r\n$8\r\ntestval1\r\n", "+OK\r\n"},
		{"*1\r\n$4\r\nECHO\r\n", "-wrong number of arguments for 'echo' command\r\n"},
		{"*2\r\n$4\r\nECHO\r\n$4\r\ntest\r\n", "$4\r\ntest\r\n"},
		{"*3\r\n$6\r\nEXISTS\r\n$2\r\nn1\r\n$2\r\nn2\r\n", ":1\r\n"},
		{"*1\r\n$4\r\nPING\r\n", "+PONG\r\n"},
		{"*3\r\n$3\r\nSET\r\n$6\r\nelpmas\r\n$8\r\necnetnes\r\n", "+OK\r\n"},
		{"*2\r\n$3\r\nGET\r\n$6\r\nelpmas\r\n", "$8\r\necnetnes\r\n"},
		{"*5\r\n$3\r\nDEL\r\n$6\r\nelpmas\r\n$4\r\ntest\r\n$2\r\nn1\r\n$2\r\nn2\r\n", ":2\r\n"},
		{"*4\r\n$3\r\nSET\r\n$6\r\nelpmas\r\n$2\r\nPX\r\n$2\r\n32\r\n", "-wrong number of arguments for 'set' command\r\n"},
		{"*5\r\n$3\r\nSET\r\n$6\r\nelpmas\r\n$8\r\necnetnes\r\n$2\r\nPX\r\n$2\r\n32\r\n", "+OK\r\n"},
	}

	tcSets := [][]TestCase{tcSet1, tcSet2}

	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+serverPort)
	require.Nil(err)

	for idx, tcSet := range tcSets {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			t.Parallel()

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			require.Nil(err)

			// Run each testCaseSet in order
			for _, tc := range tcSet {
				_, err := conn.Write([]byte(tc.input))
				require.Nil(err)

				reply := make([]byte, 1024)
				numBytesRead, err := conn.Read(reply)
				require.Nil(err)

				assert.Equal(tc.want, string(reply[:numBytesRead]))
			}

			_ = conn.Close()
		})
	}
}

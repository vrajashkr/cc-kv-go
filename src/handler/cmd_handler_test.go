package handler_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "PING"},
				},
			},
			data.SimpleString{Contents: "PONG"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "PING"},
					data.BulkString{Data: "hello world"},
				},
			},
			data.BulkString{Data: "hello world"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "ECHO"},
					data.BulkString{Data: "hello world"},
				},
			},
			data.BulkString{Data: "hello world"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "ECHO"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.SimpleString{Contents: "test"},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.SimpleString{Contents: "test"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "UNSUPPORTED"},
				},
			},
			data.Error{ErrMsg: "unsupported command UNSUPPORTED"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "3"},
				},
			},
			data.Error{ErrMsg: "NOPROTO sorry, this protocol version is not supported"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "2"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "nan"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.Integer{Value: 1},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Null{},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
					data.Integer{Value: 1},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
					data.BulkString{Data: "testVal"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.BulkString{Data: "testVal"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.SimpleString{Contents: "GET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.BulkString{Data: "GET"},
				},
			},
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "maxmemory"},
					data.BulkString{Data: "0"},
					data.BulkString{Data: "save"},
					data.BulkString{Data: ""},
					data.BulkString{Data: "appendonly"},
					data.BulkString{Data: "no"},
				},
			},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.BulkString{Data: "NEXIST"},
				},
			},
			data.Error{ErrMsg: "unsupported subcommand NEXIST for CONFIG"},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input-%s", tc.input.ToDataString()), func(t *testing.T) {
			result := handler.HandleCommand(tc.input, &storageEngine)
			assert.Equal(tc.want, result)
		})
	}
}

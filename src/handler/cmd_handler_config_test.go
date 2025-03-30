package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleConfigCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
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
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
			result := handler.HandleCommand(tc.input, &storageEngine)
			assert.Equal(tc.want, result)
		})
	}
}

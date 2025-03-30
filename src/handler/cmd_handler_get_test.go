package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleGetCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
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
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
			result := handler.HandleCommand(tc.input, &storageEngine)
			assert.Equal(tc.want, result)
		})
	}
}

package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleListOperationCommands(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
					data.BulkString{Data: "ctr1"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'lpush' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
					data.BulkString{Data: "samplelist"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
					data.SimpleString{Contents: "samplelist"},
					data.BulkString{Data: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "RPUSH"},
					data.BulkString{Data: "ctr1"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'rpush' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "RPUSH"},
					data.SimpleString{Contents: "samplelist"},
					data.BulkString{Data: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "RPUSH"},
					data.BulkString{Data: "samplelist"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'lpush' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "RPUSH"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'rpush' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
					data.BulkString{Data: "mylist"},
					data.BulkString{Data: "key1"},
					data.BulkString{Data: "key2"},
				},
			},
			data.Integer{Value: 2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "RPUSH"},
					data.BulkString{Data: "mylist"},
					data.BulkString{Data: "key1"},
					data.BulkString{Data: "key2"},
				},
			},
			data.Integer{Value: 2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "LPUSH"},
					data.BulkString{Data: "mylist"},
					data.BulkString{Data: "key1"},
					data.BulkString{Data: "key2"},
				},
			},
			data.Integer{Value: 2},
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

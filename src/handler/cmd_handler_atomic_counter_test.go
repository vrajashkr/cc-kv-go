package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleAtomicCounterCommands(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "ictr"},
				},
			},
			data.Integer{Value: 1},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "dctr"},
				},
			},
			data.Integer{Value: -1},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "ictr"},
				},
			},
			data.Integer{Value: 2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "dctr"},
				},
			},
			data.Integer{Value: -2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "nanCtr"},
					data.BulkString{Data: "testVal"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "nanCtr"},
				},
			},
			data.Error{ErrMsg: "value is not an integer or out of range"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "nanCtr"},
				},
			},
			data.Error{ErrMsg: "value is not an integer or out of range"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
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

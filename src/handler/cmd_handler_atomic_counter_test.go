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
	ch := handler.NewCommandHandler(&storageEngine)

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
			data.Error{ErrMsg: "wrong number of arguments for 'incr' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'decr' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
			result := ch.HandleCommand(tc.input)
			assert.Equal(tc.want, result)
		})
	}
}

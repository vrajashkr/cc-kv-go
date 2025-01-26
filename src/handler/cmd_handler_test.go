package handler_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
)

func TestHandleCommand(t *testing.T) {
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
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input-%s", tc.input.ToDataString()), func(t *testing.T) {
			result := handler.HandleCommand(tc.input)
			assert.Equal(tc.want, result)
		})
	}
}

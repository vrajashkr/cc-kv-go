package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()
	ch := handler.NewCommandHandler(&storageEngine)

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
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
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
			result := ch.HandleCommand(tc.input)
			assert.Equal(tc.want, result)
		})
	}
}

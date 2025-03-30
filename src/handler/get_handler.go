package handler

import (
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/get/
func handleGet(cmd data.Array, strg storage.StorageEngine) data.Message {
	if len(cmd.Elements) < 2 {
		return INVALID_CMD_ARGS
	}

	keyHolder, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	ok, val, err := strg.Get(keyHolder.Data)
	if err != nil {
		return data.Error{ErrMsg: "failed to retrieve data due to error: " + err.Error()}
	}

	if !ok {
		return data.Null{}
	}

	return data.BulkString{Data: val}
}

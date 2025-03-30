package handler

import (
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/del/
func handleDelete(cmd data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmd.Elements)
	if cmdLen < 2 {
		return INVALID_CMD_ARGS
	}

	keys := make([]string, cmdLen-1)

	for idx := range cmdLen - 1 {
		keyToCheck, ok := cmd.Elements[1+idx].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

		keys[idx] = keyToCheck.Data
	}

	result, err := strg.Delete(keys)
	if err != nil {
		return data.Error{ErrMsg: "failed to execute delete. error: " + err.Error()}
	}

	return data.Integer{Value: int64(result)}
}

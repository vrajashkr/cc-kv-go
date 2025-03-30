package handler

import (
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/exists/
func handleExists(cmd data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmd.Elements)

	keys := make([]string, cmdLen-1)

	for idx := range cmdLen - 1 {
		keyToCheck := cmd.Elements[1+idx].(data.BulkString)
		keys[idx] = keyToCheck.Data
	}

	result, err := strg.Exists(keys)
	if err != nil {
		return data.Error{ErrMsg: "failed to execute exists. error: " + err.Error()}
	}

	return data.Integer{Value: int64(result)}
}

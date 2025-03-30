package handler

import (
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/lpush/
// https://redis.io/docs/latest/commands/rpush/
func handleListPush(cmdArray data.Array, strg storage.StorageEngine, isPrepend bool) data.Message {
	cmdLen := len(cmdArray.Elements)
	if cmdLen < 3 {
		return INVALID_CMD_ARGS
	}

	listToUpdate, ok := cmdArray.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	listNameToUpdate := listToUpdate.Data

	listValues := make([]string, cmdLen-2)

	for idx := range cmdLen - 2 {
		keyToCheck, ok := cmdArray.Elements[2+idx].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

		listValues[idx] = keyToCheck.Data
	}

	res, err := strg.ListPush(listNameToUpdate, listValues, isPrepend)
	if err != nil {
		return data.Error{ErrMsg: "command failed. error: " + err.Error()}
	}

	return data.Integer{Value: res}
}

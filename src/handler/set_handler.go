package handler

import (
	"strconv"
	"time"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/set/
func handleSet(cmd data.Array, strg storage.StorageEngine) data.Message {
	numArgs := len(cmd.Elements)
	keyHolder := cmd.Elements[1].(data.BulkString)
	valueHolder := cmd.Elements[2].(data.BulkString)

	valueContents := valueHolder.Data
	expires := false
	var expiresAtTimeStampMillis int64 = -1

	// if there are more args, 2 more are expected
	if numArgs > 3 {
		if numArgs != 5 {
			return data.Error{ErrMsg: "wrong number of arguments for 'set' command"}
		}

		// handle command options
		option := cmd.Elements[3].(data.BulkString)
		optionArg := cmd.Elements[4].(data.BulkString)

		expires = true
		optionTimeInt, err := strconv.ParseInt(optionArg.Data, 10, 64)
		if err != nil {
			return data.Error{ErrMsg: "failed to set due to error: " + err.Error()}
		}

		switch option.Data {
		case CMD_SET_OPT_EX:
			// expiry time in seconds
			expiresAtTimeStampMillis = time.Now().UnixMilli() + (optionTimeInt * 1000)
		case CMD_SET_OPT_PX:
			// expiry time in milliseconds
			expiresAtTimeStampMillis = time.Now().UnixMilli() + optionTimeInt
		case CMD_SET_OPT_EXAT:
			// expiry timestamp epoch in seconds
			expiresAtTimeStampMillis = optionTimeInt * 1000
		case CMD_SET_OPT_PXAT:
			// expiry timestamp epoch in milliseconds
			expiresAtTimeStampMillis = optionTimeInt
		default:
			return data.Error{ErrMsg: "syntax error"}
		}
	}

	err := strg.Set(keyHolder.Data, valueContents, expires, expiresAtTimeStampMillis)
	if err != nil {
		return data.Error{ErrMsg: "failed to store data due to error: " + err.Error()}
	}
	return OK
}

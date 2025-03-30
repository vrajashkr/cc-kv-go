package handler

import "github.com/vrajashkr/cc-kv-go/src/data"

// https://redis.io/commands/echo/
func handleEcho(cmd data.Array) data.Message {
	if len(cmd.Elements) < 2 {
		return INVALID_CMD_ARGS
	}

	incomingContents, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}
	return incomingContents
}

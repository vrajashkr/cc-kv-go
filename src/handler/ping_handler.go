package handler

import "github.com/vrajashkr/cc-kv-go/src/data"

// https://redis.io/docs/latest/commands/ping/
func handlePing(cmd data.Array) data.Message {
	cmdLen := len(cmd.Elements)

	if cmdLen == 1 {
		return data.SimpleString{Contents: "PONG"}
	}

	return cmd.Elements[1].(data.BulkString)
}

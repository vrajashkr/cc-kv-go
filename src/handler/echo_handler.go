package handler

import "github.com/vrajashkr/cc-kv-go/src/data"

// https://redis.io/commands/echo/
func handleEcho(cmd data.Array) data.Message {
	return cmd.Elements[1].(data.BulkString)
}

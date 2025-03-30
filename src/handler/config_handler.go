package handler

import (
	"fmt"

	"github.com/vrajashkr/cc-kv-go/src/data"
)

// https://redis.io/docs/latest/commands/config-get/
func handleConfig(cmdArray data.Array) data.Message {
	if len(cmdArray.Elements) < 2 {
		return INVALID_CMD_ARGS
	}

	subCommandHolder, ok := cmdArray.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	switch subCommandHolder.Data {
	case "GET":
		return data.Array{
			Elements: []data.Message{
				data.BulkString{Data: "maxmemory"},
				data.BulkString{Data: "0"},
				data.BulkString{Data: "save"},
				data.BulkString{Data: ""},
				data.BulkString{Data: "appendonly"},
				data.BulkString{Data: "no"},
			},
		}
	default:
		return data.Error{
			ErrMsg: fmt.Sprintf("unsupported subcommand %s for %s", subCommandHolder.Data, CMD_CONFIG),
		}
	}
}

package handler

import (
	"fmt"

	"github.com/vrajashkr/cc-kv-go/src/data"
)

const (
	CMD_PING = "PING"
)

// https://redis.io/docs/latest/commands/ping/
func handlePing(cmd data.Array) data.Message {
	cmdLen := len(cmd.Elements)

	if cmdLen == 1 {
		return data.SimpleString{Contents: "PONG"}
	}

	incomingContents, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return data.Error{ErrMsg: "invalid args for command"}
	}
	return incomingContents
}

func HandleCommand(msg data.Message) data.Message {
	cmdArray, ok := msg.(data.Array)
	if !ok {
		return data.Error{
			ErrMsg: "invalid format for command",
		}
	}

	firstCmd, ok := cmdArray.Elements[0].(data.BulkString)
	if !ok {
		return data.Error{
			ErrMsg: "invalid format for command",
		}
	}

	var result data.Message
	switch firstCmd.Data {
	case CMD_PING:
		result = handlePing(cmdArray)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}

	return result
}

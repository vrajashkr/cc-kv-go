package handler

import "github.com/vrajashkr/cc-kv-go/src/data"

const (
	CMD_PING = "PING"
)

func handlePing(_ data.Array) data.Message {
	return data.SimpleString{
		Contents: "PONG",
	}
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
		result = data.Error{ErrMsg: "unsupported command"}
	}

	return result
}

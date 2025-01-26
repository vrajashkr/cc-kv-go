package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/vrajashkr/cc-kv-go/src/data"
)

const (
	CMD_PING  = "PING"
	CMD_HELLO = "HELLO"
)

var INVALID_CMD_FMT = data.Error{ErrMsg: "invalid format for command"}
var INVALID_CMD_ARGS = data.Error{ErrMsg: "invalid args for command"}

// https://redis.io/docs/latest/commands/ping/
func handlePing(cmd data.Array) data.Message {
	cmdLen := len(cmd.Elements)

	if cmdLen == 1 {
		return data.SimpleString{Contents: "PONG"}
	}

	incomingContents, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}
	return incomingContents
}

// https://redis.io/docs/latest/commands/hello/
func handleHello(cmd data.Array) data.Message {
	versionEntry, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_FMT
	}
	versionRequested, err := strconv.Atoi(versionEntry.Data)
	if err != nil {
		slog.Error("invalid command args: " + err.Error())
		return INVALID_CMD_ARGS
	}
	if versionRequested != 2 {
		return data.Error{ErrMsg: "NOPROTO sorry, this protocol version is not supported"}
	}
	return data.SimpleString{Contents: "OK"}
}

func HandleCommand(msg data.Message) data.Message {
	cmdArray, ok := msg.(data.Array)
	if !ok {
		return INVALID_CMD_FMT
	}

	firstCmd, ok := cmdArray.Elements[0].(data.BulkString)
	if !ok {
		return INVALID_CMD_FMT
	}

	var result data.Message
	switch firstCmd.Data {
	case CMD_PING:
		result = handlePing(cmdArray)
	case CMD_HELLO:
		result = handleHello(cmdArray)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}

	return result
}

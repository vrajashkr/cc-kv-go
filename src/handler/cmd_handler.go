package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

const (
	CMD_PING  = "PING"
	CMD_HELLO = "HELLO"
	CMD_ECHO  = "ECHO"
	CMD_SET   = "SET"
	CMD_GET   = "GET"
)

var INVALID_CMD_FMT = data.Error{ErrMsg: "invalid format for command"}
var INVALID_CMD_ARGS = data.Error{ErrMsg: "invalid args for command"}
var OK = data.SimpleString{Contents: "OK"}

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
	return OK
}

// https://redis.io/docs/latest/commands/set/
func handleSet(cmd data.Array, strg storage.StorageEngine) data.Message {
	if len(cmd.Elements) < 3 {
		return INVALID_CMD_ARGS
	}

	keyHolder, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	valueHolder, ok := cmd.Elements[2].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	err := strg.Set(keyHolder.Data, valueHolder.Data)
	if err != nil {
		return data.Error{ErrMsg: "failed to store data due to error: " + err.Error()}
	}
	return OK
}

// https://redis.io/docs/latest/commands/get/
func handleGet(cmd data.Array, strg storage.StorageEngine) data.Message {
	if len(cmd.Elements) < 2 {
		return INVALID_CMD_ARGS
	}

	keyHolder, ok := cmd.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	ok, val, err := strg.Get(keyHolder.Data)
	if err != nil {
		return data.Error{ErrMsg: "failed to retrieve data due to error: " + err.Error()}
	}

	if !ok {
		return data.Null{}
	}

	return data.BulkString{Data: val}
}

func HandleCommand(msg data.Message, strg storage.StorageEngine) data.Message {
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
	case CMD_ECHO:
		result = handleEcho(cmdArray)
	case CMD_SET:
		result = handleSet(cmdArray, strg)
	case CMD_GET:
		result = handleGet(cmdArray, strg)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}

	return result
}

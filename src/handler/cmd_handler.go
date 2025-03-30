package handler

import (
	"fmt"
	"log/slog"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

const (
	CMD_PING         = "PING"
	CMD_HELLO        = "HELLO"
	CMD_ECHO         = "ECHO"
	CMD_SET          = "SET"
	CMD_SET_OPT_EX   = "EX"
	CMD_SET_OPT_EXAT = "EXAT"
	CMD_SET_OPT_PX   = "PX"
	CMD_SET_OPT_PXAT = "PXAT"
	CMD_GET          = "GET"
	CMD_CONFIG       = "CONFIG"
	CMD_EXISTS       = "EXISTS"
	CMD_DELETE       = "DEL"
	CMD_INCR         = "INCR"
	CMD_DECR         = "DECR"
	CMD_LPUSH        = "LPUSH"
	CMD_RPUSH        = "RPUSH"
)

var (
	INVALID_CMD_FMT  = data.Error{ErrMsg: "invalid format for command"}
	INVALID_CMD_ARGS = data.Error{ErrMsg: "invalid args for command"}
	OK               = data.SimpleString{Contents: "OK"}
)

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
	case CMD_CONFIG:
		result = handleConfig(cmdArray)
	case CMD_EXISTS:
		result = handleExists(cmdArray, strg)
	case CMD_DELETE:
		result = handleDelete(cmdArray, strg)
	case CMD_INCR:
		result = handleAtomicUnitDelta(cmdArray, strg, false)
	case CMD_DECR:
		result = handleAtomicUnitDelta(cmdArray, strg, true)
	case CMD_LPUSH:
		result = handleListPush(cmdArray, strg, true)
	case CMD_RPUSH:
		result = handleListPush(cmdArray, strg, false)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}
	slog.Info(fmt.Sprintf("responding with %q", result.ToDataString()))
	return result
}

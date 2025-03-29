package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

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
)

var (
	INVALID_CMD_FMT  = data.Error{ErrMsg: "invalid format for command"}
	INVALID_CMD_ARGS = data.Error{ErrMsg: "invalid args for command"}
	OK               = data.SimpleString{Contents: "OK"}
)

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
		slog.Error("invalid command args", "error", err.Error())
		return INVALID_CMD_ARGS
	}
	if versionRequested != 2 {
		return data.Error{ErrMsg: "NOPROTO sorry, this protocol version is not supported"}
	}
	return OK
}

// https://redis.io/docs/latest/commands/set/
func handleSet(cmd data.Array, strg storage.StorageEngine) data.Message {
	numArgs := len(cmd.Elements)

	if numArgs < 3 {
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

	valueContents := valueHolder.Data
	expires := false
	var expiresAtTimeStampMillis int64 = -1

	// if there are more args, 2 more are expected
	if numArgs > 3 {
		if numArgs != 5 {
			return INVALID_CMD_ARGS
		}

		// handle command options
		option, ok := cmd.Elements[3].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

		optionArg := cmd.Elements[4].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

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
			return INVALID_CMD_FMT
		}
	}

	err := strg.Set(keyHolder.Data, valueContents, expires, expiresAtTimeStampMillis)
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

// https://redis.io/docs/latest/commands/exists/
func handleExists(cmd data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmd.Elements)
	if cmdLen < 2 {
		return INVALID_CMD_ARGS
	}

	keys := make([]string, cmdLen-1)

	for idx := range cmdLen - 1 {
		keyToCheck, ok := cmd.Elements[1+idx].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

		keys[idx] = keyToCheck.Data
	}

	result, err := strg.Exists(keys)
	if err != nil {
		return data.Error{ErrMsg: "failed to execute exists. error: " + err.Error()}
	}

	return data.Integer{Value: int64(result)}
}

// https://redis.io/docs/latest/commands/del/
func handleDelete(cmd data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmd.Elements)
	if cmdLen < 2 {
		return INVALID_CMD_ARGS
	}

	keys := make([]string, cmdLen-1)

	for idx := range cmdLen - 1 {
		keyToCheck, ok := cmd.Elements[1+idx].(data.BulkString)
		if !ok {
			return INVALID_CMD_ARGS
		}

		keys[idx] = keyToCheck.Data
	}

	result, err := strg.Delete(keys)
	if err != nil {
		return data.Error{ErrMsg: "failed to execute delete. error: " + err.Error()}
	}

	return data.Integer{Value: int64(result)}
}

// https://redis.io/docs/latest/commands/incr/
func handleIncr(cmdArray data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmdArray.Elements)
	if cmdLen < 2 {
		return INVALID_CMD_ARGS
	}

	key, ok := cmdArray.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	res, err := strg.AtomicDelta(key.Data, 1)
	if err != nil {
		slog.Error("failed to execute increment", "error", err.Error())
		return data.Error{ErrMsg: "value is not an integer or out of range"}
	}

	return data.Integer{Value: res}
}

// https://redis.io/docs/latest/commands/decr/
func handleDecr(cmdArray data.Array, strg storage.StorageEngine) data.Message {
	cmdLen := len(cmdArray.Elements)
	if cmdLen < 2 {
		return INVALID_CMD_ARGS
	}

	key, ok := cmdArray.Elements[1].(data.BulkString)
	if !ok {
		return INVALID_CMD_ARGS
	}

	res, err := strg.AtomicDelta(key.Data, -1)
	if err != nil {
		slog.Error("failed to execute decrement", "error", err.Error())
		return data.Error{ErrMsg: "value is not an integer or out of range"}
	}

	return data.Integer{Value: res}
}

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
		result = handleIncr(cmdArray, strg)
	case CMD_DECR:
		result = handleDecr(cmdArray, strg)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}
	slog.Info(fmt.Sprintf("responding with %q", result.ToDataString()))
	return result
}

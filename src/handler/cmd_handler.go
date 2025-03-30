package handler

import (
	"fmt"
	"log/slog"
	"strings"

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

// for commands that have a minimum arg count, an entry is added to this map.
// if there is no entry for that command, it is assumed that there is no minimum argument count for it.
var CMD_MIN_ARGS = map[string]int{
	CMD_INCR:   1,
	CMD_DECR:   1,
	CMD_CONFIG: 1,
	CMD_DELETE: 1,
	CMD_ECHO:   1,
	CMD_EXISTS: 1,
	CMD_GET:    1,
	CMD_LPUSH:  2,
	CMD_RPUSH:  2,
	CMD_SET:    2,
}

func validateCommand(cmd data.Array) error {
	// check that all the entries are BulkString
	for _, element := range cmd.Elements {
		_, ok := element.(data.BulkString)
		if !ok {
			return fmt.Errorf("invalid format for command")
		}
	}

	command := cmd.Elements[0].(data.BulkString).Data

	// check that the command has the correct minimum number of args
	numArgs := CMD_MIN_ARGS[command]

	if len(cmd.Elements) < numArgs+1 {
		return fmt.Errorf("wrong number of arguments for '%s' command", strings.ToLower(command))
	}

	return nil
}

type CommandHandler struct {
	strgEngine storage.StorageEngine
}

func NewCommandHandler(storageEngine storage.StorageEngine) CommandHandler {
	return CommandHandler{
		strgEngine: storageEngine,
	}
}

func (ch CommandHandler) ServeInput(rawData []byte) string {
	rawStr := string(rawData)
	rawStrLen := len(rawStr)
	slog.Debug("received message", "msg", rawStr)
	numProcessedChars := 0
	result := ""

	for {
		numChars, parsedMsg, err := data.ProcessMessageString(rawStr[numProcessedChars:])
		if err != nil {
			result += data.Error{ErrMsg: err.Error()}.ToDataString()
		} else {
			result += ch.HandleCommand(parsedMsg).ToDataString()
		}
		numProcessedChars += numChars
		if numProcessedChars == rawStrLen {
			break
		}
	}
	slog.Debug("response", "resp", result)
	return result
}

func (ch CommandHandler) HandleCommand(msg data.Message) data.Message {
	cmdArray, ok := msg.(data.Array)
	if !ok {
		return INVALID_CMD_FMT
	}

	firstCmd, ok := cmdArray.Elements[0].(data.BulkString)
	if !ok {
		return INVALID_CMD_FMT
	}

	err := validateCommand(cmdArray)
	if err != nil {
		return data.Error{ErrMsg: err.Error()}
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
		result = handleSet(cmdArray, ch.strgEngine)
	case CMD_GET:
		result = handleGet(cmdArray, ch.strgEngine)
	case CMD_CONFIG:
		result = handleConfig(cmdArray)
	case CMD_EXISTS:
		result = handleExists(cmdArray, ch.strgEngine)
	case CMD_DELETE:
		result = handleDelete(cmdArray, ch.strgEngine)
	case CMD_INCR:
		result = handleAtomicUnitDelta(cmdArray, ch.strgEngine, false)
	case CMD_DECR:
		result = handleAtomicUnitDelta(cmdArray, ch.strgEngine, true)
	case CMD_LPUSH:
		result = handleListPush(cmdArray, ch.strgEngine, true)
	case CMD_RPUSH:
		result = handleListPush(cmdArray, ch.strgEngine, false)
	default:
		result = data.Error{
			ErrMsg: fmt.Sprintf("unsupported command %s", firstCmd.Data),
		}
	}
	return result
}

package handler

import (
	"log/slog"
	"strconv"

	"github.com/vrajashkr/cc-kv-go/src/data"
)

// https://redis.io/docs/latest/commands/hello/
func handleHello(cmd data.Array) data.Message {
	versionEntry := cmd.Elements[1].(data.BulkString)

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

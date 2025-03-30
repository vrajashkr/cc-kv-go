package handler

import (
	"log/slog"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

// https://redis.io/docs/latest/commands/incr/
// https://redis.io/docs/latest/commands/decr/
func handleAtomicUnitDelta(cmdArray data.Array, strg storage.StorageEngine, isDecrement bool) data.Message {
	key := cmdArray.Elements[1].(data.BulkString)

	delta := int64(1)
	if isDecrement {
		delta = int64(-1)
	}

	res, err := strg.AtomicDelta(key.Data, delta)
	if err != nil {
		slog.Error("failed to execute atomic counter change", "error", err.Error())
		return data.Error{ErrMsg: "value is not an integer or out of range"}
	}

	return data.Integer{Value: res}
}

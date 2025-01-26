package handler

import (
	"log/slog"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func ServeInput(rawData []byte, storage storage.StorageEngine) string {
	rawStr := string(rawData)
	_, parsedMsg, err := data.ProcessMessageString(rawStr)
	slog.Info("received message: " + rawStr)
	if err != nil {
		return data.Error{
			ErrMsg: err.Error(),
		}.ToDataString()
	}
	return HandleCommand(parsedMsg, storage).ToDataString()
}

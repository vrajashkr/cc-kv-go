package handler

import (
	"fmt"
	"log/slog"

	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func ServeInput(rawData []byte, storage storage.StorageEngine) string {
	rawStr := string(rawData)
	rawStrLen := len(rawStr)
	slog.Info(fmt.Sprintf("received message: %q", rawStr))
	numProcessedChars := 0
	result := ""

	for {
		numChars, parsedMsg, err := data.ProcessMessageString(rawStr[numProcessedChars:])
		if err != nil {
			result += data.Error{ErrMsg: err.Error()}.ToDataString()
		} else {
			result += HandleCommand(parsedMsg, storage).ToDataString()
		}
		numProcessedChars += numChars
		if numProcessedChars == rawStrLen {
			break
		}
	}

	return result
}

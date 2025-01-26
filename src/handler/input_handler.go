package handler

import "github.com/vrajashkr/cc-kv-go/src/data"

func ServeInput(rawData []byte) string {
	rawStr := string(rawData)
	_, parsedMsg, err := data.ProcessMessageString(rawStr)
	if err != nil {
		return data.Error{
			ErrMsg: err.Error(),
		}.ToDataString()
	}
	return HandleCommand(parsedMsg).ToDataString()
}

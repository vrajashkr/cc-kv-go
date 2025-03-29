package handler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
	"github.com/vrajashkr/cc-kv-go/src/handler"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestHandleCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "PING"},
				},
			},
			data.SimpleString{Contents: "PONG"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "PING"},
					data.BulkString{Data: "hello world"},
				},
			},
			data.BulkString{Data: "hello world"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "ECHO"},
					data.BulkString{Data: "hello world"},
				},
			},
			data.BulkString{Data: "hello world"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "ECHO"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.SimpleString{Contents: "test"},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.SimpleString{Contents: "test"},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "UNSUPPORTED"},
				},
			},
			data.Error{ErrMsg: "unsupported command UNSUPPORTED"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "3"},
				},
			},
			data.Error{ErrMsg: "NOPROTO sorry, this protocol version is not supported"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "2"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "HELLO"},
					data.BulkString{Data: "nan"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.Integer{Value: 1},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Null{},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
					data.Integer{Value: 1},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
					data.BulkString{Data: "testVal"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.BulkString{Data: "testVal"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "EXISTS"},
					data.BulkString{Data: "testKey"},
					data.BulkString{Data: "testNoKey"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Integer{Value: 2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "EXISTS"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "EXISTS"},
					data.Integer{Value: 12},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DEL"},
					data.BulkString{Data: "testKey"},
					data.BulkString{Data: "testNoKey"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Integer{Value: 1},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DEL"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DEL"},
					data.Integer{Value: 12},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "ictr"},
				},
			},
			data.Integer{Value: 1},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "dctr"},
				},
			},
			data.Integer{Value: -1},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "ictr"},
				},
			},
			data.Integer{Value: 2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "dctr"},
				},
			},
			data.Integer{Value: -2},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "nanCtr"},
					data.BulkString{Data: "testVal"},
				},
			},
			data.SimpleString{Contents: "OK"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.BulkString{Data: "nanCtr"},
				},
			},
			data.Error{ErrMsg: "value is not an integer or out of range"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.BulkString{Data: "nanCtr"},
				},
			},
			data.Error{ErrMsg: "value is not an integer or out of range"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "INCR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "DECR"},
					data.SimpleString{Contents: "ctr1"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.SimpleString{Contents: "GET"},
				},
			},
			data.Error{ErrMsg: "invalid args for command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.BulkString{Data: "GET"},
				},
			},
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "maxmemory"},
					data.BulkString{Data: "0"},
					data.BulkString{Data: "save"},
					data.BulkString{Data: ""},
					data.BulkString{Data: "appendonly"},
					data.BulkString{Data: "no"},
				},
			},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "CONFIG"},
					data.BulkString{Data: "NEXIST"},
				},
			},
			data.Error{ErrMsg: "unsupported subcommand NEXIST for CONFIG"},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input-%s", tc.input.ToDataString()), func(t *testing.T) {
			result := handler.HandleCommand(tc.input, &storageEngine)
			assert.Equal(tc.want, result)
		})
	}
}

func TestHandleSetWithTimeOptions(t *testing.T) {
	assert := assert.New(t)

	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input                data.Array
		sleepDuration        time.Duration
		want                 data.Message
		immediateFetchKey    string
		immediateFetchResult string
		insertCurrentTime    bool
		currentTimeUnit      string
	}{
		{
			input: data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "timedEX"},
					data.BulkString{Data: "timedEXVal"},
					data.BulkString{Data: "EX"},
					data.BulkString{Data: "1"},
				},
			},
			sleepDuration:        1 * time.Second,
			want:                 data.SimpleString{Contents: "OK"},
			immediateFetchKey:    "timedEX",
			immediateFetchResult: "timedEXVal",
			insertCurrentTime:    false,
		},
		{
			input: data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "timedPX"},
					data.BulkString{Data: "timedPXVal"},
					data.BulkString{Data: "PX"},
					data.BulkString{Data: "5"},
				},
			},
			sleepDuration:        6 * time.Millisecond,
			want:                 data.SimpleString{Contents: "OK"},
			immediateFetchKey:    "timedPX",
			immediateFetchResult: "timedPXVal",
			insertCurrentTime:    false,
		},
		{
			input: data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "timedEXAT"},
					data.BulkString{Data: "timedEXATVal"},
					data.BulkString{Data: "EXAT"},
				},
			},
			sleepDuration:        1 * time.Second,
			want:                 data.SimpleString{Contents: "OK"},
			immediateFetchKey:    "timedEXAT",
			immediateFetchResult: "timedEXATVal",
			insertCurrentTime:    true,
			currentTimeUnit:      "s",
		},
		{
			input: data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "timedPXAT"},
					data.BulkString{Data: "timedPXATVal"},
					data.BulkString{Data: "PXAT"},
				},
			},
			sleepDuration:        5 * time.Millisecond,
			want:                 data.SimpleString{Contents: "OK"},
			immediateFetchKey:    "timedPXAT",
			immediateFetchResult: "timedPXATVal",
			insertCurrentTime:    true,
			currentTimeUnit:      "ms",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
			if tc.insertCurrentTime {
				timeToInsert := time.Now().Add(tc.sleepDuration)
				timeValueToInsert := timeToInsert.Unix()
				if tc.currentTimeUnit == "ms" {
					timeValueToInsert = timeToInsert.UnixMilli()
				}
				tc.input.Elements = append(tc.input.Elements, data.BulkString{Data: fmt.Sprintf("%d", timeValueToInsert)})
			}

			result := handler.HandleCommand(tc.input, &storageEngine)
			assert.Equal(tc.want, result)

			fetchCmd := data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "GET"},
					data.BulkString{Data: tc.immediateFetchKey},
				},
			}

			preWaitFetchResult := handler.HandleCommand(fetchCmd, &storageEngine)
			assert.Equal(data.BulkString{Data: tc.immediateFetchResult}, preWaitFetchResult)

			time.Sleep(tc.sleepDuration)

			// post wait, no data is expected
			postWaitFetchResult := handler.HandleCommand(fetchCmd, &storageEngine)
			assert.Equal(data.Null{}, postWaitFetchResult)
		})
	}
}

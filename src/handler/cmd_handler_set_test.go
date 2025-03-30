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

func TestHandleSetCommand(t *testing.T) {
	storageEngine := storage.NewMapStorageEngine()

	testCases := []struct {
		input data.Message
		want  data.Message
	}{
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'set' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
				},
			},
			data.Error{ErrMsg: "wrong number of arguments for 'set' command"},
		},
		{
			data.Array{
				Elements: []data.Message{
					data.BulkString{Data: "SET"},
					data.BulkString{Data: "testKey"},
					data.Integer{Value: 1},
				},
			},
			data.Error{ErrMsg: "invalid format for command"},
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
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.input.ToDataString(), func(t *testing.T) {
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

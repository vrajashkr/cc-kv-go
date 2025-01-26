package data_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
)

func TestProcessMessageStringWithInvalidStrings(t *testing.T) {
	testCases := []string{
		"",
		"hello",
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc), func(t *testing.T) {
			charsConsumedCount, _, err := data.ProcessMessageString(tc)
			assert.NotNil(err)
			assert.Equal(0, charsConsumedCount)
		})
	}
}

func TestProcessMessageStringWithSimpleString(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"+OK\r\n", "OK"},
		{"+Hello world 123\r\n", "Hello world 123"},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			charsConsumedCount, result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			sS, ok := result.(data.SimpleString)
			assert.True(ok)
			assert.Equal(tc.want, sS.Contents)
			assert.Equal(len(tc.input), charsConsumedCount)
		})
	}
}

func TestSimpleStringToDataString(t *testing.T) {
	testCases := []struct {
		want  string
		input string
	}{
		{"+OK\r\n", "OK"},
		{"+Hello world 123\r\n", "Hello world 123"},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			msg := data.SimpleString{
				Contents: tc.input,
			}
			assert.Equal(tc.want, msg.ToDataString())
		})
	}
}

func TestProcessMessageStringWithError(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"-Error Message\r\n", "Error Message"},
		{"-error123\r\n", "error123"},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			charsConsumedCount, result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			e, ok := result.(data.Error)
			assert.True(ok)
			assert.Equal(tc.want, e.ErrMsg)
			assert.Equal(len(tc.input), charsConsumedCount)
		})
	}
}

func TestErrorToDataString(t *testing.T) {
	testCases := []struct {
		want  string
		input string
	}{
		{"-Error Message\r\n", "Error Message"},
		{"-error123\r\n", "error123"},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			msg := data.Error{
				ErrMsg: tc.input,
			}
			assert.Equal(tc.want, msg.ToDataString())
		})
	}
}

func TestProcessMessageStringWithInteger(t *testing.T) {
	testCases := []struct {
		input string
		want  int64
	}{
		{":123\r\n", 123},
		{":+123\r\n", 123},
		{":-123\r\n", -123},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			charsConsumedCount, result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			e, ok := result.(data.Integer)
			assert.True(ok)
			assert.Equal(tc.want, e.Value)
			assert.Equal(len(tc.input), charsConsumedCount)
		})
	}
}

func TestProcessMessageStringWithIncorrectInteger(t *testing.T) {
	testCases := []string{
		":+nonsense\r\n",
		":??\r\n",
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc), func(t *testing.T) {
			charsConsumedCount, _, err := data.ProcessMessageString(tc)
			assert.NotNil(err)
			assert.Equal(0, charsConsumedCount)
		})
	}
}

func TestIntegerToDataString(t *testing.T) {
	testCases := []struct {
		want  string
		input int64
	}{
		{":123\r\n", 123},
		{":-123\r\n", -123},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %d", tc.input), func(t *testing.T) {
			msg := data.Integer{
				Value: tc.input,
			}
			assert.Equal(tc.want, msg.ToDataString())
		})
	}
}

func TestProcessMessageStringWithBulkString(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"$5\r\nhello\r\n", "hello"},
		{"$15\r\nHello world 123\r\n", "Hello world 123"},
		{"$0\r\n\r\n", ""},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			charsConsumedCount, result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			bS, ok := result.(data.BulkString)
			assert.True(ok)
			assert.Equal(tc.want, bS.Data)
			assert.Equal(len(tc.input), charsConsumedCount)
		})
	}
}

func TestBulkStringToDataString(t *testing.T) {
	testCases := []struct {
		want  string
		input string
	}{
		{"$5\r\nhello\r\n", "hello"},
		{"$15\r\nHello world 123\r\n", "Hello world 123"},
		{"$0\r\n\r\n", ""},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			msg := data.BulkString{
				Data: tc.input,
			}
			assert.Equal(tc.want, msg.ToDataString())
		})
	}
}

func TestProcessMessageStringWithIncorrectBulkString(t *testing.T) {
	testCases := []string{
		"$nonsense\r\n",
		"$\r\n",
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc), func(t *testing.T) {
			charsConsumedCount, _, err := data.ProcessMessageString(tc)
			assert.NotNil(err)
			assert.Equal(0, charsConsumedCount)
		})
	}
}

func TestProcessMessageStringWithArray(t *testing.T) {
	testCases := []struct {
		input string
		want  []data.Message
	}{
		{"*3\r\n:1\r\n:2\r\n:3\r\n", []data.Message{
			data.Integer{
				Value: 1,
			},
			data.Integer{
				Value: 2,
			},
			data.Integer{
				Value: 3,
			},
		},
		},
		{"*0\r\n", []data.Message{}},
		{"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []data.Message{
			data.BulkString{
				Data: "hello",
			},
			data.BulkString{
				Data: "world",
			},
		}},
		{"*3\r\n$5\r\nhello\r\n$5\r\nworld\r\n:67\r\n", []data.Message{
			data.BulkString{
				Data: "hello",
			},
			data.BulkString{
				Data: "world",
			},
			data.Integer{
				Value: 67,
			},
		}},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			charsConsumedCount, result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			bS, ok := result.(data.Array)
			assert.True(ok)
			assert.Equal(tc.want, bS.Elements)
			assert.Equal(len(tc.input), charsConsumedCount)
		})
	}
}

func TestArrayToDataString(t *testing.T) {
	testCases := []struct {
		want  string
		input []data.Message
	}{
		{"*3\r\n:1\r\n:2\r\n:3\r\n", []data.Message{
			data.Integer{
				Value: 1,
			},
			data.Integer{
				Value: 2,
			},
			data.Integer{
				Value: 3,
			},
		},
		},
		{"*0\r\n", []data.Message{}},
		{"*3\r\n$5\r\nhello\r\n$5\r\nworld\r\n:67\r\n", []data.Message{
			data.BulkString{
				Data: "hello",
			},
			data.BulkString{
				Data: "world",
			},
			data.Integer{
				Value: 67,
			},
		}},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			msg := data.Array{
				Elements: tc.input,
			}
			assert.Equal(tc.want, msg.ToDataString())
		})
	}
}

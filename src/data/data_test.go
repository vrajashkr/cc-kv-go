package data_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/data"
)

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
			result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			sS, ok := result.(data.SimpleString)
			assert.True(ok)
			assert.Equal(tc.want, sS.Contents)
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

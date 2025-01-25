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
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc), func(t *testing.T) {
			_, err := data.ProcessMessageString(tc)
			assert.NotNil(err)
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
			result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			e, ok := result.(data.Error)
			assert.True(ok)
			assert.Equal(tc.want, e.ErrMsg)
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
			result, err := data.ProcessMessageString(tc.input)
			assert.Nil(err)

			e, ok := result.(data.Integer)
			assert.True(ok)
			assert.Equal(tc.want, e.Value)
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
			_, err := data.ProcessMessageString(tc)
			assert.NotNil(err)
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

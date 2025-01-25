package data

import "fmt"

const (
	MSG_TYPE_SIMPLE_STR = '+'
	MSG_TYPE_ERROR      = '-'
	MSG_TYPE_INT        = ':'
	MSG_TYPE_BULK_STR   = '$'
	MSG_TYPE_ARRAY      = '*'
)

type Message interface {
	ToDataString() string
}

type SimpleString struct {
	Contents string
}

func NewSimpleString(rawMsg string) (SimpleString, error) {
	return SimpleString{
		Contents: rawMsg[1 : len(rawMsg)-2],
	}, nil
}

func (s SimpleString) ToDataString() string {
	return fmt.Sprintf("%c%s\r\n", MSG_TYPE_SIMPLE_STR, s.Contents)
}

func ProcessMessageString(msg string) (Message, error) {
	if len(msg) <= 1 {
		return nil, fmt.Errorf("received empty invalid message")
	}

	var convertedMsg Message
	var err error
	switch msg[0] {
	case MSG_TYPE_SIMPLE_STR:
		{
			convertedMsg, err = NewSimpleString(msg)
		}
	}
	return convertedMsg, err
}

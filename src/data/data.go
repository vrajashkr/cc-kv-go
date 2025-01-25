package data

import (
	"fmt"
	"strconv"
)

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

type Error struct {
	ErrMsg string
}

func NewError(rawMsg string) (Error, error) {
	return Error{
		ErrMsg: rawMsg[1 : len(rawMsg)-2],
	}, nil
}

func (e Error) ToDataString() string {
	return fmt.Sprintf("%c%s\r\n", MSG_TYPE_ERROR, e.ErrMsg)
}

type Integer struct {
	Value int64
}

func NewInteger(rawMsg string) (Integer, error) {
	numericString := rawMsg[1 : len(rawMsg)-2]
	val, err := strconv.Atoi(numericString)
	if err != nil {
		return Integer{}, err
	}
	return Integer{
		Value: int64(val),
	}, nil
}

func (i Integer) ToDataString() string {
	return fmt.Sprintf("%c%d\r\n", MSG_TYPE_INT, i.Value)
}

func ProcessMessageString(msg string) (Message, error) {
	if len(msg) <= 1 {
		return nil, fmt.Errorf("received empty invalid message")
	}

	var convertedMsg Message
	var err error
	switch msg[0] {
	case MSG_TYPE_SIMPLE_STR:
		convertedMsg, err = NewSimpleString(msg)
	case MSG_TYPE_ERROR:
		convertedMsg, err = NewError(msg)
	case MSG_TYPE_INT:
		convertedMsg, err = NewInteger(msg)
	}
	return convertedMsg, err
}

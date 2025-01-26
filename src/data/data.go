package data

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MSG_TYPE_SIMPLE_STR = '+'
	MSG_TYPE_ERROR      = '-'
	MSG_TYPE_INT        = ':'
	MSG_TYPE_BULK_STR   = '$'
	MSG_TYPE_ARRAY      = '*'
	MSG_NULL_W_BULK_STR = "$-1\r\n"
	MSG_NULL_W_ARRAY    = "*-1\r\n"
)

type Message interface {
	ToDataString() string
}

type SimpleString struct {
	Contents string
}

func NewSimpleString(rawMsg string) (int, SimpleString, error) {
	firstCRLFIndex := strings.Index(rawMsg, "\r\n")
	return firstCRLFIndex + 2, SimpleString{
		Contents: rawMsg[1:firstCRLFIndex],
	}, nil
}

func (s SimpleString) ToDataString() string {
	return fmt.Sprintf("%c%s\r\n", MSG_TYPE_SIMPLE_STR, s.Contents)
}

type Error struct {
	ErrMsg string
}

func NewError(rawMsg string) (int, Error, error) {
	firstCRLFIndex := strings.Index(rawMsg, "\r\n")
	return firstCRLFIndex + 2, Error{
		ErrMsg: rawMsg[1:firstCRLFIndex],
	}, nil
}

func (e Error) ToDataString() string {
	return fmt.Sprintf("%c%s\r\n", MSG_TYPE_ERROR, e.ErrMsg)
}

type Integer struct {
	Value int64
}

func NewInteger(rawMsg string) (int, Integer, error) {
	firstCRLFIndex := strings.Index(rawMsg, "\r\n")
	numericString := rawMsg[1:firstCRLFIndex]
	val, err := strconv.Atoi(numericString)
	if err != nil {
		return 0, Integer{}, err
	}
	return firstCRLFIndex + 2, Integer{
		Value: int64(val),
	}, nil
}

func (i Integer) ToDataString() string {
	return fmt.Sprintf("%c%d\r\n", MSG_TYPE_INT, i.Value)
}

type BulkString struct {
	Data string
}

func NewBulkString(rawMsg string) (int, BulkString, error) {
	firstCRLFIndex := strings.Index(rawMsg, "\r\n")
	if firstCRLFIndex < 2 {
		return 0, BulkString{}, fmt.Errorf("invalid message format for bulk string")
	}

	strLen, err := strconv.Atoi(rawMsg[1:firstCRLFIndex])
	if err != nil {
		return 0, BulkString{}, err
	}

	dataStartIdx := firstCRLFIndex + 2
	dataEndIdx := dataStartIdx + strLen
	return dataEndIdx + 2, BulkString{
		Data: rawMsg[dataStartIdx:dataEndIdx],
	}, nil
}

func (bs BulkString) ToDataString() string {
	return fmt.Sprintf("%c%d\r\n%s\r\n", MSG_TYPE_BULK_STR, len(bs.Data), bs.Data)
}

type Array struct {
	Elements []Message
}

func NewArray(rawMsg string) (int, Array, error) {
	firstCRIndex := strings.IndexByte(rawMsg, '\r')
	if firstCRIndex < 2 {
		return 0, Array{}, fmt.Errorf("invalid message format for array")
	}

	numElements, err := strconv.Atoi(rawMsg[1:firstCRIndex])
	if err != nil {
		return 0, Array{}, err
	}

	charConsumedCount := firstCRIndex + 2
	elements := []Message{}

	for i := 0; i < numElements; i++ {
		charsCount, msg, err := ProcessMessageString(rawMsg[charConsumedCount:])
		if err != nil {
			return 0, Array{}, err
		}
		charConsumedCount += charsCount
		elements = append(elements, msg)
	}

	return charConsumedCount, Array{
		Elements: elements,
	}, nil
}

func (a Array) ToDataString() string {
	out := fmt.Sprintf("%c%d\r\n", MSG_TYPE_ARRAY, len(a.Elements))
	for _, elem := range a.Elements {
		out += elem.ToDataString()
	}
	return out
}

type Null struct{}

func NewNull() (int, Null) {
	return len(MSG_NULL_W_BULK_STR), Null{}
}

func (n Null) ToDataString() string {
	return MSG_NULL_W_BULK_STR
}

func ProcessMessageString(msg string) (int, Message, error) {
	if len(msg) <= 1 {
		return 0, nil, fmt.Errorf("received empty invalid message")
	}

	var consumedCount int
	var convertedMsg Message
	var err error
	switch msg[0] {
	case MSG_TYPE_SIMPLE_STR:
		consumedCount, convertedMsg, err = NewSimpleString(msg)
	case MSG_TYPE_ERROR:
		consumedCount, convertedMsg, err = NewError(msg)
	case MSG_TYPE_INT:
		consumedCount, convertedMsg, err = NewInteger(msg)
	case MSG_TYPE_BULK_STR:
		if msg == MSG_NULL_W_BULK_STR {
			consumedCount, convertedMsg = NewNull()
		} else {
			consumedCount, convertedMsg, err = NewBulkString(msg)
		}
	case MSG_TYPE_ARRAY:
		if msg == MSG_NULL_W_ARRAY {
			consumedCount, convertedMsg = NewNull()
		} else {
			consumedCount, convertedMsg, err = NewArray(msg)
		}
	default:
		err = fmt.Errorf("unsupported message discriminator")
	}
	return consumedCount, convertedMsg, err
}

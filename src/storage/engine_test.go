package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestMapStorageEngine(t *testing.T) {
	assert := assert.New(t)
	mse := storage.NewMapStorageEngine()
	err := mse.Set("hello", "world")
	assert.Nil(err)

	ok, val, err := mse.Get("hello")
	assert.True(ok)
	assert.Equal("world", val)
	assert.Nil(err)

	ok, val, err = mse.Get("absent")
	assert.False(ok)
	assert.Empty(val)
	assert.Nil(err)
}

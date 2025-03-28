package storage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vrajashkr/cc-kv-go/src/storage"
)

func TestMapStorageEngine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	mse := storage.NewMapStorageEngine()
	err := mse.Set("hello", "world", false, 10202)
	require.Nil(err)

	ok, val, err := mse.Get("hello")
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("world", val)

	ok, val, err = mse.Get("absent")
	assert.Nil(err)
	assert.False(ok)
	assert.Empty(val)

	err = mse.Set("timing", "result", true, time.Now().UnixMilli()+5)
	require.Nil(err)

	ok, val, err = mse.Get("timing")
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("result", val)

	time.Sleep(6 * time.Millisecond)
	ok, val, err = mse.Get("timing")
	assert.Nil(err)
	assert.False(ok)
	assert.Empty(val)
}

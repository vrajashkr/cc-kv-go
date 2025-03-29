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

	err = mse.Set("fakecounter", "world", false, 10202)
	require.Nil(err)

	ok, val, err := mse.Get("hello")
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("world", val)

	ok, val, err = mse.Get("absent")
	assert.Nil(err)
	assert.False(ok)
	assert.Empty(val)

	res, err := mse.Exists([]string{"hello", "absent", "hello", "absent"})
	assert.Nil(err)
	assert.Equal(2, res)

	res, err = mse.Delete([]string{"hello", "absent"})
	assert.Nil(err)
	assert.Equal(1, res)

	deltaResult, err := mse.AtomicDelta("ctr1", 1)
	assert.Nil(err)
	assert.Equal(int64(1), deltaResult)

	deltaResult, err = mse.AtomicDelta("ctr2", -2)
	assert.Nil(err)
	assert.Equal(int64(-2), deltaResult)

	deltaResult, err = mse.AtomicDelta("ctr1", 2)
	assert.Nil(err)
	assert.Equal(int64(3), deltaResult)

	_, err = mse.AtomicDelta("fakecounter", 1)
	assert.NotNil(err)

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

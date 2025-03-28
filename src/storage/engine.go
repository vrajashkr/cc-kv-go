package storage

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type StorageEngine interface {
	Set(key string, value string, expires bool, expiresAtTimeStampMillis int64) error
	Get(key string) (bool, string, error)
}

type DataContainer struct {
	Data      string
	Expires   bool
	ExpiresAt time.Time
}

type MapStorageEngine struct {
	store map[string]DataContainer
	mu    sync.Mutex
}

func NewMapStorageEngine() MapStorageEngine {
	return MapStorageEngine{
		store: make(map[string]DataContainer),
	}
}

func (mse *MapStorageEngine) Set(key string, value string, expires bool, expiresAtTimeStampMillis int64) error {
	mse.mu.Lock()
	defer mse.mu.Unlock()
	expiryTime := time.Now()
	slog.Info(fmt.Sprintf("current time: %d expiresAt: %d", expiryTime.UnixMilli(), expiresAtTimeStampMillis))
	if expires {
		expiryTime = time.UnixMilli(expiresAtTimeStampMillis)
	}

	mse.store[key] = DataContainer{
		Data:      value,
		Expires:   expires,
		ExpiresAt: expiryTime,
	}
	return nil
}

func (mse *MapStorageEngine) Get(key string) (bool, string, error) {
	mse.mu.Lock()
	defer mse.mu.Unlock()
	result, ok := mse.store[key]
	if !ok {
		return false, "", nil
	}

	if result.Expires && time.Since(result.ExpiresAt).Milliseconds() >= 0 {
		// entry has expired, delete from map
		delete(mse.store, key)
		return false, "", nil
	}

	return true, result.Data, nil
}

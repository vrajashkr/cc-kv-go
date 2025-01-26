package storage

import "sync"

type StorageEngine interface {
	Set(key string, value string) error
	Get(key string) (bool, string, error)
}

type MapStorageEngine struct {
	store map[string]string
	mu    sync.Mutex
}

func NewMapStorageEngine() MapStorageEngine {
	return MapStorageEngine{
		store: make(map[string]string),
	}
}

func (mse *MapStorageEngine) Set(key string, value string) error {
	mse.mu.Lock()
	defer mse.mu.Unlock()
	mse.store[key] = value
	return nil
}

func (mse *MapStorageEngine) Get(key string) (bool, string, error) {
	mse.mu.Lock()
	defer mse.mu.Unlock()
	result, ok := mse.store[key]
	if !ok {
		return false, "", nil
	}
	return true, result, nil
}

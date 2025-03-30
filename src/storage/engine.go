package storage

import (
	"time"
)

type StorageEngine interface {
	Set(key string, value string, expires bool, expiresAtTimeStampMillis int64) error
	Get(key string) (bool, string, error)
	Exists(keys []string) (int, error)
	Delete(keys []string) (int, error)
	AtomicDelta(key string, delta int64) (int64, error)
	ListPush(key string, values []string, isPrepend bool) (int64, error)
}

type DataContainer struct {
	Data      string
	Expires   bool
	ExpiresAt time.Time
}

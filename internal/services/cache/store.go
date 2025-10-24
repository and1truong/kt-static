package cache

import (
	"context"
	"time"
)

var (
	globalStore Store
	storeSet    bool
)

// Store is an interface for persistent cache storage.
type Store interface {
	// Read returns the raw data and the expiration time.
	Read(ctx context.Context, key string) (data []byte, expiry time.Time, err error)
	// Write stores the raw data and its expiration time.
	Write(ctx context.Context, key string, data []byte, expiry time.Time) error
}

// SetStore sets the global cache store for the application.
// It can only be called once. Subsequent calls will panic.
func SetStore(store Store) {
	if storeSet {
		panic("cache store is already set and cannot be set again")
	}
	globalStore = store
	storeSet = true
}

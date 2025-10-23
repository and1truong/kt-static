package util

import (
	"context"
	"sync"
	"time"
)

// inMemoryCacheEntry is the structure stored in the in-memory cache store.
type inMemoryCacheEntry struct {
	Expiry time.Time
	Data   []byte
}

// InMemoryStore implements CacheStore using an in-memory map.
// This is primarily for testing or scenarios where persistence is not required
// but the CacheStore interface needs to be satisfied.
type InMemoryStore struct {
	store map[string]inMemoryCacheEntry
	mu    sync.RWMutex
}

// NewInMemoryStore creates a new in-memory cache store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]inMemoryCacheEntry),
	}
}

// Read retrieves the raw data and expiration time from the in-memory store.
func (m *InMemoryStore) Read(ctx context.Context, key string) (data []byte, expiry time.Time, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.store[key]
	if !ok {
		return nil, time.Time{}, nil // Cache miss
	}

	return entry.Data, entry.Expiry, nil
}

// Write stores the raw data and its expiration time in the in-memory store.
func (m *InMemoryStore) Write(ctx context.Context, key string, data []byte, expiry time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = inMemoryCacheEntry{
		Expiry: expiry,
		Data:   data,
	}

	return nil
}

var (
	defaultInMemoryWriter     CacheStore
	defaultInMemoryWriterOnce sync.Once
)

// DefaultInMemoryStore returns a singleton instance of InMemoryStore.
// This is used as the default persistent store when none is configured.
func DefaultInMemoryStore() CacheStore {
	defaultInMemoryWriterOnce.Do(func() {
		defaultInMemoryWriter = NewInMemoryStore()
	})
	return defaultInMemoryWriter
}

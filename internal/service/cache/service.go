package cache

import (
	"context"
	"encoding/json"
	"time"
)

var (
	globalStore Store
	storeSet    bool
)

// SetStore sets the global cache store for the application.
// It can only be called once. Subsequent calls will panic.
func SetStore(store Store) {
	if storeSet {
		panic("cache store is already set and cannot be set again")
	}
	globalStore = store
	storeSet = true
}

// --- In-Memory Cache Implementation ---

// --- Cache Options ---

type cacheConfig struct {
	ttl     time.Duration
	noCache bool
	store   Store // Field for persistent storage
}

type CacheOption func(*cacheConfig)

// WithTTL sets the time-to-live for the cached entry.
func WithTTL(ttl time.Duration) CacheOption {
	return func(c *cacheConfig) {
		c.ttl = ttl
	}
}

// WithNoCache forces the value-generating function to be called, bypassing the cache read.
// The result will still be written to the cache unless the value-generating function returns an error.
func WithNoCache() CacheOption {
	return func(c *cacheConfig) {
		c.noCache = true
	}
}

// WithCacheStore sets a custom Store for persistent storage.
func WithCacheStore(store Store) CacheOption {
	return func(c *cacheConfig) {
		c.store = store
	}
}

// Cache is a generic function to process and cache the result of a value-generating function.
// It now relies solely on the configured Store (which defaults to an in-memory singleton)
// for all cache read and write operations.
func Cache[T any](ctx context.Context, key string, generateValue func() (T, error), options ...CacheOption) (T, error) {
	var zero T

	config := cacheConfig{}
	for _, opt := range options {
		opt(&config)
	}

	// Ensure a Store is always configured.
	if config.store == nil {
		if globalStore != nil {
			config.store = globalStore
		} else {
			// Fallback to the default in-memory store if no global store is set.
			config.store = DefaultInMemoryStore()
		}
	}

	// 1. Try to read from cache (via Store)
	if !config.noCache {
		data, expiry, err := config.store.Read(ctx, key)
		if err != nil {
			// In a real app, log this error.
		} else if data != nil {
			// Found in cache
			if expiry.IsZero() || expiry.After(time.Now()) {
				// Not expired, deserialize
				var value T
				if err := json.Unmarshal(data, &value); err == nil {
					return value, nil // Cache hit
				}
				// If deserialization fails, treat as miss and continue to generate
			}
		}
	}

	// 2. Generate new value
	value, err := generateValue()
	if err != nil {
		return zero, err
	}

	// 3. Write to cache (via Store)
	// Serialize value
	data, err := json.Marshal(value)
	if err == nil {
		var expiry time.Time
		if config.ttl > 0 {
			expiry = time.Now().Add(config.ttl)
		}

		// Write to cache
		if writeErr := config.store.Write(ctx, key, data, expiry); writeErr != nil {
			// In a real app, log this error.
		}
	}
	// Handle serialization error (e.g., log it)

	return value, nil
}

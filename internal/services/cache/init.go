package cache

import (
	"fmt"
	"time"

	"htruong/kt-crawler/internal"
)

// InitStore initializes a cache store based on the provided configuration.
// It returns the initialized Store and any error encountered.
func InitStore(cfg internal.CacheStoreConfig) (Store, error) {
	// Parse Expiry duration, although it's not used for store initialization,
	// it's part of the config and might be used later.
	if _, err := time.ParseDuration(cfg.Expiry); err != nil {
		return nil, fmt.Errorf("invalid cache expiry duration '%s': %w", cfg.Expiry, err)
	}

	switch cfg.Backend {
	case "filesystem":
		return NewFileSystemStore(cfg.Filesystem.Directory)
	case "memory":
		return NewInMemoryStore(), nil
	default:
		return nil, fmt.Errorf("unsupported cache backend: %s", cfg.Backend)
	}
}

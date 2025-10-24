package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// --- Persistent Store Interface and Implementation ---

// fileCacheEntry is the structure stored in the file system cache.
type fileCacheEntry struct {
	Expiry time.Time `json:"expiry"`
	Data   []byte    `json:"data"`
}

const defaultCacheDir = "/tmp/cache/"

// FileSystemStore implements Store using the local file system.
type FileSystemStore struct {
	dir string
}

// NewFileSystemStore creates a new file system cache store.
// If dir is empty, it defaults to /tmp/cache/.
func NewFileSystemStore(dir string) (*FileSystemStore, error) {
	if dir == "" {
		dir = defaultCacheDir
	}
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory %s: %w", dir, err)
	}
	return &FileSystemStore{dir: dir}, nil
}

func (f *FileSystemStore) Read(ctx context.Context, key string) (data []byte, expiry time.Time, err error) {
	filePath := filepath.Join(f.dir, key)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, time.Time{}, nil // Cache miss, not an error
		}
		return nil, time.Time{}, fmt.Errorf("failed to read cache file %s: %w", filePath, err)
	}

	var entry fileCacheEntry
	if err := json.Unmarshal(fileData, &entry); err != nil {
		// Corrupted file, delete it and treat as miss
		os.Remove(filePath)
		return nil, time.Time{}, fmt.Errorf("corrupted cache file for key %s: %w", key, err)
	}

	return entry.Data, entry.Expiry, nil
}

func (f *FileSystemStore) Write(ctx context.Context, key string, data []byte, expiry time.Time) error {
	filePath := filepath.Join(f.dir, key)

	entry := fileCacheEntry{
		Expiry: expiry,
		Data:   data,
	}

	fileData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry for key %s: %w", key, err)
	}

	// Write directly to the file path. Note: This is not atomic.
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file %s: %w", filePath, err)
	}

	return nil
}

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

// Store is an interface for persistent cache storage.
type Store interface {
	// Read returns the raw data and the expiration time.
	Read(ctx context.Context, key string) (data []byte, expiry time.Time, err error)
	// Write stores the raw data and its expiration time.
	Write(ctx context.Context, key string, data []byte, expiry time.Time) error
}

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
func NewFileSystemStore(dir string) *FileSystemStore {
	if dir == "" {
		dir = defaultCacheDir
	}
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		// Note: In a real application, this error should be logged or returned.
		// For this utility, we proceed and let subsequent file operations fail if the directory is unusable.
	}
	return &FileSystemStore{dir: dir}
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

	// Write to a temporary file first and then rename for atomicity
	tmpFile, err := os.CreateTemp(f.dir, key+".tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file for key %s: %w", key, err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file on error

	if _, err := tmpFile.Write(fileData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temp file for key %s: %w", key, err)
	}
	tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), filePath); err != nil {
		return fmt.Errorf("failed to rename temp file to %s: %w", filePath, err)
	}

	return nil
}

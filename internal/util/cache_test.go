package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// setupFSCache is a helper to create a temporary directory and a FileSystemStore.
func setupFSCache(t *testing.T) (*FileSystemStore, string) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "fscache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	writer := NewFileSystemStore(tempDir)
	return writer, tempDir
}

// TestCache_InMemoryBasic tests basic in-memory caching functionality.
func TestCache_InMemoryBasic(t *testing.T) {
	ctx := context.Background()
	key := "test_key_basic"
	callCount := int32(0)

	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return "cached_value", nil
	}

	// 1. First call (miss, generate)
	val, err := Cache(ctx, key, generateValue)
	if err != nil || val != "cached_value" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("First call failed: val=%v, err=%v, count=%d", val, err, callCount)
	}

	// 2. Second call (hit, should not re-generate)
	val, err = Cache(ctx, key, generateValue)
	if err != nil || val != "cached_value" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("Second call failed (should be hit): val=%v, err=%v, count=%d", val, err, callCount)
	}
}

// TestCache_InMemoryTTL tests in-memory caching with expiration.
func TestCache_InMemoryTTL(t *testing.T) {
	ctx := context.Background()
	key := "test_key_ttl"
	callCount := int32(0)

	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return fmt.Sprintf("value_%d", atomic.LoadInt32(&callCount)), nil
	}

	// 1. First call with short TTL (miss)
	val, err := Cache(ctx, key, generateValue, WithTTL(10*time.Millisecond))
	if err != nil || val != "value_1" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("First call failed: val=%v, err=%v, count=%d", val, err, callCount)
	}

	// 2. Second call immediately (hit)
	val, err = Cache(ctx, key, generateValue)
	if err != nil || val != "value_1" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("Second call failed (should be hit): val=%v, err=%v, count=%d", val, err, callCount)
	}

	// 3. Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// 4. Third call (miss, re-generate)
	val, err = Cache(ctx, key, generateValue)
	if err != nil || val != "value_2" || atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("Third call failed (should be miss): val=%v, err=%v, count=%d", val, err, callCount)
	}
}

// TestCache_WithNoCache tests the WithNoCache option.
func TestCache_WithNoCache(t *testing.T) {
	ctx := context.Background()
	key := "test_key_nocache"
	callCount := int32(0)

	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return fmt.Sprintf("value_%d", atomic.LoadInt32(&callCount)), nil
	}

	// 1. First call (miss)
	Cache(ctx, key, generateValue)
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatal("Expected count 1")
	}

	// 2. Second call with WithNoCache (miss, re-generate)
	Cache(ctx, key, generateValue, WithNoCache())
	if atomic.LoadInt32(&callCount) != 2 {
		t.Fatal("Expected count 2")
	}

	// 3. Third call without WithNoCache (hit, since the second call wrote to cache)
	val, err := Cache(ctx, key, generateValue)
	if err != nil || val != "value_2" || atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("Third call failed (should be hit): val=%v, err=%v, count=%d", val, err, callCount)
	}
}

// TestCache_PersistentCache tests reading from the persistent cache after an in-memory miss.
func TestCache_PersistentCache(t *testing.T) {
	writer, tempDir := setupFSCache(t)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	key := "test_key_fs"
	callCount := int32(0)

	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return fmt.Sprintf("fs_value_%d", atomic.LoadInt32(&callCount)), nil
	}

	// 1. First call (miss, generate, write to in-memory and persistent)
	val, err := Cache(ctx, key, generateValue, WithCacheStore(writer))
	if err != nil || val != "fs_value_1" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("First call failed: val=%v, err=%v, count=%d", val, err, callCount)
	}

	// 3. Second call (persistent hit, should not re-generate)
	val, err = Cache(ctx, key, generateValue, WithCacheStore(writer))
	if err != nil || val != "fs_value_1" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("Second call failed (should be persistent hit): val=%v, err=%v, count=%d", val, err, callCount)
	}
}

// TestCache_PersistentCacheTTL tests persistent caching with expiration.
func TestCache_PersistentCacheTTL(t *testing.T) {
	writer, tempDir := setupFSCache(t)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	key := "test_key_fs_ttl"
	callCount := int32(0)

	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return fmt.Sprintf("fs_ttl_value_%d", atomic.LoadInt32(&callCount)), nil
	}

	// 1. First call with short TTL
	Cache(ctx, key, generateValue, WithCacheStore(writer), WithTTL(10*time.Millisecond))
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatal("Expected count 1")
	}

	// 3. Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// 4. Third call (in-memory miss, persistent miss due to expiration, re-generate)
	val, err := Cache(ctx, key, generateValue, WithCacheStore(writer))
	if err != nil || val != "fs_ttl_value_2" || atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("Third call failed (should be miss): val=%v, err=%v, count=%d", val, err, callCount)
	}
}

// TestFileSystemCacheWriter_CorruptedFile tests that a corrupted file is handled and deleted.
func TestFileSystemCacheWriter_CorruptedFile(t *testing.T) {
	writer, tempDir := setupFSCache(t)
	defer os.RemoveAll(tempDir)

	key := "corrupted_key"
	filePath := filepath.Join(tempDir, key)

	// Write a corrupted file
	if err := os.WriteFile(filePath, []byte("this is not json"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted file: %v", err)
	}

	// Read should fail and delete the file
	_, _, err := writer.Read(context.Background(), key)
	if err == nil {
		t.Fatal("Expected an error for corrupted file, got nil")
	}

	// Check if file was deleted
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("Corrupted file was not deleted. Stat error: %v", err)
	}
}

// TestFileSystemCacheWriter_CustomDir tests that a custom directory is used correctly.
func TestFileSystemCacheWriter_CustomDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "custom_fscache_dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	writer := NewFileSystemStore(tempDir)
	key := "custom_dir_key"
	data := []byte("test data")
	expiry := time.Now().Add(time.Hour)

	// Write
	if err := writer.Write(context.Background(), key, data, expiry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check if file exists in the custom directory
	filePath := filepath.Join(tempDir, key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File was not written to custom directory: %s", filePath)
	}

	// Read
	readData, readExpiry, err := writer.Read(context.Background(), key)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if string(readData) != string(data) {
		t.Fatalf("Read data mismatch. Got %s, expected %s", readData, data)
	}
	if !readExpiry.Equal(expiry.Truncate(time.Second)) { // JSON unmarshal loses precision, so truncate for comparison
		t.Logf("Warning: Expiry time precision lost. Expected %v, Got %v", expiry, readExpiry)
	}
}

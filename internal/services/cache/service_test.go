package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// mockStore is a simple implementation of the Store interface for testing.
type mockStore struct {
	readCount  int32
	writeCount int32
}

func (m *mockStore) Read(ctx context.Context, key string) (data []byte, expiry time.Time, err error) {
	atomic.AddInt32(&m.readCount, 1)
	return nil, time.Time{}, nil // Always a miss for simplicity
}

func (m *mockStore) Write(ctx context.Context, key string, data []byte, expiry time.Time) error {
	atomic.AddInt32(&m.writeCount, 1)
	return nil
}

// setupFSCache is a helper to create a temporary directory and a FileSystemStore.
func setupFSCache(t *testing.T) (*FileSystemStore, string) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "fscache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	writer, err := NewFileSystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileSystemStore: %v", err)
	}
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
	_, err := Cache(ctx, key, generateValue)
	if err != nil {
		t.Fatalf("Cache failed on first call: %v", err)
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatal("Expected count 1")
	}

	// 2. Second call with WithNoCache (miss, re-generate)
	_, err = Cache(ctx, key, generateValue, WithNoCache())
	if err != nil {
		t.Fatalf("Cache failed on second call: %v", err)
	}
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
	val, err := Cache(ctx, key, generateValue, WithStore(writer))
	if err != nil || val != "fs_value_1" || atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("First call failed: val=%v, err=%v, count=%d", val, err, callCount)
	}

	// 3. Second call (persistent hit, should not re-generate)
	val, err = Cache(ctx, key, generateValue, WithStore(writer))
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
	_, err := Cache(ctx, key, generateValue, WithStore(writer), WithTTL(10*time.Millisecond))
	if err != nil {
		t.Fatalf("Cache failed on first call: %v", err)
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatal("Expected count 1")
	}

	// 3. Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// 4. Third call (in-memory miss, persistent miss due to expiration, re-generate)
	val, err := Cache(ctx, key, generateValue, WithStore(writer))
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

	writer, err := NewFileSystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileSystemStore: %v", err)
	}
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
	// JSON unmarshal loses nanosecond precision (retains microsecond precision).
	// Truncate the expected time to microseconds for a valid comparison.
	expectedExpiry := expiry.Truncate(time.Microsecond)
	if !readExpiry.Equal(expectedExpiry) {
		t.Fatalf("Read expiry time mismatch. Expected %v, Got %v", expectedExpiry, readExpiry)
	}
}

// TestSetStore_SetOnceAndUsage tests that SetStore works, is used by Cache, and panics on second call.
func TestSetStore_SetOnceAndUsage(t *testing.T) {
	// 1. Test SetStore usage in Cache
	mock := &mockStore{}

	// Check if the store is already set by a previous test. If not, set it.
	// This is necessary because SetStore modifies global state and panics on subsequent calls.
	if !storeSet {
		SetStore(mock)
	} else {
		t.Log("Global store already set by another test. Skipping SetStore call for usage test.")
	}

	ctx := context.Background()
	key := "test_global_store"
	callCount := int32(0)
	generateValue := func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return "value", nil
	}

	// First call to Cache. Should use the global store (mock).
	_, err := Cache(ctx, key, generateValue)
	if err != nil {
		t.Fatalf("Cache failed: %v", err)
	}

	// Check if the mock store was used for reading (it should be a miss, so read is called once)
	if atomic.LoadInt32(&mock.readCount) != 1 {
		t.Errorf("Expected mock store Read to be called 1 time, got %d", mock.readCount)
	}
	// Check if the mock store was used for writing
	if atomic.LoadInt32(&mock.writeCount) != 1 {
		t.Errorf("Expected mock store Write to be called 1 time, got %d", mock.writeCount)
	}

	// 2. Test SetStore panic on second call
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected SetStore to panic on second call, but it did not")
		}
	}()

	// Attempt to set the store again (should panic)
	SetStore(&mockStore{})
}

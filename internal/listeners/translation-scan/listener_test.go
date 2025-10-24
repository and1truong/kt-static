package translation_scan

import (
	"context"
	"htruong/kt-crawler/internal/services/cache"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/fetch"
	"htruong/kt-crawler/internal/services/logging/logging"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawTranslation(t *testing.T) {
	ctx := context.Background()
	url := "https://kinhthanh.httlvn.org/?v=VI1934"

	// Setup dependencies
	logger := &logging.MockLogger{}
	dispatcher := eventdispatcher.NewDispatcher(logger)
	config := fetch.FetchConfig{
		Timeout: "5s",
		Cache:   false, // Disable cache for this test
	}
	cacheStore := cache.NewInMemoryStore()

	listener := NewTranslationScanListener(dispatcher, config, cacheStore)
	event := NewTranslationScanEvent(url)

	err := listener.Handle(ctx, event)
	assert.NoError(t, err)
}

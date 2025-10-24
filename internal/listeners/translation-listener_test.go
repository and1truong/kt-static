package listeners

import (
	"context"
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
	config := fetch.Config{
		Timeout: "5s",
	}

	listener := NewTranslationScanListener(dispatcher, config)
	event := NewTranslationScanEvent(url, "VI1934")

	err := listener.Handle(ctx, event)
	assert.NoError(t, err)
}

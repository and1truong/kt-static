package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/listeners/book-scan"
	"htruong/kt-crawler/internal/listeners/chapter-listener"
	"htruong/kt-crawler/internal/listeners/store-listener"
	"htruong/kt-crawler/internal/listeners/translation-scan"
	"htruong/kt-crawler/internal/services/cache"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/logging"
)

func loadConfig(path string) (*internal.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Start with default config, then unmarshal the file content over it.
	config := internal.NewDefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the JSON configuration file")
	flag.Parse()

	var config *internal.Config
	if configPath != "" {
		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
		}
	} else {
		config = internal.NewDefaultConfig()
	}

	var (
		ctx        = context.Background()
		logger     = logging.NewLogger()
		dispatcher = eventdispatcher.NewDispatcher(logger)
	)

	// Initialize cache store
	cacheStore, err := cache.NewFileSystemStore(config.Fetch.CacheDir)
	if err != nil {
		log.Fatalf("Failed to initialize cache store: %v", err)
	}
	
	var (
		translationListener = translation_scan.NewTranslationScanListener(dispatcher, config.Fetch, cacheStore)
		bookListener        = book_scan.NewBookScanListener(dispatcher)
		chapterListener     = chapter_listener.NewChapterScanListener(dispatcher, config.Fetch, cacheStore)
		storeListener       = store_listener.NewStoreListener(config.Listeners.Store, logger)
	)

	// Register Listeners
	// ---------------------
	{
		dispatcher.Register(book_scan.BookScanEventName, bookListener)
		dispatcher.Register(translation_scan.TranslationScanEventName, translationListener)
		dispatcher.Register(chapter_listener.ChapterScanEventName, chapterListener)
		dispatcher.Register(store_listener.StoreEventName, storeListener)
	}

	// Start the crawling process by dispatching the initial event.
	// ---------------------
	{
		initialEvent := translation_scan.NewTranslationScanEvent(config.InitialURL)
		if err := dispatcher.Dispatch(ctx, initialEvent); err != nil {
			log.Fatalf("Initial dispatch failed: %v", err)
		}
	}

	logger.Info("crawler finished successfully.")
}

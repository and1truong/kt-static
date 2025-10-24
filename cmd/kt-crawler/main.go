package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/listeners"
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

	var (
		ctx        = context.Background()
		logger     = logging.NewLogger()
		dispatcher = eventdispatcher.NewDispatcher(logger)
		config     *internal.Config
	)

	if configPath != "" {
		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
		}
	} else {
		config = internal.NewDefaultConfig()
	}

	// Initialize Cache Store
	store, err := cache.InitStore(config.Cache.Default.Store)
	if err != nil {
		log.Fatalf("Failed to initialize cache store: %v", err)
	}
	cache.SetStore(store)

	var (
		translationListener = listeners.NewTranslationScanListener(dispatcher, config.Fetch)
		bookListener        = listeners.NewBookScanListener(dispatcher)
		chapterListener     = listeners.NewChapterScanListener(dispatcher, config.Fetch)
		storeListener       = listeners.NewStoreListener(config.Listeners.Store, logger)
	)

	// Register Listeners
	// ---------------------
	{
		dispatcher.Register(listeners.BookScanEventName, bookListener)
		dispatcher.Register(listeners.TranslationScanEventName, translationListener)
		dispatcher.Register(listeners.ChapterScanEventName, chapterListener)
		dispatcher.Register(listeners.StoreEventName, storeListener)
	}

	// Start the crawling process by dispatching the initial event.
	// ---------------------
	{
		parsedURL, err := url.Parse(config.InitialURL)
		if err != nil {
			log.Fatalf("Invalid initial URL: %v", err)
		}
		translationCode := parsedURL.Query().Get("v")
		if translationCode == "" {
			log.Fatalf("Translation code (v) not found in the initial URL")
		}

		initialEvent := listeners.NewTranslationScanEvent(config.InitialURL, translationCode)
		if err := dispatcher.Dispatch(ctx, initialEvent); err != nil {
			log.Fatalf("Initial dispatch failed: %v", err)
		}
	}

	logger.Info("crawler finished successfully.")
}

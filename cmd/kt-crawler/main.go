package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"

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
	app := &cli.App{
		Name:  "kt-crawler",
		Usage: "A crawler for the kt website",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to the JSON configuration file",
				Value:   "config.sample.json",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "scan",
				Usage:  "Starts the crawling process",
				Action: runScan,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runScan(c *cli.Context) error {
	configPath := c.String("config")

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
			return cli.Exit(fmt.Sprintf("Failed to load config from %s: %v", configPath, err), 1)
		}
	} else {
		config = internal.NewDefaultConfig()
	}

	// Initialize Cache Store
	store, err := cache.InitStore(config.Cache.Default.Store)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize cache store: %v", err), 1)
	}
	cache.SetStore(store)

	var (
		translationListener = listeners.NewTranslationScanListener(dispatcher, config.Fetch)
		bookListener        = listeners.NewBookScanListener(dispatcher, config.Listeners.Store, logger)
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
			return cli.Exit(fmt.Sprintf("Invalid initial URL: %v", err), 1)
		}
		translationCode := parsedURL.Query().Get("v")
		if translationCode == "" {
			return cli.Exit("Translation code (v) not found in the initial URL", 1)
		}

		initialEvent := listeners.NewTranslationScanEvent(config.InitialURL, translationCode)
		if err := dispatcher.Dispatch(ctx, initialEvent); err != nil {
			return cli.Exit(fmt.Sprintf("Initial dispatch failed: %v", err), 1)
		}
	}

	logger.Info("crawler finished successfully.")

	return nil
}

package internal

import "htruong/kt-crawler/internal/services/fetch"

// NewDefaultConfig returns a Config struct with default values.
func NewDefaultConfig() *Config {
	return &Config{
		InitialURL: "https://kinhthanh.httlvn.org/?v=VI1934",
		Fetch: fetch.Config{
			Timeout:         "30s",
			Cache:           true,
			CacheDir:        ".cache",
			CacheTtlSeconds: 86400, // 24 hours
		},
		Listeners: ListenersConfig{
			Store: StoreConfig{
				Filesystem: FilesystemConfig{
					Directory: "/tmp/output/chapters",
				},
				DryRun: false,
			},
		},
	}
}

// Config holds the application configuration.
type Config struct {
	InitialURL string          `json:"initialURL"`
	Listeners  ListenersConfig `json:"listeners"`
	Fetch      fetch.Config    `json:"fetch"`
}

type (
	FilesystemConfig struct {
		Directory string `json:"directory"`
	}

	StoreConfig struct {
		Filesystem FilesystemConfig `json:"filesystem"`
		DryRun     bool             `json:"dry"`
	}

	ListenersConfig struct {
		Store StoreConfig `json:"store"`
	}
)

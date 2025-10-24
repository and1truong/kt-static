package internal

import "htruong/kt-crawler/internal/services/fetch"

// NewDefaultConfig returns a Config struct with default values.
func NewDefaultConfig() *Config {
	return &Config{
		InitialURL: "https://kinhthanh.httlvn.org/?v=VI1934",
		Fetch: fetch.Config{
			Timeout: "30s",
		},
		Listeners: ListenersConfig{
			Store: StoreConfig{
				Filesystem: FilesystemConfig{
					Directory: "/tmp/output/chapters",
				},
				DryRun: false,
			},
		},
		Cache: CacheConfig{
			Default: DefaultCacheConfig{
				Store: CacheStoreConfig{
					Backend: "filesystem",
					Expiry:  "24h",
					Filesystem: FilesystemConfig{
						Directory: "/tmp/kt-crawler/cache",
					},
				},
			},
		},
	}
}

// Config holds the application configuration.
type Config struct {
	InitialURL string          `json:"initialURL"`
	Listeners  ListenersConfig `json:"listeners"`
	Fetch      fetch.Config    `json:"fetch"`
	Cache      CacheConfig     `json:"cache"`
}

type CacheConfig struct {
	Default DefaultCacheConfig `json:"default"`
}

type DefaultCacheConfig struct {
	Store CacheStoreConfig `json:"store"`
}

type CacheStoreConfig struct {
	Backend    string           `json:"backend"`
	Expiry     string           `json:"expiry"`
	Filesystem FilesystemConfig `json:"filesystem"`
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

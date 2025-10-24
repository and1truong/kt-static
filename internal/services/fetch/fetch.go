package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Config holds configuration for a fetch operation.
type Config struct {
	Timeout string `json:"timeout"`
}

// Fetcher is the service responsible for fetching content.
type Fetcher struct {
	config Config
}

// NewFetcher creates a new Fetcher instance.
func NewFetcher(config Config) *Fetcher {
	return &Fetcher{
		config: config,
	}
}

// Config returns the configuration the Fetcher was initialized with.
func (f *Fetcher) Config() Config {
	return f.config
}

// --- Options Pattern ---

type options struct {
	config *Config
}

// Option is a function that configures a Fetch operation.
type Option func(*options)

// WithConfig provides a custom Config for the operation.
func WithConfig(cfg Config) Option {
	return func(opts *options) {
		opts.config = &cfg
	}
}

// Fetch fetches content from a requestPath, optionally using a custom timeout.
//
// Simple call: f.Fetch(ctx, path) uses the Fetcher's default timeout.
// Configured call: f.Fetch(ctx, path, WithConfig(config)) uses the provided config.
func (f *Fetcher) Fetch(ctx context.Context, path string, opts ...Option) ([]byte, error) {
	config := f.config

	fetchOpts := &options{}
	for _, opt := range opts {
		opt(fetchOpts)
	}

	if fetchOpts.config != nil {
		config = *fetchOpts.config
	}

	_, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid requestPath: %w", err)
	}

	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout duration in config: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

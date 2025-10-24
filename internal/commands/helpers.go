package commands

import (
	"encoding/json"
	"os"

	"htruong/kt-crawler/internal"
)

func LoadConfig(path string) (*internal.Config, error) {
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

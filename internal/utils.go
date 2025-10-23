package internal

import (
	"context"
	"io"
	"net/http"
)

func Fetch(ctx context.Context, path string) ([]byte, error) {
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

	return io.ReadAll(res.Body)
}

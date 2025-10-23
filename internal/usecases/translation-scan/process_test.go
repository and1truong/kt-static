package translation_scan

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawTranslation(t *testing.T) {
	ctx := context.Background()
	url := "https://kinhthanh.httlvn.org/?v=VI1934"
	result, err := Run(ctx, url)
	assert.NoError(t, err)

	fmt.Println(result)
}

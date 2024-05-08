package book_crawler

import (
	"context"
	"fmt"
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestParseActivity(t *testing.T) {
	chap, err := ParseActivity(context.Background(), mockFetchResponse())
	require.NoError(t, err)
	require.Len(t, chap.AudioLinks, 2)
	require.Greater(t, len(chap.Blocks), 25)
	
	for _, block := range chap.Blocks {
		fmt.Println(block)
	}
}

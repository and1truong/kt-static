package book_crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestParseTitle(t *testing.T) {
	// block, chapNumber := parseTitle(selection, attrClass)
	reader := `<table><div class="title"><h1>1</h1><h3>Lời đạt và chào thăm</h3></div></table>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(reader))
	
	selection := doc.Find("div.title")
	_, chapterCount := parseTitle(selection, "title")
	
	require.NoError(t, err)
	require.Equal(t, 1, chapterCount)
}

func TestParseActivity(t *testing.T) {
	type dataset struct {
		path   string
		audio  int
		blocks int
	}
	
	items := []dataset{
		{
			path:   "resources/fixtures/chapter.vi.html",
			audio:  2,
			blocks: 25,
		},
		{
			path:   "resources/fixtures/chapter.en.html",
			audio:  0,
			blocks: 25,
		},
	}
	
	for _, item := range items {
		chap, err := BookParseActivity(context.Background(), mockFetchResponse(item.path))
		require.NoError(t, err)
		require.Len(t, chap.AudioLinks, item.audio)
		require.Greater(t, len(chap.Blocks), item.blocks)
		
		if false {
			// encode result for next activity
			chapBytes, err := json.Marshal(chap)
			require.NoError(t, err)
			
			fmt.Println(string(chapBytes), "\n ===\n ")
		}
		
		if false {
			// debug output
			fmt.Println(chap)
		}
	}
}

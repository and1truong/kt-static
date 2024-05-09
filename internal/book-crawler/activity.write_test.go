package book_crawler

import (
	"encoding/json"
	"os"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources/fixtures"
)

type writeTestData struct {
	bookRaw        []byte
	chapterRaw     []byte
	resultPath     string
	resultContains []string
}

func TestWriteResultActivity(t *testing.T) {
	items := []writeTestData{
		{
			bookRaw:    fixtures.BookSample_VI_JudeJson,
			chapterRaw: fixtures.ChapterSample_VI_Json,
			resultPath: "/build/static/VI1934/giu/1.md",
			resultContains: []string{
				"title: Giu-đe  1",
				"book: Giu-đe",
				"chapter: 1",
				"translation: VI1934",
				"slug: /giu/1/VI1934",
			},
		},
	}
	
	var book entity.BookInfo
	var chapter entity.ChapterInfo
	logs := map[string]string{}
	
	nopeWriter := ResultWriter{
		writer: func(name string, data []byte, perm os.FileMode) error {
			logs[name] = string(data)
			
			return nil
			
		},
	}
	
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestActivityEnvironment()
	env.SetTestTimeout(10 * time.Minute)
	env.RegisterActivity(nopeWriter.WriteResultActivity)
	
	for _, item := range items {
		require.NoError(t, json.Unmarshal(item.bookRaw, &book))
		require.NoError(t, json.Unmarshal(item.chapterRaw, &chapter))
		
		_, err := env.ExecuteActivity(nopeWriter.WriteResultActivity, book, chapter)
		require.NoError(t, err)
		
		log := logs[item.resultPath]
		for _, contain := range item.resultContains {
			require.Contains(t, log, contain)
		}
	}
}

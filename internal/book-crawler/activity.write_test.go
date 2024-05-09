package book_crawler

import (
	"encoding/json"
	"os"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

type writeTestData struct {
	bookPath       string
	chapterPath    string
	resultPath     string
	resultContains []string
}

func TestWriteResultActivity(t *testing.T) {
	items := []writeTestData{
		{
			bookPath:    "resources/fixtures/book-info.jude.json",
			chapterPath: "resources/fixtures/chapter.vi.json",
			resultPath:  "../../build/VI1934/giu/1.html",
			resultContains: []string{
				"title: Giu-đe  1",
				"book: Giu-đe",
				"chapter: 1",
				"translation: VI1934",
				"slug: /giu/1/VI1934",
			},
		},
	}
	
	var book BookInfo
	var chapter ChapterInfo
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
	env.RegisterActivity(nopeWriter.ActivityHandler)
	
	for _, item := range items {
		bookRaw, bookErr := os.ReadFile(item.bookPath)
		chapterRaw, chapterErr := os.ReadFile(item.chapterPath)
		require.NoError(t, bookErr)
		require.NoError(t, chapterErr)
		
		json.Unmarshal(bookRaw, &book)
		json.Unmarshal(chapterRaw, &chapter)
		
		_, err := env.ExecuteActivity(nopeWriter.ActivityHandler, book, chapter)
		require.NoError(t, err)
		
		log := logs[item.resultPath]
		for _, contain := range item.resultContains {
			require.Contains(t, log, contain)
		}
	}
}

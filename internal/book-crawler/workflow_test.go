package book_crawler

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
	
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/resources/fixtures"
)

func mockBookInfo() BookInfo {
	var book BookInfo
	var err error
	
	err = json.Unmarshal(fixtures.BookSampleJudeJson, &book)
	if err != nil {
		panic(err)
	}
	
	return book
}

func nopeWriterActivity() any {
	writer := ResultWriter{
		writer: func(name string, data []byte, perm os.FileMode) error {
			if false {
				fmt.Println("ResultWriter › write", name, string(data))
			}
			
			return nil
			
		},
	}
	
	return writer.WriteResultActivity
}

func TestWorkflow(t *testing.T) {
	// setup workflow
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.SetTestTimeout(10 * time.Minute)
	env.SetDetachedChildWait(true)
	env.SetTestTimeout(10 * time.Minute)
	env.
		OnActivity(activities.FetchActivity, mock.Anything, "https://kinhthanh.httlvn.org/doc-kinh-thanh/giu/1?v=VI1934").
		Return(fixtures.ChapterSampleViHtml, nil)
	env.RegisterActivity(BookParseActivity)
	env.RegisterActivity(nopeWriterActivity())
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, mockBookInfo())
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

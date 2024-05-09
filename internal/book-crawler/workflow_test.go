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
)

func mockBookInfo() BookInfo {
	content, err := os.ReadFile("resources/fixtures/book-info.jude.json")
	if err != nil {
		panic(err)
	}
	
	var book BookInfo
	err = json.Unmarshal(content, &book)
	if err != nil {
		panic(err)
	}
	
	return book
}

func mockFetchResponse(path string) []byte {
	mockHTML, err := os.ReadFile(path)
	
	if err != nil {
		panic(err)
	}
	
	return mockHTML
}

func nopeWriterActivity() any {
	writer := ResultWriter{
		writer: func(name string, data []byte, perm os.FileMode) error {
			fmt.Println("ResultWriter › write", name, string(data))
			
			return nil
			
		},
	}
	
	return writer.ActivityHandler
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
		Return(mockFetchResponse("resources/fixtures/chapter.vi.html"), nil)
	env.RegisterActivity(BookParseActivity)
	env.RegisterActivity(nopeWriterActivity())
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, mockBookInfo())
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

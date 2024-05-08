package book_crawler

import (
	"context"
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

func mockFetchResponse() []byte {
	mockHTML, err := os.ReadFile("resources/fixtures/chapter.html")
	
	if err != nil {
		panic(err)
	}
	
	return mockHTML
}

func TestParseActivity(t *testing.T) {
	chap, err := ParseActivity(context.Background(), mockFetchResponse())
	if nil != err {
		panic(err)
	}
	
	fmt.Println(chap)
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
		Return(mockFetchResponse(), nil)
	env.RegisterActivity(ParseActivity)
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, mockBookInfo())
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

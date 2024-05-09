package book_crawler

import (
	"fmt"
	"os"
	"testing"
	"time"
	
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources/fixtures"
)

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
		Return(fixtures.ChapterSample_VI_Html, nil)
	env.RegisterActivity(BookParseActivity)
	env.RegisterActivity(nopeWriterActivity())
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, MockBookInfo())
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

func TestWorkflowWithCodeName(t *testing.T) {
	// setup workflow
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.SetTestTimeout(10 * time.Minute)
	env.SetDetachedChildWait(true)
	env.SetTestTimeout(10 * time.Minute)
	
	// env.
	// 	OnActivity(activities.FetchActivity, mock.Anything, "https://kinhthanh.httlvn.org/doc-kinh-thanh/giu/1?v=VI1934").
	// 	Return(fixtures.ChapterSample_VI_Html, nil)
	
	env.RegisterActivity(BookParseActivity)
	env.RegisterActivity(activities.FetchActivity)
	env.RegisterActivity(nopeWriterActivity())
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, entity.BookInfo{
		Tran:     "VI1934",
		BookCode: "he",
	})
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

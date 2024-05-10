//go:build expensive_tests
// +build expensive_tests

package book_crawler

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/chapter-crawler"
	"temporal-crawler/internal/entity"
)

func TestWorkflowWithCodeName(t *testing.T) {
	type testcase struct {
		tran string
		book string
	}
	
	tests := []entity.BookInfo{
		{Tran: "VI1934", BookCode: "he"},
		{Tran: "VI1934", BookCode: "nha"},
		{Tran: "VI1934", BookCode: "2co"},
	}
	
	for _, bookInfo := range tests {
		// setup workflow
		ts := &testsuite.WorkflowTestSuite{}
		env := ts.NewTestWorkflowEnvironment()
		env.SetTestTimeout(10 * time.Minute)
		env.SetDetachedChildWait(true)
		env.SetTestTimeout(10 * time.Minute)
		env.RegisterWorkflow(chapter_crawler.ChapterCrawlerWorkflow)
		env.RegisterActivity(chapter_crawler.ChapterParseActivity)
		env.RegisterActivity(activities.FetchActivity)
		env.RegisterActivity(nopeWriterActivity())
		
		// run it
		env.ExecuteWorkflow(BookCrawlerWorkflow, bookInfo)
		
		require.True(t, env.IsWorkflowCompleted())
		require.NoError(t, env.GetWorkflowError())
		
		var result int
		require.NoError(t, env.GetWorkflowResult(&result))
	}
}

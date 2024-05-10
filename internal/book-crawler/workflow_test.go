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
	"temporal-crawler/internal/chapter-crawler"
	"temporal-crawler/internal/resources/fixtures"
)

func nopeWriterActivity() any {
	writer := chapter_crawler.NewResultWriter(
		func(name string, data []byte, perm os.FileMode) error {
			if false {
				fmt.Println("ResultWriter › write", name, string(data))
			}
			
			return nil
		},
		false,
	)
	
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
	env.RegisterWorkflow(chapter_crawler.ChapterCrawlerWorkflow)
	env.RegisterActivity(chapter_crawler.ChapterParseActivity)
	env.RegisterActivity(nopeWriterActivity())
	env.RegisterActivity(BuildBookDocusaurusIndexActivity)
	
	env.
		OnActivity(
			activities.FileWritingActivity,
			mock.Anything,
			"build/static/VI1934/giu/_category_.json",
			[]byte(`{"label":"Giu-đe ","position":65,"link":{"type":"generated-index"}}`),
			mock.Anything,
		).
		Return(nil)
	
	// run it
	env.ExecuteWorkflow(BookCrawlerWorkflow, MockBookInfo(fixtures.BookInfo_VI_JudeJson))
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	
	var result int
	require.NoError(t, env.GetWorkflowResult(&result))
}

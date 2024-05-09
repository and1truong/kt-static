package translation_crawler

import (
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources/fixtures"
)

func _TestBookCrawlerWorkflow(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.SetDetachedChildWait(true)
	env.SetTestTimeout(10 * time.Minute)
	env.
		OnActivity(activities.FetchActivity, mock.Anything, "https://kinhthanh.httlvn.org/?v=VI1934").
		Return(fixtures.TranslationSample_VI1934_Html, nil)
	env.RegisterActivity(TranslationParseActivity)
	env.ExecuteWorkflow(TranslationCrawlerWorkflow, "VI1934")
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result *entity.TranslationInfo
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "VI1934", result.Name)
	require.Equal(t, 66, len(result.Books))
	require.Equal(t, "Sáng-thế Ký", strings.Trim(result.Books[0].BookName, " "))
}

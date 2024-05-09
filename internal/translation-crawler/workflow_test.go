package translation_crawler

import (
	"os"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/activities"
)

func _TestBookCrawlerWorkflow(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.SetDetachedChildWait(true)
	env.SetTestTimeout(10 * time.Minute)
	mockHTML, err := os.ReadFile("resources/fixtures/fetch.translation.VI1934.html")
	env.OnActivity(activities.FetchActivity, mock.Anything, "https://kinhthanh.httlvn.org/?v=VI1934").Return(mockHTML, err)
	env.RegisterActivity(TranslationParseActivity)
	env.ExecuteWorkflow(TranslationCrawlerWorkflow, "VI1934")
	
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result *TranslationInfo
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "VI1934", result.Name)
	require.Equal(t, 66, len(result.Books))
	require.Equal(t, "Sáng-thế Ký", strings.Trim(result.Books[0].BookName, " "))
}

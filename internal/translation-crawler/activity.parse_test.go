package translation_crawler

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources/fixtures"
)

func TestParseActivity(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestActivityEnvironment()
	env.SetTestTimeout(10 * time.Minute)
	env.RegisterActivity(TranslationParseActivity)
	
	val, err := env.ExecuteActivity(TranslationParseActivity, "VI1934", fixtures.TranslationSample_VI1934_Html)
	require.NoError(t, err)
	
	var result *entity.TranslationInfo
	require.NoError(t, val.Get(&result))
	require.Equal(t, "VI1934", result.Name)
	require.Equal(t, 66, len(result.Books))
	require.Equal(t, "/doc-kinh-thanh/sa/1?v=VI1934", result.Books[0].Chapters[0])
	require.Equal(t, "Sáng-thế Ký", strings.Trim(result.Books[0].BookName, " "))
	
	if false {
		{
			out, _ := json.Marshal(result.Books[18])
			fmt.Println(string(out))
		}
		
		{
			out, _ := json.Marshal(result.Books[64])
			fmt.Println(string(out))
		}
	}
}

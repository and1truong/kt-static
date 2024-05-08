package internal

import (
	"go.temporal.io/sdk/workflow"
)

func TranslationCrawlerWorkflow(ctx workflow.Context, translation string) (string, error) {
	return "hello", nil
}

package translation_crawler

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflow(TranslationCrawlerWorkflow)
	w.RegisterActivity(TranslationParseActivity)
}

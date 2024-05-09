package translation_crawler

import (
	"go.temporal.io/sdk/worker"
	"temporal-crawler/internal"
	"temporal-crawler/internal/activities"
)

func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflow(internal.TranslationCrawlerWorkflow)
	w.RegisterActivity(activities.FetchActivity)
	w.RegisterActivity(TranslationParseActivity)
}

package chapter_crawler

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker, wd string) {
	w.RegisterWorkflow(ChapterCrawlerWorkflow)
	w.RegisterActivity(ChapterParseActivity)
	
	writer := NewResultWriter(nil, wd, true)
	w.RegisterActivity(writer.WriteResultActivity)
}

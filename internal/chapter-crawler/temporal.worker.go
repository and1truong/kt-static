package chapter_crawler

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflow(ChapterCrawlerWorkflow)
	w.RegisterActivity(ChapterParseActivity)
	
	writer := NewResultWriter(nil, true)
	w.RegisterActivity(writer.WriteResultActivity)
}

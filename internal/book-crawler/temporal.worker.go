package book_crawler

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker, wd string) {
	writer := NewResultWriter(nil, wd)
	w.RegisterWorkflow(BookCrawlerWorkflow)
	w.RegisterActivity(BookParseActivity)
	w.RegisterActivity(writer.WriteResultActivity)
}

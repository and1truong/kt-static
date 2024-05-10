package book_crawler

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflow(BookCrawlerWorkflow)
}

package activities

import (
	"go.temporal.io/sdk/worker"
)

func RegisterWorkflow(w worker.Worker) {

}

func RegisterActivities(w worker.Worker) {
	w.RegisterActivity(FetchActivity)
	w.RegisterActivity(FileWritingActivity)
}

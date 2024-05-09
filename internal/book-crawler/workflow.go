package book_crawler

import (
	"time"

	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
)

var (
	baseUrl = "https://kinhthanh.httlvn.org"
)

func BookCrawlerWorkflow(ctx workflow.Context, bookInfo BookInfo) (int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})

	var futures []workflow.Future

	// trigger fetch activities
	for _, chapPath := range bookInfo.Chapters {
		ft := workflow.ExecuteActivity(ctx, activities.FetchActivity, baseUrl+chapPath)
		futures = append(futures, ft)
	}

	// wait for activities to be completed
	for _, ft := range futures {
		var body []byte
		err := ft.Get(ctx, &body)
		if err != nil {
			return 0, err
		}

		// trigger ParseActivity
		var chapter *ChapterInfo
		err = workflow.ExecuteActivity(ctx, ParseActivity, body).Get(ctx, &chapter)
		if err != nil {
			return 0, err
		}
	}

	return len(futures), nil
}

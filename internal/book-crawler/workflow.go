package book_crawler

import (
	"time"
	
	"github.com/pkg/errors"
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
	
	var fetchFutures []workflow.Future
	var parseFutures []workflow.Future
	var writeFutures []workflow.Future
	
	// ============================
	// trigger fetch activities
	// ============================
	for _, chapPath := range bookInfo.Chapters {
		ft := workflow.ExecuteActivity(ctx, activities.FetchActivity, baseUrl+chapPath)
		fetchFutures = append(fetchFutures, ft)
	}
	
	// ============================
	// wait for fetching-activities to be completed
	// and trigger parsing activities
	// ============================
	for _, ft := range fetchFutures {
		var body []byte
		err := ft.Get(ctx, &body)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get fetch result")
		}
		
		// trigger BookParseActivity
		ft := workflow.ExecuteActivity(ctx, BookParseActivity, body)
		parseFutures = append(parseFutures, ft)
	}
	
	// ============================
	// wait for parsing-activities to be completed
	// and trigger writing activities
	// ============================
	writer := ResultWriter{}
	for _, ft := range parseFutures {
		var chapter ChapterInfo
		err := ft.Get(ctx, &chapter)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get chapter result")
		}
		
		ft := workflow.ExecuteActivity(ctx, writer.ActivityHandler, bookInfo, chapter)
		writeFutures = append(writeFutures, ft)
	}
	
	// wait for writing-activities to be completed
	for _, ft := range writeFutures {
		err := ft.Get(ctx, nil)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get result from writing activity")
		}
	}
	
	return len(fetchFutures), nil
}

package chapter_crawler

import (
	"time"
	
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources/translation"
)

var (
	baseUrl = "https://kinhthanh.httlvn.org"
)

func ChapterCrawlerWorkflow(ctx workflow.Context, bookInfo entity.BookInfo, chapPath string) (int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})
	
	var body []byte
	err := workflow.
		ExecuteActivity(ctx, activities.FetchActivity, baseUrl+chapPath).
		Get(ctx, &body)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get fetch result")
	}
	
	var chapter entity.ChapterInfo
	err = workflow.
		ExecuteActivity(ctx, ChapterParseActivity, translation.Translations[bookInfo.Tran], body).
		Get(ctx, &chapter)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get chapter result")
	}
	
	writer := ResultWriter{}
	err = workflow.
		ExecuteActivity(ctx, writer.WriteResultActivity, bookInfo, chapter).
		Get(ctx, nil)
	
	if err != nil {
		return 0, errors.Wrap(err, "failed to get result from writing activity")
	}
	
	return 1, nil
}

package translation_crawler

import (
	"fmt"
	"time"
	
	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
	book_crawler "temporal-crawler/internal/book-crawler"
	"temporal-crawler/internal/entity"
)

func TranslationCrawlerWorkflow(ctx workflow.Context, tran string) (map[int]int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})
	
	// ============================
	// execute FetchActivity
	// ============================
	var body []byte
	err := workflow.ExecuteActivity(ctx, activities.FetchActivity, "https://kinhthanh.httlvn.org/?v="+tran).Get(ctx, &body)
	if err != nil {
		return nil, err
	}
	
	// ============================
	// parse translation information
	// ============================
	var transInfo *entity.TranslationInfo
	err = workflow.ExecuteActivity(ctx, TranslationParseActivity, tran, body).Get(ctx, &transInfo)
	if err != nil {
		return nil, err
	}
	
	// ============================
	// start crawling books
	// ============================
	var results []workflow.ChildWorkflowFuture
	for _, bookInfo := range transInfo.Books {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf(
				"%s/%s/%s",
				workflow.GetInfo(ctx).WorkflowExecution.ID,
				tran,
				bookInfo.BookCode,
			),
		})
		
		child := workflow.ExecuteChildWorkflow(childCtx, book_crawler.BookCrawlerWorkflow, bookInfo)
		results = append(results, child)
	}
	
	// ============================
	// Waits for all child workflows to complete
	// ============================
	out := map[int]int{}
	for i, result := range results {
		chapters := 0
		if err := result.Get(ctx, &chapters); err != nil {
			return nil, err
		}
		out[i] = chapters
	}
	
	return out, nil
}

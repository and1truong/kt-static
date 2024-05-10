package book_crawler

import (
	"context"
	"fmt"
	"time"
	
	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
	chapter_crawler "temporal-crawler/internal/chapter-crawler"
	"temporal-crawler/internal/entity"
)

func BookCrawlerWorkflow(ctx workflow.Context, bookInfo entity.BookInfo) (int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 20 * time.Minute,
	})
	
	// if full bookInfo is not provided, we should rebuild them
	if len(bookInfo.Chapters) == 0 {
		if bookInfo.Tran == "" || bookInfo.BookCode == "" {
			panic("invalid book info")
		}
		
		var body []byte
		err := workflow.ExecuteActivity(ctx, activities.FetchActivity, "https://kinhthanh.httlvn.org/?v="+bookInfo.Tran).Get(ctx, &body)
		if err != nil {
			return 0, err
		}
		
		// TODO: parse book
		// func BookListParse(ctx context.Context, tran string, body []byte) (*tc.TranslationInfo, error) {
		tranInfo, err := chapter_crawler.BookListParse(context.TODO(), bookInfo.Tran, body)
		if err != nil {
			return 0, err
		}
		
		for _, item := range tranInfo.Books {
			if item.BookCode == bookInfo.BookCode {
				bookInfo = item
			}
		}
		
		if len(bookInfo.Chapters) == 0 {
			panic("book not found")
		}
	}
	
	// ============================
	// trigger chapter crawling workflow
	// ============================
	var results []workflow.ChildWorkflowFuture
	for i, chapPath := range bookInfo.Chapters {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf(
				"%s/%d",
				workflow.GetInfo(ctx).WorkflowExecution.ID,
				i+1,
			),
		})
		
		child := workflow.ExecuteChildWorkflow(childCtx, chapter_crawler.ChapterCrawlerWorkflow, bookInfo, chapPath)
		results = append(results, child)
	}
	
	// ============================
	// Waits for all child workflows to complete
	// ============================
	counter := 0
	for _, result := range results {
		val := 0
		if err := result.Get(ctx, &val); err != nil {
			return 0, err
		}
		
		counter = counter + val
	}
	
	return counter, nil
}

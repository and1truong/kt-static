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
		
		// parse book
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
	
	// ============================
	// build book index file: _category_.json
	// schema is like:
	//  {
	//      "label": "1 Cô-rinh-tô",
	//      "position": 1,
	//      "link": { "type": "generated-index" }
	// }
	// ============================
	var docusaurusIndex []byte
	var err error
	err = workflow.ExecuteActivity(ctx, BuildBookDocusaurusIndexActivity, bookInfo).Get(ctx, &docusaurusIndex)
	if err != nil {
		return 0, err
	}
	
	// write the file
	writePath := fmt.Sprintf("build/static/%s/%s/_category_.json", bookInfo.Tran, bookInfo.BookCode)
	err = workflow.ExecuteActivity(ctx, activities.FileWritingActivity, writePath, docusaurusIndex, 0644).Get(ctx, nil)
	if err != nil {
		return 0, nil
	}
	
	return counter, nil
}

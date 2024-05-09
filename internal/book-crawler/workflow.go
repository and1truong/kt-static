package book_crawler

import (
	"context"
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

func BookCrawlerWorkflow(ctx workflow.Context, bookInfo entity.BookInfo) (int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})
	
	var fetchFutures []workflow.Future
	var parseFutures []workflow.Future
	var writeFutures []workflow.Future
	
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
		tranInfo, err := BookListParse(context.TODO(), bookInfo.Tran, body)
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
		ft := workflow.ExecuteActivity(ctx, BookParseActivity, translation.Translations[bookInfo.Tran], body)
		parseFutures = append(parseFutures, ft)
	}
	
	// ============================
	// wait for parsing-activities to be completed
	// and trigger writing activities
	// ============================
	writer := ResultWriter{}
	for _, ft := range parseFutures {
		var chapter entity.ChapterInfo
		err := ft.Get(ctx, &chapter)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get chapter result")
		}
		
		ft := workflow.ExecuteActivity(ctx, writer.WriteResultActivity, bookInfo, chapter)
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

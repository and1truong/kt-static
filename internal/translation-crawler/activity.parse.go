package translation_crawler

import (
	"context"
	
	book_crawler "temporal-crawler/internal/chapter-crawler"
	"temporal-crawler/internal/entity"
)

func TranslationParseActivity(ctx context.Context, tran string, body []byte) (*entity.TranslationInfo, error) {
	return book_crawler.BookListParse(ctx, tran, body)
}

package listeners

import (
	"context"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/services/eventdispatcher"
)

const BookScanEventName = "book.scan"

type BookScanEvent struct {
	*eventdispatcher.BaseEvent
	baseURL         string
	book            internal.Book
	translationCode string
}

// NewBookScanEvent creates a new BookScanEvent.
func NewBookScanEvent(baseURL string, book internal.Book, translationCode string) *BookScanEvent {
	return &BookScanEvent{
		BaseEvent:       eventdispatcher.NewBaseEvent(BookScanEventName),
		baseURL:         baseURL,
		book:            book,
		translationCode: translationCode,
	}
}

type BookScanListener struct {
	Dispatcher *eventdispatcher.Dispatcher
}

// NewBookScanListener creates a new BookScanListener.
func NewBookScanListener(dispatcher *eventdispatcher.Dispatcher) *BookScanListener {
	return &BookScanListener{
		Dispatcher: dispatcher,
	}
}

func (l *BookScanListener) Handle(ctx context.Context, rawEvent eventdispatcher.Event) error {
	event, ok := rawEvent.(*BookScanEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", rawEvent)
	}

	if len(event.book.Chapters) == 0 {
		return fmt.Errorf("book does not contain any chapters")
	}

	// Dispatch a event for each chapter
	for _, chapterPath := range event.book.Chapters {
		crawlEvent := NewChapterScanEvent(event.baseURL, event.book, chapterPath, event.translationCode)
		if err := l.Dispatcher.Dispatch(ctx, crawlEvent); err != nil {
			return fmt.Errorf("failed to dispatch chapter crawl event for %s: %w", chapterPath, err)
		}
	}

	return nil
}

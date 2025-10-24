package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/logging"
	"os"
	"path/filepath"
	"strings"
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

type categoryFile struct {
	Label    string `json:"label"`
	Position uint   `json:"position"`
	Link     struct {
		Type string `json:"type"`
	} `json:"link"`
}

type BookScanListener struct {
	Dispatcher *eventdispatcher.Dispatcher
	Config     internal.StoreConfig
	Logger     logging.Logger
}

// NewBookScanListener creates a new BookScanListener.
func NewBookScanListener(dispatcher *eventdispatcher.Dispatcher, config internal.StoreConfig, logger logging.Logger) *BookScanListener {
	return &BookScanListener{
		Dispatcher: dispatcher,
		Config:     config,
		Logger:     logger,
	}
}

func (l *BookScanListener) Handle(ctx context.Context, rawEvent eventdispatcher.Event) error {
	event, ok := rawEvent.(*BookScanEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", rawEvent)
	}

	category := categoryFile{
		Label:    strings.TrimSpace(event.book.Name),
		Position: event.book.Number,
		Link: struct {
			Type string `json:"type"`
		}{
			Type: "generated-index",
		},
	}

	content, err := json.MarshalIndent(category, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal category file content for book %s: %w", event.book.Code, err)
	}

	// Construct the file path: [base_dir]/[translation_code]/[book_code]/_category_.json
	baseDir := l.Config.Filesystem.Directory
	translationDir := filepath.Join(baseDir, event.translationCode)
	bookDir := filepath.Join(translationDir, event.book.Code)
	filePath := filepath.Join(bookDir, "_category_.json")

	// Ensure the book directory exists
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return fmt.Errorf("failed to create book directory %s: %w", bookDir, err)
	}

	if l.Config.DryRun {
		l.Logger.Info("Dry run: Would write book category file",
			"book", event.book.Code,
			"path", filePath,
		)
	} else {
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("failed to write book category file to %s: %w", filePath, err)
		}
		l.Logger.Info("Wrote book category file",
			"book", event.book.Code,
			"path", filePath,
		)
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

package listeners

import (
	"context"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/logging"
	"os"
	"path/filepath"
)

const StoreEventName = "chapter.parsed"

// StoreEvent is dispatched after a chapter has been successfully parsed.
type StoreEvent struct {
	*eventdispatcher.BaseEvent
	Book            internal.Book
	Chapter         internal.Chapter
	TranslationCode string
}

func NewStoreEvent(book internal.Book, chapter internal.Chapter, translationCode string) *StoreEvent {
	return &StoreEvent{
		BaseEvent:       eventdispatcher.NewBaseEvent(StoreEventName),
		Book:            book,
		Chapter:         chapter,
		TranslationCode: translationCode,
	}
}

// StoreListener writes the parsed chapter content to the filesystem.
type StoreListener struct {
	Config internal.StoreConfig
	Logger logging.Logger
}

// NewStoreListener creates a new StoreListener.
func NewStoreListener(config internal.StoreConfig, logger logging.Logger) *StoreListener {
	return &StoreListener{
		Config: config,
		Logger: logger,
	}
}

func (l *StoreListener) Handle(ctx context.Context, rawEvent eventdispatcher.Event) error {
	event, ok := rawEvent.(*StoreEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", rawEvent)
	}

	// Construct the file path: [base_dir]/[translation_code]/[book_code]/[chapter_number].txt
	baseDir := l.Config.Filesystem.Directory
	translationDir := filepath.Join(baseDir, event.TranslationCode)
	bookDir := filepath.Join(translationDir, event.Book.Code)
	fileName := fmt.Sprintf("%d.txt", event.Chapter.Number)
	filePath := filepath.Join(bookDir, fileName)

	// Ensure the book directory exists
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return fmt.Errorf("failed to create book directory %s: %w", bookDir, err)
	}

	// Write the chapter content
	content := ""
	for _, block := range event.Chapter.Blocks {
		content += block.String()
	}

	if l.Config.DryRun {
		l.Logger.Info("Dry run: Would write chapter content",
			"book", event.Book.Code,
			"chapter", event.Chapter.Number,
			"path", filePath,
			"content_length", len(content),
		)
		return nil
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write chapter content to file %s: %w", filePath, err)
	}

	return nil
}

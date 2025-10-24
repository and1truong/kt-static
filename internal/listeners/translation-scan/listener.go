package translation_scan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/listeners/book-scan"
	"htruong/kt-crawler/internal/services/cache"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/fetch"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const TranslationScanEventName = "translation.scan"

// TranslationScanEvent is the event dispatched when a translation needs to be scanned.
type TranslationScanEvent struct {
	*eventdispatcher.BaseEvent
	RequestPath string
}

// NewTranslationScanEvent creates a new TranslationScanEvent.
func NewTranslationScanEvent(requestPath string) *TranslationScanEvent {
	return &TranslationScanEvent{
		BaseEvent:   eventdispatcher.NewBaseEvent(TranslationScanEventName),
		RequestPath: requestPath,
	}
}

// TranslationScanListener is the listener for translation scan events.
type TranslationScanListener struct {
	Dispatcher *eventdispatcher.Dispatcher
	Fetcher    *fetch.Fetcher
}

// NewTranslationScanListener creates a new TranslationScanListener.
func NewTranslationScanListener(dispatcher *eventdispatcher.Dispatcher, config fetch.Config, cacheStore cache.Store) *TranslationScanListener {
	return &TranslationScanListener{
		Dispatcher: dispatcher,
		Fetcher:    fetch.NewFetcher(config, cacheStore),
	}
}

type Result struct {
	testament   string
	group       string
	bookNumber  uint
	translation internal.Translation
}

func baseURL(requestURL string) (string, error) {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
}

// Handle implements the event_dispatcher.Listener interface.
func (l *TranslationScanListener) Handle(ctx context.Context, event eventdispatcher.Event) error {
	scanEvent, ok := event.(*TranslationScanEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", event)
	}

	var (
		requestURL = scanEvent.RequestPath
		body       []byte
		err        error
	)

	baseURL, err := baseURL(requestURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	body, err = l.Fetcher.Fetch(ctx, requestURL, fetch.WithConfig(l.Fetcher.Config()))
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return err
	}

	result, err := scan(doc)
	if err != nil {
		return err
	}

	for _, book := range result.translation.Books {
		event := book_scan.NewBookScanEvent(baseURL, book)
		if err := l.Dispatcher.Dispatch(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func scan(doc *goquery.Document) (*Result, error) {
	result := &Result{}

	// scan translation name
	tName := doc.Find(".book-name h1")
	if tName.Size() == 0 {
		return nil, errors.New("translation name not found")
	}

	result.translation.Name = tName.Text()

	// find list of books
	bookList := doc.Find(".book-list")
	if bookList.Size() == 0 {
		return nil, fmt.Errorf("book list not found")
	}

	// loop through left & right
	bookList.Find(".col-md-6").EachWithBreak(
		func(i int, column *goquery.Selection) bool {
			return scanColumn(result, column)
		},
	)

	return result, nil
}

func scanColumn(result *Result, column *goquery.Selection) bool {
	column.Find(".col-md-12").EachWithBreak(
		func(i int, box *goquery.Selection) bool {
			return scanGroup(result, box)
		},
	)

	return true
}

func scanGroup(result *Result, box *goquery.Selection) bool {
	if box.Find("h3").Size() > 0 {
		result.testament = box.Find("h3").Text()
		result.group = ""
	} else if box.Find("h4").Size() > 0 {
		result.group = box.Find("h4").Text()
	} else {
		bookName := box.Find("span").Text()
		result.bookNumber++
		book := internal.Book{
			Number:    result.bookNumber,
			Name:      bookName,
			Chapters:  []string{},
			Group:     result.group,
			Testament: result.testament,
		}

		box.Find(".dropdown.pull-right > ul > li > a").EachWithBreak(
			func(i int, link *goquery.Selection) bool {
				uri, exists := link.Attr("href")
				if exists && uri != "" {
					book.Chapters = append(book.Chapters, uri)
				}

				return true
			},
		)

		// find book's machine name from chapter's RequestPath
		if len(book.Chapters) > 0 && len(book.Chapters[0]) > 0 {
			// The format is like: /doc-kinh-thanh/sa/1?v=VI1934
			// Book's machine name should be: sa
			parts := strings.Split(book.Chapters[0], "/")
			book.Code = parts[2]
		}

		result.translation.Books = append(result.translation.Books, book)

		return true
	}

	return true
}

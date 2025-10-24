package chapter_listener

import (
	"bytes"
	"context"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/listeners/store-listener"
	"htruong/kt-crawler/internal/services/cache"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/fetch"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const ChapterScanEventName = "chapter.scan"

type ChapterScanEvent struct {
	*eventdispatcher.BaseEvent
	baseURL     string
	book        internal.Book
	requestPath string
}

func NewChapterScanEvent(baseURL string, book internal.Book, requestPath string) *ChapterScanEvent {
	return &ChapterScanEvent{
		BaseEvent:   eventdispatcher.NewBaseEvent(ChapterScanEventName),
		book:        book,
		baseURL:     baseURL,
		requestPath: requestPath,
	}
}

type ChapterScanListener struct {
	Dispatcher *eventdispatcher.Dispatcher
	Fetcher    *fetch.Fetcher
}

// NewChapterScanListener creates a new ChapterScanListener.
func NewChapterScanListener(dispatcher *eventdispatcher.Dispatcher, config fetch.FetchConfig, cacheStore cache.Store) *ChapterScanListener {
	return &ChapterScanListener{
		Dispatcher: dispatcher,
		Fetcher:    fetch.NewFetcher(config, cacheStore),
	}
}

func (l *ChapterScanListener) Handle(ctx context.Context, rawEvent eventdispatcher.Event) error {
	event, ok := rawEvent.(*ChapterScanEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", rawEvent)
	}

	if event.baseURL == "" {
		return fmt.Errorf("chapter scan event is missing baseURL")
	}

	// The TODO requested to assert requestPath, so we check for it.
	if event.requestPath == "" {
		return fmt.Errorf("chapter scan event is missing requestPath")
	}

	requestURL := fmt.Sprintf("%s%s", event.baseURL, event.requestPath)
	body, err := l.Fetcher.Fetch(ctx, requestURL, fetch.WithConfig(l.Fetcher.DefaultConfig()))
	if err != nil {
		return fmt.Errorf("could not fetch book indexing page: %w", err)
	}
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return fmt.Errorf("could not parse chapter page: %w", err)
	}

	if blocks, number, err := parseHtmlDoc(doc); nil != err {
		fmt.Printf("could not parse chapter page: %s\n", err)

		return err
	} else {
		audioLinks := parseAudioLinks(doc)

		chapter := internal.Chapter{
			Number:     number,
			Blocks:     blocks,
			AudioLinks: audioLinks,
		}

		parsedEvent := store_listener.NewStoreEvent(event.book, chapter)
		if err := l.Dispatcher.Dispatch(ctx, parsedEvent); err != nil {
			return fmt.Errorf("failed to dispatch chapter parsed event: %w", err)
		}
	}

	return nil
}

func parseHtmlDoc(doc *goquery.Document) ([]internal.Block, int, error) {
	blocks := []internal.Block{}
	chapNumber := 0
	doc.Find(".bible-read > div > *").EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			block, chap, ok := parseBlock(selection)
			if !ok {
				return true
			}

			blocks = append(blocks, block)
			if chap > 0 {
				chapNumber = chap
			}

			return true
		},
	)

	return blocks, chapNumber, nil
}

func parseBlock(selection *goquery.Selection) (internal.Block, int, bool) {
	attrClass, found := selection.Attr("class")
	if !found {
		return internal.Block{}, 0, false
	}

	if strings.Contains(attrClass, "title") {
		block, chapNumber := parseTitle(selection, attrClass)

		return block, chapNumber, true
	}

	return parseVerse(selection, attrClass), 0, true
}

func parseTitle(selection *goquery.Selection, attrClass string) (internal.Block, int) {
	// <h1>1</h1><h3>Lời đạt và chào thăm</h3>
	chapNumber := 0
	block := internal.Block{
		Kind:    "title",
		Classes: strings.Split(attrClass, " "),
		Content: []internal.InnerBlock{},
	}

	selection.FindMatcher(goquery.Single("h1")).Each(
		func(i int, sub *goquery.Selection) {
			if num, err := strconv.Atoi(sub.Text()); err == nil {
				chapNumber = num
			}
		},
	)

	selection.FindMatcher(goquery.Single("h3")).Each(
		func(i int, sub *goquery.Selection) {
			block.Content = append(block.Content, internal.InnerBlock{
				Content:    cleanUpString(sub.Text()),
				References: nil,
			})
		},
	)

	return block, chapNumber
}

func parseVerse(selection *goquery.Selection, attrClass string) internal.Block {
	block := internal.Block{
		Kind:       "verse",
		Classes:    strings.Split(attrClass, " "),
		References: []string{},
	}

	var txt string
	block.Number, txt, block.NewLine = cleanupInnerText(selection)
	block.Content = parseInnerBlocks(txt)

	return block
}

func cleanupInnerText(selection *goquery.Selection) (string, string, bool) {
	var num string
	var newLine bool
	txt, _ := selection.Html()

	// remove <sup>…</sup>
	selection.Find("sup").Each(
		func(i int, sup *goquery.Selection) {
			num = sup.Text()
			txt = strings.Trim(
				strings.Replace(txt, outerHTML(sup), "", 1),
				" ",
			)
		},
	)

	txt = strings.Trim(txt, "   ")

	if strings.HasSuffix(txt, "<br/>") {
		newLine = true
		txt = strings.TrimSuffix(txt, "<br/>")
	}

	return num, cleanUpString(txt), newLine
}

func cleanUpString(txt string) string {
	txt = strings.Trim(txt, "  ")
	txt = strings.Replace(txt, " ", " ", -1)
	txt = strings.Replace(txt, " ", " ", -1)
	txt = strings.Replace(txt, "Jêsus", "Giê-su", -1)
	txt = strings.Replace(txt, "Christ", "Cơ-đốc", -1)
	txt = strings.Replace(txt, "nầy", "này", -1)

	return txt
}

func parseInnerBlocks(txt string) []internal.InnerBlock {
	var blocks []internal.InnerBlock

	parts := strings.Split(txt, "<br/>")
	for _, part := range parts {
		blocks = append(blocks, parseInnerBlock(part))
	}

	return blocks
}

func parseInnerBlock(txt string) internal.InnerBlock {
	block := internal.InnerBlock{
		Content:    txt,
		References: []string{},
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(txt))
	if err != nil {
		return block
	}

	// Remove <a data-toggle="tooltip" data-placement="bottom" title="…">⚓</a>
	doc.Find("a[data-toggle]").Each(
		func(i int, ref *goquery.Selection) {
			block.Content = strings.Replace(block.Content, outerHTML(ref), "", -1)
			block.References = strings.Split(ref.AttrOr("title", ""), "; ")
		},
	)

	return block
}

func outerHTML(selection *goquery.Selection) string {
	var buf bytes.Buffer
	html.Render(&buf, selection.Nodes[0])

	return buf.String()
}

func parseAudioLinks(doc *goquery.Document) []string {
	audioLinks := []string{}
	doc.Find(".audio-collapse > div > audio").Each(
		func(i int, audio *goquery.Selection) {
			path, found := audio.Attr("src")
			if found {
				audioLinks = append(audioLinks, path)
			}
		},
	)

	return audioLinks
}

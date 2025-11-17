package listeners

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"htruong/kt-crawler/internal"
	"htruong/kt-crawler/internal/services/cache"
	"htruong/kt-crawler/internal/services/eventdispatcher"
	"htruong/kt-crawler/internal/services/fetch"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const ChapterScanEventName = "chapter.scan"

type ChapterScanEvent struct {
	*eventdispatcher.BaseEvent
	baseURL         string
	book            internal.Book
	requestPath     string
	translationCode string
}

func NewChapterScanEvent(baseURL string, book internal.Book, requestPath string, translationCode string) *ChapterScanEvent {
	return &ChapterScanEvent{
		BaseEvent:       eventdispatcher.NewBaseEvent(ChapterScanEventName),
		book:            book,
		baseURL:         baseURL,
		requestPath:     requestPath,
		translationCode: translationCode,
	}
}

type ChapterScanListener struct {
	Dispatcher *eventdispatcher.Dispatcher
	Fetcher    *fetch.Fetcher
}

// NewChapterScanListener creates a new ChapterScanListener.
func NewChapterScanListener(dispatcher *eventdispatcher.Dispatcher, config fetch.Config) *ChapterScanListener {
	return &ChapterScanListener{
		Dispatcher: dispatcher,
		Fetcher:    fetch.NewFetcher(config),
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

	cacheKey := fmt.Sprintf("chapter_scan:%s", url.QueryEscape(requestURL))
	body, err := cache.Cache(
		ctx,
		cacheKey,
		func() ([]byte, error) {
			return l.Fetcher.Fetch(ctx, requestURL, fetch.WithConfig(l.Fetcher.Config()))
		},
	)

	if err != nil {
		return fmt.Errorf("could not fetch book indexing page: %w", err)
	}
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return fmt.Errorf("could not parse chapter page: %w", err)
	}

	if requestURL == "https://kinhthanh.httlvn.org/doc-kinh-thanh/gi/6?v=VI1934" {
		fmt.Println("wip")
	}

	if blocks, number, err := l.parseHtmlDoc(doc); nil != err {
		fmt.Printf("could not parse chapter page: %s\n", err)

		return err
	} else {
		audioLinks := l.parseAudioLinks(doc)

		chapter := internal.Chapter{
			Number:     number,
			Blocks:     blocks,
			AudioLinks: audioLinks,
		}

		parsedEvent := NewStoreEvent(event.book, chapter, event.translationCode)
		if err := l.Dispatcher.Dispatch(ctx, parsedEvent); err != nil {
			return fmt.Errorf("failed to dispatch chapter parsed event: %w", err)
		}
	}

	return nil
}

func (l *ChapterScanListener) parseHtmlDoc(doc *goquery.Document) ([]internal.Block, int, error) {
	blocks := []internal.Block{}
	chapNumber := 0
	doc.Find(".bible-read > div > *").EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			block, chap, ok := l.parseBlock(selection)
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

func (l *ChapterScanListener) parseBlock(selection *goquery.Selection) (internal.Block, int, bool) {
	attrClass, found := selection.Attr("class")
	if !found {
		return internal.Block{}, 0, false
	}

	if strings.Contains(attrClass, "title") {
		block, chapNumber := l.parseTitle(selection, attrClass)

		return block, chapNumber, true
	}

	return l.parseVerse(selection, attrClass), 0, true
}

func (l *ChapterScanListener) parseTitle(selection *goquery.Selection, attrClass string) (internal.Block, int) {
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
				Content:    l.cleanUpString(sub.Text()),
				References: nil,
			})
		},
	)

	return block, chapNumber
}

func (l *ChapterScanListener) parseVerse(selection *goquery.Selection, attrClass string) internal.Block {
	block := internal.Block{
		Kind:       "verse",
		Classes:    strings.Split(attrClass, " "),
		References: []string{},
	}

	var txt string
	block.Number, txt, block.NewLine = l.cleanupInnerText(selection)
	block.Content = l.parseInnerBlocks(txt)

	return block
}

func (l *ChapterScanListener) cleanupInnerText(selection *goquery.Selection) (string, string, bool) {
	var num string
	var newLine bool
	txt, _ := selection.Html()

	// remove <sup>…</sup>
	selection.Find("sup").Each(
		func(i int, sup *goquery.Selection) {
			num = sup.Text()
			supHTML, _ := l.outerHTML(sup)
			txt = strings.Trim(
				strings.Replace(txt, supHTML, "", 1),
				" ",
			)
		},
	)

	txt = strings.Trim(txt, "   ")

	if strings.HasSuffix(txt, "<br/>") {
		newLine = true
		txt = strings.TrimSuffix(txt, "<br/>")
	}

	return num, l.cleanUpString(txt), newLine
}

func (l *ChapterScanListener) cleanUpString(txt string) string {
	txt = strings.Trim(txt, "  ")
	txt = strings.Replace(txt, " ", " ", -1)
	txt = strings.Replace(txt, " ", " ", -1)
	txt = strings.Replace(txt, "Jêsus", "Giê-su", -1)
	txt = strings.Replace(txt, "Christ", "Cơ-đốc", -1)
	txt = strings.Replace(txt, "nầy", "này", -1)

	return txt
}

func (l *ChapterScanListener) parseInnerBlocks(txt string) []internal.InnerBlock {
	var blocks []internal.InnerBlock

	parts := strings.Split(txt, "<br/>")
	for _, part := range parts {
		blocks = append(blocks, l.parseInnerBlock(part))
	}

	return blocks
}

func (l *ChapterScanListener) parseInnerBlock(txt string) internal.InnerBlock {
	block := internal.InnerBlock{
		// Content will be set later
		References: []string{},
	}

	// Wrap the fragment in a body tag to help goquery parse it correctly.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<body>" + txt + "</body>"))
	if err != nil {
		block.Content = txt // fallback
		return block
	}

	bodyHtml, _ := doc.Find("body").Html()
	block.Content = bodyHtml

	// Replace <a data-toggle="tooltip" data-placement="bottom" title="…">⚓</a> with a markdown-style link.
	doc.Find("a[data-toggle]").Each(
		func(i int, ref *goquery.Selection) {
			refHTML, _ := l.outerHTML(ref)
			title := ref.AttrOr("title", "")
			linkText := ref.Text()

			// Create the markdown-style link
			// e.g. [⚓](tooltip:{"title":"…"})
			payload := map[string]string{"title": title}
			jsonPayload, _ := json.Marshal(payload)
			markdownLink := fmt.Sprintf("[%s](tooltip:%s)", linkText, string(jsonPayload))

			// Replace the original <a> tag with the new markdown link
			block.Content = strings.Replace(block.Content, refHTML, markdownLink, 1)

			// Keep populating references
			block.References = append(block.References, strings.Split(title, "; ")...)
		},
	)

	return block
}

func (l *ChapterScanListener) outerHTML(selection *goquery.Selection) (string, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, selection.Nodes[0]); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (l *ChapterScanListener) parseAudioLinks(doc *goquery.Document) []string {
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

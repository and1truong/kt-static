package translation_scan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"htruong/kt-crawler/internal"
	"strings"
)
import "github.com/PuerkitoBio/goquery"

var (
	fetch = internal.Fetch
)

type Result struct {
	testament   string
	group       string
	bookNumber  uint
	translation internal.Translation
}

func Run(ctx context.Context, url string) (*Result, error) {
	var (
		body []byte
		err  error
	)

	body, err = fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	return scan(doc)
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
				uri, _ := link.Attr("href")
				book.Chapters = append(book.Chapters, uri)

				return true
			},
		)

		// find book's machine name from chapter's URL
		if len(book.Chapters[0]) > 0 {
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

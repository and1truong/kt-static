package translation_crawler

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	
	"github.com/PuerkitoBio/goquery"
	book_crawler "temporal-crawler/internal/book-crawler"
)

type parser struct {
	testament  string
	group      string
	bookNumber uint
	result     *TranslationInfo
}

func TranslationParseActivity(ctx context.Context, tran string, body []byte) (*TranslationInfo, error) {
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}
	
	pr := &parser{
		testament:  "",
		group:      "",
		bookNumber: 0,
		result: &TranslationInfo{
			Name:  tran,
			Books: []book_crawler.BookInfo{},
		},
	}
	
	if err := pr.parse(doc); err != nil {
		return nil, err
	}
	
	return pr.result, nil
}

// .book-list
// .col-md-6.col-sm-6 (left)
//      .col-md-12 > h3 "Cựu Ước"
//      .col-md-12 > h4 "Năm sách Môi-se"
//      .col-md-12 > h4 "Năm sách Môi-se"
//      .col-md-12 > span "Sáng Thế Ký"
//      .col-md-12 > .dropdown.pull-right > ul > li > a "1"
// .col-md-6.col-sm-6 (right)
//      h3 Tân Ước
func (p *parser) parse(doc *goquery.Document) error {
	bookList := doc.Find(".book-list")
	if bookList.Size() == 0 {
		return fmt.Errorf("book list not found")
	}
	
	// loop through left & right
	bookList.Find(".col-md-6").EachWithBreak(
		func(i int, column *goquery.Selection) bool {
			return p.parseColumn(column)
		},
	)
	
	return nil
}

func (p *parser) parseColumn(column *goquery.Selection) bool {
	column.Find(".col-md-12").EachWithBreak(
		func(i int, box *goquery.Selection) bool {
			return p.parseGroup(i, box)
		},
	)
	
	return true
}

func (p *parser) parseGroup(i int, box *goquery.Selection) bool {
	if box.Find("h3").Size() > 0 {
		p.testament = box.Find("h3").Text()
		p.group = ""
	} else if box.Find("h4").Size() > 0 {
		p.group = box.Find("h4").Text()
	} else {
		bookName := box.Find("span").Text()
		p.bookNumber++
		book := book_crawler.BookInfo{
			Tran:       p.result.Name,
			BookNumber: p.bookNumber,
			BookName:   bookName,
			Chapters:   []string{},
			Group:      p.group,
			Testament:  p.testament,
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
			book.BookCode = parts[2]
		}
		
		p.result.Books = append(p.result.Books, book)
		
		return true
	}
	
	return true
}

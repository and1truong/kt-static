package translation_crawler

import (
	"bytes"
	"context"
	"fmt"
	
	"github.com/PuerkitoBio/goquery"
	book_crawler "temporal-crawler/internal/book-crawler"
)

func ParseActivity(ctx context.Context, tran string, body []byte) (*TranslationInfo, error) {
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}
	
	testament := ""
	group := ""
	bookNumber := uint(0)
	tranInfo := &TranslationInfo{
		Name:  tran,
		Books: []book_crawler.BookInfo{},
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
	
	bookList := doc.Find(".book-list")
	if bookList.Size() == 0 {
		return nil, fmt.Errorf("book list not found")
	}
	
	bookList.Find(".col-md-6").Each(
		// loop through left & right
		func(i int, column *goquery.Selection) {
			column.Find(".col-md-12").EachWithBreak(
				func(i int, box *goquery.Selection) bool {
					if box.Find("h3").Size() > 0 {
						testament = box.Find("h3").Text()
						group = ""
						// book = ""
					} else if box.Find("h4").Size() > 0 {
						group = box.Find("h4").Text()
						// book = ""
					} else {
						bookName := box.Find("span").Text()
						bookNumber += 1
						book := book_crawler.BookInfo{
							BookNumber: bookNumber,
							BookName:   bookName,
							Chapters:   []string{},
							Group:      group,
							Testament:  testament,
						}
						
						box.Find(".dropdown.pull-right > ul > li > a").EachWithBreak(
							func(i int, link *goquery.Selection) bool {
								uri, _ := link.Attr("href")
								book.Chapters = append(book.Chapters, uri)
								
								return true
							},
						)
						
						tranInfo.Books = append(tranInfo.Books, book)
						
						return true
					}
					
					return true
				},
			)
		},
	)
	
	return tranInfo, nil
}

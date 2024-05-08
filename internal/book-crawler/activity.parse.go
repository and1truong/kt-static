package book_crawler

import (
	"bytes"
	"context"
	"fmt"
	
	"github.com/PuerkitoBio/goquery"
)

type block struct {
	Kind       string
	Content    string
	Attributes map[string][]string
}

func ParseActivity(ctx context.Context, body []byte) (*ChapterInfo, error) {
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}
	
	var chapterInfo *ChapterInfo
	
	// parse HTML content
	if content, err := parseContent(doc); nil != err {
		return nil, err
	} else {
		audioLinks := parseAudioLinks(doc)
		
		// return write(baseDir, chap, content, audioLinks)
		
		fmt.Println(content, audioLinks)
	}
	
	return chapterInfo, nil
}

func parseContent(doc *goquery.Document) (string, error) {
	return doc.Find(".bible-read > div").Html()
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

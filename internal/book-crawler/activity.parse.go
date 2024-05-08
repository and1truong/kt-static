package book_crawler

import (
	"bytes"
	"context"
	"strings"
	
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func ParseActivity(ctx context.Context, body []byte) (*ChapterInfo, error) {
	bodyReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}
	
	chap := &ChapterInfo{
		Blocks:     []Block{},
		AudioLinks: []string{},
	}
	
	// parse HTML content
	if chap.Blocks, err = parseContent(doc); nil != err {
		return nil, err
	} else {
		chap.AudioLinks = parseAudioLinks(doc)
	}
	
	return chap, nil
}

func parseContent(doc *goquery.Document) ([]Block, error) {
	blocks := []Block{}
	doc.Find(".bible-read > div > *").EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			if block, ok := parseBlock(selection); !ok {
				return true
			} else {
				blocks = append(blocks, block)
			}
			
			return true
		},
	)
	
	return blocks, nil
}

func parseBlock(selection *goquery.Selection) (Block, bool) {
	attrClass, found := selection.Attr("class")
	if !found {
		return Block{}, false
	}
	
	block := Block{
		Kind:       "verse",
		Classes:    strings.Split(attrClass, " "),
		References: []string{},
	}
	
	if strings.Contains(attrClass, "title") {
		// <h1>1</h1><h3>Lời đạt và chào thăm</h3>
		block.Kind = "title"
		
		selection.Find("h3").Each(
			func(i int, sub *goquery.Selection) {
				block.Content = sub.Text()
			},
		)
	} else {
		block.Content, _ = selection.Html()
		
		// remove <sup>…</sup>
		selection.Find("sup").Each(
			func(i int, sup *goquery.Selection) {
				block.Number = sup.Text()
				block.Content = strings.Trim(
					strings.Replace(block.Content, outerHTML(sup), "", 1),
					" ",
				)
			},
		)
		
		// Remove <a data-toggle="tooltip" data-placement="bottom" title="…">⚓</a>
		selection.Find("a[data-toggle]").Each(
			func(i int, ref *goquery.Selection) {
				block.Content = strings.Replace(block.Content, outerHTML(ref), "", 1)
				block.References = strings.Split(ref.AttrOr("title", ""), "; ")
			},
		)
		
		block.Content = strings.Trim(block.Content, "   ")
		block.Content = strings.Replace(block.Content, "Jêsus", "Giê-su", -1)
		block.Content = strings.Replace(block.Content, "Christ", "Cơ-đốc", -1)
	}
	
	return block, true
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

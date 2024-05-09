package book_crawler

import (
	"bytes"
	"context"
	"strconv"
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
		Number:     0,
		Blocks:     []Block{},
		AudioLinks: []string{},
	}
	
	// parse HTML content
	if chap.Blocks, chap.Number, err = parseHtmlDoc(doc); nil != err {
		return nil, err
	} else {
		chap.AudioLinks = parseAudioLinks(doc)
	}
	
	return chap, nil
}

func parseHtmlDoc(doc *goquery.Document) ([]Block, int, error) {
	blocks := []Block{}
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

func parseBlock(selection *goquery.Selection) (Block, int, bool) {
	attrClass, found := selection.Attr("class")
	if !found {
		return Block{}, 0, false
	}
	
	if strings.Contains(attrClass, "title") {
		block, chapNumber := parseTitle(selection, attrClass)
		
		return block, chapNumber, true
	}
	
	return parseVerse(selection, attrClass), 0, true
}

func parseTitle(selection *goquery.Selection, attrClass string) (Block, int) {
	// <h1>1</h1><h3>Lời đạt và chào thăm</h3>
	chapNumber := 0
	block := Block{
		Kind:    "title",
		Classes: strings.Split(attrClass, " "),
		Content: []InnerBlock{},
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
			block.Content = append(block.Content, InnerBlock{
				Content:    sub.Text(),
				References: nil,
			})
		},
	)
	
	return block, chapNumber
}

func parseVerse(selection *goquery.Selection, attrClass string) Block {
	block := Block{
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
	
	txt = strings.Trim(txt, "  ")
	txt = strings.Replace(txt, " ", " ", -1)
	txt = strings.Replace(txt, "Jêsus", "Giê-su", -1)
	txt = strings.Replace(txt, "Christ", "Cơ-đốc", -1)
	txt = strings.Replace(txt, "nầy", "này", -1)
	
	return num, txt, newLine
}

func parseInnerBlocks(txt string) []InnerBlock {
	var blocks []InnerBlock
	
	parts := strings.Split(txt, "<br/>")
	for _, part := range parts {
		blocks = append(blocks, parseInnerBlock(part))
	}
	
	return blocks
}

func parseInnerBlock(txt string) InnerBlock {
	block := InnerBlock{
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

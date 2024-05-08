package book_crawler

import (
	"time"
	
	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
)

var (
	baseUrl = "https://kinhthanh.httlvn.org"
)

type (
	BookInfo struct {
		BookNumber uint
		BookName   string
		Chapters   []string
		Group      string
		Testament  string
	}
	
	ChapterInfo struct {
		Blocks     []Block
		AudioLinks []string
	}
	
	Block struct {
		Kind       string
		Number     string
		Content    string // Giu-đe, tôi tớ của Đức Chúa Jêsus Christ và em Gia-cơ, đạt cho những kẻ đã được kêu gọi, được Đức Chúa Trời, là Cha, yêu thương, và được Đức Chúa Jêsus Christ giữ gìn:<a data-toggle="tooltip" data-placement="bottom" title="Mat 13:55; Mac 6:3">⚓</a>
		Classes    []string
		References []string
	}
)

func (b Block) String() string {
	out := ""
	if b.Kind == "title" {
		out += "\n"
		out += "## "
		out += b.Content
		out += "\n"
	} else {
		out += b.Number + ". "
		out += b.Content
		
		if len(b.References) > 0 {
			out += "\n\n"
			for _, ref := range b.References {
				out += "    🔴 " + ref + "\n"
			}
			out += "\n"
		}
	}
	
	return out
}

func BookCrawlerWorkflow(ctx workflow.Context, bookInfo BookInfo) (int, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})
	
	var futures []workflow.Future
	
	// trigger fetch activities
	for _, chapPath := range bookInfo.Chapters {
		// https://kinhthanh.httlvn.org/doc-kinh-thanh/sa/10?v=VI1934
		ft := workflow.ExecuteActivity(ctx, activities.FetchActivity, baseUrl+chapPath)
		futures = append(futures, ft)
	}
	
	// wait for activities to be completed
	for _, ft := range futures {
		var body []byte
		err := ft.Get(ctx, &body)
		if err != nil {
			return 0, err
		}
		
		// parse content
		var chapter *ChapterInfo
		err = workflow.ExecuteActivity(ctx, ParseActivity, body).Get(ctx, &chapter)
		if err != nil {
			return 0, err
		}
	}
	
	return 0, nil
}

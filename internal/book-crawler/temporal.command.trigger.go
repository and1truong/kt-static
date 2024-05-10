package book_crawler

import (
	_ "embed"
	"encoding/json"
	
	"github.com/spf13/cobra"
	"temporal-crawler/internal"
	"temporal-crawler/internal/entity"
)

func MockBookInfo(content []byte) entity.BookInfo {
	var book entity.BookInfo
	var err error
	
	// content := fixtures.BookInfo_VI_JudeJson
	// content := fixtures.BookInfo_VI_NahumJson
	
	err = json.Unmarshal(content, &book)
	if err != nil {
		panic(err)
	}
	
	return book
}

var TriggerCommand = &cobra.Command{
	Use:     "book",
	Short:   "Start crawling a book",
	Example: "kt-crawler start book VI1934 kh",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var result int
		
		return internal.ExecuteTemporalWorkflow(
			cmd.Context(),
			BookCrawlerWorkflow,
			&result,
			entity.BookInfo{
				Tran:     args[0],
				BookCode: args[1],
			},
		)
	},
}

package book_crawler

import (
	_ "embed"
	"encoding/json"

	"github.com/spf13/cobra"
	"temporal-crawler/internal"
	"temporal-crawler/internal/resources/fixtures"
)

func MockBookInfo() BookInfo {
	var book BookInfo
	var err error

	// content := fixtures.BookInfo_VI_JudeJson
	content := fixtures.BookInfo_VI_NahumJson

	err = json.Unmarshal(content, &book)
	if err != nil {
		panic(err)
	}

	return book
}

var TriggerCommand = &cobra.Command{
	Use:   "book",
	Short: "Start crawling a book",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result int

		return internal.ExecuteTemporalWorkflow(cmd.Context(), BookCrawlerWorkflow, &result, MockBookInfo())
	},
}

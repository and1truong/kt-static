package book_crawler

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"temporal-crawler/internal/resources/fixtures"
)

func mockBookInfo() BookInfo {
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
	Use:   "start",
	Short: "Start crawling a book",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.Dial(client.Options{})
		if err != nil {
			panic(err)
		} else {
			defer c.Close()
		}
		
		wfOptions := client.StartWorkflowOptions{TaskQueue: "kt-crawling"}
		ctx := context.Background()
		
		bookInfo := mockBookInfo()
		process, err := c.ExecuteWorkflow(ctx, wfOptions, BookCrawlerWorkflow, bookInfo)
		
		if err != nil {
			return errors.Wrap(err, "failure starting workflow")
		} else {
			log.Println("Started Workflow Execution", "WorkflowID", process.GetID(), "RunID", process.GetRunID())
		}
		
		// Wait for Workflow Execution completion.
		// This is rarely needed in real use cases as batch workflows are usually long-running.
		var result int
		err = process.Get(ctx, &result)
		if err != nil {
			return err
		}
		
		log.Println("Completed workflow", "WorkflowID", process.GetID(), "RunID", process.GetRunID(), "Result", result)
		
		return nil
	},
}

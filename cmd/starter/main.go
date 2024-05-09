package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	
	"go.temporal.io/sdk/client"
	book_crawler "temporal-crawler/internal/book-crawler"
	"temporal-crawler/internal/book-crawler/resources/fixtures"
)

func mockBookInfo() book_crawler.BookInfo {
	content := fixtures.BookInfoJudeJson
	
	var book book_crawler.BookInfo
	if err := json.Unmarshal(content, &book); err != nil {
		panic(err)
	}
	
	return book
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "batch-sliding-window",
	}
	ctx := context.Background()
	we, err := c.ExecuteWorkflow(ctx, workflowOptions, batch_sliding_window.ProcessBatchWorkflow, batch_sliding_window.ProcessBatchWorkflowInput{
		PageSize:          5,
		SlidingWindowSize: 10,
		Partitions:        3,
	})
	if err != nil {
		log.Fatalln("Failure starting workflow", err)
	}
	log.Println("Started Workflow Execution", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	
	// Wait for Workflow Execution completion.
	// This is rarely needed in real use cases as batch workflows are usually long-running.
	var result int
	err = we.Get(ctx, &result)
	if err != nil {
		panic(err)
	}
	log.Println("Completed workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID(), "Result", result)
}

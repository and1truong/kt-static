package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	
	"go.temporal.io/sdk/client"
	book_crawler "temporal-crawler/internal/book-crawler"
	"temporal-crawler/internal/resources/fixtures"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	} else {
		defer c.Close()
	}
	
	wfOptions := client.StartWorkflowOptions{TaskQueue: "kt-crawling"}
	ctx := context.Background()
	
	bookInfo := mockBookInfo()
	process, err := c.ExecuteWorkflow(ctx, wfOptions, book_crawler.BookCrawlerWorkflow, bookInfo)
	
	if err != nil {
		log.Fatalln("Failure starting workflow", err)
	} else {
		log.Println("Started Workflow Execution", "WorkflowID", process.GetID(), "RunID", process.GetRunID())
	}
	
	// Wait for Workflow Execution completion.
	// This is rarely needed in real use cases as batch workflows are usually long-running.
	var result int
	err = process.Get(ctx, &result)
	if err != nil {
		panic(err)
	}
	
	log.Println("Completed workflow", "WorkflowID", process.GetID(), "RunID", process.GetRunID(), "Result", result)
}

func mockBookInfo() book_crawler.BookInfo {
	content := fixtures.BookInfoJudeJson
	
	var book book_crawler.BookInfo
	if err := json.Unmarshal(content, &book); err != nil {
		panic(err)
	}
	
	return book
}

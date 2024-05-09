package main

import (
	"log"
	
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"temporal-crawler/internal"
	"temporal-crawler/internal/activities"
	book_crawler "temporal-crawler/internal/book-crawler"
	translation_crawler "temporal-crawler/internal/translation-crawler"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	
	w := worker.New(c, "kt-crawling", worker.Options{})
	
	{
		w.RegisterWorkflow(internal.TranslationCrawlerWorkflow)
		w.RegisterActivity(activities.FetchActivity)
		w.RegisterActivity(translation_crawler.TranslationParseActivity)
	}
	
	{
		writer := book_crawler.ResultWriter{}
		w.RegisterWorkflow(book_crawler.BookCrawlerWorkflow)
		w.RegisterActivity(book_crawler.BookParseActivity)
		w.RegisterActivity(writer.ActivityHandler)
	}
	
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

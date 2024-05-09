package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"temporal-crawler/internal/book-crawler"
	translation_crawler "temporal-crawler/internal/translation-crawler"
)

var rootCmd = &cobra.Command{
	Use:   "kt-crawler",
	Short: "KT crawler",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var TemporalWorkerComand = &cobra.Command{
	Use:   "worker",
	Short: "worker",
	Run: func(cmd *cobra.Command, args []string) {
		// The client and worker are heavyweight objects that should be created once per process.
		c, err := client.Dial(client.Options{
			HostPort: client.DefaultHostPort,
		})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()
		
		wd, _ := os.Getwd()
		w := worker.New(c, "kt-crawling", worker.Options{})
		
		translation_crawler.RegisterWorkflow(w)
		book_crawler.RegisterWorkflow(w, wd)
		
		// log.Println("Started Workflow Execution", "WorkflowID", w.GetID(), "RunID", w.GetRunID())
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}

func main() {
	rootCmd.AddCommand(book_crawler.TriggerCommand)
	rootCmd.AddCommand(TemporalWorkerComand)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

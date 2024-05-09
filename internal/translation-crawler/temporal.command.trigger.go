package translation_crawler

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"temporal-crawler/internal"
)

var TriggerCommand = &cobra.Command{
	Use:   "translation",
	Short: "Start crawling a translation",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result map[int]int
		
		err := internal.ExecuteTemporalWorkflow(cmd.Context(), TranslationCrawlerWorkflow, &result, "VI1934")
		if err != nil {
			return err
		}
		
		fmt.Println("result: ", result)
		
		return nil
	},
}

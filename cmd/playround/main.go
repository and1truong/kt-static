package main

import (
	"context"
	"fmt"
	
	"temporal-crawler/internal"
	tc "temporal-crawler/internal/translation-crawler"
)

func main() {
	result := map[int]int{}
	
	err := internal.ExecuteTemporalWorkflow(context.Background(), tc.TranslationCrawlerWorkflow, &result, "VI1934")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("result: ", result)
}

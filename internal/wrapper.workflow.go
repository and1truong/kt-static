package internal

import (
	"time"
	
	"go.temporal.io/sdk/workflow"
	"temporal-crawler/internal/activities"
)

// BackgroundCheck is your custom Workflow Definition.
func BackgroundCheck(ctx workflow.Context, param string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	})
	
	// Step 3: parse chapter
	//  {
	//     Translation: tran,
	//     BookNumber:  bookNumber,
	//     BookName:    bookName,
	//     Chapter:     uint(chapter),
	//     Group:       group,
	//     Testament:   testament,
	//     Uri:         uri,
	//  }
	// Step 4: write results to static files.
	
	// Step 1: get translations
	activity := workflow.ExecuteActivity(ctx, activities.GetTranslationsActivity)
	translations := []string{}
	if err := activity.Get(ctx, &translations); err != nil {
		return "", err
	}
	
	// Step 2: get book list
	for _, translation := range translations {
		panic(translation)
	}
	
	futures := []workflow.Future{}
	
	if true {
		panic(futures)
	}
	
	var ssnTraceResult string
	err := workflow.ExecuteActivity(ctx, activities.SSNTraceActivity, param).Get(ctx, &ssnTraceResult)
	if err != nil {
		return "", err
	}
	
	return ssnTraceResult, nil
}

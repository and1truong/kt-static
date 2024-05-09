package internal

import (
	"context"
	"log"
	
	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
	"temporal-crawler/internal/resources"
)

func TemporalClient() client.Client {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	
	return c
}

func ExecuteTemporalWorkflow[T any, R any](ctx context.Context, wf T, resul *R, args ...interface{}) error {
	c := TemporalClient()
	defer c.Close()
	
	wfOptions := client.StartWorkflowOptions{TaskQueue: resources.TaskQueueName}
	process, err := c.ExecuteWorkflow(ctx, wfOptions, wf, args...)
	
	if err != nil {
		return errors.Wrap(err, "failure starting workflow")
	} else {
		log.Println("Started Workflow Execution", "WorkflowID", process.GetID(), "RunID", process.GetRunID())
	}
	
	// Wait for Workflow Execution completion.
	// This is rarely needed in real use cases as batch workflows are usually long-running.
	var result R
	err = process.Get(ctx, &result)
	if err != nil {
		return err
	}
	
	log.Println("Completed workflow", "WorkflowID", process.GetID(), "RunID", process.GetRunID(), "Result", result)
	
	return nil
}

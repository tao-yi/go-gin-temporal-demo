package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow functions are where you configure and organize the execution of Activity functions
func HelloWorkflow(ctx workflow.Context, name string) (string, error) {
	workflow.GetLogger(ctx).Info("Starting HelloWorkflow...")

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	input := &HelloActivityInput{Name: name}

	err := workflow.ExecuteActivity(ctx, HelloActivity, *input).Get(ctx, &result)
	return result, err
}

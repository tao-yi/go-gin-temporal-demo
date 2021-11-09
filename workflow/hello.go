package workflow

import (
	"time"

	"github.com/tao-yi/go-gin-temporal-demo/activity"
	"go.temporal.io/sdk/workflow"
)

// Workflow functions are where you configure and organize the execution of Activity functions
func HelloWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, activity.HelloActivity, name).Get(ctx, &result)
	return result, err
}

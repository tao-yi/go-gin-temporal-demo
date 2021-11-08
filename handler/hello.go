package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func Hello(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		options := client.StartWorkflowOptions{
			ID:        "greeting-workflow",
			TaskQueue: HelloTaskQueue,
		}

		ctx := context.Background()
		r, err := cli.ExecuteWorkflow(ctx, options, HelloWorkflow, name)
		if err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}

		var greeting string
		if err = r.Get(ctx, &greeting); err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": greeting})
	}
}

const HelloTaskQueue = "HELLO_TASK_QUEUE"

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

package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/workflow"
	"go.temporal.io/sdk/client"
)

func Hello(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		options := client.StartWorkflowOptions{
			TaskQueue: workflow.TaskQueue,
		}

		ctx := context.Background()
		r, err := cli.ExecuteWorkflow(ctx, options, workflow.HelloWorkflow, name)
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

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/workflow"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type CreateCronJobInput struct {
	CronSchedule string `json:"cronSchedule"`
}

func CreateCronJob(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json CreateCronJobInput
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		options := client.StartWorkflowOptions{
			ID:           "my-cron-job",
			TaskQueue:    workflow.TaskQueue,
			CronSchedule: json.CronSchedule,
		}

		ctx := context.Background()
		r, err := cli.ExecuteWorkflow(ctx, options, workflow.CronJobWorkflow)
		if err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}

		c.JSON(200, gin.H{"workflowId": r.GetID()})
	}
}

func ListCronJobs(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowType := c.Query("workflowType")
		res, err := cli.ListWorkflow(c, &workflowservice.ListWorkflowExecutionsRequest{
			Query: fmt.Sprintf("WorkflowType='%s'", workflowType),
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"res": res.Executions})
	}
}

func GetCronJob(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowID := c.Param("workflowId")
		res, err := cli.DescribeWorkflowExecution(c, workflowID, "")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"res": res})
	}
}

func GetCronJobSchedule(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowID := c.Param("workflowId")
		i := cli.GetWorkflowHistory(c, workflowID, "", false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
		if !i.HasNext() {
			c.JSON(200, gin.H{"error": "no cron schedule"})
			return
		}

		event, err := i.Next()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		cronSchedule := event.GetWorkflowExecutionStartedEventAttributes().GetCronSchedule()
		c.JSON(200, gin.H{"cronSchedule": cronSchedule})
	}
}

func DeleteCronJob(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowID := c.Param("workflowId")
		err := cli.TerminateWorkflow(c, workflowID, "", "updated cron schedule")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
	}
}

type UpdateCronjobInput struct {
	WorkflowID   string `json:"workflowId"`
	CronSchedule string `json:"cronSchedule"`
}

func UpdateCronjob(cli client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UpdateCronjobInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// terminate previous workflow
		// !! note: CancelWorkflow does not work
		err := cli.TerminateWorkflow(c, input.WorkflowID, "", "cron schedule updated")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// start workflow with new cron schedule
		options := client.StartWorkflowOptions{
			ID:           input.WorkflowID,
			TaskQueue:    workflow.TaskQueue,
			CronSchedule: input.CronSchedule,
		}

		ctx := context.Background()
		r, err := cli.ExecuteWorkflow(ctx, options, workflow.CronJobWorkflow)
		if err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}

		c.JSON(200, gin.H{"workflowId": r.GetID()})
	}
}

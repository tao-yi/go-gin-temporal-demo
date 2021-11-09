package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CronJobWorkflow(ctx workflow.Context) error {
	workflow.GetLogger(ctx).Info("Starting HelloWorkflow...")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Actvitivy functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	err := workflow.ExecuteActivity(ctx, CrobJobActivity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

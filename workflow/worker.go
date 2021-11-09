package workflow

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "TASK_QUEUE"

type temporalWorker struct {
	worker.Worker
}

func New(c client.Client) *temporalWorker {
	w := worker.New(c, TaskQueue, worker.Options{})
	w.RegisterWorkflow(HelloWorkflow)
	w.RegisterWorkflow(CronJobWorkflow)
	w.RegisterActivity(HelloActivity)
	w.RegisterActivity(CrobJobActivity)

	return &temporalWorker{Worker: w}
}

func (t *temporalWorker) Start(cb ...func()) {
	go func() {
		err := t.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start Worker", err)
		}
	}()

	for _, f := range cb {
		f()
	}
}

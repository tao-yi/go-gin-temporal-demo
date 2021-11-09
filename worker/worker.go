package worker

import (
	"log"

	"github.com/tao-yi/go-gin-temporal-demo/activity"
	"github.com/tao-yi/go-gin-temporal-demo/handler"
	"github.com/tao-yi/go-gin-temporal-demo/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type temporalWorker struct {
	worker.Worker
}

func New(c client.Client) *temporalWorker {
	w := worker.New(c, handler.HelloTaskQueue, worker.Options{})
	w.RegisterWorkflow(workflow.HelloWorkflow)
	w.RegisterActivity(activity.HelloActivity)

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

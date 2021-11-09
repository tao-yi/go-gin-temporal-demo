package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/server"
	"github.com/tao-yi/go-gin-temporal-demo/workflow"
	"go.temporal.io/sdk/client"
)

var c client.Client

func init() {
	// Create the client object just once per process
	cli, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}

	c = cli
}

func main() {
	defer c.Close()

	r := gin.Default()
	w := workflow.New(c)

	r.GET("/hello", server.Hello(c))
	r.POST("/cronjob", server.CreateCronJob(c))
	r.GET("/cronjob/:workflowId/cronschedule", server.GetCronJobSchedule(c))
	r.GET("/cronjob/:workflowId", server.GetCronJob(c))
	r.GET("/cronjob", server.ListCronJobs(c))
	r.DELETE("/cronjob/:workflowId", server.DeleteCronJob(c))
	r.PUT("/cronjob", server.UpdateCronjob(c))

	w.Start()

	s := server.New(":8091", r)
	s.Start(func() {
		log.Println("the http server started...")
	})

	s.AwaitTerm(func() {
		log.Println("the http server stopped...")
	})
}

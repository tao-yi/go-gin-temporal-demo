package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/handler"
	"github.com/tao-yi/go-gin-temporal-demo/server"
	"github.com/tao-yi/go-gin-temporal-demo/worker"
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
	w := worker.New(c)

	r.GET("/hello", handler.Hello(c))
	// r.GET("/cronjob", handler.GetCronJobS)
	// r.POST("/cronjob", handler.CreateCronjob)
	// r.DELETE("/cronjob", handler.DeleteCronjob)
	// r.PUT("/cronjob", handler.UpdateCronjob)

	w.Start()

	s := server.New(":8091", r)
	s.Start(func() {
		log.Println("the http server started...")
	})

	s.AwaitTerm(func() {
		log.Println("the http server stopped...")
	})
}

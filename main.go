package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tao-yi/go-gin-temporal-demo/activity"
	"github.com/tao-yi/go-gin-temporal-demo/handler"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func temporalClient() client.Client {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	return c
}

func main() {
	r := gin.Default()
	c := temporalClient()
	defer c.Close()

	w := worker.New(c, handler.HelloTaskQueue, worker.Options{})
	w.RegisterWorkflow(handler.HelloWorkflow)
	w.RegisterActivity(activity.HelloActivity)

	r.GET("/hello", handler.Hello(c))
	// r.GET("/cronjob", handler.GetCronJobS)
	// r.POST("/cronjob", handler.CreateCronjob)
	// r.DELETE("/cronjob", handler.DeleteCronjob)
	// r.PUT("/cronjob", handler.UpdateCronjob)

	srv := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start Worker", err)
		}
	}()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to serve http: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("stopping the http server...")
	err := srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

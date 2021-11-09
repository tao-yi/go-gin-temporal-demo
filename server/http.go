package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	s *http.Server
}

func New(addr string, handler http.Handler) *server {
	return &server{
		s: &http.Server{
			Addr:    ":8091",
			Handler: handler,
		},
	}
}

func (srv *server) Start(cb ...func()) {
	go func() {
		err := srv.s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to serve http: %v\n", err)
		}
	}()

	for _, f := range cb {
		f()
	}
}

func (srv *server) AwaitTerm(cb ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("stopping the http server...")

	err := srv.s.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}

	for _, f := range cb {
		f()
	}
}

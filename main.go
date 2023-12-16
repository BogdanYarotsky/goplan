package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	l *log.Logger
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("hello")
}

func main() {
	l := log.New(os.Stdout, "planning-api", log.LstdFlags)
	sm := http.NewServeMux()
	h := &Handler{l}
	sm.Handle("/", h)

	s := &http.Server{
		Addr:    "localhost:9090",
		Handler: sm,
	}

	go func() {
		s.ListenAndServe()
	}()

	sigChannel := make(chan os.Signal, 2)
	signal.Notify(sigChannel, os.Interrupt)
	sig := <-sigChannel
	l.Printf("Received %s signal, shutting down\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
	cancel()
}

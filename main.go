package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BogdanYarotsky/goplan/domain"
	"github.com/BogdanYarotsky/goplan/handlers"
)

func main() {
	l := log.New(os.Stdout, "plan-api", log.LstdFlags)
	ps := domain.NewPlanService()
	ph := handlers.NewPlanHandler(l, ps)
	sm := http.NewServeMux()
	sm.Handle("/", ph)
	s := &http.Server{
		Addr:    "localhost:9090",
		Handler: sm,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	signal.Notify(sigChannel, signals...)
	sig := <-sigChannel
	l.Printf("Received %s signal, shutting down\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		l.Fatal(err)
	}
}

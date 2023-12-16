package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TimeSlot struct {
	id        int
	machineId int
	contentId int
	start     time.Time
	end       time.Time
}

var slots = []*TimeSlot{
	{
		1, 1, 1, time.Now(), time.Now(),
	},
}

type Handler struct {
	l *log.Logger
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func main() {
	l := log.New(os.Stdout, "plan-api", log.LstdFlags)
	sm := http.NewServeMux()
	h := &Handler{l}
	sm.Handle("/", h)
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

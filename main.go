package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TimeSlot struct {
	Id        int
	ContentId int
	Start     time.Time
	End       time.Time
}

var slots = []*TimeSlot{
	{
		1, 777, time.Now(), time.Now(),
	},
}

type PlanningHandler struct {
	l *log.Logger
}

func (h *PlanningHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(slots)
	if err != nil {
		http.Error(w, "Could not parse slots", http.StatusInternalServerError)
	}
	c, err := w.Write(json)
	if err != nil || c < len(json) {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
	}
}

func main() {
	l := log.New(os.Stdout, "plan-api ", log.LstdFlags)
	sm := http.NewServeMux()
	h := &PlanningHandler{l}
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

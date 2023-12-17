package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BogdanYarotsky/goplan/domain"
)

type getSlotsRequest struct {
	machineId int
	start     time.Time
}

type planHandler struct {
	l *log.Logger
	s *domain.PlanService
}

func NewPlanHandler(l *log.Logger, s *domain.PlanService) http.Handler {
	return &planHandler{l, s}
}

func (h *planHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	machineStr := q.Get("machine")
	startStr := q.Get("start")
	endStr := q.Get("end")

	machineId, err := strconv.Atoi(machineStr)
	if err != nil {
		http.Error(w, "Machine should be identified with number", http.StatusBadRequest)
		return
	}

	const layout = time.RFC3339
	start, err := time.Parse(layout, startStr)
	if err != nil {
		http.Error(w, "Start timestamp is not in correct format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		http.Error(w, "End timestamp is not in correct format", http.StatusBadRequest)
		return
	}

	slots, err := h.s.GetSlots(domain.MachineId(machineId), domain.TimeRange{Start: start, End: end})
	if err != nil {
		if _, ok := err.(domain.ValidationError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Could not fetch slots", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(slots)
	if err != nil {
		http.Error(w, "Could not parse slots", http.StatusInternalServerError)
		return
	}
	bc, err := w.Write(data)
	if err != nil || bc < len(data) {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return
	}
}

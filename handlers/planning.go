package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BogdanYarotsky/goplan/domain"
)

type planHandler struct {
	l *log.Logger
}

func NewPlanHandler(l *log.Logger) http.Handler {
	return &planHandler{l}
}

func (h *planHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(domain.Slots)
	if err != nil {
		http.Error(w, "Could not parse slots", http.StatusInternalServerError)
	}
	c, err := w.Write(json)
	if err != nil || c < len(json) {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
	}
}

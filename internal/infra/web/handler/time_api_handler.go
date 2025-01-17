package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TimeApiHandler struct {
	ctx context.Context
}

func NewTimeApiHandler(ctx context.Context) *TimeApiHandler {
	return &TimeApiHandler{
		ctx: ctx,
	}
}

func (d *TimeApiHandler) GetActualDate(w http.ResponseWriter, r *http.Request) {
	type dateResponse struct {
		Date      string `json:"date"`
		Greetings string `json:"greetings"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	now := time.Now()

	err := json.NewEncoder(w).Encode(dateResponse{
		Date:      now.Format("2006-01-02 15:04:05"),
		Greetings: getGreetings(now),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getGreetings(now time.Time) string {
	hour := now.Hour()

	if hour >= 0 && hour < 12 {
		return "Good morning"
	}

	if hour >= 12 && hour < 18 {
		return "Good afternoon"
	}

	return "Good night"
}

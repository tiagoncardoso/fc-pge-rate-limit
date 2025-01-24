package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type GreetingsHandler struct {
	ctx context.Context
}

func NewGreetingsHandler(ctx context.Context) *GreetingsHandler {
	return &GreetingsHandler{
		ctx: ctx,
	}
}

func (d *GreetingsHandler) GetDateAndGreetings(w http.ResponseWriter, r *http.Request) {
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

func (d *GreetingsHandler) GetJapaneseDateAndGreetings(w http.ResponseWriter, r *http.Request) {
	type dateResponse struct {
		Date      string `json:"date"`
		Greetings string `json:"greetings"`
	}

	japanLocation, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	japanTime := time.Now().In(japanLocation)

	err = json.NewEncoder(w).Encode(dateResponse{
		Date:      japanTime.Format("2006-01-02 15:04:05"),
		Greetings: getGreetings(japanTime),
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

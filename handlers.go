package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "alive")
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]

	if user_answers, exists := GlobalAnswers[user_id]; exists {
		json.NewEncoder(w).Encode(user_answers)
	} else {
		json.NewEncoder(w).Encode(&Answers{})
	}
}

func SaveCorrectAnswerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]
	updated_answers_db := SaveAnswer(GlobalAnswers, user_id, true)
	GlobalAnswers = *updated_answers_db
	w.WriteHeader(201)
}

func SaveIncorrectAnswerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, user_id, false))
	w.WriteHeader(201)
}

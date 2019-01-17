package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestIfSavesCorrectAnswer(t *testing.T) {
	// Arrange
	user_id := "UserID"
	answers := CreateAnswersDb()

	// Act
	initial_correct_answers := GetTotalCorrectAnswers(answers, user_id)
	answers = SaveAnswer(*answers, user_id, true)
	total_correct_answers := GetTotalCorrectAnswers(answers, user_id)

	if total_correct_answers != (initial_correct_answers + 1) {
		t.Errorf("Error saving correct answer")
	}
}

func TestIfSavesIncorrectAnswer(t *testing.T) {
	// Arrange
	user_id := "UserID"
	answers := CreateAnswersDb()

	// Act
	initial_incorrect_answers := GetTotalIncorrectAnswers(answers, user_id)
	answers = SaveAnswer(*answers, user_id, false)
	total_incorrect_answers := GetTotalIncorrectAnswers(answers, user_id)

	if total_incorrect_answers != (initial_incorrect_answers + 1) {
		t.Errorf("Error saving correct answer")
	}
}

func TestGetStatsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/myuser/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/user/{user_id}/stats", GetStatsHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned error status")
	}

	response := rr.Body.Next(10000000000000000)
	answers := Answers{}

	if err := json.Unmarshal(response, &answers); err != nil {
		t.Fatal(err)
	}

	if answers.Correct != 0 {
		t.Errorf("Incorrect number of correct answers")
	}

	if answers.Incorrect != 0 {
		t.Errorf("Incorrect number of incorrect answers")
	}
}

func TestSaveCorrectAnswerHandler(t *testing.T) {
	// Arrange
	user_id := "myuser"
	correct_answers := GetTotalCorrectAnswers(&GlobalAnswers, user_id)

	// Act
	req, err := http.NewRequest("POST", "/user/"+user_id+"/correct_answer", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/user/{user_id}/correct_answer", SaveCorrectAnswerHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != 201 {
		t.Errorf("handler status != 201")
	}

	if GlobalAnswers[user_id].Correct != (correct_answers + 1) {
		t.Errorf("correct answer was not saved")
	}
}

func TestSaveIncorrectAnswerHandler(t *testing.T) {
	// Arrange
	user_id := "myuser"
	incorrect_answers := GetTotalIncorrectAnswers(&GlobalAnswers, user_id)

	// Act
	req, err := http.NewRequest("POST", "/user/"+user_id+"/incorrect_answer", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/user/{user_id}/incorrect_answer", SaveIncorrectAnswerHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != 201 {
		t.Errorf("handler status != 201")
	}

	if GlobalAnswers[user_id].Incorrect != (incorrect_answers + 1) {
		t.Errorf("incorrect answer was not saved")
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "alive"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

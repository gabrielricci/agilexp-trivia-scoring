package main

import (
	"encoding/json"
	"log"
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

func TestIfItUpdatesScoreWhenSavingAnswers(t *testing.T) {
	// Arrange
	user_id := "UserID"
	answers := CreateAnswersDb()

	// Act
	answers = SaveAnswer(*answers, user_id, true)
	answers = SaveAnswer(*answers, user_id, true)
	answers = SaveAnswer(*answers, user_id, false)
	answers = SaveAnswer(*answers, user_id, false)

	// Assert
	score := (*answers)[user_id].Score
	if score != float32(50) {
		t.Errorf("Wrong score calculated. Expected: %f, got %f", float32(50), score)
	}
}

func TestIfItSetsUserIdOnAnswers(t *testing.T) {
	// Arrange
	user_id := "UserID123"
	answers := CreateAnswersDb()

	// Act
	answers = SaveAnswer(*answers, user_id, true)

	// Assert
	saved_user_id := (*answers)[user_id].UserId
	if saved_user_id != user_id {
		t.Errorf("User id was not saved")
	}
}

func TestIfItGeneratesRanking(t *testing.T) {
	// Arrange
	user1 := "user1"
	user2 := "user2"
	user3 := "user3"
	answers := CreateAnswersDb()

	// Act
	answers = SaveAnswer(*answers, user3, true)
	answers = SaveAnswer(*answers, user3, true)
	answers = SaveAnswer(*answers, user3, true)
	answers = SaveAnswer(*answers, user2, true)
	answers = SaveAnswer(*answers, user2, true)
	answers = SaveAnswer(*answers, user1, true)

	// Assert
	ranking := GetRanking(answers)

	if len(ranking) != 3 {
		t.Errorf("Ranking is messed up")
	}
	if ranking[0].UserId != user3 {
		t.Errorf("User 3 is not on first place")
	}
	if ranking[1].UserId != user2 {
		t.Errorf("User 2 is not on first place")
	}
	if ranking[2].UserId != user1 {
		t.Errorf("User 1 is not on first place")
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

func TestRankingHandler(t *testing.T) {
	GlobalAnswers = make(UserAnswers)

	// Act
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user3", true))
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user3", true))
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user3", true))
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user2", true))
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user2", true))
	GlobalAnswers = *(SaveAnswer(GlobalAnswers, "user1", true))

	// Act
	req, err := http.NewRequest("GET", "/ranking", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/ranking", RankingHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("error status returned")
	}

	response := rr.Body.Next(10000000000000000)
	var ranking []Answers

	if err := json.Unmarshal(response, &ranking); err != nil {
		t.Fatal(err)
	}

	log.Println(ranking)

	if len(ranking) != 3 {
		t.Errorf("Ranking is messed up")
	}
	if ranking[0].UserId != "user3" {
		t.Errorf("User 3 is not on first place")
	}
	if ranking[1].UserId != "user2" {
		t.Errorf("User 2 is not on first place")
	}
	if ranking[2].UserId != "user1" {
		t.Errorf("User 1 is not on first place")
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

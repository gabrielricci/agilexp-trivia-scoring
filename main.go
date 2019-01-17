package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var GlobalAnswers = make(UserAnswers)

func CreateAnswersDb() *UserAnswers {
	db := make(UserAnswers)
	return &db
}

func SaveAnswer(answers UserAnswers, user_id string, is_correct bool) *UserAnswers {
	user_answers, exists := answers[user_id]
	if !exists {
		user_answers = Answers{}
	}

	if is_correct {
		user_answers.Correct++
	} else {
		user_answers.Incorrect++
	}

	answers[user_id] = user_answers
	return &answers
}

func GetTotalCorrectAnswers(answers *UserAnswers, user_id string) int {
	if user_answers, exists := (*answers)[user_id]; exists {
		return user_answers.Correct
	} else {
		return 0
	}

}

func GetTotalIncorrectAnswers(answers *UserAnswers, user_id string) int {
	if user_answers, exists := (*answers)[user_id]; exists {
		return user_answers.Incorrect
	} else {
		return 0
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/{user_id}/stats", GetStatsHandler)
	router.HandleFunc("/user/{user_id}/correct_answer", SaveCorrectAnswerHandler)
	router.HandleFunc("/user/{user_id}/incorrect_answer", SaveIncorrectAnswerHandler)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// with error handling
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

package main

import (
	"log"
	"net/http"
	"os"
	"sort"

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
		user_answers = Answers{
			UserId: user_id,
		}
	}

	if is_correct {
		user_answers.Correct++
	} else {
		user_answers.Incorrect++
	}

	total_questions := user_answers.Correct + user_answers.Incorrect
	score := (float32(user_answers.Correct) / float32(total_questions)) * 100
	user_answers.Score = score

	answers[user_id] = user_answers
	return &answers
}

func GetRanking(answers *UserAnswers) []Answers {
	var users_and_answers []Answers
	for _, answer := range *answers {
		users_and_answers = append(users_and_answers, answer)
	}

	sort.SliceStable(users_and_answers, func(i, j int) bool {
		user1_score := users_and_answers[i].Score
		user2_score := users_and_answers[j].Score
		user1_correct := users_and_answers[i].Correct
		user2_correct := users_and_answers[j].Correct

		if user1_score == user2_score {
			return user1_correct > user2_correct
		}

		return users_and_answers[i].Score > users_and_answers[j].Score
	})

	return users_and_answers
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
	router.HandleFunc("/ranking", RankingHandler)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// with error handling
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

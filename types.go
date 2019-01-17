package main

type UserAnswers map[string]Answers
type Answers struct {
	UserId    string  `json:"user_id"`
	Correct   int     `json:"correct_answers"`
	Incorrect int     `json:"incorrect_answers"`
	Score     float32 `json:"score"`
}

package main

type UserAnswers map[string]Answers
type Answers struct {
	Correct   int `json:"correct_answers"`
	Incorrect int `json:"incorrect_answers"`
}

package models

type GuessResponse struct {
	Success       bool   `json:"success"`
	IsOneAway     bool   `json:"isOneAway"`
	CategoryTitle string `json:"categoryTitle"`
}

package models

import "time"

type GameResult struct {
	Username string    `json:"username"`
	Outcome  string    `json:"outcome"`
	PlayedAt time.Time `json:"played_at"`
}

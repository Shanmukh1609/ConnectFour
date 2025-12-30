package models

type GameEvent struct {
	GameId          string `json:"game_id"`
	WinnerId        string `json:"winner_id"`
	IsBotWin        bool   `json:"is_bot_win"`
	Duration        int    `json:"duration_seconds"`
	WaitTime        int    `json:"wait_time_seconds"`
}
package controllers

import (
	"time"

	"github.com/Shanmukh1609/backend/models"
)

func CreateNewGame(p1, p2 *models.Player,gameId string) *models.ConnectFour {

	rowMap := make(map[int]int)
	for i := 0; i < models.ColCount; i++ {
		rowMap[i] = models.RowCount - 1
	}

	return &models.ConnectFour{
		ID:          gameId,
		Players:     []*models.Player{p1, p2},
		Rows:        rowMap,
		Status:      models.InProgress,
		CurrentTurn: 1,
		ReconTimers: make(map[string]*time.Timer),
	}
}

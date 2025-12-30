package controllers

import (
	"math"
	"github.com/Shanmukh1609/backend/models"
)

const (
	PlayerPiece = 1
	BotPiece    = 2
	Empty       = 0
	MaxDepth    = 6 // How many moves ahead the bot looks
)

// GetBestMove evaluates the board and returns the best column
// In controllers/bot.go
func GetBestMove(game *models.ConnectFour) int {
    bestScore := math.MinInt64
    bestCol := 3 

    for col := 0; col < 7; col++ {
        // IMPORTANT: Ensure your ValidMove implementation inside models 
        // matches the 0-6 indexing used here, or use getAvailableRow check:
        row := getAvailableRow(game, col)
        
        if row != -1 { // This check prevents the [-1] index panic
            // Simulate the move
            game.Board[row][col] = BotPiece
            
            // Note: depth should be decreased here or inside minimax
            score := minimax(game, MaxDepth, math.MinInt64, math.MaxInt64, false)
            
            // Undo the move
            game.Board[row][col] = Empty

            if score > bestScore {
                bestScore = score
                bestCol = col
            }
        }
    }
    // Return col + 1 if your game logic elsewhere expects 1-7
    return bestCol + 1 
}

func minimax(game *models.ConnectFour, depth int, alpha, beta int, isMaximizing bool) int {
	// 1. Terminal cases: Win, Loss, or Draw
	if game.WinningMove(BotPiece) { return 1000000 + depth }
	if game.WinningMove(PlayerPiece) { return -1000000 - depth }
	if depth == 0 || game.MoveCount == 42 { return evaluateBoard(game) }

	if isMaximizing {
        maxEval := math.MinInt64
        for col := 0; col < 7; col++ {
            row := getAvailableRow(game, col)
            if row != -1 { // CRITICAL CHECK: Ensure the column isn't full
                game.Board[row][col] = BotPiece
                eval := minimax(game, depth-1, alpha, beta, false)
                game.Board[row][col] = Empty // Backtrack
                
                maxEval = max(maxEval, eval)
                alpha = max(alpha, eval)
                if beta <= alpha { break }
            }
        }
        return maxEval
    } else {
        minEval := math.MaxInt64
        for col := 0; col < 7; col++ {
            row := getAvailableRow(game, col)
            if row != -1 { // CRITICAL CHECK: Apply to minimizing player as well
                game.Board[row][col] = PlayerPiece
                eval := minimax(game, depth-1, alpha, beta, true)
                game.Board[row][col] = Empty // Backtrack
                
                minEval = min(minEval, eval)
                beta = min(beta, eval)
                if beta <= alpha { break }
            }
        }
        return minEval
    }
}

// Simple heuristic: gives points for 2-in-a-row or 3-in-a-row
func evaluateBoard(game *models.ConnectFour) int {
    score := 0
    // Give a small bonus for pieces in the center column
    centerCol := 3
    for r := 0; r < 6; r++ {
        if game.Board[r][centerCol] == BotPiece {
            score += 3
        } else if game.Board[r][centerCol] == PlayerPiece {
            score -= 3
        }
    }
    return score
}

func getAvailableRow(game *models.ConnectFour, col int) int {
	for r := 5; r >= 0; r-- {
		if game.Board[r][col] == 0 {
			return r
		}
	}
	return -1
}


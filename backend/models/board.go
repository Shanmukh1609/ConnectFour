package models

import (
	"fmt"
	"sync"
	"time"
)

const RowCount int = 6
const ColCount int = 7

type GameStatus string

const (
	InProgress GameStatus = "IN_PROGRESS"
	Finished   GameStatus = "FINISHED"
)

type ConnectFour struct {
	ID          string                 `json:"userid"`
	Players     []*Player              `json:"players"`
	MoveCount   int                    `json:"move_count"`
	Board       [6][7]int              `json:"board"`
	Rows        map[int]int            `json:"rows"`
	CurrentTurn int                    `json:"current_turn"`
	Status      GameStatus             `json:"status"`
	Winner      int                    `json:"winner"`
	ReconTimers map[string]*time.Timer `json:"-"`
	TurnStartTime time.Time             `json:"turnStartTime"`
	Mu          sync.Mutex
}

func (board *ConnectFour)BroadcastState(){
   
	// state,err:= json.Marshal(board);
	for _,p:= range board.Players{
		if p.Conn!=nil && p.IsOnline{
         p.Conn.WriteJSON(board);
		}
	}
}

func (board *ConnectFour) DisplayBoard() {
	fmt.Printf("[")
	for _, row := range board.Board {
		for _, col := range row {
			fmt.Printf("%d\t", col)
		}
		fmt.Printf("\n")
	}
	fmt.Println("]")
}

// Determing the valid move.

func (board *ConnectFour) ValidMove(column int) bool {
	colIdx := column - 1
	if colIdx < 0 || colIdx >= ColCount {
		return false
	}
	// Check if the column is already full (Row map tracks next available row)
	return board.Rows[colIdx] >= 0
}

func (board *ConnectFour) Assign(p int, col int) {

	isValid := board.ValidMove(col)
	if isValid {
		currentRow := board.Rows[col-1]
		// board.Rows[col] = currentRow + 1
		fmt.Println(currentRow)
		board.Board[currentRow][col-1] = p
		board.MoveCount++
		board.Rows[col-1] = currentRow - 1
	} else {
		fmt.Println("Invalid move")
	}
}

// There is an issue in the valid Move.

func (b *ConnectFour) WinningMove(piece int) bool {
	var board [6][7]int = b.Board
	rowCount := len(board)
	columnCount := len(board[0])

	// Check horizontal locations for win
	for c := 0; c < columnCount-3; c++ {
		for r := 0; r < rowCount; r++ {
			if board[r][c] == piece && board[r][c+1] == piece && board[r][c+2] == piece && board[r][c+3] == piece {
				return true
			}
		}
	}

	// Check vertical locations for win
	for c := 0; c < columnCount; c++ {
		for r := 0; r < rowCount-3; r++ {
			if board[r][c] == piece && board[r+1][c] == piece && board[r+2][c] == piece && board[r+3][c] == piece {
				return true
			}
		}
	}

	// Check positively sloped diagonals
	for c := 0; c < columnCount-3; c++ {
		for r := 0; r < rowCount-3; r++ {
			if board[r][c] == piece && board[r+1][c+1] == piece && board[r+2][c+2] == piece && board[r+3][c+3] == piece {
				return true
			}
		}
	}

	// Check negatively sloped diagonals
	for c := 0; c < columnCount-3; c++ {
		for r := 3; r < rowCount; r++ {
			if board[r][c] == piece && board[r-1][c+1] == piece && board[r-2][c+2] == piece && board[r-3][c+3] == piece {
				return true
			}
		}
	}

	if b.MoveCount == 42 {
		return true
	}

	return false
}

package models

import (
	"sync"
	"github.com/gorilla/websocket"
)

type GameRoom struct {
	Players     map[*websocket.Conn]int
    P1Name      string
	GameState   *ConnectFour
	CurrentTurn int        // Track whose turn it is (1 or 2)
	Mutex       sync.Mutex // Prevents concurrent data corruption
	Ready       bool
	GameOver    bool
}

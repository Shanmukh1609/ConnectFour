package websocket

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shanmukh1609/backend/models"
	"github.com/gorilla/websocket"
)

// Configure the upgrader to allow Postman connections (CORS)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allows Postman and other external clients to connect
	},
}
var Manager = &SessionManager{
	ActivePlayers:  make(map[string]*models.ConnectFour),
	WaitingPlayers: make(map[string]*models.Player),
}

func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection attempt detected...")

	// 1. Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrade failed: %v\n", err)
		return
	}
	// defer ws.Close()

	fmt.Println("Handshake successful! Connection upgraded to WebSocket.")

	type initialData struct {
		PlayerName string  `json:"playerName"`
		GameId     *string `json:"gameId,omitempty"`
		UserId     *string `json:"userId,omitempty"`
	}

	var incomingData initialData
	errk := ws.ReadJSON(&incomingData)
	if errk != nil {
		fmt.Printf("Read error (Postman didn't send JSON yet): %v\n", err)
		// Send error back via WebSocket, not HTTP
		ws.WriteJSON(map[string]string{"error": "Initial setup JSON required"})
		return
	}

	fmt.Println("Incoming ")

	var gameId string
	if incomingData.GameId != nil {
		gameId = *incomingData.GameId
	}
	var userId string
	if incomingData.UserId != nil {
		userId = *incomingData.UserId
	}

	fmt.Println("Calling the func Join or Reconnect")

	Manager.JoinOrReconnect(userId, incomingData.PlayerName, ws, gameId)
	// 5. Connect to Game Logic
	// Note: HandleGameConnections should be updated to handle string RoomIds

	go MonitorGame(userId, ws)

	//What is the need of writing go here ?
}

func MonitorGame(userId string, conn *websocket.Conn) {
	defer func() {
		Manager.HandleDisconnect(userId)
		conn.Close()
	}()

	for {
		var move struct {
			Column int `json:"column"`
		}
		if err := conn.ReadJSON(&move); err != nil {
			fmt.Println("Error in Monitor Game", err.Error())
			break
		}
		//Verify whether it is column or not.
		HandlePlayerMove(userId, int(move.Column), conn)
	}

}

func HandlePlayerMove(userId string, col int, conn *websocket.Conn) {

	game, exists := Manager.ActivePlayers[userId]

	if exists && game.Status == models.Finished {

		Manager.mu.Lock()
		//Removing the players from the game.
		delete(Manager.ActivePlayers, game.Players[0].UserId)
		delete(Manager.ActivePlayers, game.Players[1].UserId)
		Manager.mu.Unlock()

		Manager.EndGame(game)
		// Make a call to database.

		return

	} else if !exists {
		fmt.Println("What is the issue");
		conn.WriteJSON(map[string]interface{}{"message": "other players to join"})
		return
	}

	game.Mu.Lock()

	currentTurn := game.CurrentTurn

	if game.Players[currentTurn-1].UserId != userId {

		for _, p := range game.Players {
			if p.Conn != nil && p.IsOnline {
				p.Conn.WriteJSON(map[string]interface{}{"message": "Invalid Turn"})
			}
		}
		game.Mu.Unlock()
		return
	} else if !game.ValidMove(col) {
		game.Players[currentTurn-1].Conn.WriteJSON(map[string]interface{}{"message": "Invalid Move"})
		game.Mu.Unlock()
		return
	} else {
		game.Assign(currentTurn, col)

		game.TurnStartTime = time.Now()

		if game.Players[1].IsBot {
			go Manager.TriggerBotMove(game) // Run in background so the user gets their update immediately
		}
	

		game.DisplayBoard()

		if game.WinningMove(currentTurn) {
			game.Winner = currentTurn
			game.Status = models.Finished

			game.BroadcastState()
			game.Mu.Unlock()

			Manager.mu.Lock()
			//Removing the players from the game.
			delete(Manager.ActivePlayers, game.Players[0].UserId)
			delete(Manager.ActivePlayers, game.Players[1].UserId)
			Manager.mu.Unlock()

			Manager.EndGame(game)

		} else {
			game.CurrentTurn = 3 - currentTurn
			game.BroadcastState()

			game.Mu.Unlock()
		}
	}
}

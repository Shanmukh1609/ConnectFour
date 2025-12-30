package websocket

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shanmukh1609/backend/controllers"
	"github.com/Shanmukh1609/backend/models"
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	ActivePlayers  map[string]*models.ConnectFour
	WaitingPlayers map[string]*models.Player // Why there is only single player.
	mu             sync.Mutex
}

func (s *SessionManager) JoinOrReconnect(userID string, username string, conn *websocket.Conn, gameId string) {
	s.mu.Lock()

	// Writing code to reconnect the user, if player is disconnected.

	fmt.Println("Join or Reconnect Function")
	if game, exists := s.ActivePlayers[userID]; exists && game.Status == models.InProgress {

		fmt.Println("Reconnecting the player again", username, userID)

		game.Mu.Lock()

		// How to start the timer.
		if timer, ok := game.ReconTimers[userID]; ok {
			timer.Stop()
			delete(game.ReconTimers, userID)
		}

		// Bringing the player into the game.
		for _, p := range game.Players {
			if p.UserId == userID {
				p.Conn = conn
				p.UserName = username
				p.IsOnline = true
			}
		}

		game.Mu.Unlock()

		s.mu.Unlock() // Why are we unlocking the session manager, I think we are just locking the session.

		game.BroadcastState()

		return
	}
	newPlayer := &models.Player{
		UserId:   userID,
		UserName: username,
		Conn:     conn,
		IsOnline: true,
	}

	if gameId == "" {

		id := time.Now().Format("20060102150405")

		s.WaitingPlayers[id] = newPlayer

		s.mu.Unlock()
		gameId := time.Now().Format("20060102150405")

		conn.WriteJSON(map[string]interface{}{"gameId": gameId, "message": "game id is shared"})

		time.AfterFunc(20*time.Second, func() {

			s.mu.Lock()
			defer s.mu.Unlock()

			// 1. Get the player from the map first
			waitingPlayer, exists := s.WaitingPlayers[gameId]

			// 2. Check if the entry exists AND if it belongs to the correct user
			if exists && waitingPlayer != nil && waitingPlayer.UserId == userID {
				fmt.Printf("Matchmaking timeout for %s. Spawning Bot.\n", username)

				bot := &models.Player{
					UserId:   "BOT_OPPONENT",
					UserName: "CompetitiveBot",
					IsBot:    true,
					IsOnline: true,
				}

				s.StartGame(waitingPlayer, bot, gameId)
				delete(s.WaitingPlayers, gameId) // Use delete instead of setting to nil
			}

		})

	} else { // IF the player is offline.
		opponent, exists := s.WaitingPlayers[gameId]
		if !exists || opponent == nil {
			conn.WriteJSON(map[string]string{"error": "Game session not found or expired"})
			s.mu.Unlock()
			return
		}
		delete(s.WaitingPlayers, gameId)
		s.StartGame(opponent, newPlayer, gameId)
		s.mu.Unlock()
	}

}

func (s *SessionManager) StartGame(p1, p2 *models.Player, gameId string) {
	game := controllers.CreateNewGame(p1, p2, gameId)

	s.ActivePlayers[p1.UserId] = game
	s.ActivePlayers[p2.UserId] = game

	fmt.Printf("Match Started: %s vs %s\n", p1.UserName, p2.UserName)
	game.BroadcastState()
	// Start the 45s Turn Supervisor
	go s.MonitorTurnTimeout(game)
}

func (s *SessionManager) MonitorTurnTimeout(game *models.ConnectFour) {
	for {
		time.Sleep(1 * time.Second)
		game.Mu.Lock()

		if game.Status != models.InProgress {
			game.Mu.Unlock()
			return
		}

		// Check if 45 seconds have passed since the turn started
		// Note: You'll need to add 'TurnStartTime' to your ConnectFour model
		if game.Players[0].IsOnline && game.Players[1].IsOnline && time.Since(game.TurnStartTime) >= 45*time.Second {
			fmt.Println("Turn timeout! Executing random move.")

			// Logic to find a random valid column
			var validCols []int
			for c := 0; c < 7; c++ {
				if game.ValidMove(c) {
					validCols = append(validCols, c)
				}
			}

			if len(validCols) > 0 {
				randomCol := validCols[0] // Or actual random logic
				game.Mu.Unlock()
				// Call the existing move handler
				HandlePlayerMove(game.Players[game.CurrentTurn-1].UserId, randomCol, nil)
				continue
			}
		}
		game.Mu.Unlock()
	}
}

// When the user disconnects and if he fails to reconnect and we will make bot to continue the game, if the bot could handle it.

func (s *SessionManager) HandleDisconnect(userId string) {
	s.mu.Lock()
	game, exists := s.ActivePlayers[userId]
	s.mu.Unlock()

	fmt.Println("Handle Disconnect")

	if !exists || game.Status != models.InProgress {
		return
	}

	game.Mu.Lock()
	defer game.Mu.Unlock()

	var username string
	for _, p := range game.Players {
		if userId == p.UserId {
			p.Conn = nil
			p.IsOnline = false
			p.UserName = username
		}
	}

	fmt.Printf("Player %s (ID: %s) disconnected. 30s until forfeit.\n", username, userId)

	game.BroadcastState()
	// Start the 30-second countdown for the forfeit
	game.ReconTimers[userId] = time.AfterFunc(30*time.Second, func() {
		s.ForfeitGame(game, userId)
	})

	// Notify the opponent that the player has disconnected

}

func (s *SessionManager) ForfeitGame(game *models.ConnectFour, disConnectedUserId string) {
	game.Mu.Lock()

	defer game.Mu.Unlock()

	if game.Status != models.InProgress {
		return
	}

	fmt.Printf("Forfeit: UserID %s failed to reconnect in time.\n", disConnectedUserId)
	//Winner Logic

	for i, p := range game.Players {
		if p.UserId != disConnectedUserId {
			game.Winner = i + 1
			game.Status = models.Finished
		}
	}

	game.BroadcastState()

	s.mu.Lock()

	delete(s.ActivePlayers, game.Players[0].UserId)
	delete(s.ActivePlayers, game.Players[1].UserId)

	s.mu.Unlock()
}
func (s *SessionManager) TriggerBotMove(game *models.ConnectFour) {
	if game.Status != models.InProgress {
		return
	}

	// Wait a moment so it feels natural to the human player
	fmt.Println("Triggering a Bot")
	time.Sleep(1 * time.Second)

	game.Mu.Lock()
	col := controllers.GetBestMove(game)
	currentTurn := game.CurrentTurn

	fmt.Println("Bot move", currentTurn, col)
	game.Assign(currentTurn, col)

	if game.WinningMove(currentTurn) {
		game.Winner = currentTurn
		game.Status = models.Finished
		game.BroadcastState()
		Manager.EndGame(game);
	}

	//How do we make sure the bot plays correctly.(Infra)
	game.CurrentTurn = 3 - currentTurn
	game.BroadcastState()


	game.Mu.Unlock()

}
func (s *SessionManager) EndGame(game *models.ConnectFour) {
	s.mu.Lock()
	defer s.mu.Unlock()

	time.Sleep(1 * time.Second)

	for i, p := range game.Players {
		// Calculate outcome: win, lost, or draw
		outcome := "lost"
		if game.Winner == 0 {
			outcome = "draw"
		} else if game.Winner == i+1 {
			outcome = "win"
		}

		// Save results for human players only
		fmt.Println("saving the results to db");
		if !p.IsBot {
			err := models.SaveGameResult(game.ID, p.UserId, p.UserName, outcome)
			if err != nil {
				fmt.Printf("DB Error: %v\n", err)
			}
		}

		// event := models.GameEvent{
		// 	GameId:   game.ID,
		// 	WinnerId: p.UserId,
		// 	IsBotWin: p.IsBot,
		// 	Duration: int(time.Since(game.TurnStartTime).Seconds()), // Calculate match duration
		// }
		// go analytics.PublishGameEnd(event) // Asynchronous call to not block the main loop

		// Clean up maps and close connections safely
		delete(s.ActivePlayers, p.UserId)
		if p.Conn != nil {
			p.Conn.Close() // Safe nil check for bots
		}
	}
}

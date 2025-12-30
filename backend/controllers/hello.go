package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shanmukh1609/backend/models"
	"github.com/google/uuid"
)

func HelloFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Main Point")
	json.NewEncoder(w).Encode(map[string]any{"message": "I am good"})
}

func HandleCookie(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Explicitly allow your React frontend and allow credentials
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	playerName := r.URL.Query().Get("playerName")
	fmt.Printf("Generating cookie for player: %s\n", playerName)

	var userId string
	cookie, errp := r.Cookie("connectFourUserId")

	if errp != nil {
		userId = uuid.NewString()
		newCookie := &http.Cookie{
			Name:     "connectFourUserId",
			Value:    userId,
			Expires:  time.Now().Add(2 * time.Hour),
			HttpOnly: false, // Must be false for React's document.cookie to see it
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		}
		http.SetCookie(w, newCookie)
		fmt.Printf("Assigned the Cookie to user %s\n", userId)
	} else {
		userId = cookie.Value
	}

	// Return the ID in JSON as a fallback in case JS still can't read the cookie
	json.NewEncoder(w).Encode(map[string]string{"userId": userId})
}

func LeaderBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	playerName := r.URL.Query().Get("playerName")
	fmt.Println("Leader Board started!!", playerName)

	// 1. Use Query instead of Exec to retrieve data
	query := `SELECT username, outcome, played_at FROM game_results WHERE username=$1`
	
	rows, err := models.DB.Query(query, playerName)
	if err != nil {
		fmt.Println("jk",err);
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Always close rows to prevent connection leaks

	var results []models.GameResult

	// 2. Loop through results and scan into struct
	fmt.Println(rows);

	for rows.Next() {
		var res models.GameResult
		if err := rows.Scan(&res.Username, &res.Outcome, &res.PlayedAt); err != nil {
			http.Error(w, "Scanning error", http.StatusInternalServerError)
			return
		}
		results = append(results, res)
	}

	fmt.Println(results)

	// 3. Set Header BEFORE writing status or body.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 4. Encode the slice directly to the ResponseWriter
	json.NewEncoder(w).Encode(results);
}

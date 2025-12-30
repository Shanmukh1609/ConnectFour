package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Blank import registers the driver
)

var DB *sql.DB

// InitDB initializes the global DB pool
func InitDB() {
	// 1. Get connection string from Environment (Render uses DATABASE_URL)
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback for local Docker testing
		connStr = "postgres://myuser:mypassword@localhost:5432/game_db?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 2. Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!");

	createTable();
}
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS game_results (
		id SERIAL PRIMARY KEY,
		game_id VARCHAR(50) NOT NULL,
		user_id VARCHAR(50) NOT NULL,
		username VARCHAR(100) NOT NULL,
		outcome VARCHAR(10) NOT NULL, 
		played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Database table verified/created.")
}
// SaveGameResult persists the final state to the database
func SaveGameResult(gameId, userId, username, outcome string) error {
	query := `INSERT INTO game_results (game_id, user_id, username, outcome, played_at) 
              VALUES ($1, $2, $3, $4, NOW())`
	_, err := DB.Exec(query, gameId, userId, username, outcome)
	fmt.Println("Saved the result");
	return err
}



// CREATE TABLE game_results (
//     id SERIAL PRIMARY KEY,
//     game_id VARCHAR(50) NOT NULL,
//     user_id VARCHAR(50) NOT NULL,
//     username VARCHAR(100) NOT NULL,
//     outcome VARCHAR(10) NOT NULL, -- 'win', 'lost', or 'draw'
//     played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
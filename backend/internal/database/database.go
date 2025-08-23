package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Init() (*sql.DB, error) {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// Open SQLite database
	dbPath := filepath.Join(dataDir, "terraforming_mars.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run migrations
	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	queries := []string{
		// Players table (simple, additive only)
		`CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Games table with revision system
		`CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_uuid TEXT NOT NULL, -- groups all revisions of the same game
			revision INTEGER NOT NULL DEFAULT 1,
			name TEXT NOT NULL,
			date DATE NOT NULL,
			map TEXT NOT NULL,
			generations INTEGER NOT NULL,
			expansions TEXT NOT NULL, -- JSON
			is_latest BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Game players table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS game_players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL, -- references specific games.id revision
			player_id INTEGER NOT NULL,
			corporation TEXT NOT NULL,
			terraforming_rating INTEGER DEFAULT 0,
			cities INTEGER DEFAULT 0,
			greeneries INTEGER DEFAULT 0,
			cards INTEGER DEFAULT 0,
			turmoil_points INTEGER DEFAULT 0,
			milestone_points INTEGER DEFAULT 0,
			award_points INTEGER DEFAULT 0,
			total_points INTEGER DEFAULT 0,
			FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
			FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
		)`,

		// Milestones table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS milestones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			winner_game_player_id INTEGER NULL, -- references game_players.id
			FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
			FOREIGN KEY (winner_game_player_id) REFERENCES game_players(id) ON DELETE SET NULL
		)`,

		// Awards table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS awards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
		)`,

		// Award placements table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS award_placements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			award_id INTEGER NOT NULL,
			game_player_id INTEGER NOT NULL,
			placement INTEGER NOT NULL, -- 0=none, 1=gold, 2=silver
			FOREIGN KEY (award_id) REFERENCES awards(id) ON DELETE CASCADE,
			FOREIGN KEY (game_player_id) REFERENCES game_players(id) ON DELETE CASCADE
		)`,

		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_games_uuid ON games(game_uuid)`,
		`CREATE INDEX IF NOT EXISTS idx_games_is_latest ON games(is_latest)`,
		`CREATE INDEX IF NOT EXISTS idx_games_uuid_revision ON games(game_uuid, revision)`,
		`CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_players_player_id ON game_players(player_id)`,
		`CREATE INDEX IF NOT EXISTS idx_milestones_game_id ON milestones(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_awards_game_id ON awards(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_award_placements_award_id ON award_placements(award_id)`,
		`CREATE INDEX IF NOT EXISTS idx_games_date ON games(date)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
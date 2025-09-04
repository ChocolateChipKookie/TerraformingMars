package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Init(connectionString string) (*sql.DB, error) {
	// Check if directory exists for file-based databases
	if connectionString != ":memory:" {
		dir := filepath.Dir(connectionString)
		fileInfo, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("database directory does not exist: %s\nPlease create the directory first with: mkdir -p %s", dir, dir)
		}
		if err != nil {
			return nil, fmt.Errorf("error accessing database directory %s: %v", dir, err)
		}
		if !fileInfo.IsDir() {
			return nil, fmt.Errorf("database path %s exists but is not a directory", dir)
		}
	}

	// Open SQLite database with provided connection string
	// Can be ":memory:" for in-memory, or a file path
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Set WAL mode (except for in-memory databases)
	if connectionString != ":memory:" {
		if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
			return nil, err
		}
	}

	// Run migrations
	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	queries := []string{
		// Player table (simple, additive only)
		`CREATE TABLE IF NOT EXISTS player (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			password_hash TEXT,
			role TEXT NOT NULL DEFAULT 'player' CHECK (role IN ('admin', 'user', 'player')),
			created_by INTEGER,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES player(id) ON DELETE SET NULL
		) STRICT`,

		// Game ID sequence table
		`CREATE TABLE IF NOT EXISTS game_sequence (
			previous_game_id INTEGER PRIMARY KEY DEFAULT 0
		) STRICT`,
		
		// Initialize game sequence if empty
		`INSERT OR IGNORE INTO game_sequence (previous_game_id) VALUES (0)`,

		// Game table with revision system
		`CREATE TABLE IF NOT EXISTS game (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			revision INTEGER NOT NULL DEFAULT 1,
			name TEXT NOT NULL,
			date TEXT NOT NULL,
			map TEXT NOT NULL,
			generations INTEGER NOT NULL,
			expansions TEXT NOT NULL,
			note TEXT,
			created_by INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (game_id, revision),
			FOREIGN KEY (created_by) REFERENCES player(id) ON DELETE CASCADE
		) STRICT`,

		// Game player table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS game_player (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			player_id INTEGER NOT NULL,
			corporation TEXT NOT NULL,
			terraforming_rating INTEGER NOT NULL DEFAULT 0,
			cities INTEGER NOT NULL DEFAULT 0,
			greeneries INTEGER NOT NULL DEFAULT 0,
			cards INTEGER NOT NULL DEFAULT 0,
			turmoil_points INTEGER NOT NULL DEFAULT 0,
			milestone_points INTEGER NOT NULL DEFAULT 0,
			award_points INTEGER NOT NULL DEFAULT 0,
			total_points INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE,
			FOREIGN KEY (player_id) REFERENCES player(id) ON DELETE CASCADE
		) STRICT`,

		// Milestone table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS milestone (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			winner_game_player_id INTEGER,
			FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE,
			FOREIGN KEY (winner_game_player_id) REFERENCES game_player(id) ON DELETE SET NULL
		) STRICT`,

		// Award table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS award (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
		) STRICT`,

		// Award placement table (linked to specific game revision)
		`CREATE TABLE IF NOT EXISTS award_placement (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			award_id INTEGER NOT NULL,
			game_player_id INTEGER NOT NULL,
			placement INTEGER NOT NULL CHECK (placement IN (1, 2)),
			FOREIGN KEY (award_id) REFERENCES award(id) ON DELETE CASCADE,
			FOREIGN KEY (game_player_id) REFERENCES game_player(id) ON DELETE CASCADE
		) STRICT`,

		// Game image table (stored as base64)
		`CREATE TABLE IF NOT EXISTS game_image (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			image_data BLOB NOT NULL,
			mime_type TEXT NOT NULL,
			display_order INTEGER NOT NULL DEFAULT 0,
			uploaded_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
		) STRICT`,

		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_game_game_id_revision ON game(game_id, revision)`,
		`CREATE INDEX IF NOT EXISTS idx_game_player_game_id ON game_player(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_player_player_id ON game_player(player_id)`,
		`CREATE INDEX IF NOT EXISTS idx_milestone_game_id ON milestone(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_award_game_id ON award(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_award_placement_award_id ON award_placement(award_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_date ON game(date)`,
		`CREATE INDEX IF NOT EXISTS idx_game_image_game_id ON game_image(game_id, display_order)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
package models

type Game struct {
	ID          int         `json:"-" db:"id"`                    // Internal primary key
	GameID      int         `json:"id" db:"game_id"`              // User-facing stable game ID
	Revision    int         `json:"revision" db:"revision"`       // Revision of game edit
	Name        string      `json:"name" db:"name"`               // Self explanitory
	Date        string      `json:"date" db:"date"`               // ISO date string YYYY-MM-DD
	Map         *string     `json:"map" db:"map"`                 // Optional for legacy games
	Generations *int        `json:"generations" db:"generations"` // Optional for legacy games
	Expansions  *Expansions `json:"expansions" db:"expansions"`   // Optional for legacy games
	Note        *string     `json:"note" db:"note"`               // Optional text note
	LegacyMode  bool        `json:"legacy_mode" db:"legacy_mode"`
	CreatedBy   int         `json:"created_by" db:"created_by"`
	CreatedAt   string      `json:"created_at" db:"created_at"` // ISO Datetime string
}

type PlayerRole string

const (
	RoleAdmin  PlayerRole = "admin"
	RoleUser   PlayerRole = "user"
	RolePlayer PlayerRole = "player"
)

func (r PlayerRole) IsValid() bool {
	return r == RoleAdmin || r == RoleUser || r == RolePlayer
}

// Placement represents the placement in an award (1st or 2nd place)
type Placement int

const (
	PlacementFirst  Placement = 1 // Gold
	PlacementSecond Placement = 2 // Silver
)

func (p Placement) IsValid() bool {
	return p == PlacementFirst || p == PlacementSecond
}

type Player struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	PasswordHash *string    `json:"-" db:"password_hash"` // Never return in JSON
	Role         PlayerRole `json:"role" db:"role"`
	CreatedBy    *int       `json:"created_by" db:"created_by"` // Foreign key to player.id, nil for initial admin
	CreatedAt    string     `json:"created_at" db:"created_at"` // ISO datetime string
	UpdatedAt    string     `json:"updated_at" db:"updated_at"` // ISO datetime string
}

type PlayerExtendedInfo struct {
	Player           Player `json:"player"`
	TotalGamesPlayed int    `json:"total_games_played"`
	TotalGamesWon    int    `json:"total_games_won"`
}

type GamePlayer struct {
	ID                 int     `json:"id" db:"id"`
	GameID             int     `json:"game_id" db:"game_id"`
	PlayerID           int     `json:"player_id" db:"player_id"`
	Corporation        *string `json:"corporation" db:"corporation"` // nil for legacy games
	TerraformingRating int     `json:"terraforming_rating" db:"terraforming_rating"`
	Cities             int     `json:"cities" db:"cities"`
	Greeneries         int     `json:"greeneries" db:"greeneries"`
	Cards              int     `json:"cards" db:"cards"`
	TurmoilPoints      int     `json:"turmoil_points" db:"turmoil_points"`
	MilestonePoints    int     `json:"milestone_points" db:"milestone_points"`
	AwardPoints        int     `json:"award_points" db:"award_points"`
	TotalPoints        int     `json:"total_points" db:"total_points"`
}

type Milestone struct {
	ID                 int    `json:"id" db:"id"`
	GameID             int    `json:"game_id" db:"game_id"`
	Name               string `json:"name" db:"name"`
	WinnerGamePlayerID *int   `json:"winner_game_player_id" db:"winner_game_player_id"` // references GamePlayer.ID
}

type Award struct {
	ID     int    `json:"id" db:"id"`
	GameID int    `json:"game_id" db:"game_id"`
	Name   string `json:"name" db:"name"`
}

type AwardPlacement struct {
	ID           int       `json:"id" db:"id"`
	AwardID      int       `json:"award_id" db:"award_id"`
	GamePlayerID int       `json:"game_player_id" db:"game_player_id"`
	Placement    Placement `json:"placement" db:"placement"`
}

type GameImage struct {
	ID           int    `json:"id" db:"id"`
	GameID       int    `json:"game_id" db:"game_id"`
	ImageData    []byte `json:"image_data" db:"image_data"`
	MimeType     string `json:"mime_type" db:"mime_type"`
	DisplayOrder int    `json:"display_order" db:"display_order"`
	UploadedAt   string `json:"uploaded_at" db:"uploaded_at"`
}

// Image metadata without actual image
type GameImageMeta struct {
	ID           int    `json:"id" db:"id"`
	DisplayOrder int    `json:"display_order" db:"display_order"`
	MimeType     string `json:"mime_type" db:"mime_type"`
	UploadedAt   string `json:"uploaded_at" db:"uploaded_at"`
}

// Expansions type for JSON storage in SQLite
type Expansions map[string]bool

// GameWithDetails combines game with related data for API responses
type GameWithDetails struct {
	Game        Game             `json:"game"`
	GamePlayers []GamePlayer     `json:"game_players"`
	Players     []Player         `json:"players"`
	Milestones  []Milestone      `json:"milestones"`
	Awards      []Award          `json:"awards"`
	Placements  []AwardPlacement `json:"award_placements"`
	Images      []GameImageMeta  `json:"images"`
}

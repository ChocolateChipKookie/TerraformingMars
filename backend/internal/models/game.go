package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Game struct {
	ID        int       `json:"id" db:"id"`
	GameUUID  string    `json:"game_uuid" db:"game_uuid"`
	Revision  int       `json:"revision" db:"revision"`
	Name      string    `json:"name" db:"name"`
	Date      time.Time `json:"date" db:"date"`
	Map       string    `json:"map" db:"map"`
	Generations int     `json:"generations" db:"generations"`
	Expansions  Expansions `json:"expansions" db:"expansions"`
	IsLatest    bool      `json:"is_latest" db:"is_latest"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Player struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type GamePlayer struct {
	ID                 int    `json:"id" db:"id"`
	GameID             int    `json:"game_id" db:"game_id"`
	PlayerID           int    `json:"player_id" db:"player_id"`
	Corporation        string `json:"corporation" db:"corporation"`
	TerraformingRating int    `json:"terraforming_rating" db:"terraforming_rating"`
	Cities             int    `json:"cities" db:"cities"`
	Greeneries         int    `json:"greeneries" db:"greeneries"`
	Cards              int    `json:"cards" db:"cards"`
	TurmoilPoints      int    `json:"turmoil_points" db:"turmoil_points"`
	MilestonePoints    int    `json:"milestone_points" db:"milestone_points"`
	AwardPoints        int    `json:"award_points" db:"award_points"`
	TotalPoints        int    `json:"total_points" db:"total_points"`
}

type Milestone struct {
	ID                  int    `json:"id" db:"id"`
	GameID              int    `json:"game_id" db:"game_id"`
	Name                string `json:"name" db:"name"`
	WinnerGamePlayerID  *int   `json:"winner_game_player_id" db:"winner_game_player_id"` // references GamePlayer.ID
}

type Award struct {
	ID     int    `json:"id" db:"id"`
	GameID int    `json:"game_id" db:"game_id"`
	Name   string `json:"name" db:"name"`
}

type AwardPlacement struct {
	ID           int `json:"id" db:"id"`
	AwardID      int `json:"award_id" db:"award_id"`
	GamePlayerID int `json:"game_player_id" db:"game_player_id"`
	Placement    int `json:"placement" db:"placement"` // 0=none, 1=gold, 2=silver
}

// Expansions type for JSON storage in SQLite
type Expansions map[string]bool

// Implement driver.Valuer interface for storing in database
func (e Expansions) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Implement sql.Scanner interface for reading from database
func (e *Expansions) Scan(value interface{}) error {
	if value == nil {
		*e = make(Expansions)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, e)
	case string:
		return json.Unmarshal([]byte(v), e)
	default:
		return fmt.Errorf("cannot scan %T into Expansions", value)
	}
}

// GameWithDetails combines game with related data for API responses
type GameWithDetails struct {
	Game        Game             `json:"game"`
	GamePlayers []GamePlayer     `json:"game_players"`
	Players     []Player         `json:"players"`
	Milestones  []Milestone      `json:"milestones"`
	Awards      []Award          `json:"awards"`
	Placements  []AwardPlacement `json:"award_placements"`
}

// CreateGameRequest represents the request body for creating a game
type CreateGameRequest struct {
	Name        string     `json:"name"`
	Date        string     `json:"date"` // Will be parsed to time.Time
	Map         string     `json:"map"`
	Generations int        `json:"generations"`
	Expansions  Expansions `json:"expansions"`
	Players     []struct {
		Name               string `json:"name"`
		Corporation        string `json:"corporation"`
		TerraformingRating int    `json:"terraforming_rating"`
		Cities             int    `json:"cities"`
		Greeneries         int    `json:"greeneries"`
		Cards              int    `json:"cards"`
		TurmoilPoints      int    `json:"turmoil_points"`
	} `json:"players"`
	Milestones []struct {
		Name                  string `json:"name"`
		WinnerGamePlayerIndex *int   `json:"winner_game_player_index"` // Index in the Players array
	} `json:"milestones"`
	Awards []struct {
		Name       string `json:"name"`
		Placements []struct {
			PlayerIndex int `json:"player_index"`
			Placement   int `json:"placement"`
		} `json:"placements"`
	} `json:"awards"`
}
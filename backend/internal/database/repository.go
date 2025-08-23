package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	
	"terraforming-mars-backend/internal/models"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreatePlayer creates a new player
func (r *Repository) CreatePlayer(name string) (*models.Player, error) {
	result, err := r.db.Exec("INSERT INTO players (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Player{
		ID:   int(id),
		Name: name,
	}, nil
}

// GetPlayerByName retrieves a player by name
func (r *Repository) GetPlayerByName(name string) (*models.Player, error) {
	var player models.Player
	err := r.db.QueryRow("SELECT id, name FROM players WHERE name = ?", name).Scan(&player.ID, &player.Name)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// GetAllPlayers retrieves all players
func (r *Repository) GetAllPlayers() ([]models.Player, error) {
	rows, err := r.db.Query("SELECT id, name FROM players ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

// CreateGame creates a new game with all related data
func (r *Repository) CreateGame(req models.CreateGameRequest) (*models.GameWithDetails, error) {
	// Validate that we have matching players and scores
	if len(req.Players) != len(req.PlayerScores) {
		return nil, fmt.Errorf("validation error: number of players (%d) must match number of player scores (%d). Each player must have a corresponding score entry", 
			len(req.Players), len(req.PlayerScores))
	}

	// Validate that we have at least one player
	if len(req.Players) == 0 {
		return nil, fmt.Errorf("validation error: at least one player is required to create a game")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v (expected YYYY-MM-DD)", err)
	}

	// Marshal expansions to JSON
	expansionsJSON, err := json.Marshal(req.Expansions)
	if err != nil {
		return nil, err
	}

	// Create game with first revision
	gameUUID := uuid.New().String()
	gameResult, err := tx.Exec(`
		INSERT INTO games (game_uuid, revision, name, date, map, generations, expansions, is_latest)
		VALUES (?, 1, ?, ?, ?, ?, ?, true)
	`, gameUUID, req.Name, date, req.Map, req.Generations, string(expansionsJSON))
	if err != nil {
		return nil, err
	}

	gameID, err := gameResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create game_players entries
	var gamePlayers []models.GamePlayer
	for i, playerReq := range req.Players {
		// Get player by name (must exist)
		var playerID int
		err = tx.QueryRow("SELECT id FROM players WHERE name = ?", playerReq.Name).Scan(&playerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("player '%s' not found. Please create the player first", playerReq.Name)
			}
			return nil, fmt.Errorf("error finding player '%s': %v", playerReq.Name, err)
		}

		// Get scores for this player
		scores := req.PlayerScores[i]

		// Create game_player entry
		gamePlayerResult, err := tx.Exec(`
			INSERT INTO game_players (
				game_id, player_id, corporation, 
				terraforming_rating, cities, greeneries, cards, turmoil_points,
				milestone_points, award_points, total_points
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0)
		`, gameID, playerID, playerReq.Corporation, 
			scores.TerraformingRating, scores.Cities, scores.Greeneries, 
			scores.Cards, scores.TurmoilPoints)
		if err != nil {
			return nil, err
		}

		gamePlayerID, err := gamePlayerResult.LastInsertId()
		if err != nil {
			return nil, err
		}

		gamePlayers = append(gamePlayers, models.GamePlayer{
			ID:                 int(gamePlayerID),
			GameID:             int(gameID),
			PlayerID:           playerID,
			Corporation:        playerReq.Corporation,
			TerraformingRating: scores.TerraformingRating,
			Cities:             scores.Cities,
			Greeneries:         scores.Greeneries,
			Cards:              scores.Cards,
			TurmoilPoints:      scores.TurmoilPoints,
			MilestonePoints:    0,
			AwardPoints:        0,
			TotalPoints:        0,
		})
	}

	// Create milestones
	var milestones []models.Milestone
	for _, milestoneReq := range req.Milestones {
		// Validate winner_game_player_index if provided
		var winnerGamePlayerID *int
		if milestoneReq.WinnerGamePlayerIndex != nil {
			// WinnerGamePlayerIndex in the request is the index of the player in the req.Players array
			if *milestoneReq.WinnerGamePlayerIndex < 0 || *milestoneReq.WinnerGamePlayerIndex >= len(gamePlayers) {
				return nil, fmt.Errorf("invalid winner_game_player_index %d for milestone '%s': must be between 0 and %d", 
					*milestoneReq.WinnerGamePlayerIndex, milestoneReq.Name, len(gamePlayers)-1)
			}
			winnerGamePlayerID = &gamePlayers[*milestoneReq.WinnerGamePlayerIndex].ID
		}

		result, err := tx.Exec(`
			INSERT INTO milestones (game_id, name, winner_game_player_id)
			VALUES (?, ?, ?)
		`, gameID, milestoneReq.Name, winnerGamePlayerID)
		if err != nil {
			return nil, fmt.Errorf("error creating milestone '%s': %v", milestoneReq.Name, err)
		}

		milestoneID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		milestones = append(milestones, models.Milestone{
			ID:                 int(milestoneID),
			GameID:             int(gameID),
			Name:               milestoneReq.Name,
			WinnerGamePlayerID: winnerGamePlayerID,
		})
	}

	// Create awards and their placements
	var awards []models.Award
	var placements []models.AwardPlacement
	for _, awardReq := range req.Awards {
		// Create the award
		result, err := tx.Exec(`
			INSERT INTO awards (game_id, name)
			VALUES (?, ?)
		`, gameID, awardReq.Name)
		if err != nil {
			return nil, fmt.Errorf("error creating award '%s': %v", awardReq.Name, err)
		}

		awardID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		awards = append(awards, models.Award{
			ID:     int(awardID),
			GameID: int(gameID),
			Name:   awardReq.Name,
		})

		// Create award placements
		for _, placementReq := range awardReq.Placements {
			// Validate player index
			if placementReq.PlayerIndex < 0 || placementReq.PlayerIndex >= len(gamePlayers) {
				return nil, fmt.Errorf("invalid player_index %d for award '%s': must be between 0 and %d", 
					placementReq.PlayerIndex, awardReq.Name, len(gamePlayers)-1)
			}

			// Validate placement value (0=none, 1=gold, 2=silver)
			if placementReq.Placement < 0 || placementReq.Placement > 2 {
				return nil, fmt.Errorf("invalid placement %d for award '%s': must be 0 (none), 1 (gold), or 2 (silver)", 
					placementReq.Placement, awardReq.Name)
			}

			// Only create placement if it's not "none" (0)
			if placementReq.Placement > 0 {
				placementResult, err := tx.Exec(`
					INSERT INTO award_placements (award_id, game_player_id, placement)
					VALUES (?, ?, ?)
				`, awardID, gamePlayers[placementReq.PlayerIndex].ID, placementReq.Placement)
				if err != nil {
					return nil, fmt.Errorf("error creating award placement: %v", err)
				}

				placementID, err := placementResult.LastInsertId()
				if err != nil {
					return nil, err
				}

				placements = append(placements, models.AwardPlacement{
					ID:           int(placementID),
					AwardID:      int(awardID),
					GamePlayerID: gamePlayers[placementReq.PlayerIndex].ID,
					Placement:    placementReq.Placement,
				})
			}
		}
	}

	// Calculate and update points for each game player
	for i := range gamePlayers {
		gp := &gamePlayers[i]
		
		// Calculate milestone points (5 points per milestone won)
		milestonePoints := 0
		for _, milestone := range milestones {
			if milestone.WinnerGamePlayerID != nil && *milestone.WinnerGamePlayerID == gp.ID {
				milestonePoints += 5
			}
		}
		
		// Calculate award points (5 for gold, 2 for silver)
		awardPoints := 0
		for _, placement := range placements {
			if placement.GamePlayerID == gp.ID {
				switch placement.Placement {
				case 1: // Gold
					awardPoints += 5
				case 2: // Silver
					awardPoints += 2
				}
			}
		}
		
		// Calculate total points
		totalPoints := gp.TerraformingRating + gp.Cities + gp.Greeneries + 
			gp.Cards + gp.TurmoilPoints + milestonePoints + awardPoints
		
		// Update the game_player record with calculated points
		_, err := tx.Exec(`
			UPDATE game_players 
			SET milestone_points = ?, award_points = ?, total_points = ?
			WHERE id = ?
		`, milestonePoints, awardPoints, totalPoints, gp.ID)
		if err != nil {
			return nil, fmt.Errorf("error updating points for game_player %d: %v", gp.ID, err)
		}
		
		// Update the struct for consistency
		gp.MilestonePoints = milestonePoints
		gp.AwardPoints = awardPoints
		gp.TotalPoints = totalPoints
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the created game with all its details
	return r.GetGameByID(int(gameID))
}

// GetGameByID retrieves a specific game revision by ID with all related data
func (r *Repository) GetGameByID(gameID int) (*models.GameWithDetails, error) {
	// Get the game
	var game models.Game
	var expansionsJSON string
	err := r.db.QueryRow(`
		SELECT id, game_uuid, revision, name, date, map, generations, expansions, is_latest, created_at
		FROM games 
		WHERE id = ?
	`, gameID).Scan(
		&game.ID, &game.GameUUID, &game.Revision, &game.Name, &game.Date,
		&game.Map, &game.Generations, &expansionsJSON, &game.IsLatest, &game.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", gameID)
		}
		return nil, err
	}

	// Parse expansions JSON
	if err := json.Unmarshal([]byte(expansionsJSON), &game.Expansions); err != nil {
		return nil, fmt.Errorf("error parsing expansions JSON: %v", err)
	}

	// Get game players and their associated player info
	rows, err := r.db.Query(`
		SELECT 
			gp.id, gp.game_id, gp.player_id, gp.corporation,
			gp.terraforming_rating, gp.cities, gp.greeneries, gp.cards,
			gp.turmoil_points, gp.milestone_points, gp.award_points, gp.total_points,
			p.id, p.name
		FROM game_players gp
		JOIN players p ON gp.player_id = p.id
		WHERE gp.game_id = ?
		ORDER BY gp.total_points DESC
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gamePlayers []models.GamePlayer
	var players []models.Player
	playersSeen := make(map[int]bool)
	
	for rows.Next() {
		var gp models.GamePlayer
		var p models.Player
		err := rows.Scan(
			&gp.ID, &gp.GameID, &gp.PlayerID, &gp.Corporation,
			&gp.TerraformingRating, &gp.Cities, &gp.Greeneries, &gp.Cards,
			&gp.TurmoilPoints, &gp.MilestonePoints, &gp.AwardPoints, &gp.TotalPoints,
			&p.ID, &p.Name,
		)
		if err != nil {
			return nil, err
		}
		
		gamePlayers = append(gamePlayers, gp)
		
		// Add player to list if not already added
		if !playersSeen[p.ID] {
			players = append(players, p)
			playersSeen[p.ID] = true
		}
	}

	// Get milestones
	milestoneRows, err := r.db.Query(`
		SELECT id, game_id, name, winner_game_player_id
		FROM milestones 
		WHERE game_id = ?
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer milestoneRows.Close()

	var milestones []models.Milestone
	for milestoneRows.Next() {
		var m models.Milestone
		err := milestoneRows.Scan(&m.ID, &m.GameID, &m.Name, &m.WinnerGamePlayerID)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}

	// Get awards
	awardRows, err := r.db.Query(`
		SELECT id, game_id, name
		FROM awards 
		WHERE game_id = ?
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer awardRows.Close()

	var awards []models.Award
	for awardRows.Next() {
		var a models.Award
		err := awardRows.Scan(&a.ID, &a.GameID, &a.Name)
		if err != nil {
			return nil, err
		}
		awards = append(awards, a)
	}

	// Get award placements
	placementRows, err := r.db.Query(`
		SELECT ap.id, ap.award_id, ap.game_player_id, ap.placement
		FROM award_placements ap
		JOIN awards a ON ap.award_id = a.id
		WHERE a.game_id = ?
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer placementRows.Close()

	var placements []models.AwardPlacement
	for placementRows.Next() {
		var p models.AwardPlacement
		err := placementRows.Scan(&p.ID, &p.AwardID, &p.GamePlayerID, &p.Placement)
		if err != nil {
			return nil, err
		}
		placements = append(placements, p)
	}

	return &models.GameWithDetails{
		Game:        game,
		GamePlayers: gamePlayers,
		Players:     players,
		Milestones:  milestones,
		Awards:      awards,
		Placements:  placements,
	}, nil
}